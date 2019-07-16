// Copyright 2019 The Meerkat Authors
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
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"os"
)

type FileType byte

const (
	MagicNumber = "mk"
	Version     = "01"

	PostingListV1 FileType = 0
	StringIndexV1 FileType = 1
	SkipListV1    FileType = 2
	RowStoreV1    FileType = 3
	RowStoreIDXV1 FileType = 4
	SegmentInfo   FileType = 5
)

type baseBinaryWriter struct {
	writer *bufio.Writer
	Offset int
}

type BinaryWriter struct {
	baseBinaryWriter
	file *os.File
}

type BufferBinaryWriter struct {
	baseBinaryWriter
	Buffer *bytes.Buffer
}

func NewBinaryWriter(name string) (*BinaryWriter, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	bw := &BinaryWriter{
		file: f,
		baseBinaryWriter: baseBinaryWriter{
			writer: bufio.NewWriter(f),
			Offset: 0},
	}
	return bw, nil
}

func NewBufferBinaryWriter() (*BufferBinaryWriter, error) {
	var b bytes.Buffer
	bw := &BufferBinaryWriter{
		Buffer: &b,
		baseBinaryWriter: baseBinaryWriter{
			writer: bufio.NewWriter(&b),
			Offset: 0},
	}
	return bw, nil
}

func (bw *baseBinaryWriter) WriteValue(v interface{}, info *segment.FieldInfo) error {
	var err error
	switch info.Type {
	case segment.FieldTypeInt:
		err = bw.WriteVarUint64(v.(uint64))
	case segment.FieldTypeText:
		err = bw.WriteString(v.(string))
	case segment.FieldTypeKeyword:
		err = bw.WriteString(v.(string))
	case segment.FieldTypeTimestamp:
		err = bw.WriteVarUint64(v.(uint64))
	case segment.FieldTypeFloat:
		err = bw.WriteFixedUint64(math.Float64bits(v.(float64)))
	}
	return err
}

func (bw *BufferBinaryWriter) Flush() error {
	err := bw.writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (bw *BinaryWriter) Close() error {
	err := bw.writer.Flush()
	if err != nil {
		return err
	}

	err = bw.file.Sync()
	if err != nil {
		return err
	}
	err = bw.file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (br *BinaryWriter) WriteHeader(fileType FileType) error {

	b := []byte(MagicNumber)

	_, err := br.Write(b)

	if err != nil {
		return err
	}

	err = br.WriteByte(byte(fileType))

	if err != nil {
		return err
	}

	return nil
}

func (br *BinaryWriter) WritePageHeader(page *inmem.Page) error {

	err := br.WriteByte(byte(page.Enc))

	if err != nil {
		return err
	}

	err = br.WriteVarUInt(page.Offset)

	if err != nil {
		return err
	}

	err = br.WriteVarUInt(page.Total)

	if err != nil {
		return err
	}

	err = br.WriteVarUInt(page.StartID)

	if err != nil {
		return err
	}

	err = br.WriteVarUInt(page.PayloadSize)

	if err != nil {
		return err
	}

	return nil
}

func (bw *baseBinaryWriter) WriteByte(x byte) error {
	bw.writer.WriteByte(x)
	bw.Offset++
	return nil
}

func (bw *baseBinaryWriter) WriteVarInt(i int) error {
	return bw.WriteZigzag64(i)
}

func (bw *baseBinaryWriter) WriteVarUInt(i int) error {
	return bw.WriteVarUint64(uint64(i))
}

// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (bw *baseBinaryWriter) WriteVarUint64(x uint64) error {
	for x >= 1<<7 {
		bw.writer.WriteByte(uint8(x&0x7f | 0x80))
		x >>= 7
		bw.Offset += 1
	}
	bw.writer.WriteByte(uint8(x))
	bw.Offset += 1
	return nil
}

// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (bw *baseBinaryWriter) EncodedVarintLen(x uint64) uint64 {
	var s uint64 = 0
	for x >= 1<<7 {
		x >>= 7
		s += 1
	}
	s += 1
	return s
}

func (bw *baseBinaryWriter) WriteFixedInt(i int) error {
	return bw.WriteFixedUint64(uint64(i))
}

// EncodeFixed64 writes a 64-bit integer to the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (bw *baseBinaryWriter) WriteFixedUint64(x uint64) error {
	bw.writer.WriteByte(uint8(x))
	bw.writer.WriteByte(uint8(x >> 8))
	bw.writer.WriteByte(uint8(x >> 16))
	bw.writer.WriteByte(uint8(x >> 24))
	bw.writer.WriteByte(uint8(x >> 32))
	bw.writer.WriteByte(uint8(x >> 40))
	bw.writer.WriteByte(uint8(x >> 48))
	bw.writer.WriteByte(uint8(x >> 56))
	bw.Offset += 8
	return nil
}

// EncodeFixed32 writes a 32-bit integer to the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (bw *baseBinaryWriter) WriteFixedUint32(x uint64) error {
	bw.writer.WriteByte(uint8(x))
	bw.writer.WriteByte(uint8(x >> 8))
	bw.writer.WriteByte(uint8(x >> 16))
	bw.writer.WriteByte(uint8(x >> 24))
	bw.Offset += 4
	return nil
}

// EncodeZigzag64 writes a zigzag-encoded 64-bit integer
// to the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (bw *baseBinaryWriter) WriteZigzag64(x int) error {
	// use signed number to get arithmetic right shift.
	bw.Offset += 8
	return bw.WriteVarUint64((uint64(x) << 1) ^ uint64(x>>63))
}

// EncodeRawBytes writes a count-delimited byte buffer to the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (bw *baseBinaryWriter) WriteBytes(b []byte) error {
	err := bw.WriteVarUInt(len(b))
	if err != nil {
		return err
	}
	_, err = bw.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (bw *baseBinaryWriter) Write(b []byte) (int, error) {
	nn, err := bw.writer.Write(b)
	if err != nil {
		return nn, err
	}
	bw.Offset += nn
	return nn, err
}

// EncodeStringBytes writes an encoded string to the Buffer.
// This is the format used for the proto2 string type.
func (bw *baseBinaryWriter) WriteString(s string) error {
	l := len(s)
	bw.WriteVarUInt(l)
	bw.writer.WriteString(s)
	bw.Offset += l
	return nil
}

// errOverflow is returned when an integer is too large to be represented.
var errOverflow = errors.New("proto: integer overflow")
var errUnknFileType = errors.New("unknown file type")

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

// ReadHeader read the file header returning the file type.
func (br *BinaryReader) ReadHeader() (FileType, error) {

	l := len(MagicNumber)

	if len(br.bytes) < (l + 1) {
		return 0, io.ErrUnexpectedEOF
	}

	fMagic := string(br.bytes[:l])

	if fMagic != MagicNumber {
		return 0, errUnknFileType
	}

	br.Offset += l + 1

	return FileType(br.bytes[l]), nil

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

func (br *BinaryReader) ReadPageHeader() (*inmem.Page, error) {

	p := new(inmem.Page)
	var err error

	b := br.ReadByte()
	p.Enc = inmem.Encoding(b)

	p.Offset, err = br.ReadVarUInt()
	if err != nil {
		return nil, err
	}

	p.Total, err = br.ReadVarUInt()
	if err != nil {
		return nil, err
	}
	p.StartID, err = br.ReadVarUInt()
	if err != nil {
		return nil, err
	}
	p.PayloadSize, err = br.ReadVarUInt()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (br *BinaryReader) ReadByte() byte {
	b := br.bytes[br.Offset]
	br.Offset++
	return b
}

// ReadVarUInt64 reads a varint-encoded integer from the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (br *BinaryReader) ReadVarUInt64() (x uint64, err error) {
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
	i, err := br.ReadZigzag64()
	return int(i), err
}

func (br *BinaryReader) ReadVarUInt() (int, error) {
	i, err := br.ReadVarUInt64()
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

// ReadZigZag64 reads a zigzag-encoded 64-bit integer
// from the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (br *BinaryReader) ReadZigzag64() (x uint64, err error) {
	x, err = br.ReadVarUInt64()
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
	x, err = br.ReadVarUInt64()
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

	n, err := br.ReadVarUInt()

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

func (br *BinaryReader) ReadValue(info *segment.FieldInfo) (interface{}, error) {
	switch info.Type {
	case segment.FieldTypeInt:
		return br.ReadVarUInt64()
	case segment.FieldTypeText:
		return br.ReadString()
	case segment.FieldTypeKeyword:
		return br.ReadString()
	case segment.FieldTypeTimestamp:
		return br.ReadVarUInt64()
	case segment.FieldTypeFloat:
		return br.ReadFixed64()
	}
	return nil, errors.New(fmt.Sprintf("info.Type %d not found ", info.Type))
}
