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
	Flush() error
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

type PageIndexWriter interface {
	IndexWriter
	IndexPages(vec IntVector)
}

type PageWriter interface {
	Flushable
	//WritePage(vect ByteSliceVector)
	WritePage(page []byte, endRid uint32) error
}

type ValidityIndexWriter interface {
	IndexWriter
	Index(rid []uint32)
}

type Encoder interface {
	Flushable
	FlushPages() error
	Type() EncodingType
}

type IntEncoder interface {
	Encoder
	Encode(vec IntVector) error
}

type IntDecoder interface {
	Decode(page []byte, buf []int) ([]int, error)
}

type UintEncoder interface {
	Encoder
	Encode(vec IntVector) error
}

type FloatEncoder interface {
	Encoder
	Encode(vec FloatVector) error
}

type ByteSliceEncoder interface {
	Encoder
	Encode(vec ByteSliceVector) error
}

type ByteSliceDecoder interface {
	Decode(page []byte, data []byte, offsets []int) ([]byte, []int, error)
}

type ColumnWriter interface {
	Write() error
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
