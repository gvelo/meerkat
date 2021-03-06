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
	buf     *io.Buffer
	offsets []int
}

func NewByteSlicePlainDecoder() *ByteSlicePlainDecoder {
	return &ByteSlicePlainDecoder{
		buf:     &io.Buffer{},
		offsets: make([]int, 1024*4),
	}
}

func (d *ByteSlicePlainDecoder) Decode(block []byte) ([]byte, []int) {

	// TODO(gvelo) dec methods should be thread safe thread safe

	d.offsets = d.offsets[0:cap(d.offsets)]

	d.buf.SetBuf(block)

	// read the len of the offsets slice
	l := d.buf.ReadUVarIntAsInt()

	// grow the offset buffer if needed
	if l > len(d.offsets) {
		d.offsets = make([]int, l*2)
	}

	for i := 0; i < l; i++ {
		d.offsets[i] = d.buf.ReadUVarIntAsInt()
	}

	d.offsets = d.offsets[:l]

	// decode offsets
	DeltaDecode(d.offsets)

	// read data
	data := d.buf.Free()

	return data, d.offsets

}
