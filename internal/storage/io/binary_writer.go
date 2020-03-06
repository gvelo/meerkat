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
	"encoding/binary"
	"os"
)

const (
	bufSize = 1024 * 1024
)

type BinaryWriter struct {
	file   *os.File
	writer *bufio.Writer
	offset int
	buf    []byte
}

func NewBinaryWriter(name string) (*BinaryWriter, error) {

	f, err := os.Create(name)

	if err != nil {
		return nil, err
	}

	bw := &BinaryWriter{
		file:   f,
		writer: bufio.NewWriterSize(f, bufSize),
		offset: 0,
		buf:    make([]byte, binary.MaxVarintLen64),
	}

	return bw, nil
}

func (bw *BinaryWriter) Close() {

	err := bw.writer.Flush()

	if err != nil {
		panic(err)
	}

	err = bw.file.Sync()

	if err != nil {
		panic(err)
	}

	err = bw.file.Close()

	if err != nil {
		panic(err)
	}

}

func (bw *BinaryWriter) Offset() int {
	return bw.offset
}

func (bw *BinaryWriter) WriteByte(x byte) {
	err := bw.writer.WriteByte(x)
	if err != nil {
		panic(err)
	}
	bw.offset++
}

// WriteUvarint write an int as an unsigned varint
func (bw *BinaryWriter) WriteUvarint(i int) {
	bw.WriteUVarint64(uint64(i))
}

// unsigned
func (bw *BinaryWriter) WriteUVarint64(x uint64) {

	n := binary.PutUvarint(bw.buf, x)

	n, err := bw.writer.Write(bw.buf[:n])

	if err != nil {
		panic(err)
	}

	bw.offset += n

}

func (bw *BinaryWriter) WriteFixedInt(i int) {
	bw.WriteFixedUint64(uint64(i))
}

func (bw *BinaryWriter) WriteFixedUint64(x uint64) {

	binary.LittleEndian.PutUint64(bw.buf, x)

	_, err := bw.writer.Write(bw.buf[:8])

	if err != nil {
		panic(err)
	}

	bw.offset += 8

}

func (bw *BinaryWriter) WriteBytes(b []byte) {

	bw.WriteUvarint(len(b))
	bw.WriteRaw(b)

}

func (bw *BinaryWriter) WriteRaw(b []byte) {

	n, err := bw.writer.Write(b)

	if err != nil {
		panic(err)
	}

	bw.offset += n

}

func (bw *BinaryWriter) Write(b []byte) (int, error) {

	n, err := bw.writer.Write(b)

	if err != nil {
		panic(err)
	}

	bw.offset += n

	return n, nil

}

func (bw *BinaryWriter) WriteString(s string) {

	l := len(s)

	bw.WriteUvarint(l)

	_, err := bw.writer.WriteString(s)

	if err != nil {
		panic(err)
	}

	bw.offset += l

}

func (bw *BinaryWriter) WriteUVarIntSlice(slice []int) {

	for _, i := range slice {
		bw.WriteUvarint(i)
	}

}
