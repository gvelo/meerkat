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

package storage

import (
	"meerkat/internal/buffer"
	"meerkat/internal/schema"
	"meerkat/internal/storage/io"
)

type Flushable interface {
	Flush()
}
type IndexWriter interface {
	Flushable
	Cardinality() int
}

type IntIndexWriter interface {
	IndexWriter
	Index(vector IntVector)
}

type UintIndexWriter interface {
	IndexWriter
	Index(vector IntVector)
}

type FloatIndexWriter interface {
	IndexWriter
	Index(vector FloatVector)
}

type ByteSliceIndexWriter interface {
	IndexWriter
	Index(vector ByteSliceVector)
}

type BlockIndexWriter interface {
	Flushable
	IndexBlock(block []byte, baseRID uint32)
}

type BlockWriter interface {
	Flushable
	WriteBlock(block []byte, baseRid uint32)
}

type ValidityIndexWriter interface {
	IndexWriter
	Index(rid []uint32)
}

type Encoder interface {
	Flushable
	FlushBlocks()
	Type() EncodingType
}

type IntEncoder interface {
	Encoder
	Encode(vec IntVector)
}

type IntDecoder interface {
	Decode(block []byte, buf []int) []int
}

type UintEncoder interface {
	Encoder
	Encode(vec IntVector)
}

type FloatEncoder interface {
	Encoder
	Encode(vec FloatVector)
}

type ByteSliceEncoder interface {
	Encoder
	Encode(vec ByteSliceVector)
}

type ByteSliceDecoder interface {
	Decode(block []byte, data []byte, offsets []int) ([]byte, []int)
}

type ColumnWriter interface {
	Write()
}

func NewColumWriter(fieldType schema.FieldType, buf buffer.Buffer, perm []int, bw *io.BinaryWriter) ColumnWriter {
	return nil
}

func NewTSColumnWriter(buf *buffer.IntBuffer, perm []int, bw *io.BinaryWriter) ColumnWriter {

	// TODO: here should be the logic of the column writer factory.
	//  The factory will build the src using an appropriate page size
	//  to feed the encoder.

	return nil

}
