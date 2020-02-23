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
	"meerkat/internal/storage"
	"meerkat/internal/storage/io"
)

const (
	maxSlicesPerBlock = 1024 * 2
)

type ByteSlicePlainEncoder struct {
	bw        storage.BlockWriter
	buf       *io.EncoderBuffer
	offsetBuf []int
}

func NewByteSlicePlainEncodeer(bw storage.BlockWriter) *ByteSlicePlainEncoder {
	return &ByteSlicePlainEncoder{
		bw:        bw,
		buf:       io.NewEncoderBuffer(64 * 1024),
		offsetBuf: make([]int, maxSlicesPerBlock),
	}
}

func (e *ByteSlicePlainEncoder) Flush() error {
	return nil
}

func (e *ByteSlicePlainEncoder) FlushBlocks() error {
	return e.bw.Flush()
}

func (e *ByteSlicePlainEncoder) Type() storage.EncodingType {
	return storage.Plain
}

func (e *ByteSlicePlainEncoder) Encode(vec storage.ByteSliceVector) error {

	size := binary.MaxVarintLen64*(vec.Len()+2) + len(vec.Data())

	e.buf.Reset(size)

	DeltaEncode(vec.Offsets(), e.offsetBuf)

	// left enough room at the beginning of the block to write the
	// block length encoded as uvarint
	e.buf.WriteVarUintSliceAt(binary.MaxVarintLen64, e.offsetBuf[:vec.Len()])

	// write the vector data
	e.buf.WriteBytes(vec.Data())

	blockSize := e.buf.Len() - binary.MaxVarintLen64

	offset := binary.MaxVarintLen64 - io.SizeUVarint(uint64(blockSize))

	e.buf.WriteUvarintAt(offset, blockSize)

	block := e.buf.Bytes()[offset:]

	return e.bw.WriteBlock(block, vec.Rid()[0])

}
