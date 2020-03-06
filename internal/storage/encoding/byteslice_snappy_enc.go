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
	buf       *io.EncoderBuffer
	offsetBuf []int
}

func NewByteSliceSnappyEncodeer(bw BlockWriter) *ByteSliceSnappyEncoder {
	return &ByteSliceSnappyEncoder{
		bw:        bw,
		buf:       io.NewEncoderBuffer(64 * 1024),
		offsetBuf: make([]int, maxSlicesPerBlock),
	}
}

func (e *ByteSliceSnappyEncoder) Flush() {
}

func (e *ByteSliceSnappyEncoder) FlushBlocks() {
}

func (e *ByteSliceSnappyEncoder) Type() EncodingType {
	return Snappy
}

func (e *ByteSliceSnappyEncoder) Encode(v colval.ByteSliceColValues) {

	// make sure that the buffer has enough space to accommodate
	// the offsets slice plus the encoded data. We need to avoid
	// allocation inside the snappy encoder.
	size := binary.MaxVarintLen64*(v.Len()+2) + snappy.MaxEncodedLen(len(v.Data()))

	e.buf.Reset(size)

	DeltaEncode(v.Offsets(), e.offsetBuf)

	// left enough room at the beginning of the block to write the
	// block length encoded as uvarint
	e.buf.WriteVarUintSliceAt(binary.MaxVarintLen64, e.offsetBuf[:v.Len()])

	// as we have enough room in the dst buffer the encoder will not
	// allocate a new slice. See MaxEncodedLen(srcLen int) int
	r := snappy.Encode(e.buf.Free(), v.Data())

	e.buf.SetLen(e.buf.Len() + len(r))

	blockSize := e.buf.Len() - binary.MaxVarintLen64

	offset := binary.MaxVarintLen64 - io.SizeUVarint(uint64(blockSize))

	e.buf.WriteUvarintAt(offset, blockSize)

	block := e.buf.Bytes()[offset:]

	e.bw.WriteBlock(block, v.Rid()[0])

}
