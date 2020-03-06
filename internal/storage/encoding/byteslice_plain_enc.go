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
	"meerkat/internal/storage/colval"
	"meerkat/internal/storage/io"
)

const (
	maxSlicesPerBlock = 1024 * 32
)

type ByteSlicePlainEncoder struct {
	bw        BlockWriter
	buf       *io.EncoderBuffer
	offsetBuf []int
}

func NewByteSlicePlainEncodeer(bw BlockWriter) *ByteSlicePlainEncoder {
	return &ByteSlicePlainEncoder{
		bw:        bw,
		buf:       io.NewEncoderBuffer(64 * 1024),
		offsetBuf: make([]int, maxSlicesPerBlock),
	}
}

func (e *ByteSlicePlainEncoder) Flush() {
}

func (e *ByteSlicePlainEncoder) FlushBlocks() {
}

func (e *ByteSlicePlainEncoder) Type() EncodingType {
	return Plain
}

func (e *ByteSlicePlainEncoder) Encode(v colval.ByteSliceColValues) {

	size := binary.MaxVarintLen64*(v.Len()+2) + len(v.Data())

	e.buf.Reset(size)

	DeltaEncode(v.Offsets(), e.offsetBuf)

	// left enough room at the beginning of the block to write the
	// block length encoded as uvarint
	e.buf.WriteVarUintSliceAt(binary.MaxVarintLen64, e.offsetBuf[:v.Len()])

	// write the vector data
	e.buf.WriteBytes(v.Data())

	blockSize := e.buf.Len() - binary.MaxVarintLen64

	offset := binary.MaxVarintLen64 - io.SizeUVarint(uint64(blockSize))

	e.buf.WriteUvarintAt(offset, blockSize)

	block := e.buf.Bytes()[offset:]

	e.bw.WriteBlock(block, v.Rid()[0])

}
