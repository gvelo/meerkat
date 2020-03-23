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

const (
	blockLen     = 1024 * 8  // numeric block length
	txtBlockSize = 1024 * 64 // txt block size in bytes
	vectorLen    = 1024 * 64 // the numeric vector len ( move to executor pakage ? )
)

type ColumnWriter interface {
	Write()
}

func NewColumWriter(fieldType schema.FieldType, buf buffer.Buffer, perm []int, bw *io.BinaryWriter) ColumnWriter {

	blkIdx := index.NewBlockIndexWriter(bw)
	blkWriter := NewBlockWriter(bw, blkIdx)

	var validity index.ValidityIndexWriter

	if buf.Nullable() {
		validity = index.NewValidityBitmapIndex(bw)
	}

	switch fieldType {
	case schema.FieldType_INT:
		// TODO(gvelo): for plain encoded scalars blockindex is not necessary.
		//              offset can be computed using the RID plus the
		//              column's value width.
		src := bufcolval.NewIntBufColSource(buf.(*buffer.IntBuffer), blockLen, perm)
		enc := encoding.NewIntPlainEncoder(blkWriter)

		return NewIntColumnWriter(schema.FieldType_INT, src, enc, nil, blkIdx, validity, bw)

	case schema.FieldType_STRING:
		src := bufcolval.NewByteSliceBufColSource(buf.(*buffer.ByteSliceBuffer), txtBlockSize, perm)
		enc := encoding.NewByteSliceSnappyEncodeer(blkWriter)
		return NewByteSliceColumnWriter(schema.FieldType_STRING, src, enc, nil, blkIdx, validity, bw)

	default:
		panic("unknown fieldType")

	}

}

func NewTSColumnWriter(buf *buffer.IntBuffer, bw *io.BinaryWriter) ColumnWriter {

	// plain 8k pages
	src := bufcolval.NewTsBufColSource(buf, blockLen)
	blkIdx := index.NewBlockIndexWriter(bw)
	blkWriter := NewBlockWriter(bw, blkIdx)
	enc := encoding.NewIntPlainEncoder(blkWriter)
	cw := NewIntColumnWriter(schema.FieldType_TIMESTAMP, src, enc, nil, blkIdx, nil, bw)

	return cw

}
