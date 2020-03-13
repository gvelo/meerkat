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
	"meerkat/internal/storage/io"
)

type ByteSlicePlainDecoder struct {
	buf        *io.DecoderBuffer
	offsetsBuf []int
}

func NewByteSlicePlainDecoder() *ByteSlicePlainDecoder {
	return &ByteSlicePlainDecoder{
		buf:        io.NewDecoderBuffer(),
		offsetsBuf: make([]int, maxSlicesPerBlock),
	}
}

func (d *ByteSlicePlainDecoder) Decode(block []byte, data []byte, offsets []int) ([]byte, []int) {

	// TODO(gvelo) dec methods should be thread safe thread safe

	d.buf.SetBytes(block)

	// discard the block length
	_ = d.buf.ReadUvarint()

	// read the offsets
	ol := d.buf.ReadVarUintSlice(d.offsetsBuf)

	// decode offsets
	DeltaDecode(d.offsetsBuf[:ol], offsets)

	// read data
	dl := d.buf.ReadBytes(data)

	return data[:dl], offsets[:ol]

}
