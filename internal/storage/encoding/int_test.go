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
	"math/rand"
	"meerkat/internal/storage"
	"testing"
)

type IntEncFactory func(bw storage.BlockWriter) storage.IntEncoder

func TestIntEncoding(t *testing.T) {

	testIntEncoding(t, func(bw storage.BlockWriter) storage.IntEncoder {
		return NewIntPlainEncoder(bw)
	}, NewIntPlainDecoder())

}

func testIntEncoding(t *testing.T, f IntEncFactory, d storage.IntDecoder) {

	s := 1024

	bw := &blockWriterMock{}

	v := createRandomIntVec(s)

	e := f(bw)

	err := e.Encode(v)

	if err != nil {
		t.Error(err)
		return
	}

	data := make([]int, len(v.Data())*2)

	data, err = d.Decode(bw.block, data)
	assert.Equal(t, v.Values(), data, "decoded data doesn't match")

}

func createRandomIntVec(size int) storage.IntVector {

	var data []int
	var rid []uint32

	for i := 0; i < size; i++ {
		data = append(data, rand.Int())
		rid = append(rid, uint32(i))
	}

	return storage.NewIntVector(data, rid)

}
