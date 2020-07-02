// Copyright 2019 The Meerkat Authors
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
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage/encoding"
	"meerkat/internal/storage/vector"
	"time"
)

type SegmentRegistry interface {
	Segment(indexId []byte, from *time.Time, to *time.Time) []*Segment
}

type Column interface {
	Encoding() encoding.EncodingType
	Validity() *roaring.Bitmap
	HasNulls() bool
	Stats() *Stats
}

type Int64Column interface {
	Column
	Index() Int64Index
	Reader() Int64ColumnReader
	Iterator() Int64Iterator
}

type Int64ColumnReader interface {
	// TODO(gvelo): hint the reader about index use.
	//  ie. avoid index use in low selectivity search.
	Read(rid []uint32) vector.Int64Vector
}

type Float64Column interface {
	Column
	Index() FloatIndex
	Read(pos []uint32) vector.Float64Vector
	Iterator() Float64Iterator
}

type Float64ColumnReader interface {
	Read(pos []uint32) vector.Float64Vector
}

type ByteSliceColumn interface {
	Column
	Dict() ByteSliceDict
	Index() ByteSliceIndex
	DictEncReader() Int64ColumnReader
	Reader() ByteSliceReader
	Iterator() ByteSliceIterator
	DictEncIterator() Int64Iterator
}

type ByteSliceReader interface {
	Read(pos []uint32) vector.ByteSliceVector
}

type TextColumn interface {
	Column
	Index() ByteSliceIndex
	Reader() ByteSliceReader
	Iterator() ByteSliceIterator
}

type TimeColumn interface {
	Column
	Index() TimeIndex
	Reader() Int64ColumnReader
	Iterator() Int64Iterator
}

type Iterator interface {
	HasNext() bool
}

type Int64Iterator interface {
	Iterator
	Next() vector.Int64Vector
}

type Float64Iterator interface {
	Iterator
	Next() vector.Float64Vector
}

type ByteSliceIterator interface {
	Iterator
	Next() vector.ByteSliceVector
}

type IntDict interface {
	DecodeInt(id int) int
}

type FloatDict interface {
	DecodeFloat(id int) float64
}

type ByteSliceDict interface {
	DecodeByteSlice(i int) []byte
}

type ByteSliceIndex interface {
	Regex(s []byte) *roaring.Bitmap
	Prefix(s []byte) *roaring.Bitmap
	Search(s []byte) *roaring.Bitmap
}

type Int64Index interface {
	Eq(i int64) *roaring.Bitmap
	Ne(i int64) *roaring.Bitmap
	Gt(i int64) *roaring.Bitmap
	Ge(i int64) *roaring.Bitmap
	Lt(i int64) *roaring.Bitmap
	Le(i int64) *roaring.Bitmap
}

type FloatIndex interface {
	Eq(f float64) *roaring.Bitmap
	Ne(f float64) *roaring.Bitmap
	Gt(f float64) *roaring.Bitmap
	Ge(f float64) *roaring.Bitmap
	Lt(f float64) *roaring.Bitmap
	Le(f float64) *roaring.Bitmap
}

type TimeIndex interface {
	TimeRange(start int, end int) (startPos, endPos int)
	TimeRangeAsBitmap(start int, end int) *roaring.Bitmap
}

type Stats struct {
	Size        int // no compressed
	Cardinality int
	Compresed   int // size compressed
	Max         interface{}
	Min         interface{}
}
