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
	"meerkat/internal/utils"
)

type BinaryReader struct {
	bytes  []byte
	offset int
	size   int
}

func NewBinaryReader(b []byte) *BinaryReader {
	return &BinaryReader{
		bytes:  b,
		offset: 0,
		size:   len(b),
	}
}

// SetBytes set a new buffer and reset the internal
// state. Handy method for reuse readers.
func (br *BinaryReader) SetBytes(b []byte) {
	br.bytes = b
	br.offset = 0
	br.size = len(b)
}

func (br *BinaryReader) Bytes() []byte {
	return br.bytes
}

func (br *BinaryReader) Offset() int {
	return br.offset
}

func (br *BinaryReader) SetOffset(o int) {
	br.offset = o
}

func (br *BinaryReader) ReadByte() byte {
	b := br.bytes[br.offset]
	br.offset++
	return b
}

func (br *BinaryReader) ReadUVarint() int {

	return int(br.ReadUVarint64())

}

func (br *BinaryReader) ReadUVarint64() uint64 {

	i, n := binary.Uvarint(br.bytes[br.offset:])

	if n <= 0 {
		panic("error reading uvarint")
	}

	br.offset += n

	return i

}

func (br *BinaryReader) ReadVarint() int {

	return int(br.ReadVarint64())

}

func (br *BinaryReader) ReadVarint64() int64 {

	i, n := binary.Varint(br.bytes[br.offset:])

	if n <= 0 {
		panic("error reading varint")
	}

	br.offset += n

	return i

}

func (br *BinaryReader) ReadFixed64() int {

	i := binary.LittleEndian.Uint64(br.bytes[br.offset:])

	br.offset += 8

	return int(i)

}

func (br *BinaryReader) ReadBytes() []byte {

	n := br.ReadUVarint()

	end := br.offset + n

	buf := br.bytes[br.offset:end]

	br.offset += n

	return buf

}

func (br *BinaryReader) SliceFrom(offset int) []byte {
	return br.bytes[offset:]
}

func (br *BinaryReader) ReadString() string {

	buf := br.ReadBytes()

	return utils.ByteSlice2String(buf)

}

func (br *BinaryReader) ReadVarUintSlice() []int {

	l := br.ReadUVarint()

	b := make([]int, l)

	for i := 0; i < l; i++ {
		b[i] = br.ReadUVarint()
	}

	return b

}
