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
	"meerkat/internal/storage/colval"
	"testing"
)

type IntEncFactory func(bw BlockWriter) Int64Encoder

func TestIntEncoding(t *testing.T) {

	testCases := []struct {
		name    string
		factory IntEncFactory
		decoder Int64Decoder
	}{
		{
			name: "plain_int_encoding",
			factory: func(bw BlockWriter) Int64Encoder {
				return NewInt64PlainEncoder(bw)
			},
			decoder: NewInt64PlainDecoder(),
		},
		{
			name: "dd_int_encoding",
			factory: func(bw BlockWriter) Int64Encoder {
				return NewInt64DdEncoder(bw)
			},
			decoder: NewInt64DdDecoder(),
		},
	}

	// RUN TC
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testIntEncoding(t, tc.factory, tc.decoder)
		})
	}

}

func testIntEncoding(t *testing.T, f IntEncFactory, d Int64Decoder) {

	s := 1024

	bw := &blockWriterMock{}

	v := createTestIntColVal(s)

	e := f(bw)

	e.Encode(v)

	data := make([]int64, len(v.Values()))

	data = d.Decode(bw.block, data)

	assert.Equal(t, v.Values(), data, "decoded data doesn't match")

}

func createTestIntColVal(size int) colval.Int64ColValues {

	var data []int64
	var rid []uint32
	top := 100
	for i := 0; i < size; i++ {
		data = append(data, int64(top-i))
		rid = append(rid, uint32(i))
	}

	return colval.NewInt64ColValues(data, rid)

}
