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
	"bufio"
	"os"
)

type BinaryWriter struct {
	file   *os.File
	writer *bufio.Writer
	Offset int
}

func NewBinaryWriter(name string) (*BinaryWriter, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	bw := &BinaryWriter{
		file:   f,
		writer: bufio.NewWriter(f),
		Offset: 0,
	}
	return bw, nil
}

func (br *BinaryWriter) Close() error {
	err := br.writer.Flush()
	if err != nil {
		return err
	}
	err = br.file.Sync()
	if err != nil {
		return err
	}
	err = br.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (br *BinaryWriter) WriteByte(x byte) error {
	br.writer.WriteByte(x)
	br.Offset++
	return nil
}

func (br *BinaryWriter) WriteVarInt(i int) error {
	return br.WriteVarUint64(uint64(i))
}

func (br *BinaryWriter) WriteVarUInt32(i uint32) error {
	return br.WriteVarUint64(uint64(i))
}

// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (br *BinaryWriter) WriteVarUint64(x uint64) error {
	for x >= 1<<7 {
		br.writer.WriteByte(uint8(x&0x7f | 0x80))
		x >>= 7
		br.Offset += 1
	}
	br.writer.WriteByte(uint8(x))
	br.Offset += 1
	return nil
}

// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (br *BinaryWriter) EncodedVarintLen(x uint64) uint64 {
	var s uint64 = 0
	for x >= 1<<7 {
		x >>= 7
		s += 1
	}
	s += 1
	return s
}

func (br *BinaryWriter) WriteFixedInt(i int) error {
	return br.WriteFixedUint64(uint64(i))
}

// EncodeFixed64 writes a 64-bit integer to the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (br *BinaryWriter) WriteFixedUint64(x uint64) error {
	br.writer.WriteByte(uint8(x))
	br.writer.WriteByte(uint8(x >> 8))
	br.writer.WriteByte(uint8(x >> 16))
	br.writer.WriteByte(uint8(x >> 24))
	br.writer.WriteByte(uint8(x >> 32))
	br.writer.WriteByte(uint8(x >> 40))
	br.writer.WriteByte(uint8(x >> 48))
	br.writer.WriteByte(uint8(x >> 56))
	br.Offset += 8
	return nil
}

// EncodeFixed32 writes a 32-bit integer to the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (br *BinaryWriter) WriteFixedUint32(x uint64) error {
	br.writer.WriteByte(uint8(x))
	br.writer.WriteByte(uint8(x >> 8))
	br.writer.WriteByte(uint8(x >> 16))
	br.writer.WriteByte(uint8(x >> 24))
	br.Offset += 4
	return nil
}

// EncodeZigzag64 writes a zigzag-encoded 64-bit integer
// to the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (br *BinaryWriter) WriteZigzag64(x uint64) error {
	// use signed number to get arithmetic right shift.
	br.Offset += 8
	return br.WriteVarUint64(uint64((x << 1) ^ uint64(int64(x)>>63)))
}

// EncodeZigzag32 writes a zigzag-encoded 32-bit integer
// to the Buffer.
// This is the format used for the sint32 protocol buffer type.
func (br *BinaryWriter) WriteZigzag32(x uint64) error {
	// use signed number to get arithmetic right shift.
	br.Offset += 4
	return br.WriteVarUint64(uint64((uint32(x) << 1) ^ uint32((int32(x) >> 31))))
}

// EncodeRawBytes writes a count-delimited byte buffer to the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (br *BinaryWriter) WriteBytes(b []byte) error {
	err := br.WriteVarInt(len(b))
	if err != nil {
		return err
	}
	_, err = br.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (br *BinaryWriter) Write(b []byte) (int, error) {
	nn, err := br.writer.Write(b)
	if err != nil {
		return nn, err
	}
	br.Offset += nn
	return nn, err
}

// EncodeStringBytes writes an encoded string to the Buffer.
// This is the format used for the proto2 string type.
func (br *BinaryWriter) WriteString(s string) error {
	l := len(s)
	br.WriteVarInt(l)
	br.writer.WriteString(s)
	br.Offset += l
	return nil
}
