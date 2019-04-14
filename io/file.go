package io

import (
	"bufio"
	"errors"
	"eventdb/io/mmap"
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
	RowStoreV1    FileType = 2
)

type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
	size   int
}

func newFileWriter(name string) (*FileWriter, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	fw := &FileWriter{
		file:   f,
		writer: bufio.NewWriter(f),
		size:   0,
	}
	return fw, nil
}

func (fw *FileWriter) Close() error {
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

// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (fw *FileWriter) WriteEncodedVarint(x uint64) error {
	for x >= 1<<7 {
		fw.writer.WriteByte(uint8(x&0x7f | 0x80))
		x >>= 7
		fw.size += 1
	}
	fw.writer.WriteByte(uint8(x))
	fw.size += 1
	return nil
}

// EncodeFixed64 writes a 64-bit integer to the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (fw *FileWriter) WriteEncodedFixed64(x uint64) error {
	fw.writer.WriteByte(uint8(x))
	fw.writer.WriteByte(uint8(x >> 8))
	fw.writer.WriteByte(uint8(x >> 16))
	fw.writer.WriteByte(uint8(x >> 24))
	fw.writer.WriteByte(uint8(x >> 32))
	fw.writer.WriteByte(uint8(x >> 40))
	fw.writer.WriteByte(uint8(x >> 48))
	fw.writer.WriteByte(uint8(x >> 56))
	fw.size += 8
	return nil
}

// EncodeFixed32 writes a 32-bit integer to the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (fw *FileWriter) WriteEncodedFixed32(x uint64) error {
	fw.writer.WriteByte(uint8(x))
	fw.writer.WriteByte(uint8(x >> 8))
	fw.writer.WriteByte(uint8(x >> 16))
	fw.writer.WriteByte(uint8(x >> 24))
	fw.size += 4
	return nil
}

// EncodeZigzag64 writes a zigzag-encoded 64-bit integer
// to the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (fw *FileWriter) WriteEncodedZigzag64(x uint64) error {
	// use signed number to get arithmetic right shift.
	fw.size += 8
	return fw.WriteEncodedVarint(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}

// EncodeZigzag32 writes a zigzag-encoded 32-bit integer
// to the Buffer.
// This is the format used for the sint32 protocol buffer type.
func (fw *FileWriter) WriteEncodedZigzag32(x uint64) error {
	// use signed number to get arithmetic right shift.
	fw.size += 4
	return fw.WriteEncodedVarint(uint64((uint32(x) << 1) ^ uint32((int32(x) >> 31))))
}

// EncodeRawBytes writes a count-delimited byte buffer to the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (fw *FileWriter) WriteEncodedRawBytes(b []byte) error {
	fw.WriteEncodedVarint(uint64(len(b)))
	for i := 0; i < len(b); i++ {
		fw.writer.WriteByte(b[i])
		fw.size += 1
	}
	return nil
}

// EncodeStringBytes writes an encoded string to the Buffer.
// This is the format used for the proto2 string type.
func (fw *FileWriter) WriteEncodedStringBytes(s string) error {
	fw.WriteEncodedVarint(uint64(len(s)))
	for i := 0; i < len(s); i++ {
		fw.writer.WriteByte(s[i])
		fw.size += 1
	}
	return nil
}

type FileReader struct {
	bytes  []byte
	Offset int
}

func newFileReader(name string) (*FileReader, error) {

	b, err := mmap.Map(name)

	if err != nil {
		return nil, err
	}

	fr := &FileReader{
		bytes:  b,
		Offset: 0,
	}
	return fr, nil
}

// errOverflow is returned when an integer is too large to be represented.
var errOverflow = errors.New("proto: integer overflow")

func (fr *FileReader) decodeVarintSlow() (x uint64, err error) {
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

// DecodeVarint reads a varint-encoded integer from the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (fr *FileReader) DecodeVarint() (x uint64, err error) {
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
func (fr *FileReader) DecodeFixed64() (x uint64, err error) {
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
func (fr *FileReader) DecodeFixed32() (x uint64, err error) {
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
func (fr *FileReader) DecodeZigzag64() (x uint64, err error) {
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
func (fr *FileReader) DecodeZigzag32() (x uint64, err error) {
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
func (fr *FileReader) DecodeRawBytes(alloc bool) (buf []byte, err error) {
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

// DecodeStringBytes reads an encoded string from the Buffer.
// This is the format used for the proto2 string type.
func (fr *FileReader) DecodeStringBytes() (s string, err error) {
	buf, err := fr.DecodeRawBytes(false)
	if err != nil {
		return
	}
	return string(buf), nil
}

// Close close and unmap the file.
func (fr *FileReader) Close() error {
	return mmap.UnMap(fr.bytes)
}
