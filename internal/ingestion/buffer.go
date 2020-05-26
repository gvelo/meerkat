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
	"meerkat/internal/jsoningester/ingestionpb"
	"meerkat/internal/schema"
	iobuff "meerkat/internal/storage/io"
)

type TSBuffer struct {
	buf []int64
	len int
}

func NewTSBuffer(cap int) *TSBuffer {
	return &TSBuffer{
		buf: make([]int64, cap),
		len: 0,
	}
}

func (b *TSBuffer) Append(i int64) {
	b.buf[b.len] = i
	b.len++
}

func (b *TSBuffer) Reserve(size int) {
	if b.len+size > len(b.buf) {
		buf := make([]int64, (len(b.buf)+size)*2)
		copy(buf, b.buf)
		b.buf = buf
	}
}

func (b *TSBuffer) Values() []int64 {
	return b.buf[:b.len]
}

func (b *TSBuffer) Len() int {
	return b.len
}

type ByteSliceSparseBuffer struct {
	buf     []byte
	offsets []uint32
	rowNums []uint32
	len     int
	size    int
}

func NewByteSliceSparseBuffer(numOfRows int, size int) *ByteSliceSparseBuffer {
	return &ByteSliceSparseBuffer{
		buf:     make([]byte, size),
		offsets: make([]uint32, numOfRows),
		rowNums: make([]uint32, numOfRows),
		len:     0,
		size:    0,
	}
}

func (b *ByteSliceSparseBuffer) Reserve(l int, size int) {

	if l+b.len > len(b.rowNums) {
		c := (len(b.rowNums) + l) * 2
		rowNums := make([]uint32, c)
		offsets := make([]uint32, c)
		copy(rowNums, b.rowNums[:b.len])
		copy(offsets, b.offsets[:b.len])
		b.rowNums = rowNums
		b.offsets = offsets
	}

	if b.size+size > len(b.buf) {
		buf := make([]byte, (len(b.buf)+size)*2)
		copy(buf, b.buf[:b.size])
		b.buf = buf
	}

}

func (b *ByteSliceSparseBuffer) Append(row uint32, v []byte) {
	b.rowNums[b.len] = row
	copy(b.buf[b.size:], v)
	b.size += len(v)
	b.offsets[b.len] = uint32(b.size) // TODO(gvelo): check b.size range
	b.len++
}

func (b *ByteSliceSparseBuffer) ToDenseBuffer(l int) *ByteSliceDenseBuffer {

	// the buffer is actually dense.
	if l == b.len {

		valids := make([]bool, l)

		// TODO(gvelo): memset ??
		for i := 0; i < l; i++ {
			valids[i] = true
		}

		return &ByteSliceDenseBuffer{
			buf:      b.buf,
			offsets:  b.offsets,
			size:     b.size,
			len:      b.len,
			valid:    valids,
			hasNulls: false,
		}

	}

	valids := make([]bool, l)
	offsets := make([]uint32, l)

	for i := 0; i < b.len; i++ {

		rowNum := b.rowNums[i]

		valids[rowNum] = true

		offsets[rowNum] = b.offsets[i]

		if rowNum > 0 {
			var start uint32
			if i > 0 {
				start = b.offsets[i-1]
			}
			offsets[rowNum-1] = start
		}

	}

	return &ByteSliceDenseBuffer{
		buf:      b.buf,
		offsets:  offsets,
		size:     b.size,
		len:      l,
		valid:    valids,
		hasNulls: true,
	}

}

type ByteSliceDenseBuffer struct {
	buf      []byte
	offsets  []uint32
	size     int
	len      int
	valid    []bool // TODO(gvelo): should we use a bitmap here ?
	hasNulls bool
}

func (b *ByteSliceDenseBuffer) Reserve(l int, size int) {

	if l+b.len > len(b.offsets) {
		c := (len(b.offsets) + l) * 2
		offsets := make([]uint32, c)
		copy(offsets, b.offsets[:b.len])
		b.offsets = offsets
	}

	if b.size+size > len(b.buf) {
		buf := make([]byte, (len(b.buf)+size)*2)
		copy(buf, b.buf[:b.size])
		b.buf = buf
	}

}

func (b *ByteSliceDenseBuffer) Append(value []byte) {
	copy(b.buf[b.size:], value)
	b.size += len(value)
	b.offsets[b.len] = uint32(b.size) // TODO(gvelo): check b.size range
	b.len++
}

func (b *ByteSliceDenseBuffer) Value(rowNum uint32) []byte {

	if int(rowNum) >= b.len {
		panic(" index out of range ")
	}

	var start uint32

	if rowNum > 0 {
		start = b.offsets[rowNum-1]
	}

	return b.buf[start:b.offsets[rowNum]]

}

func (b *ByteSliceDenseBuffer) Valids() []bool {
	return b.valid[:b.len]
}

func (b *ByteSliceDenseBuffer) HasNulls() bool {
	return b.hasNulls
}

type Float64SparseBuffer struct {
	buf     []float64
	rownums []uint32
	len     int
}

func NewFloat64SparseBuffer(len int) *Float64SparseBuffer {
	return &Float64SparseBuffer{
		buf:     make([]float64, len),
		rownums: make([]uint32, len),
	}
}

func (b *Float64SparseBuffer) Append(rowNum uint32, value float64) {
	b.buf[b.len] = value
	b.rownums[b.len] = rowNum
	b.len++
}

func (b *Float64SparseBuffer) Reserve(l int) {

	if b.len+l > len(b.buf) {
		c := (len(b.buf) + l) * 2
		buf := make([]float64, c)
		rowNums := make([]uint32, c)
		copy(buf, b.buf[:b.len])
		copy(rowNums, b.rownums[:b.len])
		b.buf = buf
		b.rownums = rowNums
	}

}

func (b *Float64SparseBuffer) ToDenseBuffer(l int) *Float64DenseBuffer {

	if l == b.len {

		valids := make([]bool, l)

		// TODO(gvelo): memset ??
		for i := 0; i < l; i++ {
			valids[i] = true
		}

		return &Float64DenseBuffer{
			buf:      b.buf,
			len:      b.len,
			valids:   valids,
			hasNulls: false,
		}
	}

	buf := make([]float64, l)
	valids := make([]bool, l)

	for i, f := range b.buf[:b.len] {
		buf[b.rownums[i]] = f
		valids[b.rownums[i]] = true
	}

	return &Float64DenseBuffer{
		buf:      buf,
		len:      l,
		valids:   valids,
		hasNulls: true,
	}

}

type Float64DenseBuffer struct {
	buf      []float64
	len      int ``
	valids   []bool
	hasNulls bool
}

func (b *Float64DenseBuffer) Append(value float64) {
	b.buf[b.len] = value
	b.len++
}

func (b *Float64DenseBuffer) Reserve(l int) {
	if b.len+l > len(b.buf) {
		c := (len(b.buf) + l) * 2
		buf := make([]float64, c)
		copy(buf, b.buf[:b.len])
		b.buf = buf
	}
}

func (b *Float64DenseBuffer) Valids() []bool {
	return b.valids[:b.len]
}

func (b *Float64DenseBuffer) Values() []float64 {
	return b.buf[:b.len]
}

type TableBufferX interface {
	// Schema() curretnly all the ingestion is schemaless.
	TableName() string
	PartitionId() int
	Columns() []*ingestionpb.Column
	TSBuffer() TSBuffer
	ColBuffer(colName string) interface{} // Always dense buffers
	Add(columns []*ingestionpb.Column, partition *ingestionpb.Partition)
	Len()
}

type ColumnBuffer struct {
	col  *ingestionpb.Column
	buff interface{}
}

type TableBuffer struct {
	partitionID uint64
	tableName   string
	columns     map[string]*ColumnBuffer
	// num of rows in the table.
	len int
}

func NewTableBuffer(tableName string, partitionID uint64) *TableBuffer {
	return &TableBuffer{
		partitionID: partitionID,
		tableName:   tableName,
		len:         0,
	}
}

func (b *TableBuffer) Columns() map[string]*ColumnBuffer {
	return b.columns
}

func (b *TableBuffer) Append(partition *ingestionpb.Partition) {

	if b.partitionID != partition.Id {
		panic("error appending to table buffer: wrong partitionID")
	}

	colBuffers := make([]interface{}, len(partition.Columns))

	for _, column := range partition.Columns {

		if c, found := b.columns[column.Name]; found {

			adaptType(c, column)

			colBuffers[column.Idx] = c.buff

			switch buff := c.buff.(type) {
			case *TSBuffer:
				buff.Reserve(int(column.Len))
			case *ByteSliceSparseBuffer:
				buff.Reserve(int(column.Len), int(column.ColSize))
			default:
				panic("invalid buffer type")
			}
			continue
		}

		var buff interface{}

		if column.Name == "_TS" {
			buff = NewTSBuffer(int(column.Len))
		} else {
			buff = NewByteSliceSparseBuffer(int(column.Len), int(column.ColSize))
		}

		colBuffers[column.Idx] = buff

		newCol := &ColumnBuffer{
			col:  column,
			buff: buff,
		}

		b.columns[column.Name] = newCol

	}

	// leemos la data

	if len(partition.Data) == 0 {
		panic("error: invalid data len")
	}

	dataBuff := iobuff.NewBufferWithData(partition.Data)
	var rowNum uint32

	// read the first TS

	colIdx := dataBuff.ReadUVarIntAsInt()

	if colIdx != 0 {
		panic("invalid column index")
	}

	ts := dataBuff.ReadFixedUInt64()
	tsBuff := colBuffers[0].(*TSBuffer)
	tsBuff.Append(int64(ts))

	for dataBuff.GetPos() < len(partition.Data) {

		colIdx := dataBuff.ReadUVarIntAsInt()

		// TS
		if colIdx == 0 {
			rowNum++
			ts := dataBuff.ReadFixedUInt64()
			tsBuff := colBuffers[0].(*TSBuffer)
			tsBuff.Append(int64(ts))
			continue
		}

		colBuff := colBuffers[colIdx].(*ByteSliceSparseBuffer)
		colBuff.Append(rowNum, dataBuff.ReadBytes())

	}

}

func adaptType(current *ColumnBuffer, new *ingestionpb.Column) {
	if current.col.Type != new.Type {
		current.col.Type = schema.ColumnType_STRING
	}
}
