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
	"meerkat/internal/schema"
	"meerkat/internal/storage"
	"testing"
)

// Select * From Index Where c1 = 12 and ts Between  F1 AND F2 limit 100
// Logs
//  | where _ts > ago(8h)
//  | where log contains "Error"
/*
func TestQueryKusto(t *testing.T) {
	// a := assert.New(t)

	var s storage.ColumnFinder

	ctx := NewContext(s)
	ctx.Value(ColumnIndexToColumnName, []string{"F", "C1", "C2", "C3"})


	// TODO check this ugly parameter
	op1 := NewTimeColumnScanOperator(ctx, Between, 11, 11, "_ts", 1000,false)
	op2 := NewStringColumnScanOperator(ctx, Contains, "Error", "message",1000, false)

	op3 := NewBinaryBitmapOperator(ctx, And, op1, op2) // ts > 1 AND ts < 2
	op4 := NewIndexScanOperator(ctx, Eq, 12, "C1")     // C1 == 12

	op5 := NewBinaryBitmapOperator(ctx, And, op3, op4) // ts > 1 AND ts < 2 AND C1 == 12

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

}*/

// Select * From Index Where c1 = 12 and ts Between  F1 AND F2 limit 100
// C1, F1 indexed
// C2 = string
// C3 = float
func TestQuery1(t *testing.T) {
	// a := assert.New(t)

	var s storage.ColumnFinder

	fields := []field{
		{
			name:     "c1",
			t:        schema.FieldType_FLOAT,
			nullable: false,
		},
		{
			name:     "c2",
			t:        schema.FieldType_INT,
			nullable: false,
		},
		{
			name:     "c2",
			t:        schema.FieldType_STRING,
			nullable: false,
		},
	}

	ii := createIndexInfo("Logs", fields...)

	ctx := NewContext(s, ii, 100)
	print(ctx)
	/*
		op1 := NewIndexScanOperator(ctx, Lt, 1, "ts") // ts > 1
		op2 := NewIndexScanOperator(ctx, Gt, 2, "ts") // ts < 2

		op3 := NewBinaryBitmapOperator(ctx, And, op1, op2) // ts > 1 AND ts < 2
		op4 := NewIndexScanOperator(ctx, Eq, 12, "C1")     // C1 == 12

		op5 := NewBinaryBitmapOperator(ctx, And, op3, op4) // ts > 1 AND ts < 2 AND C1 == 12

		op6 := NewReaderOperator(ctx, op5, "F")
		op7 := NewReaderOperator(ctx, op5, "C1")
		op8 := NewReaderOperator(ctx, op5, "C2")
		op9 := NewReaderOperator(ctx, op5, "C3")

		op10 := NewBufferOperator(ctx, []VectorOperator{op6, op7, op8, op9})

		op11 := NewMaterializeOperator(ctx, op10, nil)

		// TODO: Check where to put this operator. it should be set in the Query node.
		// op12 := NewLimitOperator(ctx, op11, 100)

		// op12.Next() // Should return the values [F] [C1] [C2] [C3]
		// limits the request.
		op11.Next()
		// op5 := Decompress
	*/

}

// Hacer un agrupado!
// Select max(C3), C1, C2 From Index Where ts Between F1 AND F2
// Group by C1, C2
// C1, ts indexed
func TestQuery2(t *testing.T) {
	/* a := assert.New(t)

	var s storage.ColumnFinder

	fields := []field{
		{
			name:     "c1",
			t:        schema.FieldType_FLOAT,
			nullable: false,
		},
		{
			name:     "c2",
			t:        schema.FieldType_INT,
			nullable: false,
		},
		{
			name:     "c2",
			t:        schema.FieldType_STRING,
			nullable: false,
		},
	}

	ii := createIndexInfo("Logs", fields...)

	ctx := NewContext(s, ii)

	ctx.Value(ColumnIndexToColumnName, []string{"_ts"}) // este operador va a ser el mas complejo, materializa bufferea, ect ect tiene que saber que hacer en todos los campos
	// quizas los histogramas y demas los podemos mapear como liteners o cosas por el estilo.
	// tenenmos que ver si pasar un contexto, para abajo para compltar cosas, ejemplo Limit. ...
	op1 := NewIndexScanOperator(ctx, Lt, 1, "_ts") // _ts > 1
	op2 := NewIndexScanOperator(ctx, Gt, 2, "_ts") // _ts < 2

	op3 := NewBinaryBitmapOperator(ctx, And, op1, op2) // _ts > 1 AND _ts < 2

	col := []string{"_ts", "F", "C1", "C2", "C4"}
	vo := make([]VectorOperator, 0)
	for _, it := range col[1:] {
		op := NewReaderOperator(ctx, op3, it)
		vo = append(vo, op)
	}

	op10 := NewBufferOperator(ctx, vo , nil)

	agList := make([]Aggregation, 0)
	agList = append(agList, Aggregation{Max, 3})

	 NewHashAggregateOperator(ctx, op10, agList, []int{1, 2})
	/*
	op12 := NewMaterialize(ctx, op11)

	op12.Next() // Should return the values [F] [C1] [C2] [C3]
	// limits the request.
	*/
	// op5 := Decompress

}
