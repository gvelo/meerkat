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
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage"
	"meerkat/internal/utils"
	"testing"
)

type EncFactory func(pw storage.PageWriter) storage.ByteSliceEncoder

type pageWriterMock struct {
	page []byte
	rid  uint32
}

func (w *pageWriterMock) Flush() error {
	return nil
}

func (w *pageWriterMock) WritePage(page []byte, endRid uint32) error {
	w.page = page
	w.rid = endRid
	return nil
}

func Test(t *testing.T) {

	t.Run("plain", func(t *testing.T) {
		testByteSliceEnc(t, func(pw storage.PageWriter) storage.ByteSliceEncoder {
			return NewByteSlicePlainEncodeer(pw)
		}, NewByteSlicePlainDecoder())
	})

	t.Run("snappy", func(t *testing.T) {
		testByteSliceEnc(t, func(pw storage.PageWriter) storage.ByteSliceEncoder {
			return NewByteSliceSnappyEncodeer(pw)
		}, NewByteSliceSnappyDecoder())
	})

}

func testByteSliceEnc(t *testing.T, ef EncFactory, d storage.ByteSliceDecoder) {

	s := 1024

	pw := &pageWriterMock{}

	v := createRandomSliceVec(s)

	e := ef(pw)

	err := e.Encode(v)

	if err != nil {
		t.Error(err)
		return
	}

	data := make([]byte, len(v.Data())*2)
	offsets := make([]int, s)

	data, offsets, err = d.Decode(pw.page, data, offsets)
	assert.Equal(t, v.Data(), data, "decoded data doesn't match")
	assert.Equal(t, v.Offsets(), offsets, "decoded offsets doesn't match")

}

func createRandomSliceVec(size int) storage.ByteSliceVector {

	var data []byte
	var offsets []int
	var rid []uint32

	for i := 0; i < size; i++ {
		data = append(data, utils.RandomString(50)...)
		offsets = append(offsets, len(data))
		rid = append(rid, uint32(i))
	}

	return storage.NewByteSliceVector(rid, data, offsets)

}
