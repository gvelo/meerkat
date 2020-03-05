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
	"meerkat/internal/util/sliceutil"
)

const (
	// 8k pages
	pageSize = 8 * 1024

	// 682 items per leaf node (8k/(size uint32 + size int64)) 1 int64 remainder
	leafNodeLen = pageSize / (8 + 4)

	// 2048 items per node (8k/size uint32)
	nodeLen = pageSize / 4

	ridLeafSize    = 4 * leafNodeLen
	offsetLeafSize = 8 * leafNodeLen

	// try to keep chunks of blocks 8k aligned on disk.
	maxChunkSize = 8 * 1024
)

var padding [pageSize]byte

// nodeWriter writes a node into a index page
type nodeWriter func(start, end int)

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

func (w *blockIndexWriter) Flush() {

	if w.chunkSize != 0 {
		w.appendBlock()
	}

	w.writeLevels()

	w.writeMetadata()

}

func (w *blockIndexWriter) IndexBlock(block []byte, baseRID uint32) {

	if w.chunkSize == 0 {
		w.baseRID = baseRID
		w.chunkStart = w.offset
	}

	blockSize := len(block)

	w.chunkSize += blockSize

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

func (w *blockIndexWriter) writeMetadata() {

	entry := w.bw.Offset()

	w.bw.WriteUvarint(len(w.levelOffset))

	for i := len(w.levelOffset) - 1; i >= 0; i-- {
		w.bw.WriteUvarint(w.levelOffset[i])
	}

	w.bw.WriteFixedInt(entry)

}

func (w *blockIndexWriter) writeLevels() {

	level := 0

	for {

		var nextLevel []uint32

		if level == 0 {
			nextLevel = w.writeLevel(leafNodeLen, w.leafNodeWriter)
		} else {

			if len(w.baseRIDList) == 1 {
				break
			}

			nextLevel = w.writeLevel(nodeLen, w.nodeWriter)
		}

		w.baseRIDList = nextLevel
		level++

	}

}

func (w *blockIndexWriter) leafNodeWriter(start, end int) {
	w.writeLeafNode(w.baseRIDList[start:end], w.chunkOffsetList[start:end])
}

func (w *blockIndexWriter) nodeWriter(start, end int) {
	w.writeNode(w.baseRIDList[start:end])
}

func (w *blockIndexWriter) writeLevel(itemsPerNode int, writer nodeWriter) []uint32 {

	w.levelOffset = append(w.levelOffset, w.bw.Offset())

	n := 0

	for i := 0; i < len(w.baseRIDList); i = i + itemsPerNode {

		end := i + itemsPerNode

		if end > len(w.baseRIDList) {
			end = len(w.baseRIDList)
		}

		writer(i, end)

		w.baseRIDList[n] = w.baseRIDList[i]

		n++

	}

	return w.baseRIDList[:n]

}

func (w *blockIndexWriter) writeLeafNode(rid []uint32, offsets []int) {

	ridBytes := sliceutil.U32B(rid)

	w.bw.WriteRaw(ridBytes)

	padLen := ridLeafSize - len(ridBytes)

	if padLen > 0 {
		w.bw.WriteRaw(padding[:padLen])
	}

	offsetBytes := sliceutil.I2B(offsets)

	w.bw.WriteRaw(offsetBytes)

	padLen = pageSize - (ridLeafSize + len(offsetBytes))

	w.bw.WriteRaw(padding[:padLen])

}

func (w *blockIndexWriter) writeNode(rid []uint32) {

	ridBytes := sliceutil.U32B(rid)

	w.bw.WriteRaw(ridBytes)

	padLen := pageSize - len(ridBytes)

	if padLen > 0 {
		w.bw.WriteRaw(padding[:padLen])
	}

}
