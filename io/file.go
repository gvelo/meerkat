package io

import (
	"bufio"
	"errors"
	"eventdb/io/mmap"
	"eventdb/segment"
	"fmt"
	"io"
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
	fw := &BinaryWriter{
		file:   f,
		writer: bufio.NewWriter(f),
		Offset: 0,
	}
	return fw, nil
}

func (fw *BinaryWriter) WriteValue(v interface{}, info segment.FieldInfo) {
	switch info.FieldType {
	case segment.FieldTypeInt:
		fw.WriteEncodedVarint(v.(uint64))
	case segment.FieldTypeText:
		fw.WriteEncodedStringBytes(v.(string))
	case segment.FieldTypeKeyword:
		fw.WriteEncodedStringBytes(v.(string))
	case segment.FieldTypeTimestamp:
		fw.WriteEncodedVarint(v.(uint64))
	case segment.FieldTypeFloat:
		fw.WriteEncodedFixed64(v.(uint64))
	}
}

func (fw *BinaryWriter) Close() error {
	err := fw.writer.Flush()
	if err != nil {
		return err
	}
	err = fw.file.Sync()
	if err != nil {
		return err
	}
	err = fw.file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (fw *BinaryWriter) WriteHeader(fileType FileType) error {
	b := []byte(MagicNumber)
	fw.Write(b)
	fw.WriteByte(byte(fileType))
	return nil
}

func (fw *BinaryWriter) WriteByte(x byte) error {
	fw.writer.WriteByte(x)
	fw.Offset++
	return nil
}

// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (fw *BinaryWriter) WriteEncodedVarint(x uint64) error {
	for x >= 1<<7 {
		fw.writer.WriteByte(uint8(x&0x7f | 0x80))
		x >>= 7
		fw.Offset += 1
	}
	fw.writer.WriteByte(uint8(x))
	fw.Offset += 1
	return nil
}

// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (fw *BinaryWriter) EncodedVarintLen(x uint64) uint64 {
	var s uint64 = 0
	for x >= 1<<7 {
		x >>= 7
		s += 1
	}
	s += 1
	return s
}

// EncodeFixed64 writes a 64-bit integer to the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (fw *BinaryWriter) WriteEncodedFixed64(x uint64) error {
	fw.writer.WriteByte(uint8(x))
	fw.writer.WriteByte(uint8(x >> 8))
	fw.writer.WriteByte(uint8(x >> 16))
	fw.writer.WriteByte(uint8(x >> 24))
	fw.writer.WriteByte(uint8(x >> 32))
	fw.writer.WriteByte(uint8(x >> 40))
	fw.writer.WriteByte(uint8(x >> 48))
	fw.writer.WriteByte(uint8(x >> 56))
	fw.Offset += 8
	return nil
}

// EncodeFixed32 writes a 32-bit integer to the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (fw *BinaryWriter) WriteEncodedFixed32(x uint64) error {
	fw.writer.WriteByte(uint8(x))
	fw.writer.WriteByte(uint8(x >> 8))
	fw.writer.WriteByte(uint8(x >> 16))
	fw.writer.WriteByte(uint8(x >> 24))
	fw.Offset += 4
	return nil
}

// EncodeZigzag64 writes a zigzag-encoded 64-bit integer
// to the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (fw *BinaryWriter) WriteEncodedZigzag64(x uint64) error {
	// use signed number to get arithmetic right shift.
	fw.Offset += 8
	return fw.WriteEncodedVarint(uint64((x << 1) ^ uint64(int64(x)>>63)))
}

// EncodeZigzag32 writes a zigzag-encoded 32-bit integer
// to the Buffer.
// This is the format used for the sint32 protocol buffer type.
func (fw *BinaryWriter) WriteEncodedZigzag32(x uint64) error {
	// use signed number to get arithmetic right shift.
	fw.Offset += 4
	return fw.WriteEncodedVarint(uint64((uint32(x) << 1) ^ uint32((int32(x) >> 31))))
}

// EncodeRawBytes writes a count-delimited byte buffer to the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (fw *BinaryWriter) WriteEncodedRawBytes(b []byte) error {
	fw.WriteEncodedVarint(uint64(len(b)))
	fw.Write(b)
	return nil
}

func (fw *BinaryWriter) Write(b []byte) (int, error) {
	nn, err := fw.writer.Write(b)
	if err != nil {
		return nn, err
	}
	fw.Offset += nn
	return nn, err
}

// EncodeStringBytes writes an encoded string to the Buffer.
// This is the format used for the proto2 string type.
func (fw *BinaryWriter) WriteEncodedStringBytes(s string) error {
	l := len(s)
	fw.WriteEncodedVarint(uint64(l))
	fw.writer.WriteString(s)
	fw.Offset += l
	return nil
}

type BinaryReader struct {
	bytes  []byte
	Offset int
	Size   int
}

func NewBinaryReader(name string) (*BinaryReader, error) {

	b, err := mmap.Map(name)

	if err != nil {
		return nil, err
	}

	fr := &BinaryReader{
		bytes:  b,
		Offset: 0,
		Size:   len(b),
	}
	return fr, nil
}

// errOverflow is returned when an integer is too large to be represented.
var errOverflow = errors.New("proto: integer overflow")
var errUnknFileType = errors.New("unknown file type")

// ReadHeader read the file header returning the file type.
func (fr *BinaryReader) ReadHeader() (FileType, error) {

	l := len(MagicNumber)

	if len(fr.bytes) < (l + 1) {
		return 0, io.ErrUnexpectedEOF
	}

	fMagic := string(fr.bytes[:l])

	if fMagic != MagicNumber {
		return 0, errUnknFileType
	}

	fr.Offset += l + 1

	return FileType(fr.bytes[l]), nil

}

func (fr *BinaryReader) decodeVarintSlow() (x uint64, err error) {
	for shift := uint(0); shift < 64; shift += 7 {
		if fr.Offset >= len(fr.bytes) {
			err = io.ErrUnexpectedEOF
			return
		}
		b := fr.bytes[fr.Offset]
		fr.Offset++
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			return
		}
	}

	// The number is too large to represent in a 64-bit value.
	err = errOverflow
	return
}

func (fr *BinaryReader) DecodeByte() byte {
	byte := fr.bytes[fr.Offset]
	fr.Offset++
	return byte
}

// DecodeVarint reads a varint-encoded integer from the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (fr *BinaryReader) DecodeVarint() (x uint64, err error) {
	i := fr.Offset
	if fr.Offset >= len(fr.bytes) {
		return 0, io.ErrUnexpectedEOF
	} else if fr.bytes[fr.Offset] < 0x80 {
		fr.Offset++
		return uint64(fr.bytes[i]), nil
	} else if len(fr.bytes)-fr.Offset < 10 {
		return fr.decodeVarintSlow()
	}

	var b uint64
	// we already checked the first byte
	x = uint64(fr.bytes[fr.Offset]) - 0x80
	fr.Offset++

	b = uint64(fr.bytes[fr.Offset])
	fr.Offset++
	x += b << 7
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 7

	b = uint64(fr.bytes[fr.Offset])
	fr.Offset++
	x += b << 14
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 14

	b = uint64(fr.bytes[fr.Offset])
	fr.Offset++
	x += b << 21
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 21

	b = uint64(fr.bytes[fr.Offset])
	fr.Offset++
	x += b << 28
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 28

	b = uint64(fr.bytes[fr.Offset])
	fr.Offset++
	x += b << 35
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 35

	b = uint64(fr.bytes[fr.Offset])
	fr.Offset++
	x += b << 42
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 42

	b = uint64(fr.bytes[fr.Offset])
	fr.Offset++
	x += b << 49
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 49

	b = uint64(fr.bytes[fr.Offset])
	fr.Offset++
	x += b << 56
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 56

	b = uint64(fr.bytes[fr.Offset])
	fr.Offset++
	x += b << 63
	if b&0x80 == 0 {
		goto done
	}

	return 0, errOverflow

done:
	return x, nil
}

// DecodeFixed64 reads a 64-bit integer from the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (fr *BinaryReader) DecodeFixed64() (x uint64, err error) {
	// x, err already 0
	i := fr.Offset + 8
	if i < 0 || i > len(fr.bytes) {
		err = io.ErrUnexpectedEOF
		return
	}
	fr.Offset = i

	x = uint64(fr.bytes[i-8])
	x |= uint64(fr.bytes[i-7]) << 8
	x |= uint64(fr.bytes[i-6]) << 16
	x |= uint64(fr.bytes[i-5]) << 24
	x |= uint64(fr.bytes[i-4]) << 32
	x |= uint64(fr.bytes[i-3]) << 40
	x |= uint64(fr.bytes[i-2]) << 48
	x |= uint64(fr.bytes[i-1]) << 56
	return
}

// DecodeFixed32 reads a 32-bit integer from the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (fr *BinaryReader) DecodeFixed32() (x uint64, err error) {
	// x, err already 0
	i := fr.Offset + 4
	if i < 0 || i > len(fr.bytes) {
		err = io.ErrUnexpectedEOF
		return
	}
	fr.Offset = i

	x = uint64(fr.bytes[i-4])
	x |= uint64(fr.bytes[i-3]) << 8
	x |= uint64(fr.bytes[i-2]) << 16
	x |= uint64(fr.bytes[i-1]) << 24
	return
}

// DecodeZigzag64 reads a zigzag-encoded 64-bit integer
// from the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (fr *BinaryReader) DecodeZigzag64() (x uint64, err error) {
	x, err = fr.DecodeVarint()
	if err != nil {
		return
	}
	x = (x >> 1) ^ uint64((int64(x&1)<<63)>>63)
	return
}

// DecodeZigzag32 reads a zigzag-encoded 32-bit integer
// from  the Buffer.
// This is the format used for the sint32 protocol buffer type.
func (fr *BinaryReader) DecodeZigzag32() (x uint64, err error) {
	x, err = fr.DecodeVarint()
	if err != nil {
		return
	}
	x = uint64((uint32(x) >> 1) ^ uint32((int32(x&1)<<31)>>31))
	return
}

// DecodeRawBytes reads a count-delimited byte buffer from the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (fr *BinaryReader) DecodeRawBytes(alloc bool) (buf []byte, err error) {
	n, err := fr.DecodeVarint()
	if err != nil {
		return nil, err
	}

	nb := int(n)
	if nb < 0 {
		return nil, fmt.Errorf("proto: bad byte length %d", nb)
	}
	end := fr.Offset + nb
	if end < fr.Offset || end > len(fr.bytes) {
		return nil, io.ErrUnexpectedEOF
	}

	if !alloc {
		// todo: check if can get more uses of alloc=false
		buf = fr.bytes[fr.Offset:end]
		fr.Offset += nb
		return
	}

	buf = make([]byte, nb)
	copy(buf, fr.bytes[fr.Offset:])
	fr.Offset += nb
	return
}

func (fr BinaryReader) SliceAt(offset int) []byte {

	return fr.bytes[offset:]
}

// DecodeStringBytes reads an encoded string from the Buffer.
// This is the format used for the proto2 string type.
func (fr *BinaryReader) DecodeStringBytes() (s string, err error) {
	buf, err := fr.DecodeRawBytes(false)
	if err != nil {
		return
	}
	return string(buf), nil
}

// Close close and unmap the file.
func (fr *BinaryReader) Close() error {
	return mmap.UnMap(fr.bytes)
}

func (fr *BinaryReader) ReadValue(info segment.FieldInfo) (interface{}, error) {
	switch info.FieldType {
	case segment.FieldTypeInt:
		return fr.DecodeVarint()
	case segment.FieldTypeText:
		return fr.DecodeStringBytes()
	case segment.FieldTypeKeyword:
		return fr.DecodeStringBytes()
	case segment.FieldTypeTimestamp:
		return fr.DecodeVarint()
	case segment.FieldTypeFloat:
		return fr.DecodeFixed64()
	}
	return nil, errors.New(fmt.Sprintf("info.FieldType %d not found ", info.FieldType))
}
