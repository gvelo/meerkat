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

const (
	// 8k pages
	indexPageSize = 8192

	// 682 items per leaf node (8k/(size uint32 + size int64)) 1 int64 remainder
	blocksPerLeafNode = indexPageSize / (8 + 4)

	// 2048 items per node (8k/(size uint32 )
	blocksPerNode = indexPageSize / 4

	ridLeafSize    = 4 * blocksPerLeafNode
	offsetLeafSize = 8 * blocksPerLeafNode
)

var padding [indexPageSize]byte

type nodeWriter func(start, end int) error

type pageIndexWriter struct {
	blockOffsets []int
	rid          []uint32
	maxBlockSize int
	offset       int
	blockSize    int
	firstRid     uint32
	blockStart   int
	bw           *io.BinaryWriter
	levelOffset  []int
}

func NewPageIndexWriter(bw *io.BinaryWriter) *pageIndexWriter {
	return &pageIndexWriter{
		bw:           bw,
		maxBlockSize: indexPageSize,
	}
}

func (p *pageIndexWriter) Flush() (int, error) {

	if p.blockSize != 0 {
		p.appendBlock()
	}

	err := p.writeLevels()

	if err != nil {
		return -1, err
	}

	return p.writeMetadata()

}

func (p *pageIndexWriter) IndexPages(page []byte, rid uint32) {

	if p.blockSize == 0 {
		p.firstRid = rid
		p.blockStart = p.offset
	}

	pageSize := len(page)

	p.blockSize = p.blockSize + pageSize

	if p.blockSize >= p.maxBlockSize {
		p.appendBlock()
	}

	p.offset = p.offset + pageSize

}

func (p *pageIndexWriter) appendBlock() {
	p.blockOffsets = append(p.blockOffsets, p.blockStart)
	p.rid = append(p.rid, p.firstRid)
	p.blockSize = 0
}

func (p *pageIndexWriter) writeMetadata() (int, error) {

	entry := p.bw.Offset

	err := p.bw.WriteVarInt(len(p.levelOffset))

	if err != nil {
		return 0, err
	}

	for i := len(p.levelOffset) - 1; i >= 0; i-- {

		err = p.bw.WriteVarInt(p.levelOffset[i])

		if err != nil {
			return 0, err
		}

	}

	return entry, nil

}

func (p *pageIndexWriter) writeLevels() error {

	level := 0

	for {

		var nextLevel []uint32
		var err error

		if level == 0 {
			nextLevel, err = p.writeLevel(blocksPerLeafNode, p.leafNodeWriter)
		} else {

			if len(p.rid) == 1 {
				break
			}

			nextLevel, err = p.writeLevel(blocksPerNode, p.nodeWriter)
		}

		if err != nil {
			return err
		}

		p.rid = nextLevel
		level++

	}

	return nil

}

func (p *pageIndexWriter) leafNodeWriter(start, end int) error {
	return p.writeLeafNode(p.rid[start:end], p.blockOffsets[start:end])
}

func (p *pageIndexWriter) nodeWriter(start, end int) error {
	return p.writeNode(p.rid[start:end])
}

func (p *pageIndexWriter) writeLevel(itemsPerNode int, writer nodeWriter) ([]uint32, error) {

	p.levelOffset = append(p.levelOffset, p.bw.Offset)

	n := 0

	for i := 0; i < len(p.rid); i = i + itemsPerNode {

		end := i + itemsPerNode

		if end > len(p.rid) {
			end = len(p.rid)
		}

		err := writer(i, end)

		if err != nil {
			return nil, err
		}

		p.rid[n] = p.rid[i]

		n++

	}

	return p.rid[:n], nil

}

func (p *pageIndexWriter) writeLeafNode(rid []uint32, offsets []int) error {

	ridBytes := utils.UInt32AsByte(rid)

	_, err := p.bw.Write(ridBytes)

	if err != nil {
		return err
	}

	padLen := ridLeafSize - len(ridBytes)

	if padLen > 0 {

		_, err = p.bw.Write(padding[:padLen])

		if err != nil {
			return err
		}

	}

	offsetBytes := utils.IntAsByte(offsets)

	_, err = p.bw.Write(offsetBytes)

	if err != nil {
		return err
	}

	padLen = indexPageSize - (ridLeafSize + len(offsetBytes))

	_, err = p.bw.Write(padding[:padLen])

	if err != nil {
		return err
	}

	return nil

}

func (p *pageIndexWriter) writeNode(rid []uint32) error {

	ridBytes := utils.UInt32AsByte(rid)

	_, err := p.bw.Write(ridBytes)

	if err != nil {
		return err
	}

	padLen := indexPageSize - len(ridBytes)

	if padLen > 0 {

		_, err = p.bw.Write(padding[:padLen])

		if err != nil {
			return err
		}

	}

	return nil

}
