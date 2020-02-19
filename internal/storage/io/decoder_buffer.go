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

package io

import "encoding/binary"

type DecoderBuffer struct {
	buf   []byte
	index int
}

func NewDecoderBuffer() *DecoderBuffer {
	return &DecoderBuffer{}
}

func (d *DecoderBuffer) Bytes() []byte {
	return d.buf
}

func (d *DecoderBuffer) SetBytes(b []byte) {
	d.buf = b
	d.index = 0
}

func (d *DecoderBuffer) ReadUvarint64() uint64 {

	u, n := binary.Uvarint(d.buf[d.index:])

	if n < 0 {
		panic("error reading uvarint64")
	}

	d.index = d.index + n

	return u
}

func (d *DecoderBuffer) ReadUvarint() int {
	return int(d.ReadUvarint64())
}

func (d *DecoderBuffer) ReadVarint64() int64 {

	u, n := binary.Varint(d.buf[d.index:])

	if n < 0 {
		panic("error reading varint64")
	}

	d.index = d.index + n

	return u
}

func (d *DecoderBuffer) ReadVarint() int {
	return int(d.ReadVarint64())
}

func (d *DecoderBuffer) ReadVarUintSlice(dst []int) int {

	l := d.ReadUvarint()

	for i := 0; i < l; i++ {
		dst[i] = d.ReadUvarint()
	}

	return l

}

func (d *DecoderBuffer) ReadBytes(b []byte) int {
	return copy(b, d.buf[d.index:])
}

func (d *DecoderBuffer) Remaining() []byte {
	return d.buf[d.index:]
}
