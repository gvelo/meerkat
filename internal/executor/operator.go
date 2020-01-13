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
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage"
)

// Operation represents an operation between two expressions
type Operation int

const (
	and Operation = iota
	or
	xor
	eq
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

// BitmapOperator is an Operator that produces bitmaps as output
type BitmapOperator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() *roaring.Bitmap
}

// VectorOperator is an Operator that produces Vector as output
type VectorOperator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() storage.Vector
}

// MultiVectorOperator is an Operator that produces a Vector array as output
type MultiVectorOperator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() []storage.Vector
}

// NewBinaryBitmapOperator creates a new bitmap binary operator.
func NewBinaryBitmapOperator(op Operation, left BitmapOperator, right BitmapOperator) *BinaryBitmapOperator {
	return &BinaryBitmapOperator{
		op:    op,
		left:  left,
		right: right,
	}
}

// BinaryBitmapOperator executes a binary operation between two bitmaps
// and returns a new bitmap.
type BinaryBitmapOperator struct {
	op    Operation
	left  BitmapOperator
	right BitmapOperator
}

func (op *BinaryBitmapOperator) Init() {
	op.left.Init()
	op.right.Init()
}

func (op *BinaryBitmapOperator) Destroy() {
	op.left.Destroy()
	op.right.Destroy()
}

func (op *BinaryBitmapOperator) Next() *roaring.Bitmap {

	// parallelize
	l := op.left.Next()
	r := op.left.Next()

	switch op.op {
	case and:
		return roaring.And(l, r)
	case or:
		return roaring.Or(l, r)
	case xor:
		return roaring.Xor(l, r)
	}
	//TODO: What do we return in this case?
	return nil
}
