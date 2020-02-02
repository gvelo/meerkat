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

// page size * N

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

// vect valid until next call
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

//////////////////////// float

func NewFloatColumnSource(buf *buffer.Float64Buffer, dstSize int, permMap []int) FloatColumSource {

	return &floatColumnSource{
		srcBuf:   buf.Values(),
		dstBuf:   make([]int, dstSize),
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

///////////////////////////// byteslice

func NewByteSliceColumnSource(buff buffer.ByteSliceBuffer, maxSize int, permMap []int) ByteSliceColumSource {

	return &byteSliceColumnSource{}
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

	cs.dstBuf = cs.dstBuf[:cs.maxSize]
	cs.dstOffsets = cs.dstOffsets[0:0]
	cs.rid = cs.rid[0:0]

	dstStart := 0
	dstEnd := 0

	for ; cs.pos < cs.bs.Len(); cs.pos++ {

		j := cs.permMap[cs.pos]

		if cs.hasNulls && cs.nulls[j] {
			continue
		}

		slice := cs.bs.Get(cs.pos)

		dstEnd = dstStart + len(slice)

		if dstEnd > cs.maxSize {

			// if there aren't enough room to accommodate the slice
			// in an empty buffer, expand the buffer.
			if dstStart == 0 {
				cs.dstBuf = append(cs.dstBuf, slice...)
				cs.dstOffsets = append(cs.dstOffsets, len(cs.dstBuf))
				cs.rid = append(cs.rid, uint32(cs.pos))
			}

			break

		}

		copy(cs.dstBuf[dstStart:dstEnd], slice)
		cs.dstOffsets = append(cs.dstOffsets, dstEnd)
		cs.rid = append(cs.rid, uint32(cs.pos))
		dstStart = dstEnd

	}

	return byteSliceVector{
		rid:     cs.rid,
		data:    cs.dstBuf[:dstEnd],
		offsets: cs.dstOffsets,
	}

}
