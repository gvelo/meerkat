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

type input struct {
	validity [][]uint64
	length   []int
	values   interface{}
}

type queryTestCase struct {
	fieldName string
	name      string
	batch     int
	in        input
	out       interface{}
	op        ComparisonOperation
	value     interface{}
}

func TestQueryStringScanOperators(t *testing.T) {

	testCases := []struct {
		fieldName string
		name      string
		batch     int
		in        input
		out       interface{}
		op        ComparisonOperation
		value     interface{}
	}{
		{
			fieldName: "col",
			name:      "Check batch 3  input not Null",
			batch:     3,
			in: input{
				values: [][]string{{"bbb", "aaa1", "aaa2", "aaa3", "aaa4", "bbbb"}, {"aaa01", "aaa02", "aaa03", "aaa04"}},
			},
			out:   [][]uint32{{1, 2, 3}, {4, 6, 7}, {8, 9}},
			op:    Contains,
			value: "a",
		},
		{
			fieldName: "col",
			name:      "Check batch 3 input Null",
			batch:     3,
			in: input{
				validity: [][]uint64{{63}, {15}}, // all valid
				length:   []int{6, 6},
				values:   [][]string{{"bbb", "aaa1", "aaa2", "aaa3", "aaa4", "bbbb"}, {"aaa01", "aaa02", "aaa03", "aaa04", "23234", "234234"}},
			},
			out:   [][]uint32{{1, 2, 3}, {4, 6, 7}, {8, 9}},
			op:    Contains,
			value: "a",
		},
	}

	// RUN TC
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctx := NewContext(nil, tc.batch)

			op1 := newColumnScanOperator(ctx, tc.op, tc.value.(string), tc)
			op1.Init()
			var i = 0
			n := op1.Next()
			for ; n != nil; n = op1.Next() {
				// assert.Equal(t, tc.out.cardinality[i], len(n), "length does not match.")
				for x := 0; x < len(n); x++ {
					assert.Equal(t, tc.out.([][]uint32)[i][x], n[x], "Not the same values")
				}
				i++
			}
		})
	}
}

func (tc *queryTestCase) init() error {
	return nil
}

func newColumnScanOperator(ctx Context, op ComparisonOperation, value interface{}, tc queryTestCase) Uint32Operator {
	switch v := value.(type) {
	case int:
		return NewIntColumnScanOperator(ctx, op, v, tc.fieldName, len(tc.in.validity) > 0)
	case string:
		return NewStringColumnScanOperator(ctx, op, v, tc.fieldName, len(tc.in.validity) > 0)
	}
	return nil
}

func TestQueryIntScanOperators(t *testing.T) {

	testCases := []queryTestCase{
		{
			fieldName: "col",
			name:      "Check batch 5 Not Null",
			batch:     5,
			in: input{
				validity: nil,
				values:   [][]int{{-1, 4, 5}, {43, 4, 5, 7}},
			},
			out:   [][]uint32{{1, 2, 3, 4, 5}, {6}},
			op:    Gt,
			value: 1,
		},
		{
			fieldName: "col",
			name:      "Check batch 5 Null",
			batch:     5,
			in: input{
				validity: [][]uint64{{255}, {255}}, // all valid
				length:   []int{3, 4},
				values:   [][]int{{-1, 4, 5}, {43, 4, 5, 7}},
			},
			out:   [][]uint32{{1, 2, 3, 4, 5}, {6}},
			op:    Gt,
			value: 1,
		},
		{
			fieldName: "col",
			name:      "Check batch 3 more input Null",
			batch:     3,
			in: input{
				validity: nil,
				values:   [][]int{{-1, 4, 5, 4, 4, 3}, {43, 4, 5, 7}},
			},
			out:   [][]uint32{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			op:    Gt,
			value: 1,
		},
		{
			fieldName: "col",
			name:      "Check batch 3 more input Not Null",
			batch:     3,
			in: input{
				validity: [][]uint64{{255}, {255}}, // all valid
				length:   []int{6, 4},
				values:   [][]int{{-1, 4, 5, 4, 4, 3}, {43, 4, 5, 7}},
			},
			out:   [][]uint32{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			op:    Gt,
			value: 1,
		},
		{
			fieldName: "col",
			name:      "Check batch 3 Null",
			batch:     3,
			in: input{
				validity: nil,
				values:   [][]int{{-1, 4, 5}, {43, 4, 5, 7}},
			},
			out:   [][]uint32{{1, 2, 3}, {4, 5, 6}},
			op:    Gt,
			value: 1,
		},
		{
			fieldName: "col",
			name:      "Check batch 3 Not Null",
			batch:     3,
			in: input{
				validity: [][]uint64{{255}, {255}}, // all valid
				length:   []int{3, 4},
				values:   [][]int{{-1, 4, 5}, {43, 4, 5, 7}},
			},
			out:   [][]uint32{{1, 2, 3}, {4, 5, 6}},
			op:    Gt,
			value: 1,
		},
		{
			fieldName: "col",
			name:      "Check batch 10 Null one vec",
			batch:     10,
			in: input{
				validity: nil,
				values:   [][]int{{-1, 4, 5, 43, 4, 5, 7}},
			},
			out:   [][]uint32{{1, 2, 3, 4, 5, 6}},
			op:    Gt,
			value: 1,
		},

		{
			fieldName: "col",
			name:      "Check batch 10 Not Null",
			batch:     10,
			in: input{
				validity: nil,
				length:   []int{10, 10},
				values:   [][]int{{-1, 5, 5, 33, 51, 54, 34, 32, 23, 32}, {-1, 5, 5, 33, 51, 54, 34, 32, 2, 3}},
			},
			out:   [][]uint32{{1, 2, 3, 4, 5, 6, 7, 11, 12, 13}, {14, 15, 16, 17}},
			op:    Gt,
			value: 1,
		},
		{
			fieldName: "col",
			name:      "Check batch 10 Null",
			batch:     10,
			in: input{
				validity: [][]uint64{{255}}, // all valid but last 2 0000 1111 1111
				length:   []int{10},         // TODO: revisar esto tambien en este caso podriamos meter un null en el medio.
				values:   [][]int{{-1, 5, 5, 33, 51, 54, 34, 32, 33232, 22323233}},
			},
			out:   [][]uint32{{1, 2, 3, 4, 5, 6, 7}},
			op:    Gt,
			value: 1,
		},
	}

	// RUN TC
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctx := NewContext(nil, tc.batch)

			op1 := newColumnScanOperator(ctx, tc.op, tc.value.(int), tc)
			op1.Init()
			var i = 0
			n := op1.Next()
			for ; n != nil; n = op1.Next() {
				// assert.Equal(t, tc.out.cardinality[i], len(n), "length does not match.")
				for x := 0; x < len(n); x++ {
					assert.Equal(t, tc.out.([][]uint32)[i][x], n[x], "Not the same values")
				}
				i++
			}
		})
	}
}
