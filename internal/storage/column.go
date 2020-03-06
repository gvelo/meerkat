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
	"fmt"
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
	Segment(indexId []byte, from time.Time, to time.Time) []Segment
}

type LocalSegmentRegistry struct {
}

func (sr LocalSegmentRegistry) Segment(indexId []byte, from time.Time, to time.Time) []Segment {
	return nil
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
	Read(pos []uint32) (Vector, error)
	Stats() *Stats
}

type IntColumn interface {
	Column
	Dict() IntDict
	Index() IntIndex
	Iterator() IntIterator
}

type FloatColumn interface {
	Column
	Dict() FloatDict
	Index() FloatIndex
	Iterator() FloatIterator
}

type StringColumn interface {
	Column
	Dict() ByteSliceDict
	Index() ByteSliceIndex
	ReadDictEnc(pos []uint32) (IntVector, error)
	DictEncodedIterator() IntIterator
	Iterator() ByteSliceIterator
}

type TextColumn interface {
	Column
	Index() ByteSliceIndex
	Iterator() ByteSliceIterator
}

type TimeColumn interface {
	Column
	Index() TimeIndex
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
	Compressed  int
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
	AppendValue(i int)
}

type FloatVector interface {
	Vector
	ValuesAsFloat() []float64
	AppendValue(i float64)
}

type ByteSliceVector interface {
	Vector
	Data() []byte
	Offsets() []int
	Get(i int) []byte
	AppendValue(i []byte)
}

func NewIntVector() IntVector {
	return &intVector{
		rid: make([]uint32, 0),
		vec: make([]int, 0),
	}
}

func NewIntVectorFromSlice(values []int) IntVector {
	rid := make([]int, len(values))
	for i, _ := range values {
		rid[i] = i
	}
	return intVector{
		rid: make([]uint32, 0),
		vec: values,
	}
}

type intVector struct {
	vec []int
	rid []uint32
}

func (v intVector) AppendValue(i int) {
	v.vec = append(v.vec, i)
	fmt.Printf("2 %p , %v\n ", v, v)
}

func (v intVector) Len() int {
	return len(v.vec)
}

func (v intVector) Rid() []uint32 {
	return v.rid
}

func (v intVector) ValuesAsInt() []int {
	return v.vec
}

type floatVector struct {
	vec []float64
	rid []uint32
}

func NewFloatVector() FloatVector {
	return floatVector{
		rid: make([]uint32, 0),
		vec: make([]float64, 0),
	}
}

func NewFloatVectorFromSlice(values []float64) FloatVector {
	rid := make([]uint32, len(values))
	for i, _ := range values {
		rid[i] = uint32(i)
	}
	return floatVector{
		rid: rid,
		vec: values,
	}
}

func (v floatVector) AppendValue(i float64) {
	v.vec = append(v.vec, i)
}

func (v floatVector) Len() int {
	return len(v.vec)
}

func (v floatVector) Rid() []uint32 {
	return v.rid
}

func (v floatVector) ValuesAsFloat() []float64 {
	return v.vec
}

func NewByteSliceVector() ByteSliceVector {
	return &byteSliceVector{
		rid:     make([]uint32, 0),
		data:    make([]byte, 0),
		offsets: make([]int, 0),
	}
}

func NewByteSliceVectorSlice(values [][]byte) ByteSliceVector {
	rid := make([]uint32, len(values))
	offsets := make([]int, len(values))
	data := make([]byte, 0, len(values)*5) // mas o menos el maximo de los strings..
	for i, _ := range values {
		rid[i] = uint32(i)
		data = append(data, values[i]...)
		offsets[i] = len(data)
	}
	return byteSliceVector{
		rid:     rid,
		data:    data,
		offsets: offsets,
	}
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

func (v byteSliceVector) AppendValue(i []byte) {
	v.data = append(v.data, i...)
	v.offsets = append(v.offsets, len(i))
	fmt.Printf("2 %p , %v\n ", v, v)

}

func (v byteSliceVector) Get(i int) []byte {

	var start int

	if i > 0 {
		start = v.offsets[i-1]
	}

	return v.data[start:v.offsets[i]]

}
