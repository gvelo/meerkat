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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
	"meerkat/internal/util/sliceutil"
	"strings"
)

func selectStringOpFn(op ComparisonOperation) func(x []byte, y string) bool {
	var v func(x []byte, y string) bool
	switch op {
	case Eq:
		v = func(x []byte, y string) bool {
			return sliceutil.BS2S(x) == y
		}
	case Gt:
		v = func(x []byte, y string) bool {
			return sliceutil.BS2S(x) > y
		}
	case Ge:
		v = func(x []byte, y string) bool {
			return sliceutil.BS2S(x) >= y
		}
	case Le:
		v = func(x []byte, y string) bool {
			return sliceutil.BS2S(x) <= y
		}
	case Lt:
		v = func(x []byte, y string) bool {
			return sliceutil.BS2S(x) < y
		}
	case Ne:
		v = func(x []byte, y string) bool {
			return sliceutil.BS2S(x) != y
		}
	case Contains:
		v = func(x []byte, y string) bool {
			return strings.Contains(sliceutil.BS2S(x), y)
		}
	case IsNull:
		v = nil
	default:
		panic("Operator Not found.")
	}
	return v
}

// NewColumnScanOperator creates a ColumnScanOperator
func NewStringColumnScanOperator(ctx Context, op ComparisonOperation, value string, fieldName string, size int, nullable bool) Uint32Operator {
	v := &StringColumnScanOperator{
		ctx:   ctx,
		opFn:  selectStringOpFn(op),
		value: value,
		fn:    fieldName,
		sz:    size,
		log:   log.With().Str("src", "IntColumnScanOperator").Logger(),
	}
	if nullable {
		v.processFn = v.processNullVector
	} else {
		v.processFn = v.processVector
	}
	return v
}

type StringColumnScanOperator struct {
	ctx           Context
	opFn          func(x []byte, y string) bool
	fn            string
	value         string
	sz            int
	iterator      storage.BinaryIterator
	lastRid       uint32
	resultLeft    []uint32
	lastCheckedId int
	lastValuePos  int
	lastVector    *vector.ByteSliceVector
	processFn     func(lastValuePos, lastCheckedId int, r []uint32) []uint32
	log           zerolog.Logger
}

func (op *StringColumnScanOperator) Next() []uint32 {
	r := make([]uint32, 0, op.sz)

	if op.lastVector != nil {

		r = append(r, op.resultLeft...)
		r = op.processVector(op.lastValuePos, op.lastCheckedId, r)
		if len(r) >= op.sz {
			op.resultLeft = r[op.sz:]
			return r[:op.sz]
		}

	}

	for op.iterator.HasNext() {

		v := op.iterator.Next()
		op.lastVector = &v

		r = op.processVector(op.lastValuePos, op.lastCheckedId, r)

		if len(r) >= op.sz {
			op.resultLeft = r[op.sz:]
			return r[:op.sz]
		}

	}

	op.lastVector = nil
	op.lastCheckedId = 0
	op.lastValuePos = 0
	if len(r) > 0 {
		return r
	} else {
		return nil
	}
}

func (op *StringColumnScanOperator) Init() {
	c := op.ctx.Segment().Col(op.fn).(storage.StringColumn)
	op.iterator = c.Iterator()
}

func (op *StringColumnScanOperator) Destroy() {
}

func (op *StringColumnScanOperator) processVector(lastValuePos, lastCheckedId int, r []uint32) []uint32 {

	i := lastValuePos
	x := lastCheckedId

	for ; x < op.lastVector.Len(); x++ {
		if op.opFn(op.lastVector.Get(x), op.value) {
			r = append(r, op.lastRid)
			i++
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

func (op *StringColumnScanOperator) processNullVector(lastValuePos, lastCheckedId int, r []uint32) []uint32 {

	i := lastValuePos
	x := lastCheckedId

	for ; x < op.lastVector.Len(); x++ {
		if op.lastVector.IsValid(x) {
			if op.opFn(op.lastVector.Get(x), op.value) {
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

func selectIntOpFn(op ComparisonOperation) func(x, y int) bool {
	var v func(x, y int) bool
	switch op {
	case Eq:
		v = func(x, y int) bool {
			return x == y
		}
	case Gt:
		v = func(x, y int) bool {
			return x > y
		}
	case Ge:
		v = func(x, y int) bool {
			return x >= y
		}
	case Le:
		v = func(x, y int) bool {
			return x <= y
		}
	case Lt:
		v = func(x, y int) bool {
			return x < y
		}
	case Ne:
		v = func(x, y int) bool {
			return x != y
		}
	case IsNull:
		v = nil
	default:
		panic("Operator Not found.")
	}
	return v
}

// NewColumnScanOperator creates a ColumnScanOperator
func NewIntColumnScanOperator(ctx Context, op ComparisonOperation, value int, fieldName string, size int, nullable bool) Uint32Operator {

	v := &IntColumnScanOperator{
		ctx:   ctx,
		opFn:  selectIntOpFn(op),
		value: value,
		fn:    fieldName,
		sz:    size,
		log:   log.With().Str("src", "IntColumnScanOperator").Logger(),
	}
	if nullable {
		v.processFn = v.processNullVector
	} else {
		v.processFn = v.processVector
	}
	return v
}

type IntColumnScanOperator struct {
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
	lastVector    *vector.IntVector
	processFn     func(lastValuePos, lastCheckedId int, r []uint32) []uint32
	log           zerolog.Logger
}

func (op *IntColumnScanOperator) Init() {
	c := op.ctx.Segment().Col(op.fn).(storage.IntColumn)
	op.iterator = c.Iterator()
}

func (op *IntColumnScanOperator) Destroy() {
}

func (op *IntColumnScanOperator) processVector(lastValuePos, lastCheckedId int, r []uint32) []uint32 {

	i := lastValuePos
	x := lastCheckedId

	for ; x < op.lastVector.Len(); x++ {
		if op.opFn(op.lastVector.Values()[x], op.value) {
			r = append(r, op.lastRid)
			i++
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

func (op *IntColumnScanOperator) processNullVector(lastValuePos, lastCheckedId int, r []uint32) []uint32 {

	i := lastValuePos
	x := lastCheckedId

	for ; x < op.lastVector.Len(); x++ {
		if op.lastVector.IsValid(x) {
			if op.opFn(op.lastVector.Values()[x], op.value) {
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

func (op *IntColumnScanOperator) Next() []uint32 {

	r := make([]uint32, 0, op.sz)

	if op.lastVector != nil {

		r = append(r, op.resultLeft...)
		r = op.processVector(op.lastValuePos, op.lastCheckedId, r)
		if len(r) >= op.sz {
			op.resultLeft = r[op.sz:]
			return r[:op.sz]
		}

	}

	for op.iterator.HasNext() {

		v := op.iterator.Next()
		op.lastVector = &v

		r = op.processVector(op.lastValuePos, op.lastCheckedId, r)

		if len(r) >= op.sz {
			op.resultLeft = r[op.sz:]
			return r[:op.sz]
		}

	}

	op.lastVector = nil
	op.lastCheckedId = 0
	op.lastValuePos = 0
	if len(r) > 0 {
		return r
	} else {
		return nil
	}

}

func selectTimeOpFn(op ComparisonOperation) func(x, y, z int) bool {
	var v func(x, y, z int) bool
	switch op {
	case Eq:
		v = func(x, y, z int) bool {
			return x == y
		}
	case Gt:
		v = func(x, y, z int) bool {
			return x > y
		}
	case Between:
		v = func(x, y, z int) bool {
			return x > y
		}
	case Ge:
		v = func(x, y, z int) bool {
			return x >= y
		}
	case Le:
		v = func(x, y, z int) bool {
			return x <= y
		}
	case Lt:
		v = func(x, y, z int) bool {
			return x < y
		}
	case Ne:
		v = func(x, y, z int) bool {
			return x != y
		}
	case IsNull:
		v = nil
	default:
		panic("Operator Not found.")
	}
	return v
}

// NewTimeColumnScanOperator creates a TimeColumnScanOperator
func NewTimeColumnScanOperator(ctx Context, op ComparisonOperation, valueFrom, valueTo int, fieldName string, size int, nullable bool) Uint32Operator {

	v := &TimeColumnScanOperator{
		ctx:     ctx,
		opFn:    selectTimeOpFn(op),
		value:   valueFrom,
		valueTo: valueTo,
		fn:      fieldName,
		sz:      size,
		log:     log.With().Str("src", "TimeColumnScanOperator").Logger(),
	}

	return v
}

type TimeColumnScanOperator struct {
	ctx           Context
	opFn          func(x, y, z int) bool
	fn            string
	value         int
	valueTo       int
	sz            int
	iterator      storage.IntIterator
	lastRid       uint32
	resultLeft    []uint32
	lastCheckedId int
	lastValuePos  int
	lastVector    *vector.IntVector
	processFn     func(lastValuePos, lastCheckedId int, r []uint32) []uint32
	log           zerolog.Logger
}

func (op *TimeColumnScanOperator) Init() {
	c := op.ctx.Segment().Col(op.fn).(storage.IntColumn)
	op.iterator = c.Iterator()
}

func (op *TimeColumnScanOperator) Destroy() {
}

func (op *TimeColumnScanOperator) processVector(lastValuePos, lastCheckedId int, r []uint32) []uint32 {

	i := lastValuePos
	x := lastCheckedId

	for ; x < op.lastVector.Len(); x++ {
		if op.opFn(op.lastVector.Values()[x], op.value, op.valueTo) {
			r = append(r, op.lastRid)
			i++
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

func (op *TimeColumnScanOperator) Next() []uint32 {

	r := make([]uint32, 0, op.sz)

	if op.lastVector != nil {

		r = append(r, op.resultLeft...)
		r = op.processVector(op.lastValuePos, op.lastCheckedId, r)
		if len(r) >= op.sz {
			op.resultLeft = r[op.sz:]
			return r[:op.sz]
		}

	}

	for op.iterator.HasNext() {

		v := op.iterator.Next()
		op.lastVector = &v

		r = op.processVector(op.lastValuePos, op.lastCheckedId, r)

		if len(r) >= op.sz {
			op.resultLeft = r[op.sz:]
			return r[:op.sz]
		}

	}

	op.lastVector = nil
	op.lastCheckedId = 0
	op.lastValuePos = 0
	if len(r) > 0 {
		return r
	} else {
		return nil
	}

}