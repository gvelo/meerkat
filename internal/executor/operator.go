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
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage"
	"sync"

	//"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
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

// TODO(sebad): check these operators.

func NewBufferOperator(ctx Context, child Uint32Operator, filter []string) Operator {
	return &BufferOperator{
		ctx: ctx,
		//children: children,
	}
}

// BufferOperator reads all positions in the bitmap
type BufferOperator struct {
	ctx      Context
	vKeys    [][]byte
	children []VectorOperator
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
	}
}

// MaterializeOperator operator
type MaterializeOperator struct {
	ctx    Context
	child  Uint32Operator
	filter []string
	cols   []interface{}
}

func (op *MaterializeOperator) Init() {
	op.child.Init()
	op.cols = make([]interface{}, 0)
	if op.filter != nil {
		for _, it := range op.filter {
			op.cols = append(op.cols, op.ctx.Segment().Col(it))
		}
	} else {
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
	// TODO(sebad) put a TimeOut.
	var wg sync.WaitGroup

	n := op.child.Next()

	if n == nil {
		return nil
	}

	for ; n != nil; n = op.child.Next() {

		for i := 0; i < len(op.cols); i++ {
			switch c := op.cols[i].(type) {
			case storage.IntColumn:
				wg.Add(1)
				go func(res []interface{}, wg sync.WaitGroup) {
					println("Runn 0")
					res = append(res, c.Reader().Read(n))
					wg.Done()
				}(res, wg)
			case storage.StringColumn:
				wg.Add(1)
				go func(res []interface{}, wg sync.WaitGroup) {
					println("Runn 1")
					res = append(res, c.Reader().Read(n))
					wg.Done()
				}(res, wg)
			case storage.FloatColumn:
				wg.Add(1)
				go func(res []interface{}, wg sync.WaitGroup) {
					defer wg.Done()
					println("Runn 2")
					res = append(res, c.Reader().Read(n))
				}(res, wg)
			case storage.TimeColumn:
				wg.Add(1)
				println("Runn 3")
				go func(res []interface{}, wg sync.WaitGroup) {
					defer wg.Done()
					res = append(res, c.Reader().Read(n))
				}(res, wg)
			}

		}

		wg.Wait()

	}

	if len(res) == 0 {
		return nil
	}
	return res

}
