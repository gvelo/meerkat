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

type EncodingType int

const (
	Plain EncodingType = iota
	Dict
	DictRleBitPacked
	DeltaBitPacked
	Snappy
)

type Column interface {
	Encoding() EncodingType
	Validity() *roaring.Bitmap
	HasNulls() bool
	Stats() *Stats
}

type IntColumn interface {
	Column
	Dict() IntDict
	Index() IntIndex
	Read(pos []uint32) (IntVector, error)
	Iterator() IntIterator
}

type FloatColumn interface {
	Column
	Dict() FloatDict
	Index() FloatIndex
	Read(pos []uint32) (FloatVector, error)
	Iterator() FloatIterator
}

type StringColumn interface {
	Column
	Dict() ByteSliceDict
	Index() ByteSliceIndex
	ReadDictEnc(pos []uint32) (IntVector, error)
	Read(pos []uint32) (ByteSliceVector, error)
	DictEncodedIterator() IntIterator
	Iterator() ByteSliceIterator
}

type TextColumn interface {
	Column
	Index() ByteSliceIndex
	Read(pos []uint32) (ByteSliceVector, error)
	Iterator() ByteSliceIterator
}

type TimeColumn interface {
	Column
	Index() TimeIndex
	Read(pos []uint32) (IntVector, error)
	Iterator() IntIterator
}

type Iterator interface {
	HasNext() bool
}

type IntIterator interface {
	Iterator
	Next() (IntVector, error)
}

type FloatIterator interface {
	Iterator
	Next() (FloatVector, error)
}

type ByteSliceIterator interface {
	Iterator
	Next() (ByteSliceVector, error)
}

type IntDict interface {
	DecodeInt(id int) (int, error)
}

type FloatDict interface {
	DecodeFloat(id int) (float64, error)
}

type ByteSliceDict interface {
	DecodeByteSlice(i int) ([]byte, error)
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

type Vector interface {
	Len() int
	Rid() []uint32
}

type IntVector interface {
	Vector
	ValuesAsInt() []int
}

type FloatVector interface {
	Vector
	ValuesAsFloat() []float64
}

type ByteSliceVector interface {
	Vector
	Data() []byte
	Offsets() []int
	Get(i int) []byte
}

type intVector struct {
	vect []int
	rid  []uint32
}

func (v intVector) Len() int {
	return len(v.vect)
}

func (v intVector) Rid() []uint32 {
	return v.rid
}

func (v intVector) ValuesAsInt() []int {
	return v.vect
}

type floatVector struct {
	vect []float64
	rid  []uint32
}

func (v floatVector) Len() int {
	return len(v.vect)
}

func (v floatVector) Rid() []uint32 {
	return v.rid
}

func (v floatVector) ValuesAsFloat() []float64 {
	return v.vect
}

type byteSliceVector struct {
	rid     []uint32
	data    []byte
	offsets []int
}

func (v byteSliceVector) Len() int {
	return len(v.offsets)
}

func (v byteSliceVector) Rid() []uint32 {
	return v.rid
}

func (v byteSliceVector) Data() []byte {
	return v.data
}

func (v byteSliceVector) Offsets() []int {
	return v.offsets
}

func (v byteSliceVector) Get(i int) []byte {

	var start int

	if i > 0 {
		start = v.offsets[i-1]
	}

	return v.data[start:v.offsets[i]]

}
