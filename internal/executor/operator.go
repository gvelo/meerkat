// Copyright 2019 The Meerkat Authors
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
	"meerkat/internal/store2"
)

type Operation int

const (
	and Operation = iota
	or
	xor
	eq
)

type Operator interface {
	Init()
	Close()
}

type BitmapOperator interface {
	Operator
	Next() *roaring.Bitmap
}

type VectorOperator interface {
	Operator
	Next() store2.Vector
}

func NewBinaryBitmapOperator(op Operation, left BitmapOperator, right BitmapOperator) *BinaryBitmapOperator {
	return &BinaryBitmapOperator{
		op:    op,
		left:  left,
		right: right,
	}
}

type BinaryBitmapOperator struct {
	op    Operation
	left  BitmapOperator
	right BitmapOperator
}

func (op *BinaryBitmapOperator) Init() {
	// do nothing for now
}

func (op *BinaryBitmapOperator) Close() {
	// nothing to release here
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

}
