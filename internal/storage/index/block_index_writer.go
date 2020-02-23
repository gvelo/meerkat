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
	pageSize = 8192

	// 682 items per leaf node (8k/(size uint32 + size int64)) 1 int64 remainder
	leafNodeLen = pageSize / (8 + 4)

	// 2048 items per node (8k/size uint32)
	nodeLen = pageSize / 4

	ridLeafSize    = 4 * leafNodeLen
	offsetLeafSize = 8 * leafNodeLen

	// try to keep chunks of blocks 8k aligned on disk.
	maxChunkSize = 8192
)

var padding [pageSize]byte

// nodeWriter writes a node into a index page
type nodeWriter func(start, end int) error

type blockIndexWriter struct {
	chunkOffsetList []int
	baseRIDList     []uint32
	offset          int
	chunkSize       int
	baseRID         uint32
	chunkStart      int
	bw              *io.BinaryWriter
	levelOffset     []int
}

func NewBlockIndexWriter(bw *io.BinaryWriter) *blockIndexWriter {
	return &blockIndexWriter{
		bw: bw,
	}
}

func (w *blockIndexWriter) Flush() (int, error) {

	if w.chunkSize != 0 {
		w.appendBlock()
	}

	err := w.writeLevels()

	if err != nil {
		return -1, err
	}

	return w.writeMetadata()

}

func (w *blockIndexWriter) IndexBlock(block []byte, baseRID uint32) {

	if w.chunkSize == 0 {
		w.baseRID = baseRID
		w.chunkStart = w.offset
	}

	blockSize := len(block)

	w.chunkSize = w.chunkSize + blockSize

	if w.chunkSize >= maxChunkSize {
		w.appendBlock()
	}

	w.offset = w.offset + blockSize

}

func (w *blockIndexWriter) appendBlock() {
	w.chunkOffsetList = append(w.chunkOffsetList, w.chunkStart)
	w.baseRIDList = append(w.baseRIDList, w.baseRID)
	w.chunkSize = 0
}

func (w *blockIndexWriter) writeMetadata() (int, error) {

	entry := w.bw.Offset

	err := w.bw.WriteVarInt(len(w.levelOffset))

	if err != nil {
		return 0, err
	}

	for i := len(w.levelOffset) - 1; i >= 0; i-- {

		err = w.bw.WriteVarInt(w.levelOffset[i])

		if err != nil {
			return 0, err
		}

	}

	return entry, nil

}

func (w *blockIndexWriter) writeLevels() error {

	level := 0

	for {

		var nextLevel []uint32
		var err error

		if level == 0 {
			nextLevel, err = w.writeLevel(leafNodeLen, w.leafNodeWriter)
		} else {

			if len(w.baseRIDList) == 1 {
				break
			}

			nextLevel, err = w.writeLevel(nodeLen, w.nodeWriter)
		}

		if err != nil {
			return err
		}

		w.baseRIDList = nextLevel
		level++

	}

	return nil

}

func (w *blockIndexWriter) leafNodeWriter(start, end int) error {
	return w.writeLeafNode(w.baseRIDList[start:end], w.chunkOffsetList[start:end])
}

func (w *blockIndexWriter) nodeWriter(start, end int) error {
	return w.writeNode(w.baseRIDList[start:end])
}

func (w *blockIndexWriter) writeLevel(itemsPerNode int, writer nodeWriter) ([]uint32, error) {

	w.levelOffset = append(w.levelOffset, w.bw.Offset)

	n := 0

	for i := 0; i < len(w.baseRIDList); i = i + itemsPerNode {

		end := i + itemsPerNode

		if end > len(w.baseRIDList) {
			end = len(w.baseRIDList)
		}

		err := writer(i, end)

		if err != nil {
			return nil, err
		}

		w.baseRIDList[n] = w.baseRIDList[i]

		n++

	}

	return w.baseRIDList[:n], nil

}

func (w *blockIndexWriter) writeLeafNode(rid []uint32, offsets []int) error {

	ridBytes := utils.UInt32AsByte(rid)

	_, err := w.bw.Write(ridBytes)

	if err != nil {
		return err
	}

	padLen := ridLeafSize - len(ridBytes)

	if padLen > 0 {

		_, err = w.bw.Write(padding[:padLen])

		if err != nil {
			return err
		}

	}

	offsetBytes := utils.IntAsByte(offsets)

	_, err = w.bw.Write(offsetBytes)

	if err != nil {
		return err
	}

	padLen = pageSize - (ridLeafSize + len(offsetBytes))

	_, err = w.bw.Write(padding[:padLen])

	if err != nil {
		return err
	}

	return nil

}

func (w *blockIndexWriter) writeNode(rid []uint32) error {

	ridBytes := utils.UInt32AsByte(rid)

	_, err := w.bw.Write(ridBytes)

	if err != nil {
		return err
	}

	padLen := pageSize - len(ridBytes)

	if padLen > 0 {

		_, err = w.bw.Write(padding[:padLen])

		if err != nil {
			return err
		}

	}

	return nil

}
