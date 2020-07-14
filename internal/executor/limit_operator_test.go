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
	"github.com/magiconair/properties/assert"
	"meerkat/internal/storage/vector"
	"testing"
)

func sl(intVector vector.Int64Vector) vector.Int64Vector {
	intVector.SetLen(intVector.Cap())
	return intVector
}

func TestLimitOperator(t *testing.T) {

	testCases := []struct {
		name  string
		v     MultiVectorOperator
		exp   interface{}
		limit int
	}{
		{
			name: "Limit for 1 Vectors",
			v: &fakeMultiVectorOperator{
				vec: [][]interface{}{
					{
						sl(vector.NewInt64Vector([]int64{1, 2, 3, 4, 4, 6, 6, 6}, nil)),
					},
				},
				idx: 0,
			},
			exp:   []int{1},
			limit: 1,
		},
		{
			name: "Limit for 5 Vectors",
			v: &fakeMultiVectorOperator{
				vec: [][]interface{}{
					{
						sl(vector.NewInt64Vector([]int64{1, 2, 3, 4, 4, 6, 6, 6}, nil)),
					},
					{
						sl(vector.NewInt64Vector([]int64{8, 8, 8, 8, 8, 9, 10, 10}, nil)),
					},
					{
						sl(vector.NewInt64Vector([]int64{8, 8, 8, 8, 8, 9, 10, 10}, nil)),
					},
					{
						sl(vector.NewInt64Vector([]int64{8, 8, 8, 8, 8, 9, 10, 10}, nil)),
					},
				},
				idx: 0,
			},
			exp:   []int{8, 8, 8, 6},
			limit: 30,
		},
	}

	// RUN TC
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run(tc.name, func(t *testing.T) {

				op1 := NewLimitOperator(nil, tc.v, tc.limit)
				op1.Init()
				var i = 0
				n := op1.Next()
				total := 0
				for ; n != nil; n = op1.Next() {
					l := getLen(n[0].(interface{}))
					assert.Equal(t, l, tc.exp.([]int)[i], "Not the same values")
					total = total + l
					i++
				}

				assert.Equal(t, total, tc.limit, "Not the same values")
			})
		})
	}

}
