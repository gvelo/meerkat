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
	"meerkat/internal/storage/colval"
	"testing"
)

type IntEncFactory func(bw BlockWriter) Int64Encoder

func TestIntEncoding(t *testing.T) {

	testIntEncoding(t, func(bw BlockWriter) Int64Encoder {
		return NewInt64PlainEncoder(bw)
	}, NewInt64PlainDecoder())

}

func testIntEncoding(t *testing.T, f IntEncFactory, d Int64Decoder) {

	s := 1024

	bw := &blockWriterMock{}

	v := createRandomIntColVal(s)

	e := f(bw)

	e.Encode(v)

	data := make([]int, len(v.Values())*8*2)

	data = d.Decode(bw.block, data)

	assert.Equal(t, v.Values(), data, "decoded data doesn't match")

}

func createRandomIntColVal(size int) colval.IntColValues {

	var data []int
	var rid []uint32

	for i := 0; i < size; i++ {
		data = append(data, rand.Int())
		rid = append(rid, uint32(i))
	}

	return colval.NewIntColValues(data, rid)

}
