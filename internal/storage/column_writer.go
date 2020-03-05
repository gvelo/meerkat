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

//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@column_types.tmpldata column_writer.gen.go.tmpl

package storage

import (
	"meerkat/internal/buffer"
	"meerkat/internal/schema"
	"meerkat/internal/storage/bufcolval"
	"meerkat/internal/storage/encoding"
	"meerkat/internal/storage/index"
	"meerkat/internal/storage/io"
)

type ColumnWriter interface {
	Write()
}

func NewColumWriter(fieldType schema.FieldType, buf buffer.Buffer, perm []int, bw *io.BinaryWriter) ColumnWriter {

	blkIdx := index.NewBlockIndexWriter(bw)
	blkWriter := NewBlockWriter(bw, blkIdx)

	switch fieldType {
	case schema.FieldType_INT:
		src := bufcolval.NewIntBufColSource(buf.(*buffer.IntBuffer), 8*1024, perm)
		enc := encoding.NewIntPlainEncoder(blkWriter)
		return NewIntColumnWriter(schema.FieldType_INT, src, enc, nil, blkIdx, nil, bw)

	case schema.FieldType_STRING:
		src := bufcolval.NewByteSliceBufColSource(buf.(*buffer.ByteSliceBuffer), 64*1024, perm)
		enc := encoding.NewByteSliceSnappyEncodeer(blkWriter)
		return NewByteSliceColumnWriter(schema.FieldType_STRING, src, enc, nil, blkIdx, nil, bw)

	default:
		panic("unknown fieldType")

	}

}

func NewTSColumnWriter(buf *buffer.IntBuffer, bw *io.BinaryWriter) ColumnWriter {

	// plain 8k pages
	src := bufcolval.NewTsBufColSource(buf, 8*1024)
	blkIdx := index.NewBlockIndexWriter(bw)
	blkWriter := NewBlockWriter(bw, blkIdx)
	enc := encoding.NewIntPlainEncoder(blkWriter)
	cw := NewIntColumnWriter(schema.FieldType_TIMESTAMP, src, enc, nil, blkIdx, nil, bw)

	return cw

}
