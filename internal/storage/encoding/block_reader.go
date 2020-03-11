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

import "meerkat/internal/storage/io"

type Block struct {
	// the block bytes
	bytes []byte
	// the num of values on the block
	l int
}

func (b *Block) Bytes() []byte {
	return b.Bytes()
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
	br  *io.BinaryReader
	end int
}

func (r *ByteSliceBlockReader) ReadBlock(offset int) Block {
	r.br.SetOffset(offset)
	return r.readBlock()
}

func (r *ByteSliceBlockReader) HasNext() bool {
	return r.br.Offset() < r.end
}

func (r *ByteSliceBlockReader) Next() Block {
	return r.readBlock()
}

func (r *ByteSliceBlockReader) readBlock() Block {
	b := Block{}
	s := r.br.ReadUVarint()
	b.l = r.br.ReadUVarint()
	b.bytes = r.br.ReadSlice(r.br.Offset(), r.br.Offset()+s)
	if r.br.Offset() > r.end {
		panic("read out of column bounds")
	}
	return b
}

func NewByteSliceBlockReader(start int, end int, bytes []byte) *ByteSliceBlockReader {

	br := io.NewBinaryReader(bytes)
	br.SetOffset(start)

	b := &ByteSliceBlockReader{
		end: end,
		br:  br,
	}

	return b
}

type ScalarPlainBlockReader struct {
	br        *io.BinaryReader
	end       int
	blockSize int
	blockLen  int
}

func NewScalarPlainBlockReader(start int,
	end int,
	bytes []byte,
	blockLen int) *ScalarPlainBlockReader {

	br := io.NewBinaryReader(bytes)
	br.SetOffset(start)

	b := &ScalarPlainBlockReader{
		br:        br,
		end:       end,
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
	return r.br.Offset() < r.end
}

func (r *ScalarPlainBlockReader) readBlock() Block {

	if r.br.Offset() == r.end {
		panic("read out of column bounds")
	}

	b := Block{
		l: r.blockLen,
	}

	blockEnd := r.br.Offset() + r.blockSize

	if blockEnd > r.end {
		blockEnd = r.end
		b.l = (blockEnd - r.br.Offset())
		if (b.l % 8) != 0 {
			panic("error reading block")
		}
	}

	b.bytes = r.br.ReadSlice(r.br.Offset(), blockEnd)

	return b
}
