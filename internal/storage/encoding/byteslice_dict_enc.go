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

package encoding

import (
	"encoding/binary"
	"meerkat/internal/storage/colval"
	"meerkat/internal/storage/io"
	"meerkat/internal/util/sliceutil"
)

type ByteSliceDictEncoder struct {
	bw          BlockWriter
	buf         *io.Buffer
	keyBuffer   *io.Buffer
	w           *io.BinaryWriter
	dict        map[string]int64
	lastKey     int64
	dictStart   int
	rowTotal    uint64
	dictOffsets []int
}

func NewByteSliceDictEncoder(bw BlockWriter, w *io.BinaryWriter) *ByteSliceDictEncoder {
	return &ByteSliceDictEncoder{
		bw:          bw,
		w:           w,
		lastKey:     0,
		dict:        make(map[string]int64),
		keyBuffer:   io.NewBuffer(64 * 1042),
		dictStart:   0,
		dictOffsets: make([]int, 0),
	}
}

func (e *ByteSliceDictEncoder) Flush() {

	dictValuesStart := 0
	e.w.WriteRaw(e.keyBuffer.Data())

	dictKeyStart := len(e.keyBuffer.Data())
	offsets := sliceutil.I2B(e.dictOffsets)
	e.w.WriteRaw(offsets)

	entry := dictKeyStart + len(offsets)
	e.w.WriteUVarint64(uint64(dictValuesStart))
	e.w.WriteUVarint64(uint64(dictKeyStart))

	e.w.WriteFixedUint64(uint64(entry))

}

func (e *ByteSliceDictEncoder) FlushBlocks() {

}

func (e *ByteSliceDictEncoder) Type() Type {
	return Dict
}

func (e *ByteSliceDictEncoder) Encode(v colval.ByteSliceColValues) {

	e.buf = io.NewBuffer(v.Len()*binary.MaxVarintLen64 + binary.MaxVarintLen64)

	// left enough room at the beginning of the block to write the
	// header ( block length encoded as uvarint )
	e.buf.Pos(binary.MaxVarintLen64)

	// if
	if len(v.Data()) > e.keyBuffer.Available() {
		e.keyBuffer.Grow((e.keyBuffer.Cap() + len(v.Data())) * 2)
	}

	for i := 0; i < v.Len(); i++ {

		key := string(v.Get(i))
		id, ok := e.dict[key]
		e.rowTotal++
		if ok {
			e.buf.WriteVarInt64(id)
		} else {
			e.dict[key] = e.lastKey
			e.buf.WriteVarInt64(e.lastKey)
			e.keyBuffer.WriteBytes([]byte(key))
			e.dictOffsets = append(e.dictOffsets, e.keyBuffer.GetPos())
			e.lastKey++
		}

	}

	blockEnd := e.buf.GetPos()

	blockSize := blockEnd - binary.MaxVarintLen64

	headerOffset := binary.MaxVarintLen64 - io.SizeUVarint(uint64(blockSize))

	e.buf.PutIntAsUVarInt(headerOffset, blockSize)

	block := e.buf.Buf()[headerOffset:blockEnd]

	e.bw.WriteBlock(block, v.Rid()[0])

}
