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

import "github.com/RoaringBitmap/roaring"

// NewBinaryBitmapOperator creates a new bitmap binary operator.
func NewBinaryBitmapOperator(ctx Context, op BinaryOperation, left BitmapOperator, right BitmapOperator) *BinaryBitmapOperator {
	return &BinaryBitmapOperator{
		ctx:   ctx,
		op:    op,
		left:  left,
		right: right,
	}
}

// BinaryBitmapOperator executes a binary operation between two bitmaps
// and returns a new bitmap.
type BinaryBitmapOperator struct {
	ctx   Context
	op    BinaryOperation
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
	panic("Operator not supported")
}
