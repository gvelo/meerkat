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
	"meerkat/internal/storage/io"
	"meerkat/internal/utils"
)

type pageIndexReader struct {
	levelOffsets []int
	leafLevel    int
	br           *io.BinaryReader
}

func NewPageIndexReader(br *io.BinaryReader) *pageIndexReader {
	return &pageIndexReader{
		br: br,
	}
}

func (p *pageIndexReader) read() error {

	var err error

	p.levelOffsets, err = p.br.ReadVarUintSlice()

	if err != nil {
		return err
	}

	p.leafLevel = len(p.levelOffsets) - 1

	return nil

}

func (p *pageIndexReader) lookup(level int, blockNum int, rid uint32) (uint32, int) {

	if level == p.leafLevel {

		rids, offsets := p.readLeaf(blockNum)

		f := findFloor(rids, rid)
		return rids[f], offsets[f]

	}

	rids := p.readNode(level, blockNum)

	f := findFloor(rids, rid)

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

func (p *pageIndexReader) readLeaf(node int) ([]uint32, []int) {

	pageBase := p.levelOffsets[p.leafLevel] + (node * indexPageSize)

	page := p.br.Bytes()[pageBase : pageBase+indexPageSize-1]

	rid := utils.BytesAsUInt32(page[0:ridLeafSize])
	offsets := utils.BytesAsInt(page[ridLeafSize : ridLeafSize+offsetLeafSize])

	return rid, offsets

}

func (p *pageIndexReader) readNode(level int, node int) []uint32 {

	pageBase := p.levelOffsets[level] + (node * indexPageSize)

	page := p.br.Bytes()[pageBase : pageBase+indexPageSize-1]

	rid := utils.BytesAsUInt32(page)

	return rid

}

func (p *pageIndexReader) Lookup(rid uint32) (uint32, int) {

	return p.lookup(0, 0, rid)

}
