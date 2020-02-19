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

import (
	"encoding/binary"
)

type EncoderBuffer struct {
	buf []byte
	len int
}

func NewEncoderBuffer(cap int) *EncoderBuffer {
	return &EncoderBuffer{
		buf: make([]byte, cap),
		len: 0,
	}
}

func (b *EncoderBuffer) Bytes() []byte {
	return b.buf[0:b.len]
}

func (b *EncoderBuffer) Reset(size int) {

	if cap(b.buf) < size {
		b.buf = make([]byte, size)
	}

	b.buf = b.buf[0:size]
	b.len = 0

}

func (b *EncoderBuffer) WriteUvarint(x int) {
	b.WriteUvarint64(uint64(x))
}

func (b *EncoderBuffer) WriteUvarintAt(offset int, x int) {
	b.WriteUvarint64At(offset, uint64(x))
}

func (b *EncoderBuffer) WriteVarint(x int) {
	b.WriteVarint64(int64(x))
}

func (b *EncoderBuffer) WriteVarintAt(offset int, x int) {
	b.WriteVarint64At(offset, int64(x))
}

func (b *EncoderBuffer) WriteUvarint64(x uint64) {
	b.len = b.len + binary.PutUvarint(b.buf[b.len:], x)
}

func (b *EncoderBuffer) WriteUvarint64At(offset int, x uint64) {

	n := binary.PutUvarint(b.buf[offset:], x)

	if b.len < offset {
		b.len = offset + n
	}

}

func (b *EncoderBuffer) WriteVarint64(x int64) {
	b.len = b.len + binary.PutVarint(b.buf[b.len:], x)
}

func (b *EncoderBuffer) WriteVarint64At(offset int, x int64) {

	n := binary.PutVarint(b.buf[offset:], x)

	if b.len < n {
		b.len = n
	}

}

func (b *EncoderBuffer) WriteBytes(bytes []byte) {

	copy(b.buf[b.len:], bytes)

	b.len = b.len + len(bytes)

}

func (b *EncoderBuffer) WriteVarUintSliceAt(offset int, x []int) {

	b.WriteUvarintAt(offset, len(x))

	for _, v := range x {
		b.WriteUvarint(v)
	}

}

func (b *EncoderBuffer) Len() int {
	return b.len
}

// SizeVarint returns the varint encoding size of an integer.
func SizeUVarint(x uint64) int {
	switch {
	case x < 1<<7:
		return 1
	case x < 1<<14:
		return 2
	case x < 1<<21:
		return 3
	case x < 1<<28:
		return 4
	case x < 1<<35:
		return 5
	case x < 1<<42:
		return 6
	case x < 1<<49:
		return 7
	case x < 1<<56:
		return 8
	case x < 1<<63:
		return 9
	}
	return 10
}
