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

package intake

import (
	"encoding/binary"
	"errors"
	"math"
)

const (
	BINARY byte = iota
	INT
	FLOAT
)

const (
	intSize     = 8
	colTypeSize = 1
	MaxVarintLen64 = 10
)

var ErrBufferFull = errors.New("bufio: buffer full")

// RowBuffer
type Buffer struct {
	buf     []byte
	offsets []int32
	types   []byte
	row     int
	offset  int
}

func NewBuffer(cap int, numOfRows int, types []byte) *Buffer {

	if len(types) > 256 {
		panic("too many columns")
	}

	t := make([]byte, 0, 256)
	t = append(t, types...)

	return &Buffer{
		buf:     make([]byte, cap),
		offsets: make([]int32, numOfRows),
		types:   t,
	}
}

func (b *Buffer) NewRow() {
	b.offsets[b.row] = int32(b.offset)
	b.row++
}

func (b *Buffer) AddInt(colIdx byte, value int64) error {

	if b.free() < binary.MaxVarintLen64+  {
		return ErrBufferFull
	}

	b.addColType(colIdx, INT)

	b.buf[b.offset] = colIdx
	b.offset++
	n := binary.PutVarint(b.buf, value)
	b.offset += n

	return nil

}

func (b *Buffer) AddFloat(colIdx byte, value float64) error {

	if b.free() < 8+1 {
		return ErrBufferFull
	}

	b.addColType(colIdx, INT)

	b.buf[b.offset] = colIdx
	b.offset++
	i := math.Float64bits(value)
	binary.LittleEndian.PutUint64(i)
	b.offset += 8

	return nil

}

func (b *Buffer) AddBinary(colIdx byte, value []byte) error {

	if b.free() < binary.MaxVarintLen64+len(value) {
		return ErrBufferFull
	}

	b.addColType(colIdx, BINARY)
	b.buf[b.offset] = colIdx
	b.offset++
	n := binary.PutVarint(b.buf, len(value))
	b.offset += int32(n)
	copy(b.buf, value)
	b.offset += int32(len(value))
	return nil

}

func (b *Buffer) GetInt(row int, colIndex byte) int {

}

func (b *Buffer) addColType(colIdx byte, colType byte) {

	if colIdx < byte(len(b.types)) {
		if b.types[colIdx] != colType {
			panic("column type doesn't match")
		}
	}

	b.types = append(b.types, colType)

}
func (b *Buffer) free() int {
	return len(b.buf) - int(b.offset)
}
