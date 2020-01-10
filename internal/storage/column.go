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

type Encoding int

const (
	Plain Encoding = iota
	Dict
	DictRleBitPacked
	DeltaBitPacked
	Snappy
)

type Column interface {
	Encoding() Encoding
	Validity() *roaring.Bitmap
	HasNulls() bool
	Stats() *Stats
}

type IntColumn interface {
	Column
	Dict() IntDict
	Index() IntIndex
	Read(pos []int) (IntVector, error)
}

type ByteColumn interface {
	Column
	Dict() ByteArrayDict
	Index() ByteArrayIndex
	ReadEnc(pos []int) (IntVector, error)
	Read(pos []int) (ByteArrayVector, error)
	IteratorEnc() IntIterator
	Iterator() ByteArrayIterator
}

type Iterator interface {
	HasNext() bool
}

type IntIterator interface {
	Iterator
	Next() (IntVector, error)
}

type ByteArrayIterator interface {
	Iterator
	Next() (ByteArrayVector, error)
}

type IntDict interface {
	DecodeInt(id int) (int, error)
}

type ByteArrayDict interface {
	DecodeByteArray(i int) ([]byte, error)
}

type ByteArrayIndex interface {
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

type Stats struct {
	Len         int
	Size        int
	Cardinality int
	Compresed   int
	Max         interface{}
	Min         interface{}
}

type Vector interface {
	Len() int
	HasNulls() bool
	Pos() []int
	ValuesAsBytes() []byte
	PosAsBytes() []byte
}

type IntVector interface {
	Vector
	ValuesAsInt() []int
}

type ByteArrayVector interface {
	Vector
	ValuesAsSlide() [][]byte
}