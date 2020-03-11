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

package bufcolval

//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@../scalar_types.tmpldata bufcolval.gen.go.tmpl

import (
	"meerkat/internal/buffer"
	"meerkat/internal/storage/colval"
)

// TODO(gvelo) add a max amount of string per page ie 2048.
func NewByteSliceBufColSource(buff *buffer.ByteSliceBuffer, maxSize int, permMap []int) *ByteSliceBufColSource {

	return &ByteSliceBufColSource{
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

type ByteSliceBufColSource struct {
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

func (cs *ByteSliceBufColSource) HasNext() bool {
	return cs.pos < cs.bs.Len()
}

func (cs *ByteSliceBufColSource) HasNulls() bool {
	return cs.bs.Nullable()
}

func (cs *ByteSliceBufColSource) Next() colval.ByteSliceColValues {

	cs.dstOffsets = cs.dstOffsets[0:0]
	cs.rid = cs.rid[0:0]

	size := 0

	for ; cs.pos < cs.bs.Len(); cs.pos++ {

		j := cs.permMap[cs.pos]

		//TODO: fix
		if cs.HasNulls() && false {
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

	return colval.NewByteSliceColValues(cs.dstBuf[:size], cs.rid, cs.dstOffsets)

}

func NewTsBufColSource(buf *buffer.IntBuffer, dstSize int) *TSBufColSource {
	return &TSBufColSource{
		dstSize: dstSize,
		srcBuf:  buf.Values(),
		rid:     make([]uint32, dstSize),
	}
}

type TSBufColSource struct {
	srcBuf  []int
	rid     []uint32
	start   int
	end     int
	dstSize int
	pos     int
}

func (cs *TSBufColSource) HasNext() bool {
	return cs.end < len(cs.srcBuf)
}

func (cs *TSBufColSource) HasNulls() bool {
	// TS is never null
	return false
}

func (cs *TSBufColSource) Next() colval.IntColValues {

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

	return colval.NewIntColValues(cs.srcBuf[cs.start:cs.end], cs.rid[0:dstLen])

}
