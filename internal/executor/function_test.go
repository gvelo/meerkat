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

package executor

import (
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage/vector"
	"testing"
)

func TestUnaryOperator(t *testing.T) {

	testCases := []struct {
		name string
		fn   Fn
		v    VectorOperator
		exp  interface{}
	}{
		{
			name: "fStrLen",
			fn:   &fStrLen{},
			v: &fakeVectorOperator{
				vec: []interface{}{
					vector.NewByteSliceVectorFromByteArray([][]byte{[]byte("COSO"), []byte("COSO1"), []byte("COSO12")}, nil),
				},
				idx: 0,
			},
			exp: []int{4, 5, 6},
		},
		{
			name: "fStrLen",
			fn:   &fStrLen{},
			v: &fakeVectorOperator{
				vec: []interface{}{
					vector.NewByteSliceVectorFromByteArray([][]byte{[]byte("COSO"), []byte("COSO1"), []byte("COSO12")}, nil),
				},
				idx: 0,
			},
			exp: []int{4, 5, 6},
		},
	}

	// RUN TC
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			op1 := NewUnaryOperator(nil, tc.v, tc.fn)
			op1.Init()
			var i = 0
			n := op1.Next()
			total := 0
			for ; n != nil; n = op1.Next() {
				l := getLen(n)
				// TODO: make a function to check.
				assert.Equal(t, l, tc.exp.([]int)[i], "Not the same values")
				total = total + l
				i++
			}

		})
	}

}
