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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewBinaryBitmapOperator creates a new bitmap binary operator.
func NewBinaryBitmapOperator(ctx Context, op BinaryOperation, left BitmapOperator, right BitmapOperator) *BinaryBitmapOperator {
	return &BinaryBitmapOperator{
		ctx:   ctx,
		op:    op,
		left:  left,
		right: right,
		log:   log.With().Str("src", "BinaryBitmapOperator").Logger(),
	}
}

// BinaryBitmapOperator executes a binary operation between two bitmaps
// and returns a new bitmap.
type BinaryBitmapOperator struct {
	ctx   Context
	op    BinaryOperation
	left  BitmapOperator
	right BitmapOperator
	log   zerolog.Logger
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

	l := op.left.Next()
	r := op.right.Next()

	switch op.op {
	case And:
		return roaring.And(l, r)
	case Or:
		return roaring.Or(l, r)
	case Xor:
		return roaring.Xor(l, r)
	}
	panic("Operator not supported")
}

type BinaryUint32Operator struct {
	ctx        Context
	op         BinaryOperation
	left       Uint32Operator
	right      Uint32Operator
	sz         int
	remainingL []uint32
	remainingR []uint32
	log        zerolog.Logger
}

// NewBinaryBitmapOperator creates a new bitmap binary operator.
func NewBinaryUint32Operator(ctx Context, op BinaryOperation, left Uint32Operator, right Uint32Operator) *BinaryUint32Operator {
	return &BinaryUint32Operator{
		ctx:   ctx,
		op:    op,
		left:  left,
		right: right,
		log:   log.With().Str("src", "BinaryUint32Operator").Logger(),
	}
}

func (op *BinaryUint32Operator) Init() {
	op.left.Init()
	op.right.Init()
}

func (op *BinaryUint32Operator) Destroy() {
	op.left.Destroy()
	op.right.Destroy()
}

func (op *BinaryUint32Operator) and(l, r []uint32) []uint32 {

	res := make([]uint32, 0, op.ctx.Sz())

	if len(l) == 0 || len(r) == 0 {
		return nil
	}

	x := 0
	i := 0
	for ; i < len(l) && x < len(r) && len(res) < op.ctx.Sz(); i++ {

		if l[i] == r[x] {
			res = append(res, l[i])
			x++
		} else {

			if l[i] > r[x] {
				for l[i] > r[x] {
					x++
				}
				if l[i] == r[x] {
					res = append(res, l[i])
					x++
				}
			}

		}

	}

	if len(res) == op.ctx.Sz() {

		if len(l) > i {
			op.remainingL = l[i:]

		} else {
			op.remainingL = nil
		}

		if len(r) > x {
			op.remainingR = r[x:]
			return res
		} else {
			op.remainingR = nil
		}

		return res
	}

	op.remainingL = nil
	op.remainingR = nil
	return res
}

func (op *BinaryUint32Operator) or(l, r []uint32) []uint32 {

	res := make([]uint32, 0, op.ctx.Sz())

	if len(l) == 0 || len(r) == 0 {

		if len(l) > op.ctx.Sz() {
			op.remainingL = l[op.ctx.Sz():]
			return l[:op.ctx.Sz()]
		}

		if len(r) > op.ctx.Sz() {
			op.remainingR = r[op.ctx.Sz():]
			return r[:op.ctx.Sz()]
		}

	}

	x := 0
	i := 0
	for ; i < len(l) && x < len(r) && len(res) < op.ctx.Sz(); i++ {
		if l[i] < r[x] {
			res = append(res, l[i])
			continue
		}

		if l[i] == r[x] {
			res = append(res, l[i])
			x++
			continue
		}

		if l[i] > r[x] {
			res = append(res, r[x])
			x++
			continue
		}
	}

	if len(res) == op.ctx.Sz() {

		if len(l) > i {
			op.remainingL = l[i:]
		} else {
			op.remainingL = nil
		}

		if len(r) > x {
			op.remainingR = r[x:]
			return res
		} else {
			op.remainingR = nil
		}

		return res
	}

	op.remainingL = nil
	op.remainingR = nil
	return res
}

func (op *BinaryUint32Operator) xor(l, r []uint32) []uint32 {

	res := make([]uint32, 0, op.ctx.Sz())

	if len(l) == 0 || len(r) == 0 {

		if len(l) > op.ctx.Sz() {
			op.remainingL = l[op.ctx.Sz():]
			return l[:op.ctx.Sz()]
		}

		if len(r) > op.ctx.Sz() {
			op.remainingR = r[op.ctx.Sz():]
			return r[:op.ctx.Sz()]
		}

	}

	x := 0
	i := 0
	for ; i < len(l) && x < len(r) && len(res) < op.ctx.Sz(); i++ {
		if l[i] < r[x] {
			res = append(res, l[i])
			continue
		}

		if l[i] == r[x] {
			x++
			continue
		}

		if l[i] > r[x] {
			res = append(res, r[x])
			x++
			continue
		}
	}

	if len(res) == op.ctx.Sz() {

		if len(l) > i {
			op.remainingL = l[i:]

		} else {
			op.remainingL = nil
		}

		if len(r) > x {
			op.remainingR = r[x:]
			return res
		} else {
			op.remainingR = nil
		}

		return res
	}

	op.remainingL = nil
	op.remainingR = nil
	return res
}

func (op *BinaryUint32Operator) Next() []uint32 {
	var res []uint32

NEXT:
	l := op.left.Next()
	r := op.right.Next()

	if len(l) > 0 && len(op.remainingL) > 0 {
		l = append(op.remainingL, l...)
	}

	if len(l) == 0 && len(op.remainingL) > 0 {
		l = op.remainingL
	}

	if len(r) > 0 && len(op.remainingR) > 0 {
		r = append(op.remainingR, r...)
	}

	if len(r) == 0 && len(op.remainingR) > 0 {
		r = op.remainingR
	}

	if res == nil && len(l) == 0 && len(r) == 0 {
		return nil
	}

	switch op.op {
	case And:
		res = append(res, op.and(l, r)...)
	case Or:
		res = append(res, op.or(l, r)...)
	case Xor:
		res = append(res, op.xor(l, r)...)
	default:
		log.Error().Msgf("Operator not supported")
	}

	if len(res) < op.ctx.Sz() {

		if l == nil && r == nil {
			return res
		}
		goto NEXT
	}

	return res

}
