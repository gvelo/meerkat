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

// NewBinaryBitmapOperator creates a new bitmap binary operator.
func NewIntIndexScanOperator(op Operation, value int, idx storage.IntIndex) BitmapOperator {
	return &IntIndexScanOperator{
		op:  op,
		idx: idx,
	}
}

// IndexScanOperator executes a binary operation between two bitmaps
// and returns a new bitmap.
type IntIndexScanOperator struct {
	op  Operation
	idx storage.Column
}

func (op *IntIndexScanOperator) Init() {

}

func (op *IntIndexScanOperator) Destroy() {

}

func (op *IntIndexScanOperator) Next() *roaring.Bitmap {

	//TODO: What do we return in this case?
	return nil
}

// NewColumnScanOperator creates a ColumnScanOperator
func NewColumnScanOperator(p []int, c interface{}) MultiVectorOperator {
	return &ColumnScanOperator{
		p: p,
		c: c,
	}
}

//
// ColumnScanOperator takes a array of positions and a condition
// it scans a non indexed column, search for that condition in the positions
// provided and returns 2 Vectors:
//
// 1 . position vector that meet the conditions
// 2 . values that meet the conditions
//
// if the arrays of positions is empty it scan all
//
type ColumnScanOperator struct {
	c interface{}
}

func (op *ColumnScanOperator) Init() {
}

func (op *ColumnScanOperator) Destroy() {
}

func (op *ColumnScanOperator) Next() []storage.Vector {

	//TODO: What do we return in this case?
	return nil
}

// ByteArrayScanOperator scans a non indexed column, search for the []pos
// registers and returns the bitmap that meets that condition.
//
type ByteArrayScanOperator struct {
	pos []int
}

func (op *ByteArrayScanOperator) Init() {
}

func (op *ByteArrayScanOperator) Destroy() {
}

func (op *ByteArrayScanOperator) Next() *roaring.Bitmap {
	//TODO: What do we return in this case?
	return nil
}
