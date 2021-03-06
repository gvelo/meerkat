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

package jsoningester

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/twmb/murmur3"
	hash2 "hash"
	"io"
	"meerkat/internal/cluster"
	"meerkat/internal/ingestion"
	"meerkat/internal/storage"
	"strconv"
	"time"
)

const (
	TSColName = "_ts"
)

type ParserError struct {
	Line   int
	Column string
	Error  string
}

type column struct {
	idx         int
	colType     storage.ColumnType
	keepParsing bool
	size        int
	len         int
}

func NewTable(name string) *Table {

	t := &Table{
		name:       name,
		partitions: make(map[int]*Partition),
	}

	return t

}

type Table struct {
	name       string
	partitions map[int]*Partition
}

func (t *Table) partition(partNum int) *Partition {

	if partition, ok := t.partitions[partNum]; ok {
		return partition
	}

	partition := &Partition{
		columns: make(map[string]*column),
		writer:  ingestion.NewRowSetWriter(4 * 1024),
	}

	partition.columns[TSColName] = &column{
		idx:         0,
		colType:     storage.ColumnType_TIMESTAMP,
		keepParsing: false,
	}

	t.partitions[partNum] = partition

	return partition

}

type Partition struct {
	colIdx    int
	columns   map[string]*column
	numOfRows int
	writer    *ingestion.RowSetWriter
}

func (p *Partition) column(name string) *column {

	if col, ok := p.columns[name]; ok {
		return col
	}

	p.colIdx++

	col := &column{
		idx:         p.colIdx,
		keepParsing: true,
	}

	p.columns[name] = col

	return col

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
	scanner := bufio.NewScanner(br)

	line := 0
	ingestedRows := 0

	for scanner.Scan() {

		line++

		decoder := json.NewDecoder(bytes.NewReader(scanner.Bytes()))
		decoder.UseNumber()

		m := make(map[string]interface{})

		err := decoder.Decode(&m)

		if err != nil {
			ingestionErrors = append(ingestionErrors, ParserError{
				Line:  line,
				Error: err.Error(),
			})
			continue
		}

		var ts int64

		if i, ok := m[TSColName]; ok {

			ts, err = parseTS(i)

			if err != nil {
				ingestionErrors = append(ingestionErrors, ParserError{
					Line:   line,
					Error:  err.Error(),
					Column: TSColName,
				})
				continue
			}

		} else {
			ts = time.Now().UnixNano()
		}

		partition := table.partition(ing.getPartition(ts, numOfPartitions))

		tsCol := partition.columns[TSColName]
		partition.writer.WriteFixedUInt64(0, uint64(ts))
		tsCol.len++
		partition.numOfRows++

		delete(m, TSColName)

		for colName, colValue := range m {

			col := partition.column(colName)

			switch v := colValue.(type) {
			case string:
				if col.keepParsing {
					_, err := time.Parse(time.RFC3339Nano, v)
					if err != nil {
						col.keepParsing = false
					} else {
						col.colType = storage.ColumnType_TIMESTAMP
					}
				}
				partition.writer.WriteString(col.idx, v)
				col.size += len(v)
				col.len++
			case json.Number:
				s := string(v)
				if col.keepParsing {
					_, err := strconv.Atoi(s)
					if err != nil {
						col.keepParsing = false
						col.colType = storage.ColumnType_STRING
					} else {
						col.colType = storage.ColumnType_INT64
					}

				}
				partition.writer.WriteString(col.idx, s)
				col.size += len(s)
				col.len++
			default:
				panic("unknown type")
			}

		}

		ingestedRows++

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

type Ingester interface {
	Ingest(stream io.Reader, tableName string) []ParserError
}

func NewIngester(rpc IngesterRpc, nodeReg cluster.NodeRegistry, bufferReg ingestion.BufferRegistry) Ingester {
	return &ingester{
		rpc:            rpc,
		nodeReg:        nodeReg,
		indexBufferReg: bufferReg,
	}
}

type ingester struct {
	nodeReg        cluster.NodeRegistry
	rpc            IngesterRpc
	indexBufferReg ingestion.BufferRegistry
}

func (ing *ingester) Ingest(stream io.Reader, tableName string) []ParserError {

	// TODO(gvelo): all this ingest implementation is for testing purposes.

	// TODO(gvelo) add a partition router.

	nodes := ing.nodeReg.Nodes([]string{cluster.Joining, cluster.Ready}, true)
	numOfPartitions := len(nodes) + 1 // num of members plus local node.

	parser := NewParser()

	table, ingestedRows, pErr := parser.Parse(stream, tableName, numOfPartitions)

	if ingestedRows == 0 {
		return pErr
	}

	pbTable := CreatePBTable(table)

	// first partition goes to the local node
	ing.indexBufferReg.AddToBuffer(&ingestion.Table{
		Name:       tableName,
		Partitions: pbTable.Partitions[:1],
	})

	for i, node := range nodes {

		err := ing.rpc.SendRequest(context.TODO(), node.Id(), &ingestion.IngestionRequest{
			Table: &ingestion.Table{
				Name:       tableName,
				Partitions: pbTable.Partitions[i+1 : i+2],
			}})

		if err != nil {
			panic(err)
		}

	}

	return pErr

}

func CreatePBTable(table *Table) *ingestion.Table {

	var pbPartitions []*ingestion.Partition

	for id, partition := range table.partitions {

		p := &ingestion.Partition{
			Id:   uint64(id),
			Data: partition.writer.Buf.Data(),
		}

		for colName, col := range partition.columns {
			pbCol := &ingestion.Column{
				Idx:     uint64(col.idx),
				Name:    colName,
				ColSize: uint64(col.size),
				Len:     uint64(col.len),
				Type:    col.colType,
			}
			p.Columns = append(p.Columns, pbCol)
		}

		pbPartitions = append(pbPartitions, p)

	}

	return &ingestion.Table{
		Name:       table.name,
		Partitions: pbPartitions,
	}

}
