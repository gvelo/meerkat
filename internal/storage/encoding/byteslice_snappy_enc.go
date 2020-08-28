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

package encoding

import (
	"encoding/binary"
	"github.com/golang/snappy"
	"meerkat/internal/storage/colval"
	"meerkat/internal/storage/io"
)

type ByteSliceSnappyEncoder struct {
	bw        BlockWriter
	buf       *io.Buffer
	offsetBuf []int
}

func NewByteSliceSnappyEncoder(bw BlockWriter) *ByteSliceSnappyEncoder {
	return &ByteSliceSnappyEncoder{
		bw:        bw,
		buf:       io.NewBuffer(64 * 1024),
		offsetBuf: make([]int, 1024*4),
	}
}

func (e *ByteSliceSnappyEncoder) Flush() {
}

func (e *ByteSliceSnappyEncoder) FlushBlocks() {
}

func (e *ByteSliceSnappyEncoder) Type() Type {
	return Snappy
}

func (e *ByteSliceSnappyEncoder) Encode(v colval.ByteSliceColValues) {

	// make sure that the buffer has enough space to accommodate
	// the offsets slice plus the encoded data. We need to avoid
	// allocation inside the snappy encoder.
	// ( header size + offset slice size ) * MaxVarintLen64 + enc data size
	size := binary.MaxVarintLen64*(v.Len()+2) + snappy.MaxEncodedLen(len(v.Data()))

	if size > e.buf.Cap() {
		// TODO(gvelo) check this grow policy.
		e.buf = io.NewBuffer(size + size/2)
	}

	e.buf.Reset()

	if v.Len() > len(e.offsetBuf) {
		e.offsetBuf = make([]int, v.Len()*2)
	}

	DeltaEncode(v.Offsets(), e.offsetBuf)

	// TODO(gvelo): framming is responsibility of the blockwriter so
	//  move this code outside this struct.

	// left enough room at the beginning of the block to write the
	// header ( block length encoded as uvarint )
	e.buf.Pos(binary.MaxVarintLen64)

	e.buf.WriteVarUintSlice(e.offsetBuf[:v.Len()])

	// as we have enough room in the dst buffer the encoder will not
	// allocate a new slice. See MaxEncodedLen(srcLen int) int
	r := snappy.Encode(e.buf.Free(), v.Data())

	blockEnd := e.buf.GetPos() + len(r)

	blockSize := blockEnd - binary.MaxVarintLen64

	headerOffset := binary.MaxVarintLen64 - io.SizeUVarint(uint64(blockSize))

	e.buf.PutIntAsUVarInt(headerOffset, blockSize)

	block := e.buf.Buf()[headerOffset:blockEnd]

	e.bw.WriteBlock(block, v.Rid()[0])

}
