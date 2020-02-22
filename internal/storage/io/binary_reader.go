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
	"errors"
	"fmt"
	"io"
)

// errOverflow is returned when an integer is too large to be represented.
var errOverflow = errors.New("proto: integer overflow")

type BinaryReader struct {
	bytes  []byte
	Offset int
	Size   int
}

func NewBinaryReader(b []byte) *BinaryReader {
	return &BinaryReader{
		bytes:  b,
		Offset: 0,
		Size:   len(b),
	}
}

// SetBytes set a new buffer and reset the internal
// state. Handy method for reuse readers.
func (br *BinaryReader) SetBytes(b []byte) {
	br.bytes = b
	br.Offset = 0
	br.Size = len(b)
}

func (br *BinaryReader) Bytes() []byte {
	return br.bytes
}

func (br *BinaryReader) decodeVarintSlow() (x uint64, err error) {
	for shift := uint(0); shift < 64; shift += 7 {
		if br.Offset >= len(br.bytes) {
			err = io.ErrUnexpectedEOF
			return
		}
		b := br.bytes[br.Offset]
		br.Offset++
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			return
		}
	}

	// The number is too large to represent in a 64-bit value.
	err = errOverflow
	return
}

func (br *BinaryReader) ReadByte() byte {
	b := br.bytes[br.Offset]
	br.Offset++
	return b
}

// ReadVarint64 reads a varint-encoded integer from the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (br *BinaryReader) ReadVarint64() (x uint64, err error) {
	i := br.Offset
	if br.Offset >= len(br.bytes) {
		return 0, io.ErrUnexpectedEOF
	} else if br.bytes[br.Offset] < 0x80 {
		br.Offset++
		return uint64(br.bytes[i]), nil
	} else if len(br.bytes)-br.Offset < 10 {
		return br.decodeVarintSlow()
	}

	var b uint64
	// we already checked the first byte
	x = uint64(br.bytes[br.Offset]) - 0x80
	br.Offset++

	b = uint64(br.bytes[br.Offset])
	br.Offset++
	x += b << 7
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 7

	b = uint64(br.bytes[br.Offset])
	br.Offset++
	x += b << 14
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 14

	b = uint64(br.bytes[br.Offset])
	br.Offset++
	x += b << 21
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 21

	b = uint64(br.bytes[br.Offset])
	br.Offset++
	x += b << 28
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 28

	b = uint64(br.bytes[br.Offset])
	br.Offset++
	x += b << 35
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 35

	b = uint64(br.bytes[br.Offset])
	br.Offset++
	x += b << 42
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 42

	b = uint64(br.bytes[br.Offset])
	br.Offset++
	x += b << 49
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 49

	b = uint64(br.bytes[br.Offset])
	br.Offset++
	x += b << 56
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 56

	b = uint64(br.bytes[br.Offset])
	br.Offset++
	x += b << 63
	if b&0x80 == 0 {
		goto done
	}

	return 0, errOverflow

done:
	return x, nil
}

func (br *BinaryReader) ReadVarInt() (int, error) {
	i, err := br.ReadVarint64()
	return int(i), err
}

// ReadFixed64 reads a 64-bit integer from the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (br *BinaryReader) ReadFixed64() (x uint64, err error) {
	// x, err already 0
	i := br.Offset + 8
	if i < 0 || i > len(br.bytes) {
		err = io.ErrUnexpectedEOF
		return
	}
	br.Offset = i

	x = uint64(br.bytes[i-8])
	x |= uint64(br.bytes[i-7]) << 8
	x |= uint64(br.bytes[i-6]) << 16
	x |= uint64(br.bytes[i-5]) << 24
	x |= uint64(br.bytes[i-4]) << 32
	x |= uint64(br.bytes[i-3]) << 40
	x |= uint64(br.bytes[i-2]) << 48
	x |= uint64(br.bytes[i-1]) << 56
	return
}

// ReadFixed32 reads a 32-bit integer from the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (br *BinaryReader) ReadFixed32() (x uint64, err error) {
	// x, err already 0
	i := br.Offset + 4
	if i < 0 || i > len(br.bytes) {
		err = io.ErrUnexpectedEOF
		return
	}
	br.Offset = i

	x = uint64(br.bytes[i-4])
	x |= uint64(br.bytes[i-3]) << 8
	x |= uint64(br.bytes[i-2]) << 16
	x |= uint64(br.bytes[i-1]) << 24
	return
}

// ReadZigzag64 reads a zigzag-encoded 64-bit integer
// from the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (br *BinaryReader) ReadZigzag64() (x uint64, err error) {
	x, err = br.ReadVarint64()
	if err != nil {
		return
	}
	x = (x >> 1) ^ uint64((int64(x&1)<<63)>>63)
	return
}

// ReadZigzag32 reads a zigzag-encoded 32-bit integer
// from  the Buffer.
// This is the format used for the sint32 protocol buffer type.
func (br *BinaryReader) ReadZigzag32() (x uint64, err error) {
	x, err = br.ReadVarint64()
	if err != nil {
		return
	}
	x = uint64((uint32(x) >> 1) ^ uint32((int32(x&1)<<31)>>31))
	return
}

// ReadBytes reads a count-delimited byte buffer from the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (br *BinaryReader) ReadBytes() (buf []byte, err error) {

	n, err := br.ReadVarInt()

	if err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, fmt.Errorf("proto: bad byte length %d", n)
	}
	end := br.Offset + n
	if end < br.Offset || end > len(br.bytes) {
		return nil, io.ErrUnexpectedEOF
	}

	buf = br.bytes[br.Offset:end]
	br.Offset += n
	return

}

func (br *BinaryReader) SliceAt(offset int) []byte {
	return br.bytes[offset:]
}

// ReadString reads an encoded string from the Buffer.
// This is the format used for the proto2 string type.
func (br *BinaryReader) ReadString() (s string, err error) {
	buf, err := br.ReadBytes()
	if err != nil {
		return
	}
	return string(buf), nil
}

func (br *BinaryReader) ReadVarUintSlice() (s []int, err error) {

	l, err := br.ReadVarInt()

	if err != nil {
		return nil, err
	}

	b := make([]int, l)

	for i := 0; i < l; i++ {

		b[i], err = br.ReadVarInt()

		if err != nil {
			return nil, err
		}

	}

	return b, nil

}
