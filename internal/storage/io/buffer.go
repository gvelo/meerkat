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
	"math"
)

func NewBuffer(cap int) *Buffer {
	return &Buffer{
		buf: make([]byte, cap),
		pos: 0,
	}
}

type Buffer struct {
	buf []byte
	pos int
}

func (b *Buffer) Cap() int {
	return len(b.buf)
}
func (b *Buffer) Pos(p int) {
	b.pos = p
}

func (b *Buffer) GetPos() int {
	return b.pos
}

func (b *Buffer) SetBuf(buf []byte) {
	b.buf = buf
	b.pos = 0
}

func (b *Buffer) Reset() {
	b.pos = 0
}

func (b *Buffer) Buf() []byte {
	return b.buf
}

func (b *Buffer) Data() []byte {
	return b.buf[:b.pos]
}

func (b *Buffer) WriteByteSlice(v []byte) {
	n := b.PutByteSlice(b.pos, v)
	b.pos += n
}

func (b *Buffer) WriteBytes(v []byte) {
	n := copy(b.buf[b.pos:], v)
	if n < len(v) {
		panic("buffer full`")
	}
	b.pos += n
}

func (b *Buffer) WriteVarInt(v int) {
	n := b.PutVarInt(b.pos, v)
	b.pos += n
}

func (b *Buffer) WriteIntAsUVarInt(v int) {
	n := b.PutIntAsUVarInt(b.pos, v)
	b.pos += n
}

func (b *Buffer) WriteVarInt64(v int64) {
	n := b.PutVarInt64(b.pos, v)
	b.pos += n
}

func (b *Buffer) WriteUVarInt64(v uint64) {
	n := b.PutUVarInt64(b.pos, v)
	b.pos += n
}

func (b *Buffer) WriteFloat64(v float64) {
	n := b.PutFloat64(b.pos, v)
	b.pos += n
}

func (b *Buffer) WriteVarUintSlice(v []int) {
	b.WriteIntAsUVarInt(len(v))
	for _, i := range v {
		b.WriteIntAsUVarInt(i)
	}
}

func (b *Buffer) Free() []byte {
	return b.buf[b.pos:]
}

func (b *Buffer) Available() int {
	return len(b.buf) - b.pos
}

func (b *Buffer) PutByte(pos int, v byte) {
	b.buf[pos] = v
}

func (b *Buffer) PutByteSlice(pos int, v []byte) int {
	n := b.PutIntAsUVarInt(pos, len(v))
	c := copy(b.buf[pos+n:], v)
	if c != len(v) {
		panic("buffer full")
	}
	return n + c
}

func (b *Buffer) PutVarInt(pos int, v int) int {
	return binary.PutVarint(b.buf[pos:], int64(v))
}

func (b *Buffer) PutIntAsUVarInt(pos int, v int) int {
	return binary.PutUvarint(b.buf[pos:], uint64(v))
}

func (b *Buffer) PutVarInt64(pos int, v int64) int {
	return binary.PutVarint(b.buf[pos:], v)
}

func (b *Buffer) PutUVarInt64(pos int, v uint64) int {
	return binary.PutUvarint(b.buf[pos:], v)
}

func (b *Buffer) PutFloat64(pos int, v float64) int {
	binary.LittleEndian.PutUint64(b.buf[pos:], math.Float64bits(v))
	return 8
}

func (b *Buffer) ReadByte() byte {
	r := b.buf[b.pos]
	b.pos++
	return r
}

func (b *Buffer) ReadBytes() []byte {
	r, n := b.Bytes(b.pos)
	b.pos += n
	return r
}

func (b *Buffer) ReadVarInt() int {
	r, n := b.VarInt(b.pos)
	b.pos += n
	return r
}

func (b *Buffer) ReadUVarInt() uint {
	r, n := b.UVarInt(b.pos)
	b.pos += n
	return r
}

func (b *Buffer) ReadUVarIntAsInt() int {
	r, n := b.UVarInt(b.pos)
	b.pos += n
	return int(r)
}

func (b *Buffer) Byte(pos int) byte {
	return b.buf[pos]
}

func (b *Buffer) Bytes(pos int) ([]byte, int) {
	l, n := binary.Uvarint(b.buf[pos:])
	if n <= 0 {
		panic("invalid Uvarint")
	}
	start := pos + n
	end := start + int(l)
	return b.buf[start:end], int(l) + n
}

func (b *Buffer) VarInt(pos int) (int, int) {
	r, n := binary.Varint(b.buf[pos:])
	if n <= 0 {
		panic("invalid Varint")
	}
	return int(r), n
}

func (b *Buffer) UVarInt(pos int) (uint, int) {
	r, n := binary.Uvarint(b.buf[pos:])
	if n <= 0 {
		panic("invalid Uvarint")
	}
	return uint(r), n
}

func (b *Buffer) UVarIntAsInt(pos int) (int, int) {
	r, n := binary.Uvarint(b.buf[pos:])
	if n <= 0 {
		panic("invalid Uvarint")
	}
	return int(r), n
}
