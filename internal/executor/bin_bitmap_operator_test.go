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
	"testing"
)

type binQueryTestCase struct {
	name string
	l    Uint32Operator
	r    Uint32Operator
	exp  expected
	op   BinaryOperation
	sz   int
}

type fakeUint32Op struct {
	i int
	v [][]uint32
}

func (f *fakeUint32Op) Init() {
}

func (f *fakeUint32Op) Destroy() {

}

func (f *fakeUint32Op) Next() []uint32 {
	if f.i == len(f.v) {
		return nil
	}
	v := f.v[f.i]
	f.i++
	return v
}

func TestNewBinaryUint32Operator(t *testing.T) {

	testCases := []binQueryTestCase{
		{
			name: "compare Xor 1",
			op:   Xor,
			l: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10}},
			},
			r: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{3, 5, 7}, {8, 9, 10}},
			},
			sz: 1,
			exp: expected{
				[][]uint32{{1}, {2}, {4}, {6}},
			},
		},
		{
			name: "compare Xor 3",
			op:   Xor,
			l: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10}},
			},
			r: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{3, 5, 7}, {8, 9, 10}},
			},
			sz: 3,
			exp: expected{
				[][]uint32{{1, 2, 4}, {6}},
			},
		},
		{
			name: "compare Or 3",
			op:   Or,
			l: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10}},
			},
			r: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{3, 5, 7}, {8, 9, 10}},
			},
			sz: 3,
			exp: expected{
				[][]uint32{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}},
			},
		},
		{
			name: "compare And 1",
			op:   And,
			l: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10}},
			},
			r: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{3, 5, 7}, {8, 9, 10}},
			},
			sz: 1,
			exp: expected{
				[][]uint32{{3}, {5}, {7}, {8}, {9}, {10}},
			},
		},

		{
			name: "compare OR 1",
			op:   Or,
			l: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
			},
			r: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{3, 4, 5, 7, 8, 9, 11, 15, 18}},
			},
			sz: 30,
			exp: expected{
				[][]uint32{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
			},
		},
		{
			name: "compare AND 1",
			op:   And,
			l: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
			},
			r: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{3, 4, 5, 7, 8, 9, 11, 15, 18}},
			},
			sz: 20,
			exp: expected{
				[][]uint32{{3, 4, 5, 7, 8, 9, 11, 15, 18}},
			},
		},
		{
			name: "compare Or 1",
			op:   Or,
			l: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10}},
			},
			r: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{3, 5, 7}, {8, 9, 10}},
			},
			sz: 1,
			exp: expected{
				[][]uint32{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}},
			},
		},
		{
			name: "compare And 3",
			op:   And,
			l: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{3, 5, 7}},
			},
			r: &fakeUint32Op{
				i: 0,
				v: [][]uint32{{1, 2, 3, 4, 5, 6, 7}},
			},
			sz: 3,
			exp: expected{
				[][]uint32{{3, 5, 7}},
			},
		},
	}
	// RUN TC
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctx := NewContext(nil, nil, tc.sz)
			op1 := NewBinaryUint32Operator(ctx, tc.op, tc.l, tc.r)
			op1.Init()
			var i = 0
			n := op1.Next()
			for ; n != nil; n = op1.Next() {
				for x := 0; x < len(n); x++ {
					assert.Equal(t, tc.exp.values.([][]uint32)[i][x], n[x], "Not the same values")
				}
				i++
			}
		})
	}
}
