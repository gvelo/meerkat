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
	"meerkat/internal/util/sliceutil"
)

type ByteSliceDictDecoder struct {
	buf        *io.Buffer
	b          []byte
	dictValues []byte
	offsets    []int
}

func NewByteSliceDictDecoder(b []byte, eb io.Bounds) *ByteSliceDictDecoder {

	r := io.NewBinaryReader(b[eb.Start:eb.End])
	r.SetOffset(r.Size() - 8)

	footer := r.ReadFixed64()
	r.SetOffset(footer)

	r.ReadUVarint64() // start
	startOffsets := r.ReadUVarint64()

	return &ByteSliceDictDecoder{
		buf:        &io.Buffer{},
		b:          b,
		dictValues: r.ReadSlice(0, int(startOffsets)),
		offsets:    sliceutil.B2I(r.ReadSlice(int(startOffsets), footer)),
	}
}

func (d *ByteSliceDictDecoder) Decode(block []byte) ([]byte, []int) {

	v := make([]byte, 0)

	offsets := make([]int, 0)
	actOffset := 0

	b := io.NewBinaryReader(block)
	dict := io.NewBinaryReader(d.dictValues)

	b.SetOffset(0)

	for i := 0; b.Offset() < b.Size(); i++ {

		idx := b.ReadVarint64()
		start := 0
		if idx > 0 {
			start = d.offsets[idx-1]
		}
		end := d.offsets[idx]

		value := dict.ReadSlice(start, end)

		v = append(v, value...)

		actOffset = actOffset + len(value)
		offsets = append(offsets, actOffset)
	}

	return v, offsets
}
