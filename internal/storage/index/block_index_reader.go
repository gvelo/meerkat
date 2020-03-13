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

package index

import (
	"fmt"
	"meerkat/internal/storage/io"
	"meerkat/internal/util/sliceutil"
)

// thread safe
type blockIndexReader struct {
	levelOffsets []int
	leafLevel    int
	bounds       io.Bounds
	b            []byte
}

func NewBlockIndexReader(b []byte, bounds io.Bounds) *blockIndexReader {

	br := io.NewBinaryReader(b)
	br.SetOffset(bounds.End - 8)
	br.SetOffset(br.ReadFixed64())

	levelOffsets := br.ReadVarUintSlice()
	leafLevel := len(levelOffsets) - 1

	fmt.Println("level offsets ", levelOffsets)

	return &blockIndexReader{
		b:            b,
		bounds:       bounds,
		levelOffsets: levelOffsets,
		leafLevel:    leafLevel,
	}
}

func (p *blockIndexReader) lookup(level int, pos int, rid uint32) (uint32, int) {

	if level == p.leafLevel {

		ridList, offsets := p.readLeaf(pos)

		f := findFloor(ridList, rid)
		return ridList[f], offsets[f]

	}

	ridList := p.readNode(level, pos)

	f := findFloor(ridList, rid)

	level++

	return p.lookup(level, f, rid)

}

func findFloor(s []uint32, rid uint32) int {

	var i int
	var v uint32

	for i, v = range s {

		if i != 0 && v == 0 {
			return i - 1
		}

		if v > rid {
			return i - 1
		}

	}

	return i
}

func (p *blockIndexReader) readLeaf(pos int) ([]uint32, []int) {

	pageOffset := p.levelOffsets[p.leafLevel] + (pos * pageSize)

	page := p.b[pageOffset : pageOffset+pageSize-1]

	ridList := sliceutil.B2U32(page[0:ridLeafSize])
	offsets := sliceutil.B2I(page[ridLeafSize : ridLeafSize+offsetLeafSize])

	return ridList, offsets

}

func (p *blockIndexReader) readNode(level int, pos int) []uint32 {

	pageOffset := p.levelOffsets[level] + (pos * pageSize)

	page := p.b[pageOffset : pageOffset+pageSize-1]

	ridList := sliceutil.B2U32(page)

	return ridList

}

func (p *blockIndexReader) Lookup(rid uint32) (uint32, int) {

	return p.lookup(0, 0, rid)

}
