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
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
)

func selectIntOpFn(op ComparisonOperation) func(x, y int) bool {
	var v func(x, y int) bool
	switch op {
	case eq:
		v = func(x, y int) bool {
			return x == y
		}
	case gt:
		v = func(x, y int) bool {
			return x > y
		}
	case ge:
		v = func(x, y int) bool {
			return x >= y
		}
	case le:
		v = func(x, y int) bool {
			return x <= y
		}
	case lt:
		v = func(x, y int) bool {
			return x < y
		}
	case ne:
		v = func(x, y int) bool {
			return x != y
		}
	case isNull:
		v = nil
	default:
		panic("Operator Not found.")
	}
	return v
}

// NewColumnScanOperator creates a ColumnScanOperator
func NewIntColumnScanOperator(ctx Context, op ComparisonOperation, value int, fieldName string, size int) Uint32Operator {

	return &IntColumnScanOperator{
		ctx:   ctx,
		opFn:  selectIntOpFn(op),
		value: value,
		fn:    fieldName,
		sz:    size,
	}
}

type IntColumnScanOperator struct {
	ctx      Context
	opFn     func(x, y int) bool
	fn       string
	value    int
	sz       int
	iterator storage.IntIterator
	scanLeft []int
	lastRid  uint32
}

func (op *IntColumnScanOperator) Init() {
	c := op.ctx.Segment().Col(op.fn).(storage.IntColumn)
	op.iterator = c.Iterator()
}

func (op *IntColumnScanOperator) Destroy() {
}

func (op *IntColumnScanOperator) processVector(src []int, r []uint32) []uint32 {

	x := 0
	for ; x < len(src) && len(r) < op.sz; x++ {
		if op.opFn(src[x], op.value) {
			r = append(r, op.lastRid)
		}
		op.lastRid++
	}

	if len(r) == op.sz {
		op.scanLeft = src[x:]
	}

	return r
}

func (op *IntColumnScanOperator) Next() []uint32 {

	r := make([]uint32, 0, op.sz)

	if len(op.scanLeft) > 0 {

		r = op.processVector(op.scanLeft, r)
		if len(r) == op.sz {
			return r
		}

	}

	for op.iterator.HasNext() {

		intVector := op.iterator.Next()
		values := intVector.Values()

		r = op.processVector(values, r)
		if len(r) == op.sz {
			return r
		}
	}

	op.scanLeft = nil
	if len(r) > 0 {
		return r
	} else {
		return nil
	}

}

// NewColumnScanOperator creates a ColumnScanOperator
func NewIntNullColumnScanOperator(ctx Context, op ComparisonOperation, value int, fieldName string, size int) Uint32Operator {

	return &IntNullColumnScanOperator{
		ctx:   ctx,
		opFn:  selectIntOpFn(op),
		value: value,
		fn:    fieldName,
		sz:    size,
	}
}

type IntNullColumnScanOperator struct {
	ctx           Context
	opFn          func(x, y int) bool
	fn            string
	value         int
	sz            int
	iterator      storage.IntIterator
	lastRid       uint32
	resultLeft    []uint32
	lastCheckedId int
	lastValuePos  int
	lastIntVector *vector.IntVector
}

func (op *IntNullColumnScanOperator) Init() {
	c := op.ctx.Segment().Col(op.fn).(storage.IntColumn)
	op.iterator = c.Iterator()
}

func (op *IntNullColumnScanOperator) Destroy() {
}

func (op *IntNullColumnScanOperator) processVector(lastValuePos, lastCheckedId int, r []uint32) []uint32 {

	i := lastValuePos
	x := lastCheckedId

	for ; x < op.lastIntVector.Len(); x++ {
		if op.lastIntVector.IsValid(x) {
			if op.opFn(op.lastIntVector.Values()[x], op.value) {
				r = append(r, op.lastRid)
				i++
			}
		}
		op.lastRid++
	}

	if len(r) >= op.sz {
		op.lastCheckedId = x
		op.lastValuePos = i
		return r
	}

	return r
}

func (op *IntNullColumnScanOperator) Next() []uint32 {

	r := make([]uint32, 0, op.sz)

	if op.lastIntVector != nil {

		r = append(r, op.resultLeft...)
		r = op.processVector(op.lastValuePos, op.lastCheckedId, r)
		if len(r) >= op.sz {
			op.resultLeft = r[op.sz:]
			return r[:op.sz]
		}

	}

	for op.iterator.HasNext() {

		intVector := op.iterator.Next()
		op.lastIntVector = &intVector

		r = op.processVector(op.lastValuePos, op.lastCheckedId, r)

		if len(r) >= op.sz {
			op.resultLeft = r[op.sz:]
			return r[:op.sz]
		}

	}

	op.lastIntVector = nil
	op.lastCheckedId = 0
	op.lastValuePos = 0
	if len(r) > 0 {
		return r
	} else {
		return nil
	}

}
