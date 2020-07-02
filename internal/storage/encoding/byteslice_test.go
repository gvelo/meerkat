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
	"fmt"
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage/colval"
	"meerkat/internal/util/testutil"
	"testing"
)

type SliceEncFactory func(bw BlockWriter) ByteSliceEncoder

type blockWriterMock struct {
	block []byte
	rid   uint32
}

func (w *blockWriterMock) Flush() {
}


func (w *blockWriterMock) Block() []byte {
	_, n := binary.Uvarint(w.block)
	return w.block[n:]
}

func (w *blockWriterMock) WriteBlock(block []byte, baseRid uint32) {
	w.block = block
	w.rid = baseRid
}

func TestByteSliceSnappyEnc(t *testing.T) {
	testByteSliceEnc(t, func(bw BlockWriter) ByteSliceEncoder {
		return NewByteSliceSnappyEncodeer(bw)
	}, NewByteSliceSnappyDecoder())
}

func TestByteSlicePlainEnc(t *testing.T) {
	testByteSliceEnc(t, func(bw BlockWriter) ByteSliceEncoder {
		return NewByteSlicePlainEncodeer(bw)
	}, NewByteSlicePlainDecoder())
}

func testByteSliceEnc(t *testing.T, ef SliceEncFactory, d ByteSliceDecoder) {

	s := 1024

	pw := &blockWriterMock{}

	v := createRandomSliceVec(s)

	e := ef(pw)

	e.Encode(v)

	data := make([]byte, len(v.Data())*2)
	offsets := make([]int, s)

	data, offsets = d.Decode(pw.Block())
	fmt.Println("exp ",v.Data())
	fmt.Println("act ",data)
	fmt.Println(len(v.Data()))
	fmt.Println(len(data))

	assert.Equal(t, v.Offsets(), offsets, "decoded offsets doesn't match")
	assert.Equal(t, v.Data(), data, "decoded data doesn't match")

}

func createRandomSliceVec(size int) colval.ByteSliceColValues {

	var data []byte
	var offsets []int
	var rid []uint32

	for i := 0; i < size; i++ {
		data = append(data, testutil.RandomString(50)...)
		offsets = append(offsets, len(data))
		rid = append(rid, uint32(i))
	}

	return colval.NewByteSliceColValues(data, rid, offsets)

}
