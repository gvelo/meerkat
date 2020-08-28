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
	"container/heap"
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage/vector"
	"testing"
)

func TestMergeOperator(t *testing.T) {

	testCases := []struct {
		name   string
		sz     int
		input  []MultiVectorOperator
		exp    [][]interface{}
		result []interface{}
	}{
		{
			name: "k3_2_Batch",
			sz:   20,
			exp: [][]interface{}{{
				vector.NewInt64Vector([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 150, 200, 250, 300, 350}, nil),
				vector.NewFloat64Vector([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 150, 200, 250, 300, 350}, nil),
				vector.NewInt64Vector([]int64{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 1, 2, 3, 4, 5, 6, 7, 8}, nil),
			}, {
				vector.NewInt64Vector([]int64{350, 400, 450, 500}, nil),
				vector.NewFloat64Vector([]float64{350, 400, 450, 500}, nil),
				vector.NewInt64Vector([]int64{5, 6, 7, 8}, nil),
			}},
			result: []interface{}{
				vector.NewInt64Vector([]int64{}, nil),
				vector.NewFloat64Vector([]float64{}, nil),
				vector.NewInt64Vector([]int64{}, nil),
			},
			input: []MultiVectorOperator{
				&fakeMultiVectorOperator{
					vec: [][]interface{}{
						{
							vector.NewInt64Vector([]int64{0, 2, 4, 6, 8, 10, 12, 14}, nil),
							vector.NewFloat64Vector([]float64{0, 2, 4, 6, 8, 10, 12, 14}, nil),
							vector.NewInt64Vector([]int64{1, 2, 3, 4, 5, 6, 7, 8}, nil),
						},
					},
				},
				&fakeMultiVectorOperator{
					vec: [][]interface{}{
						{
							vector.NewInt64Vector([]int64{1, 3, 5, 7, 9, 11, 13, 15}, nil),
							vector.NewFloat64Vector([]float64{1, 3, 5, 7, 9, 11, 13, 15}, nil),
							vector.NewInt64Vector([]int64{1, 2, 3, 4, 5, 6, 7, 8}, nil),
						},
					},
				},
				&fakeMultiVectorOperator{
					vec: [][]interface{}{
						{
							vector.NewInt64Vector([]int64{150, 200, 250, 300, 350, 400, 450, 500}, nil),
							vector.NewFloat64Vector([]float64{150, 200, 250, 300, 350, 400, 450, 500}, nil),
							vector.NewInt64Vector([]int64{1, 2, 3, 4, 5, 6, 7, 8}, nil),
						},
					},
				},
			},
		},
		{
			name: "k3_2_Batch_shorter",
			sz:   10,
			exp: [][]interface{}{{
				vector.NewInt64Vector([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, nil),
				vector.NewFloat64Vector([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, nil),
				vector.NewInt64Vector([]int64{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6}, nil),
			}, {
				vector.NewInt64Vector([]int64{12, 13, 14, 15, 150, 200, 250, 300, 350}, nil),
				vector.NewFloat64Vector([]float64{12, 13, 14, 15, 150, 200, 250, 300, 350}, nil),
				vector.NewInt64Vector([]int64{6, 7, 7, 8, 8, 1, 2, 3, 4, 5, 6, 7, 8}, nil),
			}},
			result: []interface{}{
				vector.NewInt64Vector([]int64{}, nil),
				vector.NewFloat64Vector([]float64{}, nil),
				vector.NewInt64Vector([]int64{}, nil),
			},
			input: []MultiVectorOperator{
				&fakeMultiVectorOperator{
					vec: [][]interface{}{
						{
							vector.NewInt64Vector([]int64{0, 2, 4, 6, 8, 10, 12, 14}, nil),
							vector.NewFloat64Vector([]float64{0, 2, 4, 6, 8, 10, 12, 14}, nil),
							vector.NewInt64Vector([]int64{1, 2, 3, 4, 5, 6, 7, 8}, nil),
						},
					},
				},
				&fakeMultiVectorOperator{
					vec: [][]interface{}{
						{
							vector.NewInt64Vector([]int64{1, 3, 5, 7, 9, 11, 13, 15}, nil),
							vector.NewFloat64Vector([]float64{1, 3, 5, 7, 9, 11, 13, 15}, nil),
							vector.NewInt64Vector([]int64{1, 2, 3, 4, 5, 6, 7, 8}, nil),
						},
					},
				},
				&fakeMultiVectorOperator{
					vec: [][]interface{}{
						{
							vector.NewInt64Vector([]int64{150, 200, 250}, nil),
							vector.NewFloat64Vector([]float64{150, 200, 250}, nil),
							vector.NewInt64Vector([]int64{1, 2, 3, 4, 5}, nil),
						},
					},
				},
			},
		},
		{
			name: "k2_1-2Batches",
			sz:   10,
			exp: [][]interface{}{{
				vector.NewInt64Vector([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, nil),
				vector.NewFloat64Vector([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, nil),
				vector.NewInt64Vector([]int64{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6}, nil),
			}, {
				vector.NewInt64Vector([]int64{12, 13, 14, 15, 150, 200, 250, 300, 350}, nil),
				vector.NewFloat64Vector([]float64{12, 13, 14, 15, 150, 200, 250, 300, 350}, nil),
				vector.NewInt64Vector([]int64{6, 7, 7, 8, 8, 1, 2, 3, 4, 5, 6, 7, 8}, nil),
			}},
			result: []interface{}{
				vector.NewInt64Vector([]int64{}, nil),
				vector.NewFloat64Vector([]float64{}, nil),
				vector.NewInt64Vector([]int64{}, nil),
			},
			input: []MultiVectorOperator{
				&fakeMultiVectorOperator{
					vec: [][]interface{}{
						{
							vector.NewInt64Vector([]int64{0, 2, 4, 6, 8, 10, 12, 14}, nil),
							vector.NewFloat64Vector([]float64{0, 2, 4, 6, 8, 10, 12, 14}, nil),
							vector.NewInt64Vector([]int64{1, 2, 3, 4, 5, 6, 7, 8}, nil),
						},
						{
							vector.NewInt64Vector([]int64{0, 2, 4, 6, 8, 10, 12, 14}, nil),
							vector.NewFloat64Vector([]float64{0, 2, 4, 6, 8, 10, 12, 14}, nil),
							vector.NewInt64Vector([]int64{1, 2, 3, 4, 5, 6, 7, 8}, nil),
						},
					},
				},
				&fakeMultiVectorOperator{
					vec: [][]interface{}{
						{
							vector.NewInt64Vector([]int64{150, 200, 250}, nil),
							vector.NewFloat64Vector([]float64{150, 200, 250}, nil),
							vector.NewInt64Vector([]int64{1, 2, 3, 4, 5}, nil),
						},
					},
				},
			},
		},
	}

	// RUN TC
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctx := NewContext(nil, tc.sz)
			op1 := NewMergeOperator(ctx, tc.input, tc.result)
			op1.Init()
			n := op1.Next()
			total := 0
			idx := 0
			for ; n != nil; n = op1.Next() {

				l := getLen(n[0].(interface{}))

				for i := 0; i < l; i++ {
					// TODO check other values.
					vv := tc.exp[idx][0].(vector.Int64Vector)
					vp := &vv
					ts := vp.Get(i)

					vv1 := tc.exp[idx][0].(vector.Int64Vector)
					vp1 := &vv1
					ts1 := vp1.Get(i)

					assert.Equal(t, ts, ts1, "Not the same values")

				}
				idx++
				total = total + l
			}
		})
	}

}

func TestMinHeap(t *testing.T) {

	h := &MinHeap{
		&Item{
			listFrom: 1,
			value:    3,
		},
		&Item{
			listFrom: 1,
			value:    0,
		},
		&Item{
			listFrom: 1,
			value:    2,
		},
	}

	heap.Init(h)
	heap.Push(h, &Item{
		listFrom: 5,
		value:    1,
	})

	i := 0
	for h.Len() > 0 {
		it := h.Pop().(*Item)
		assert.Equal(t, int64(i), it.value)
		i++
	}

}
