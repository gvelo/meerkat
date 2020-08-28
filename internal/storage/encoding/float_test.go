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

type FloatEncFactory func(bw BlockWriter) Float64Encoder

func TestFloatEncoding(t *testing.T) {

	testCases := []struct {
		name    string
		factory FloatEncFactory
		decoder Float64Decoder
	}{
		{
			name: "plain_float_encoding",
			factory: func(bw BlockWriter) Float64Encoder {
				return NewFloat64PlainEncoder(bw)
			},
			decoder: NewFloat64PlainDecoder(),
		},
		{
			name: "xor_float_encoding",
			factory: func(bw BlockWriter) Float64Encoder {
				return NewFloat64XorEncoder(bw)
			},
			decoder: NewFloat64XorDecoder(),
		},
	}

	// RUN TC
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testFloatEncoding(t, tc.factory, tc.decoder)
		})
	}

}

func testFloatEncoding(t *testing.T, f FloatEncFactory, d Float64Decoder) {

	s := 1024

	bw := &blockWriterMock{}

	v := createRandomFloatColVal(s)

	e := f(bw)

	e.Encode(v)

	data := make([]float64, len(v.Values())*8*2)

	data = d.Decode(bw.block, data)

	assert.Equal(t, v.Values(), data, "decoded data doesn't match")

}

func createRandomFloatColVal(size int) colval.Float64ColValues {

	var data []float64
	var rid []uint32

	for i := 0; i < size; i++ {
		data = append(data, rand.Float64())
		rid = append(rid, uint32(i))
	}

	return colval.NewFloat64ColValues(data, rid)

}
