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

package encoding

import (
	"meerkat/internal/storage/io"
)

type Block struct {
	// the block bytes
	bytes []byte
	// the num of values on the block
	l int
}

func (b *Block) Bytes() []byte {
	return b.bytes
}

func (b *Block) Len() int {
	return b.l
}

type BlockReader interface {
	ReadBlock(offset int) Block
	Next() Block
	HasNext() bool
}

type ByteSliceBlockReader struct {
	br     *io.BinaryReader
	bounds io.Bounds
}

func (r *ByteSliceBlockReader) ReadBlock(offset int) Block {
	r.br.SetOffset(offset)
	return r.readBlock()
}

func (r *ByteSliceBlockReader) HasNext() bool {
	return r.br.Offset() < r.bounds.End
}

func (r *ByteSliceBlockReader) Next() Block {
	return r.readBlock()
}

func (r *ByteSliceBlockReader) readBlock() Block {

	b := Block{}

	size := r.br.ReadUVarint()
	start := r.br.Offset()
	b.l = r.br.ReadUVarint()
	b.bytes = r.br.ReadSlice(start, start+size)

	if r.br.Offset() > r.bounds.End {
		panic("read out of column bounds")
	}

	return b

}

func NewByteSliceBlockReader(bytes []byte, bounds io.Bounds) *ByteSliceBlockReader {

	br := io.NewBinaryReader(bytes)
	br.SetOffset(bounds.Start)

	b := &ByteSliceBlockReader{
		bounds: bounds,
		br:     br,
	}

	return b

}

type ScalarPlainBlockReader struct {
	br        *io.BinaryReader
	bounds    io.Bounds
	blockSize int
	blockLen  int
}

func NewScalarPlainBlockReader(
	bytes []byte,
	bounds io.Bounds,
	blockLen int) *ScalarPlainBlockReader {

	br := io.NewBinaryReader(bytes)
	br.SetOffset(bounds.Start)

	b := &ScalarPlainBlockReader{
		br:        br,
		bounds:    bounds,
		blockSize: blockLen * 8,
		blockLen:  blockLen,
	}

	return b

}

func (r *ScalarPlainBlockReader) ReadBlock(offset int) Block {
	r.br.SetOffset(offset)
	return r.readBlock()
}

func (r *ScalarPlainBlockReader) Next() Block {
	return r.readBlock()
}

func (r *ScalarPlainBlockReader) HasNext() bool {
	return r.br.Offset() < r.bounds.End
}

func (r *ScalarPlainBlockReader) readBlock() Block {

	if r.br.Offset() >= r.bounds.End {
		panic("read out of column bounds")
	}

	b := Block{
		l: r.blockLen,
	}

	blockEnd := r.br.Offset() + r.blockSize

	if blockEnd > r.bounds.End {

		blockEnd = r.bounds.End

		size := blockEnd - r.br.Offset()

		if (size % 8) != 0 {
			panic("error reading block")
		}

		b.l = size / 8

	}

	b.bytes = r.br.ReadSlice(r.br.Offset(), blockEnd)

	return b

}
