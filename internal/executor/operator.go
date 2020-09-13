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

// TODO: hacer un operador por tipo asi sacamos los switchs.

import (
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
	"meerkat/internal/util/sliceutil"
	"sync"
	"time"
)

// BinaryOperation represents an operation between two expressions
type BinaryOperation int

const TsIndex = 0

const (
	And BinaryOperation = iota
	Or
	Xor
)

// Operation represents an operation between two expressions
type Operation int

// ComparisonOperation represents an comparison Between two expressions
type ComparisonOperation int

const (
	Eq ComparisonOperation = iota
	Lt
	Le
	Gt
	Ge
	Ne
	Contains
	Between
	IsNull
	Rex
	Pref
)

// Operator represents an Operator of a physical plan execution.
// Operators take one or more inputs and produce an output in the form
// of vectors, bitmaps or some other type.
type Operator interface {
	// Init initializes the Operator acquiring the required resources.
	// Init will call the init method on all it's input operators.
	Init()

	// Destroy the Operator releasing all the acquired resources.
	// Destroy will cascade calling the Destroy method on all it's
	// children operators.
	Destroy()
}

// ColumnOperator represents an Operator which acts to a column
// GetName returns the column name.
type ColumnOperator interface {
	GetName() []byte
}

// BitmapOperator is an Operator that produces bitmaps as output
type BitmapOperator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() *roaring.Bitmap
}

// Uint32Operator is an Operator that []uint32 bitmaps as output
type Uint32Operator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() []uint32
}

// VectorOperator is an Operator that produces Vector as output
type VectorOperator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() interface{}
}

// MultiVectorOperator is an Operator that produces a Vector array as output
type MultiVectorOperator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() []interface{}
}

// MultiVectorOperator is an Operator that produces a Vector array as output
type StringOperator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	Next() [][]string
}

func NewMaterializeOperator(ctx Context, child Uint32Operator, filter []string) *MaterializeOperator {
	return &MaterializeOperator{
		ctx:    ctx,
		child:  child,
		filter: filter,
		log:    log.With().Str("src", "MaterializeOperator").Logger(),
	}
}

// MaterializeOperator operator
type MaterializeOperator struct {
	ctx    Context
	child  Uint32Operator
	filter []string
	cols   []interface{}
	log    zerolog.Logger
}

func (op *MaterializeOperator) Init() {
	op.child.Init()
	op.cols = make([]interface{}, 0)
}

func (op *MaterializeOperator) Destroy() {
	op.child.Destroy()
}

func (op *MaterializeOperator) Next() []interface{} {

	res := make([]interface{}, 0)

	for x := 0; x < len(op.cols); x++ {
		res = append(res, nil)
	}

	// TODO(sebad) put a TimeOut.
	var wg sync.WaitGroup
	wg.Add(len(op.cols))
	n := op.child.Next()

	if n == nil {
		return nil
	}

	for ; n != nil; n = op.child.Next() {

		for i := 0; i < len(op.cols); i++ {
			switch c := op.cols[i].(type) {
			case storage.Int64Column:
				go func(x int) {
					res[x] = c.Reader().Read(n)
					wg.Done()
				}(i)
			case storage.ByteSliceColumn:
				go func(x int) {
					res[x] = c.Reader().Read(n)
					wg.Done()
				}(i)
			case storage.Float64Column:
				go func(x int) {
					res[x] = c.Read(n)
					wg.Done()
				}(i)
			case storage.TimeColumn:
				go func(x int) {
					res[x] = c.Reader().Read(n)
					wg.Done()
				}(i)
			}

		}

		wg.Wait()

	}

	if len(res) == 0 {
		return nil
	}

	return res

}

// TODO: pass a writer to write the table...
func NewColumnToRowOperator(ctx Context, child MultiVectorOperator) *ColumnToRowOperator {
	op := &ColumnToRowOperator{
		ctx:   ctx,
		child: child,
		log:   log.With().Str("src", "ColumnToRowOperator").Logger(),
	}
	op.IntFormat = "%d"
	op.FloatFormat = "%9.2f"
	op.TimeFormat = "2006-01-02T15:04:05"
	return op
}

// ColumnToRowOperator reads all positions in the bitmap
type ColumnToRowOperator struct {
	ctx         Context
	child       MultiVectorOperator
	FloatFormat string
	IntFormat   string
	TimeFormat  string
	log         zerolog.Logger
}

func (op *ColumnToRowOperator) Init() {
	op.child.Init()
}

func (op *ColumnToRowOperator) Destroy() {
	op.child.Destroy()
}

func (op *ColumnToRowOperator) Next() [][]string {

	n := op.child.Next()

	if len(n) == 0 {
		return nil
	}

	vv := n[TsIndex].(vector.Int64Vector)
	l := vv.Len()

	res := make([][]string, 0, l)
	for i := 0; i < l; i++ {
		row := make([]string, 0, len(n))
		for x := 0; x < len(n); x++ {
			row = append(row, op.get(i, x, n))
			fmt.Printf(" %s ", op.get(i, x, n))
		}
		fmt.Println("")
		res = append(res, row)
	}

	return res
}

func (op *ColumnToRowOperator) get(i, x int, v []interface{}) string {
	switch t := v[x].(type) {
	case vector.Int64Vector:
		if t.HasNulls() && !t.IsValid(i) {
			return "null"
		}
		if x == 0 { // 0 is always ts.
			time := time.Unix(0, int64(t.Values()[i]))
			return fmt.Sprintf(time.Format(op.TimeFormat))
		} else {
			return fmt.Sprintf(op.IntFormat, t.Values()[i])
		}
	case vector.Float64Vector:
		if t.HasNulls() && !t.IsValid(i) {
			return "null"
		}
		return fmt.Sprintf(op.FloatFormat, t.Values()[i])
	case vector.ByteSliceVector:
		if t.HasNulls() && !t.IsValid(i) {
			return "null"
		}
		return fmt.Sprintf("'%v'", sliceutil.BS2S(t.Get(i)))
	case vector.BoolVector:
		if t.HasNulls() && !t.IsValid(i) {
			return "null"
		}
		return fmt.Sprintf("%v", t.Values()[i])
	default:
		op.log.Error().Err(fmt.Errorf("type %v not found", t))
	}
	return ""
}

//TODO: check it
func getLen(n interface{}) int {
	switch v := n.(type) {
	case vector.Int64Vector:
		return v.Len()
	case vector.BoolVector:
		return v.Len()
	case vector.Float64Vector:
		return v.Len()
	case vector.ByteSliceVector:
		return v.Len()
	default:
		panic(fmt.Errorf("type %v not found", v))
	}

	return 0
}
