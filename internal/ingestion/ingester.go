// Copyright 2020 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ingestion

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/twmb/murmur3"
	hash2 "hash"
	"io"
	"meerkat/internal/cluster"
	"meerkat/internal/ingestion/ingestionpb"
	"meerkat/internal/schema"
	iobuff "meerkat/internal/storage/io"
	"strconv"
	"time"
)

//go:generate protoc  -I . -I ../../build/proto/ -I ../../internal/schema/   --plugin ../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc:. ./ingestionpb/ingester.proto

const (
	EOR       = 0x00
	TSColName = "_ts"
)

type ParserError struct {
	Line   int
	Column string
	Error  string
}

type Column struct {
	idx         int
	colType     schema.ColumnType
	keepParsing bool
}

func NewTable(name string) *Table {

	t := &Table{
		columns:    make(map[string]*Column),
		name:       name,
		partitions: make(map[int]*Partition),
		colIdx:     1, // 0 reserved for EOR
	}

	t.columns[TSColName] = &Column{
		idx:         1,
		colType:     schema.ColumnType_TIMESTAMP,
		keepParsing: false,
	}

	return t

}

type Table struct {
	columns    map[string]*Column
	name       string
	partitions map[int]*Partition
	colIdx     int
	numOfRows  int
}

func (t *Table) column(name string) *Column {

	if c, ok := t.columns[name]; ok {
		return c
	}

	t.colIdx++

	c := &Column{
		idx: t.colIdx,
	}

	t.columns[name] = c

	for _, p := range t.partitions {
		p.colSize = append(p.colSize, 0)
		p.colLen = append(p.colLen, 0)
	}

	return c

}

func (t *Table) partition(p int) *Partition {

	if partition, ok := t.partitions[p]; ok {
		return partition
	}

	partition := &Partition{
		colSize: make([]uint64, t.colIdx),
		colLen:  make([]uint64, t.colIdx),
		writer:  NewRowSetWriter(4 * 1024),
	}

	t.partitions[p] = partition

	return partition

}

type Partition struct {
	colSize []uint64
	colLen  []uint64
	writer  *RowSetWriter
}

type RowSetWriter struct {
	Buf *iobuff.Buffer
}

func NewRowSetWriter(cap int) *RowSetWriter {
	return &RowSetWriter{
		Buf: iobuff.NewBuffer(cap),
	}
}

func (rs *RowSetWriter) WriteString(colId int, str string) {

	rs.Reserve(len(str) + binary.MaxVarintLen64*2)

	rs.Buf.WriteIntAsUVarInt(colId)
	rs.Buf.WriteString(str)

}

func (rs *RowSetWriter) WriteInt(colId int, i int) {

	rs.Reserve(binary.MaxVarintLen64 * 2)

	rs.Buf.WriteIntAsUVarInt(colId)
	rs.Buf.WriteVarInt(i)

}

func (rs *RowSetWriter) WriteIntAsUVarInt(colId int, i int) {

	rs.Reserve(binary.MaxVarintLen64 * 2)

	rs.Buf.WriteIntAsUVarInt(colId)
	rs.Buf.WriteIntAsUVarInt(i)

}

func (rs *RowSetWriter) WriteFixedUInt64(colId int, i uint64) {

	rs.Reserve(binary.MaxVarintLen64 * 2)

	rs.Buf.WriteIntAsUVarInt(colId)
	rs.Buf.WriteFixedUInt64(i)

}

func (rs *RowSetWriter) WriteFloat(colId int, f float64) {

	rs.Reserve(binary.MaxVarintLen64 * 2)

	rs.Buf.WriteIntAsUVarInt(colId)
	rs.Buf.WriteFloat64(f)

}

func (rs *RowSetWriter) WriteEOR() {
	rs.Reserve(1)
	rs.Buf.WriteByte(EOR)
}

func (rs *RowSetWriter) Reserve(size int) {
	if size > rs.Buf.Available() {
		rs.Buf.Grow((rs.Buf.Cap() + size) * 2)
	}
}

func (rs *RowSetWriter) grow(newSize int) {
	rs.Buf.Grow(newSize)
}

func NewParser() *Parser {
	return &Parser{
		hash: murmur3.New64(),
	}
}

type Parser struct {
	hash hash2.Hash64
}

func (ing *Parser) Parse(reader io.Reader, tableName string, numOfPartitions int) (*Table, int, []ParserError) {

	var ingestionErrors []ParserError

	table := NewTable(tableName)

	br := bufio.NewReader(reader)
	decoder := json.NewDecoder(br)
	decoder.UseNumber()

	line := 0
	ingestedRows := 0

	for decoder.More() {

		line++

		m := make(map[string]interface{})

		err := decoder.Decode(&m)

		if err != nil {
			ingestionErrors = append(ingestionErrors, ParserError{
				Line:  line,
				Error: err.Error(),
			})
			// TODO(gvelo): seems like json.Decoder cannot recover
			//  from errors on the stream. We need a parser that
			//  can be able to recover errors skipping new lines.
			//continue
			return table, ingestedRows, ingestionErrors
		}

		var ts int64

		if i, ok := m[TSColName]; ok {

			ts, err = parseTS(i)

			if err != nil {
				ingestionErrors = append(ingestionErrors, ParserError{
					Line:  line,
					Error: err.Error(),
				})
				// TODO(gvelo): see prev TODO.
				//continue
				return table, ingestedRows, ingestionErrors
			}

		} else {
			ts = time.Now().UnixNano()
		}

		partition := table.partition(ing.getPartition(ts, numOfPartitions))
		tsCol := table.columns[TSColName]
		partition.writer.WriteFixedUInt64(tsCol.idx, uint64(ts))
		partition.colLen[tsCol.idx]++

		for colName, colValue := range m {

			if colName == TSColName {
				continue
			}

			col := table.column(colName)

			switch v := colValue.(type) {
			case string:
				if col.keepParsing {
					_, err := time.Parse(time.RFC3339Nano, v)
					if err != nil {
						col.keepParsing = false
					} else {
						col.colType = schema.ColumnType_TIMESTAMP
					}
				}
				partition.writer.WriteString(col.idx, v)
				partition.colSize[col.idx] += uint64(len(v))
				partition.colLen[col.idx]++
			case json.Number:
				s := string(v)
				if col.keepParsing {
					_, err := strconv.ParseFloat(s, 64)
					if err != nil {
						col.keepParsing = false
						col.colType = schema.ColumnType_STRING
					} else {
						col.colType = schema.ColumnType_REAL
					}

				}
				partition.writer.WriteString(col.idx, s)
				partition.colSize[col.idx] += uint64(len(s))
				partition.colLen[col.idx]++
			default:
				panic("unknown type")
			}

		}
		ingestedRows++
		partition.writer.WriteEOR()

	}

	return table, ingestedRows, ingestionErrors

}

func (ing *Parser) getPartition(ts int64, numOfPartiton int) int {
	ing.hash.Reset()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(ts))
	_, err := ing.hash.Write(b)
	if err != nil {
		panic(err)
	}
	h := ing.hash.Sum64()
	return int(h % uint64(numOfPartiton))
}

func parseTS(i interface{}) (int64, error) {

	s, ok := i.(string)

	if !ok {
		return -1, fmt.Errorf("cannot parse [%v] to date", i)
	}

	ts, err := strconv.Atoi(s)

	if err == nil {
		return int64(ts), nil
	}

	t, err := time.Parse(time.RFC3339Nano, s)

	if err != nil {
		return -1, err
	}

	return t.UnixNano(), nil

}

func NewInester(localNodeName string, conReg cluster.ConnRegistry, bufferReg IngestBufferRegistry) *Ingester {
	return &Ingester{
		conReg:         conReg,
		localNodeName:  localNodeName,
		indexBufferReg: bufferReg,
	}
}

type Ingester struct {
	conReg         cluster.ConnRegistry
	localNodeName  string
	indexBufferReg IngestBufferRegistry
}

func (ing *Ingester) Ingest(stream io.Reader, tableName string) []ParserError {

	// TODO(gvelo): only for testing purposes. Remove.
	m := ing.conReg.Members()
	numOfPartitions := len(m)

	parser := NewParser()

	table, ingestedRows, pErr := parser.Parse(stream, tableName, numOfPartitions)

	if ingestedRows == 0 {
		return pErr
	}

	var pbColumns []*ingestionpb.Column

	for name, column := range table.columns {

		c := &ingestionpb.Column{
			Name: name,
			Idx:  uint64(column.idx),
			Type: column.colType,
		}

		pbColumns = append(pbColumns, c)

	}

	var pbPartitions []*ingestionpb.Partition

	for id, partition := range table.partitions {

		p := &ingestionpb.Partition{
			Id:      id,
			ColSize: partition.colSize,
			ColLen:  partition.colLen,
			Data:    partition.writer.Buf.Data(),
		}

		pbPartitions = append(pbPartitions, p)

	}

	pbTable := &ingestionpb.Table{
		Name:       tableName,
		Columns:    pbColumns
		Partitions: pbPartitions,
	}

}
