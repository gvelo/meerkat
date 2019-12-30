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

// Page Iteration
// Las paginas de una columna se pueden iterar en un fullscan o filtradas por un bitmap posicional.
// Los nulls se representan como bitmap posicional
// agregar filtrado de nulls

// Vectores
// - Continuos sin Nulls ( no poseen nulls ) : se representan en forma continua sin nuls
// - Continuos con Nulls : se representa como un array continuo mas un array de bools para validez
// - Posicionales ( poseen nulls ) los nulls se representan como un bitmap posicional
// - String son representados como vectores cuyo primer byte es el len.
// - Vectores de paginas ??? procesamiento comprimido.

// Ejemplo sum() posicional vs continuo
//
//
//          +---------(+)-------+
//          |                   |
//          |                   |
//   +------+-------+    +------+-------+
//   |  pos | value |    |  pos | value |
//   +------+-------+    +------+-------+
//   |  132 |   10  |    |    0 |  112  |
//   |  156 | 2334  |    |   15 |   23  |
//   | 1234 |   11  |    |  132 |  345  |
//   | 1344 |   12  |    | 1344 |  654  |
//   +------+-------+    +------+-------+

import (
	"github.com/RoaringBitmap/roaring"
	"time"
)

type Segment interface {
	Index() string
	From() time.Time
	To() time.Time
	Rows() int
	Cols() map[string]Column
}

type SegmentRegistry interface {
	Segment(indexId string, from *time.Time, to *time.Time) []Segment
}

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
	Validity() *roaring.Bitmap
	Scan() PageIterator
	Page(rows *roaring.Bitmap) PageIterator
	Index() Index
	Stats() *Stats
	Dictionary() Dictionary
}

type Dictionary interface {
	String(id int) string
}

type StringIndex interface {
	Regex(s string) *roaring.Bitmap
	Gt(i int)
	Lt()
	Eq()
}

type Page interface {
	FirstRow() int
	Type() PageType
	Size() int // value count
	Len() int  // byte len
	Bytes() []byte
	Read(p []byte) (n int, err error) //nuls ????

	// bitmaps ??
	// full ?
	// nuls ??
}
type PageIterator interface {
	HasNext() bool
	Next() []Page
	Next(p []Page)
}

type Stats struct {
	Len         int
	Size        int
	Cardinality int
	Compresed   int
	Max         interface{}
	Min         interface{}
}

type Vec struct {
	values []byte
	pos    []int
}

func (v *Vec) AsInt() []int {

}

type Batch struct {
	Vec []Vec
}
