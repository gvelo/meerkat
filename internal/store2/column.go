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

package store2

import (
	"github.com/RoaringBitmap/roaring"
)

type Encoding int

const (
	Plain Encoding = iota
	Dict
	DictRleBitPacked
	DeltaBitPacked
	Snappy
)

type PageType int

const (
	Plain PageType = iota
	BitPack
	Run
)

type Column interface {
	Encoding() Encoding
	Nulls() *roaring.Bitmap
	Scan() Iterator
	Page(row int) Page
	ColumnIndex() Index
	Stats() *Stats
	Dictionary() Dictionary
}

type Page interface {
	Type() PageType
	Size() int // value count
	Len() int  // byte len
	Bytes() []byte
	toValue ??
	oneValue ??
}

type Iterator interface {
	HasNext() bool
	Next() Page
}

type Stats struct {
	Len int
	Size int
	Cardinality int
	Compresed
}