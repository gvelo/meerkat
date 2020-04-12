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
	"meerkat/internal/schema"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
	"meerkat/internal/util/sliceutil"
	"sync"
	"time"
)

// BinaryOperation represents an operation between two expressions
type BinaryOperation int

const (
	And BinaryOperation = iota
	Or
	Xor
)

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
	Next() vector.Vector
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

func NewBufferOperator(ctx Context, child Uint32Operator, filter []string) Operator {
	return &BufferOperator{
		ctx: ctx,
		//children: children,
		log: log.With().Str("src", "BufferOperator").Logger(),
	}
}

// BufferOperator reads all positions in the bitmap
type BufferOperator struct {
	ctx      Context
	vKeys    [][]byte
	children []VectorOperator
	log      zerolog.Logger
}

func (r *BufferOperator) Init() {
	r.vKeys = make([][]byte, 0)
	for i, c := range r.children {
		r.vKeys[i] = c.(ColumnOperator).GetName()
		c.Init()
	}
	r.ctx.Value(ColumnIndexToColumnName, r.vKeys)
}

func (r *BufferOperator) Destroy() {
	for _, c := range r.children {
		c.Destroy()
	}
}

func (r *BufferOperator) Next() []vector.Vector {
	op := make([]vector.Vector, 0, len(r.children))
	for i, c := range r.children {
		// Paralelize
		op[i] = c.Next()

	}
	return op
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
	if op.filter != nil {

		fp := make([]schema.Field, 0, len(op.filter))

		for _, it := range op.filter {
			f, err := op.ctx.IndexInfo().FieldByName(it)
			if err != nil {
				op.log.Err(err)
			} else {
				fp = append(fp, f)
			}
			op.cols = append(op.cols, op.ctx.Segment().Col(it))
		}
		op.ctx.SetFieldProcessed(fp)
	} else {
		op.ctx.SetFieldProcessed(op.ctx.IndexInfo().Fields)
		for _, it := range op.ctx.IndexInfo().Fields {
			op.cols = append(op.cols, op.ctx.Segment().Col(it.Name))
		}

	}
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
			case storage.IntColumn:
				go func(x int) {
					res[x] = c.Reader().Read(n)
					wg.Done()
				}(i)
			case storage.StringColumn:
				go func(x int) {
					res[x] = c.Reader().Read(n)
					wg.Done()
				}(i)
			case storage.FloatColumn:
				go func(x int) {
					res[x] = c.Reader().Read(n)
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
	l := op.len(n[0])

	res := make([][]string, 0, l)
	for i := 0; i < l; i++ {
		row := make([]string, 0, len(n))
		for x := 0; x < len(n); x++ {
			row = append(row, op.get(i, x, n))
		}
		res = append(res, row)
	}

	return res
}

func (op *ColumnToRowOperator) get(i, x int, v []interface{}) string {
	switch t := v[x].(type) {
	case vector.IntVector:
		if op.ctx.GetFieldProcessed()[x].FieldType == schema.FieldType_TIMESTAMP {
			time := time.Unix(0, int64(t.Values()[i]))
			return fmt.Sprintf(time.Format(op.TimeFormat))
		} else {
			return fmt.Sprintf(op.IntFormat, t.Values()[i])
		}
	case vector.FloatVector:
		return fmt.Sprintf(op.FloatFormat, t.Values()[i])
	case vector.ByteSliceVector:
		return fmt.Sprintf("'%v'", sliceutil.BS2S(t.Get(i)))
	case vector.UintVector:
		return fmt.Sprintf(op.IntFormat, t.Values()[i])
	case vector.BoolVector:
		return fmt.Sprintf("%v", t.Values()[i])
	default:
		log.Printf("Type %v not found.", t)
	}
	return ""
}

// TODO: QUE MIERDA?????
func (op *ColumnToRowOperator) len(v interface{}) int {
	// No entiendo porque mierda no me toma el vector.Vector que implementan todos!?????
	switch t := v.(type) {
	case vector.IntVector:
		return t.Len()
	case vector.FloatVector:
		return t.Len()
	case vector.ByteSliceVector:
		return t.Len()
	case vector.UintVector:
		return t.Len()
	case vector.BoolVector:
		return t.Len()
	default:
		log.Printf("Type %v not found.", t)
	}
	return 0
}
