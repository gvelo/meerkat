package io

import (
	"bufio"
	"errors"
	"eventdb/segment"
	"eventdb/segment/inmem"
	"fmt"
	"io"
	"math"
	"os"
)

type FileType byte

const (
	MagicNumber = "ed"
	Version     = "01"

	PostingListV1 FileType = 0
	StringIndexV1 FileType = 1
	SkipListV1    FileType = 2
	RowStoreV1    FileType = 3
	RowStoreIDXV1 FileType = 4
	SegmentInfo   FileType = 5
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

func (br *BinaryWriter) WriteValue(v interface{}, info *segment.FieldInfo) error {
	var err error
	switch info.Type {
	case segment.FieldTypeInt:
		err = br.WriteVarUint64(v.(uint64))
	case segment.FieldTypeText:
		err = br.WriteString(v.(string))
	case segment.FieldTypeKeyword:
		err = br.WriteString(v.(string))
	case segment.FieldTypeTimestamp:
		err = br.WriteVarUint64(v.(uint64))
	case segment.FieldTypeFloat:
		err = br.WriteFixedUint64(math.Float64bits(v.(float64)))
	}
	return err
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

	err = br.WriteVarInt(page.Total)

	if err != nil {
		return err
	}

	err = br.WriteVarInt(page.StartID)

	if err != nil {
		return err
	}

	err = br.WriteVarInt(page.PayloadSize)

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

func (br *BinaryReader) ReadValue(info *segment.FieldInfo) (interface{}, error) {
	switch info.Type {
	case segment.FieldTypeInt:
		return br.ReadVarint64()
	case segment.FieldTypeText:
		return br.ReadString()
	case segment.FieldTypeKeyword:
		return br.ReadString()
	case segment.FieldTypeTimestamp:
		return br.ReadVarint64()
	case segment.FieldTypeFloat:
		return br.ReadFixed64()
	}
	return nil, errors.New(fmt.Sprintf("info.Type %d not found ", info.Type))
}
