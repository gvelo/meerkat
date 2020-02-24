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

package storage

import (
	"meerkat/internal/buffer"
)

type ColumnSource interface {
	HasNext() bool
	HasNulls() bool
}

type IntColumSource interface {
	ColumnSource
	Next() IntVector
}

type UintColumSource interface {
	ColumnSource
	Next() IntVector
}

type FloatColumSource interface {
	ColumnSource
	Next() FloatVector
}

type ByteSliceColumSource interface {
	ColumnSource
	Next() ByteSliceVector
}

type UUIDColumSource interface {
	ColumnSource
	Next() UUIDVector
}

func NewIntColumnSource(buf *buffer.IntBuffer, dstSize int, permMap []int) IntColumSource {

	return &intColumnSource{
		srcBuf:   buf.Int(),
		dstBuf:   make([]int, dstSize),
		nulls:    buf.Nulls(),
		rid:      make([]uint32, dstSize),
		permMap:  permMap,
		hasNulls: buf.Nullable(),
	}

}

type intColumnSource struct {
	srcBuf   []int
	dstBuf   []int
	nulls    []bool
	rid      []uint32
	permMap  []int
	hasNulls bool
	pos      int
}

func (cs *intColumnSource) HasNext() bool {
	return cs.pos < len(cs.srcBuf)
}

func (cs *intColumnSource) HasNulls() bool {
	return cs.hasNulls
}

// The underlying array point to an internal buffer that will be
// overwritten by a subsequent call to Next().
func (cs *intColumnSource) Next() IntVector {

	var i int

	for i = 0; i < len(cs.dstBuf) && cs.pos < len(cs.srcBuf); i++ {

		j := cs.permMap[cs.pos]

		if cs.hasNulls && cs.nulls[j] {
			i--
		} else {
			cs.dstBuf[i] = cs.srcBuf[j]
			cs.rid[i] = uint32(cs.pos)
		}

		cs.pos++

	}

	return intVector{
		vect: cs.dstBuf[:i],
		rid:  cs.rid[:i],
	}

}

func NewTsColumnSource(buf *buffer.IntBuffer, dstSize int) *tsColumnSource {
	return &tsColumnSource{
		dstSize: dstSize,
		srcBuf:  buf.Values(),
		rid:     make([]uint32, dstSize),
	}
}

type tsColumnSource struct {
	srcBuf  []int
	rid     []uint32
	start   int
	end     int
	dstSize int
	pos     int
}

func (cs *tsColumnSource) HasNext() bool {
	return cs.end < len(cs.srcBuf)
}

func (cs *tsColumnSource) HasNulls() bool {
	return false
}

func (cs *tsColumnSource) Next() IntVector {

	cs.start = cs.end
	cs.end = cs.start + cs.dstSize
	dstLen := cs.dstSize

	if cs.end > len(cs.srcBuf) {
		cs.end = len(cs.srcBuf)
		dstLen = cs.end - cs.start
	}

	for i := 0; i < dstLen; i++ {
		cs.rid[i] = uint32(cs.pos)
		cs.pos++
	}

	return intVector{
		vect: cs.srcBuf[cs.start:cs.end],
		rid:  cs.rid[0:dstLen],
	}

}

func NewFloatColumnSource(buf *buffer.Float64Buffer, dstSize int, permMap []int) FloatColumSource {

	return &floatColumnSource{
		srcBuf:   buf.Values(),
		dstBuf:   make([]float64, dstSize),
		nulls:    buf.Nulls(),
		rid:      make([]uint32, dstSize),
		permMap:  permMap,
		hasNulls: buf.Nullable(),
	}

}

type floatColumnSource struct {
	srcBuf   []float64
	dstBuf   []float64
	nulls    []bool
	rid      []uint32
	permMap  []int
	hasNulls bool
	pos      int
}

func (cs *floatColumnSource) HasNext() bool {
	return cs.pos < len(cs.srcBuf)
}

func (cs *floatColumnSource) HasNulls() bool {
	return cs.hasNulls
}

// vect valid until next call
func (cs *floatColumnSource) Next() FloatVector {

	var i int

	for i = 0; i < len(cs.dstBuf) && cs.pos < len(cs.srcBuf); i++ {

		j := cs.permMap[cs.pos]

		if cs.hasNulls && cs.nulls[j] {
			i--
		} else {
			cs.dstBuf[i] = cs.srcBuf[j]
			cs.rid[i] = uint32(cs.pos)
		}

		cs.pos++

	}

	return floatVector{
		vect: cs.dstBuf[:i],
		rid:  cs.rid[:i],
	}

}

// TODO(gvelo) add a max amount of string per page ie 2048.
func NewByteSliceColumnSource(buff *buffer.ByteSliceBuffer, maxSize int, permMap []int) ByteSliceColumSource {

	return &byteSliceColumnSource{
		bs:         buff,
		maxSize:    maxSize,
		dstBuf:     make([]byte, maxSize),
		dstOffsets: make([]int, 0, maxSize/8),
		nulls:      buff.Nulls(),
		rid:        make([]uint32, 0, maxSize/8),
		permMap:    permMap,
		hasNulls:   buff.Nullable(),
	}
}

type byteSliceColumnSource struct {
	bs         *buffer.ByteSliceBuffer
	maxSize    int
	dstBuf     []byte
	dstOffsets []int
	nulls      []bool
	rid        []uint32
	permMap    []int
	hasNulls   bool
	pos        int
}

func (cs *byteSliceColumnSource) HasNext() bool {
	return cs.pos < cs.bs.Len()
}

func (cs *byteSliceColumnSource) HasNulls() bool {
	return cs.bs.Nullable()
}

func (cs *byteSliceColumnSource) Next() ByteSliceVector {

	cs.dstOffsets = cs.dstOffsets[0:0]
	cs.rid = cs.rid[0:0]

	size := 0

	for ; cs.pos < cs.bs.Len(); cs.pos++ {

		j := cs.permMap[cs.pos]

		if cs.hasNulls && cs.nulls[j] {
			continue
		}

		slice := cs.bs.Get(j)

		available := cs.maxSize - size

		if len(slice) > available {

			// if there aren't enough room to accommodate the slice
			// in an empty buffer, allocate a new one.
			if size == 0 {
				if len(slice) > cap(cs.dstBuf) {
					cs.dstBuf = make([]byte, len(slice))
				}
				copy(cs.dstBuf[0:], slice)
				size = len(slice)
				cs.dstOffsets = append(cs.dstOffsets, size)
				cs.rid = append(cs.rid, uint32(cs.pos))
				cs.pos++
				break
			}

			break

		}

		copy(cs.dstBuf[size:], slice)
		size = size + len(slice)
		cs.dstOffsets = append(cs.dstOffsets, size)
		cs.rid = append(cs.rid, uint32(cs.pos))

	}

	return byteSliceVector{
		rid:     cs.rid,
		data:    cs.dstBuf[:size],
		offsets: cs.dstOffsets,
	}

}

// TODO(gvelo) uuid are fixed lenght. Create a specialized source.

func NewUUIDColumnSource(buff *buffer.UUIDBuffer, dstLen int, permMap []int) UUIDColumSource {

	return &uuidColumnSource{
		bs:       buff,
		dstLen:   dstLen,
		dstBuf:   make([]byte, dstLen*16),
		nulls:    buff.Nulls(),
		rid:      make([]uint32, 0, dstLen),
		permMap:  permMap,
		hasNulls: buff.Nullable(),
	}
}

type uuidColumnSource struct {
	bs       *buffer.UUIDBuffer
	dstLen   int
	dstBuf   []byte
	nulls    []bool
	rid      []uint32
	permMap  []int
	hasNulls bool
	pos      int
}

func (cs *uuidColumnSource) HasNext() bool {
	return cs.pos < cs.bs.Len()
}

func (cs *uuidColumnSource) HasNulls() bool {
	return cs.bs.Nullable()
}

func (cs *uuidColumnSource) Next() UUIDVector {

	cs.rid = cs.rid[0:0]

	l := 0

	for ; cs.pos < cs.bs.Len(); cs.pos++ {

		j := cs.permMap[cs.pos]

		if cs.hasNulls && cs.nulls[j] {
			continue
		}

		uuid := cs.bs.Get(j)

		copy(cs.dstBuf[l:], uuid)
		cs.rid = append(cs.rid, uint32(cs.pos))
		l += 16

		if l == len(cs.dstBuf) {
			break
		}

	}

	return byteSliceVector{
		rid:  cs.rid,
		data: cs.dstBuf[:l],
	}

}
