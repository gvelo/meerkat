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

type Segment interface {
	IndexName() string
	IndexID() []byte
	From() time.Time
	To() time.Time
	Rows() int
	Col(id []byte) Column
}

type SegmentRegistry interface {
	Segment(indexId []byte, from *time.Time, to *time.Time) []Segment
}

type Col interface {
	Encoding() encoding.EncodingType
	Validity() *roaring.Bitmap
	HasNulls() bool
	Read(pos []uint32) vector.Vec
	Idx() interface{}
	Iterator()

}

type Column interface {
	Encoding() encoding.EncodingType
	Validity() *roaring.Bitmap
	HasNulls() bool
	Stats() *Stats
}

type IntColumn interface {
	Column
	Dict() IntDict
	Index() IntIndex
	// TODO(gvelo): hint the reader about index use.
	//  ie. avoid index use in low selectivity search.
	Read(pos []uint32) vector.IntVector
	Iterator() IntIterator
}

type FloatColumn interface {
	Column
	Dict() FloatDict
	Index() FloatIndex
	Read(pos []uint32) vector.FloatVector
	Iterator() FloatIterator
}

type StringColumn interface {
	Column
	Dict() ByteSliceDict
	Index() ByteSliceIndex
	ReadDictEnc(pos []uint32) vector.IntVector
	Read(pos []uint32) vector.ByteSliceVector
	DictEncodedIterator() IntIterator
	Iterator() ByteSliceIterator
}

type TextColumn interface {
	Column
	Index() ByteSliceIndex
	Read(pos []uint32) vector.ByteSliceVector
	Iterator() ByteSliceIterator
}

type TimeColumn interface {
	Column
	Index() TimeIndex
	Read(pos []uint32) vector.IntVector
	Iterator() IntIterator
}

type Iterator interface {
	HasNext() bool
}

type IntIterator interface {
	Iterator
	Next() vector.IntVector
}

type FloatIterator interface {
	Iterator
	Next() vector.FloatVector
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

type IntIndex interface {
	Eq(i int) *roaring.Bitmap
	Ne(i int) *roaring.Bitmap
	Gt(i int) *roaring.Bitmap
	Ge(i int) *roaring.Bitmap
	Lt(i int) *roaring.Bitmap
	Le(i int) *roaring.Bitmap
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
