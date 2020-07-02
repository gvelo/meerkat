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
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
	"testing"
)

func TestSortOperator_build_diff_array(t *testing.T) {

	v := vector.NewIntVector([]int{5, 6, 6, 55, 22}, []uint64{4})
	v.SetLen(5)

	v1 := vector.NewIntVector([]int{5, 6, 6, 55, 22}, nil)
	v1.SetLen(5)

	v2 := vector.NewIntVector([]int{1, 2, 3, 5, 8}, nil)
	v2.SetLen(5)

	v3 := vector.NewIntVector([]int{1, 3, 2, 8, 5}, nil)
	v3.SetLen(5)

	testCases := []struct {
		name        string
		order       []int
		expFalseIdx []int
		v           vector.IntVector
	}{
		{
			name:        "Test 1 w/nulls",
			order:       []int{0, 1, 2, 3, 4},
			v:           v,
			expFalseIdx: []int{1, 4},
		}, {
			name:        "Test 1 w/o nulls",
			order:       []int{0, 1, 2, 3, 4},
			v:           v1,
			expFalseIdx: []int{2},
		}, {
			name:        "Test 3",
			order:       []int{0, 1, 2, 3, 4},
			v:           v2,
			expFalseIdx: []int{},
		}, {
			name:        "Test 4",
			order:       []int{0, 2, 1, 4, 3},
			v:           v3,
			expFalseIdx: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			b := make([]bool, len(tc.order))
			buildDiffIntVector(tc.v, tc.order, b)

			fId := 0
			for i, it := range b {
				if fId < len(tc.expFalseIdx) && tc.expFalseIdx[fId] == i {
					assert.Equal(t, it, false)
					fId++
				} else {
					assert.Equal(t, it, true)
				}
			}

		})
	}

}

func TestSortOperator(t *testing.T) {

	// Set up child
	vec := make([][]interface{}, 0)

	lv := make([]interface{}, 0)
	timeArr, _ := createTimeArray("17-05-2020 19:18:00 -03:00", 5, "1m")
	v := vector.NewIntVector(timeArr, nil)
	v.SetLen(5)
	lv = append(lv, v)

	v1 := createNotNullStringVector([]string{"Hola1", "Hola2", "Hola3", "Hola4", "Hola5"}, nil)
	v1.SetLen(5)
	lv = append(lv, v1)

	v2 := vector.NewIntVector([]int{1, 2, 1, 2, 8}, nil)
	v2.SetLen(5)
	lv = append(lv, v2)

	v3 := vector.NewIntVector([]int{10, 1, 5, 1, 8}, nil)
	v3.SetLen(5)
	lv = append(lv, v3)

	vec = append(vec, lv)

	f := &fakeMultiVectorOperator{
		vec: vec,
	}
	// Set up cf
	sMap := make(map[string]storage.Column)
	sMap["_ts"] = &fakeIntColumn{}
	sMap["hola"] = &fakeStringColumn{}
	sMap["n1"] = &fakeIntColumn{}
	sMap["n2"] = &fakeIntColumn{}

	// Create ctx
	ctx := NewContext(nil, 100)

	opts := []SortOpt{
		{
			"n1",
			true,
		}, {
			"n2",
			false,
		},
	}
	op := NewSortOperator(ctx, f, opts)
	op2 := NewColumnToRowOperator(ctx, op)

	op2.Init()
	for op2.Next() != nil {
	}

}

func TestSortOperator2(t *testing.T) {

	// Set up child
	vec := make([][]interface{}, 0)

	lv := make([]interface{}, 0)
	timeArr, _ := createTimeArray("17-05-2020 19:18:00 -03:00", 5, "1m")
	v := vector.NewIntVector(timeArr, nil)
	v.SetLen(5)
	lv = append(lv, v)

	v1 := createNotNullStringVector([]string{"Hola1", "Hola2", "Hola3", "Hola4", "Hola5"}, nil)
	v1.SetLen(5)
	lv = append(lv, v1)

	v2 := vector.NewIntVector([]int{1, 2, 1, 2, 8}, []uint64{6})
	v2.SetLen(5)
	lv = append(lv, v2)

	v3 := vector.NewIntVector([]int{10, 1, 5, 1, 8}, []uint64{6})
	v3.SetLen(5)
	lv = append(lv, v3)

	vec = append(vec, lv)

	f := &fakeMultiVectorOperator{
		vec: vec,
	}
	// Set up cf
	sMap := make(map[string]storage.Column)
	sMap["_ts"] = &fakeIntColumn{}
	sMap["hola"] = &fakeStringColumn{}
	sMap["n1"] = &fakeIntColumn{}
	sMap["n2"] = &fakeIntColumn{}

	// Create ctx
	ctx := NewContext(nil, 100)

	opts := []SortOpt{
		{
			"n1",
			true,
		}, {
			"n2",
			false,
		},
	}
	op := NewSortOperator(ctx, f, opts)
	op2 := NewColumnToRowOperator(ctx, op)

	op2.Init()
	for op2.Next() != nil {
	}

}
