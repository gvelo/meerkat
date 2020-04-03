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
	"meerkat/internal/storage"
	"testing"
)

type expected struct {
	cardinality []int
	values      interface{}
}

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
	out       expected
	op        ComparisonOperation
	value     interface{}
}

func (tc *queryTestCase) init() error {
	return nil
}

func createColFinder(in interface{}) storage.ColumnFinder {

	m := make(map[string]storage.Column)
	m["intFieldId"] = NewFakeIntColumn(in)
	s := NewFakeColFinder(m)
	return s
}

func newColumnScanOperator(ctx Context, op ComparisonOperation, value interface{}, tc queryTestCase) Uint32Operator {

	if len(tc.in.validity) > 0 {
		return NewIntNullColumnScanOperator(ctx, op, value.(int), tc.fieldName, tc.batch)
	} else {
		return NewIntColumnScanOperator(ctx, op, value.(int), tc.fieldName, tc.batch)
	}

}

func TestQueryScanOperators(t *testing.T) {

	testCases := []queryTestCase{
		{
			fieldName: "intFieldId",
			name:      "Check batch 5",
			batch:     5,
			in: input{
				validity: nil,
				values:   [][]int{{-1, 4, 5}, {43, 4, 5, 7}},
			},
			out: expected{
				cardinality: []int{2, 3},
				values:      [][]uint32{{1, 2, 3, 4, 5}, {6}},
			},
			op:    gt,
			value: 1,
		},
		{
			fieldName: "intFieldId",
			name:      "Check batch 3 more input",
			batch:     3,
			in: input{
				validity: nil,
				values:   [][]int{{-1, 4, 5, 4, 4, 3}, {43, 4, 5, 7}},
			},
			out: expected{
				cardinality: []int{2, 3},
				values:      [][]uint32{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			},
			op:    gt,
			value: 1,
		},

		{
			fieldName: "intFieldId",
			name:      "Check batch 3",
			batch:     3,
			in: input{
				validity: nil,
				values:   [][]int{{-1, 4, 5}, {43, 4, 5, 7}},
			},
			out: expected{
				cardinality: []int{2, 3},
				values:      [][]uint32{{1, 2, 3}, {4, 5, 6}},
			},
			op:    gt,
			value: 1,
		},
		{
			fieldName: "intFieldId",
			name:      "Check batch 10",
			batch:     10,
			in: input{
				validity: nil,
				values:   [][]int{{-1, 4, 5, 43, 4, 5, 7}},
			},
			out: expected{
				cardinality: []int{2, 3},
				values:      [][]uint32{{1, 2, 3, 4, 5, 6}},
			},
			op:    gt,
			value: 1,
		},
	}

	// RUN TC
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			if err := tc.init(); err != nil {
				t.Fatal(err)
			}

			ctx := NewContext(createColFinder(tc.in))

			op1 := newColumnScanOperator(ctx, tc.op, tc.value, tc)
			op1.Init()
			var i = 0
			n := op1.Next()
			for ; n != nil; n = op1.Next() {
				// assert.Equal(t, tc.out.cardinality[i], len(n), "length does not match.")
				for x := 0; x < len(n); x++ {
					assert.Equal(t, n[x], tc.out.values.([][]uint32)[i][x], "Not the same values")
				}
				i++
			}
		})
	}
}

// Select * From Index Where c1 = 12 and ts Between  F1 AND F2 limit 100
// C1, F1 indexed
// C2 = string
// C3 = float
func TestQuery1(t *testing.T) {
	// a := assert.New(t)

	var s storage.ColumnFinder

	ctx := NewContext(s)
	ctx.Value(ColumnIndexToColumnName, []string{"F", "C1", "C2", "C3"})

	op1 := NewIndexScanOperator(ctx, lt, 1, "ts") // ts > 1
	op2 := NewIndexScanOperator(ctx, gt, 2, "ts") // ts < 2

	op3 := NewBinaryBitmapOperator(ctx, and, op1, op2) // ts > 1 AND ts < 2
	op4 := NewIndexScanOperator(ctx, eq, 12, "C1")     // C1 == 12

	op5 := NewBinaryBitmapOperator(ctx, and, op3, op4) // ts > 1 AND ts < 2 AND C1 == 12

	op6 := NewReaderOperator(ctx, op5, "F")
	op7 := NewReaderOperator(ctx, op5, "C1")
	op8 := NewReaderOperator(ctx, op5, "C2")
	op9 := NewReaderOperator(ctx, op5, "C3")

	op10 := NewBufferOperator(ctx, []VectorOperator{op6, op7, op8, op9})

	op11 := NewMaterialize(ctx, op10)

	// TODO: Check where to put this operator. it should be set in the Query node.
	// op12 := NewLimitOperator(ctx, op11, 100)

	// op12.Next() // Should return the values [F] [C1] [C2] [C3]
	// limits the request.
	op11.Next()
	// op5 := Decompress

}

// Hacer un agrupado!
// Select max(C3), C1, C2 From Index Where ts Between F1 AND F2
// Group by C1, C2
// C1, ts indexed
func TestQuery2(t *testing.T) {
	// a := assert.New(t)

	var s storage.ColumnFinder

	ctx := NewContext(s)
	ctx.Value(ColumnIndexToColumnName, []string{"_ts"})

	op1 := NewIndexScanOperator(ctx, lt, 1, "_ts") // _ts > 1
	op2 := NewIndexScanOperator(ctx, gt, 2, "_ts") // _ts < 2

	op3 := NewBinaryBitmapOperator(ctx, and, op1, op2) // _ts > 1 AND _ts < 2

	// este operador va a ser el mas complejo, materializa bufferea, ect ect tiene que saber que hacer en todos los campos
	// quizas los histogramas y demas los podemos mapear como liteners o cosas por el estilo.
	// tenenmos que ver si pasar un contexto, para abajo para compltar cosas, ejemplo Limit. ...

	col := []string{"_ts", "F", "C1", "C2", "C4"}
	vo := make([]VectorOperator, 0)
	for _, it := range col[1:] {
		op := NewReaderOperator(ctx, op3, it)
		vo = append(vo, op)
	}

	op10 := NewBufferOperator(ctx, vo)

	agList := make([]Aggregation, 0)
	agList = append(agList, Aggregation{Max, 3})

	op11 := NewHashAggregateOperator(ctx, op10, agList, []int{1, 2})

	op12 := NewMaterialize(ctx, op11)

	op12.Next() // Should return the values [F] [C1] [C2] [C3]
	// limits the request.

	// op5 := Decompress

}
