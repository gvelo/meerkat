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

//func NewColumWriterOld(columnType ColumnType, buf buffer.Buffer, perm []int, bw *io.BinaryWriter) ColumnWriter {
//
//	blkIdx := index.NewBlockIndexWriter(bw)
//	blkWriter := NewBlockWriter(bw, blkIdx)
//
//	var validity index.ValidityIndexWriter
//
//	if buf.Nullable() {
//		validity = index.NewValidityBitmapIndex(bw)
//	}
//
//	switch columnType {
//	case ColumnType_INT64:
//		// TODO(gvelo): for plain encoded scalars blockindex is not necessary.
//		//              offset can be computed using the RID plus the
//		//              column's value width.
//		src := bufcolval.NewIntBufColSource(buf.(*buffer.IntBuffer), blockLen, perm)
//		enc := encoding.NewInt64PlainEncoder(blkWriter)
//
//		return NewInt64ColumnWriter(ColumnType_INT64, src, enc, nil, blkIdx, validity, bw)
//
//	case ColumnType_STRING:
//		src := bufcolval.NewByteSliceBufColSource(buf.(*buffer.ByteSliceBuffer), txtBlockSize, perm)
//		enc := encoding.NewByteSliceSnappyEncodeer(blkWriter)
//		return NewByteSliceColumnWriter(ColumnType_STRING, src, enc, nil, blkIdx, validity, bw)
//
//	default:
//		panic("unknown columnType")
//
//	}
//
//}

//func NewTSColumnWriter(buf *buffer.IntBuffer, bw *io.BinaryWriter) ColumnWriter {
//
//	// plain 8k pages
//	src := bufcolval.NewTsBufColSource(buf, blockLen)
//	blkIdx := index.NewBlockIndexWriter(bw)
//	blkWriter := NewBlockWriter(bw, blkIdx)
//	enc := encoding.NewInt64PlainEncoder(blkWriter)
//	cw := NewInt64ColumnWriter(ColumnType_TIMESTAMP, src, enc, nil, blkIdx, nil, bw)
//
//	return cw
//
//}

func NewColumnWriter(info ColumnSourceInfo, segmentSrc SegmentSource, bw *io.BinaryWriter) ColumnWriter {

	switch info.ColumnType {

	case ColumnType_TIMESTAMP:

		blkIdx := index.NewBlockIndexWriter(bw)
		blkWriter := NewBlockWriter(bw, blkIdx)
		enc := encoding.NewInt64PlainEncoder(blkWriter)
		src := segmentSrc.ColumnSource(info.Name, blockLen).(Int64ColumnSource)

		return NewInt64ColumnWriter(ColumnType_TIMESTAMP,
			src,
			enc,
			nil,
			blkIdx,
			nil, bw)

	case ColumnType_DATETIME:
		panic("not implemented yet")

	case ColumnType_BOOL:
		panic("not implemented yet")

	case ColumnType_INT32:
		panic("not implemented yet")

	case ColumnType_INT64:

		// TODO(gvelo): for plain encoded numeric columns, blockindex is not necessary.
		// offset can be computed using the RID plus the column's value width.

		blkIdx := index.NewBlockIndexWriter(bw)
		blkWriter := NewBlockWriter(bw, blkIdx)
		src := segmentSrc.ColumnSource(info.Name, blockLen).(Int64ColumnSource)
		enc := encoding.NewInt64PlainEncoder(blkWriter)
		var validity index.ValidityIndexWriter

		if info.Nullable {
			validity = index.NewValidityBitmapIndex(bw)
		}

		return NewInt64ColumnWriter(ColumnType_INT64, src, enc, nil, blkIdx, validity, bw)

	case ColumnType_FLOAT64:
		panic("not implemented yet")

	case ColumnType_STRING:
		blkIdx := index.NewBlockIndexWriter(bw)
		blkWriter := NewBlockWriter(bw, blkIdx)
		src := segmentSrc.ColumnSource(info.Name, txtBlockSize).(ByteSliceColumnSource)
		enc := encoding.NewByteSliceSnappyEncodeer(blkWriter)
		var validity index.ValidityIndexWriter

		if info.Nullable {
			validity = index.NewValidityBitmapIndex(bw)
		}

		return NewByteSliceColumnWriter(ColumnType_STRING, src, enc, nil, blkIdx, validity, bw)

	case ColumnType_DYNAMIC:
		panic("not implemented yet")

	case ColumnType_GUID:
		panic("not implemented yet")

	default:
		panic("unknown column type")
	}

}
