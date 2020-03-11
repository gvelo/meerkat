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
	"meerkat/internal/storage/vector"
)

type ColumnSource interface {
	HasNext() bool
	HasNulls() bool
}

type IntColumSource interface {
	ColumnSource
	Next() vector.IntVector
}

type UintColumSource interface {
	ColumnSource
	Next() vector.IntVector
}

type FloatColumSource interface {
	ColumnSource
	Next() vector.FloatVector
}

type ByteSliceColumSource interface {
	ColumnSource
	Next() vector.ByteSliceVector
}

func NewIntColumnSource(buf *buffer.IntBuffer, dstSize int, permMap []int) IntColumSource {

	return &intColumnSource{
		srcBuf:   buf.Values(),
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
	nulls    []uint64
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
func (cs *intColumnSource) Next() vector.IntVector {

	var i int

	for i = 0; i < len(cs.dstBuf) && cs.pos < len(cs.srcBuf); i++ {

		j := cs.permMap[cs.pos]

		// TODO:(sebad) sacar cuando este implementado.
		if cs.HasNulls() {
			i--
		} else {
			cs.dstBuf[i] = cs.srcBuf[j]
			cs.rid[i] = uint32(cs.pos)
		}

		cs.pos++

	}

	return vector.NewIntVector(cs.dstBuf[:i], []uint64{})

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

func (cs *tsColumnSource) Next() vector.IntVector {

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

	return vector.NewIntVector(cs.srcBuf[cs.start:cs.end], []uint64{})
}

func NewFloatColumnSource(buf *buffer.FloatBuffer, dstSize int, permMap []int) FloatColumSource {

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
	nulls    []uint64
	permMap  []int
	rid      []uint32
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
func (cs *floatColumnSource) Next() vector.FloatVector {

	var i int

	for i = 0; i < len(cs.dstBuf) && cs.pos < len(cs.srcBuf); i++ {

		j := cs.permMap[cs.pos]

		// TODO:(sebad) sacar cuando este implementado.
		if cs.HasNulls() {
			i--
		} else {
			cs.dstBuf[i] = cs.srcBuf[j]
			cs.rid[i] = uint32(cs.pos)
		}

		cs.pos++

	}

	return vector.NewFloatVector(cs.dstBuf[:i], []uint64{})

}

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
	nulls      []uint64
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

func (cs *byteSliceColumnSource) Next() vector.ByteSliceVector {

	cs.dstOffsets = cs.dstOffsets[0:0]
	cs.rid = cs.rid[0:0]

	size := 0

	for ; cs.pos < cs.bs.Len(); cs.pos++ {

		j := cs.permMap[cs.pos]

		// TODO:(sebad) sacar cuando este implementado.
		if cs.HasNulls() {
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

	return vector.NewByteSliceVector(cs.dstBuf[:size], nil, cs.dstOffsets)

}
