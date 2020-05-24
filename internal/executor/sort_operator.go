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

//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@../storage/scalar_types.tmpldata sort_operator.gen.go.tmpl

package executor

import (
	"bytes"
	"fmt"
	"github.com/psilva261/timsort"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
	"meerkat/internal/storage/vector"
	"sort"
)

type SortOpt struct {
	keyName  string
	orderAsc bool
}

func NewSortOperator(ctx Context, child MultiVectorOperator, opts []SortOpt) *SortOperator {
	return &SortOperator{
		ctx:      ctx,
		child:    child,
		sortOpts: opts,
		log:      log.With().Str("src", "SortOperator").Logger(),
	}
}

// SortOperator
type SortOperator struct {
	ctx            Context
	child          MultiVectorOperator // (Positions to review)
	sortOpts       []SortOpt
	sorted         [][]byte
	batchProc      int
	order          []int
	slicesAppended []interface{}
	log            zerolog.Logger
}

// TODO(sebad): This operator should spill over disk ... some day.
func (op *SortOperator) Init() {
	op.child.Init()

	colIds := make([]int, 0, len(op.sortOpts))
	for _, it := range op.sortOpts {
		_, id, err := op.ctx.GetFieldProcessed().FindField(it.keyName)
		if err != nil {
			log.Error().Err(err)
		}

		colIds = append(colIds, id)
	}

	n := op.child.Next()
	op.slicesAppended = n

	op.order = make([]int, 0)
	i := 0

	l1 := getLen(n[0])
	for ; i < l1; i++ {
		op.order = append(op.order, i)
	}

	// second batch.
	n = op.child.Next()
	for ; n != nil; n = op.child.Next() {
		l1 := getLen(n[0])
		for ; i < l1; i++ {
			op.order = append(op.order, i)
		}
		op.appendSlices(n)
	}

	// this is the sorting rational
	// we have these cols for sort:
	//  c1 c2
	//  1  7
	//  2  6
	//  1  3
	//  2  4
	//
	//  then we should sort the 1 first
	//
	//  1  7
	//  1  3
	//  2  6
	//  2  4
	//
	//  then we should sort the following columns in the
	//  case we have partitions with the previous keys like this
	//
	//  1 3
	//  1 7
	//  2 4
	//  2 6

	so := getVectorSorter(op.slicesAppended[colIds[0]], op.order, op.sortOpts[0].orderAsc)
	// Sort the 1st vector.
	timsort.TimSort(so)

	if len(colIds) > 1 {

		b := make([]bool, len(op.order))
		sv := make([]int, 16)

		for i := 1; i < len(colIds); i++ {
			// Create the dif array from prev col.
			buildDiffArray(op.slicesAppended[colIds[i-1]], op.order, b)
			sv = createPartitions(b, sv[:0])
			sortPartitions(op.slicesAppended[colIds[i]], op.order, op.sortOpts[i].orderAsc, sv)
		}
	}

}

func sortPartitions(v interface{}, order []int, asc bool, p []int) {

	if len(p) < 1 {
		log.Error().Err(fmt.Errorf("invalid partitions list %v", p))
	}

	for i, ps := range p {
		var pe int
		if i == len(p)-1 {
			pe = len(order)
		} else {
			pe = p[i+1]
		}
		so := getVectorSorter(v, order[ps:pe], asc)
		timsort.TimSort(so)
	}

}

//
func createPartitions(b []bool, sel []int) []int {
	l := len(b)
	for i := 0; i < l; i++ {
		if b[i] {
			sel = append(sel, i)
		}
	}
	return sel
}

func buildDiffArray(so interface{}, o []int, b []bool) {
	switch t := so.(type) {
	case vector.IntVector:
		buildDiffIntVector(t, o, b)
	case vector.ByteSliceVector:
		buildDiffByteSliceVector(t, o, b)
	case vector.FloatVector:
		buildDiffFloatVector(t, o, b)
	case vector.BoolVector:
		panic(log.Error().Err(fmt.Errorf("wtf")))
	default:
		log.Error().Msg("No found.")
	}
}

func buildDiffIntVector(colVec vector.IntVector, order []int, b []bool) {
	var lastVal int
	var lastValNull bool
	b[0] = true
	if colVec.HasNulls() {
		for outputIdx, checkIdx := range order {
			null := !colVec.IsValid(checkIdx)
			if null {
				if !lastValNull {
					// The current value is null while the previous was not.
					b[outputIdx] = true
				}
			} else {
				v := colVec.Get(checkIdx)
				if lastValNull {
					// The previous value was null while the current is not.
					b[outputIdx] = true
				} else {
					// Neither value is null, so we must compare.
					var unique bool

					{
						var cmpResult int

						{
							a, b := int64(v), int64(lastVal)
							if a < b {
								cmpResult = -1
							} else if a > b {
								cmpResult = 1
							} else {
								cmpResult = 0
							}
						}

						unique = cmpResult != 0
					}

					b[outputIdx] = b[outputIdx] || unique
				}
				lastVal = v
			}
			lastValNull = null
		}
	} else {
		for outputIdx, checkIdx := range order {
			v := colVec.Get(checkIdx)
			var unique bool

			{
				var cmpResult int

				{
					a, b := int64(v), int64(lastVal)
					if a < b {
						cmpResult = -1
					} else if a > b {
						cmpResult = 1
					} else {
						cmpResult = 0
					}
				}

				unique = cmpResult != 0
			}

			b[outputIdx] = b[outputIdx] || unique
			lastVal = v
		}
	}
}

func buildDiffFloatVector(colVec vector.FloatVector, order []int, b []bool) {
	var lastVal float64
	var lastValNull bool
	b[0] = true
	if colVec.HasNulls() {
		for outputIdx, checkIdx := range order {
			null := !colVec.IsValid(checkIdx)
			if null {
				if !lastValNull {
					// The current value is null while the previous was not.
					b[outputIdx] = true
				}
			} else {
				v := colVec.Get(checkIdx)
				if lastValNull {
					// The previous value was null while the current is not.
					b[outputIdx] = true
				} else {
					// Neither value is null, so we must compare.
					var unique bool

					{
						var cmpResult int

						{
							a, b := v, lastVal
							if a < b {
								cmpResult = -1
							} else if a > b {
								cmpResult = 1
							} else if a == b {
								cmpResult = 0
							} else if math.IsNaN(a) {
								if math.IsNaN(b) {
									cmpResult = 0
								} else {
									cmpResult = -1
								}
							} else {
								cmpResult = 1
							}
						}

						unique = cmpResult != 0
					}

					b[outputIdx] = b[outputIdx] || unique
				}
				lastVal = v
			}
			lastValNull = null
		}
	} else {
		for outputIdx, checkIdx := range order {
			v := colVec.Get(checkIdx)
			var unique bool

			{
				var cmpResult int

				{
					a, b := v, lastVal
					if a < b {
						cmpResult = -1
					} else if a > b {
						cmpResult = 1
					} else if a == b {
						cmpResult = 0
					} else if math.IsNaN(a) {
						if math.IsNaN(b) {
							cmpResult = 0
						} else {
							cmpResult = -1
						}
					} else {
						cmpResult = 1
					}
				}

				unique = cmpResult != 0
			}

			b[outputIdx] = b[outputIdx] || unique
			lastVal = v
		}
	}
}

func buildDiffByteSliceVector(colVec vector.ByteSliceVector, order []int, b []bool) {
	var lastVal []byte
	var lastValNull bool
	b[0] = true
	if colVec.HasNulls() {
		for outputIdx, checkIdx := range order {
			null := !colVec.IsValid(checkIdx)
			if null {
				if !lastValNull {
					// The current value is null while the previous was not.
					b[outputIdx] = true
				}
			} else {
				v := colVec.Get(checkIdx)
				if lastValNull {
					// The previous value was null while the current is not.
					b[outputIdx] = true
				} else {
					// Neither value is null, so we must compare.
					var unique bool

					{
						var cmpResult int

						{
							cmpResult = bytes.Compare(v, lastVal)
						}

						unique = cmpResult != 0
					}

					b[outputIdx] = b[outputIdx] || unique
				}
				lastVal = v
			}
			lastValNull = null
		}
	} else {
		for outputIdx, checkIdx := range order {
			v := colVec.Get(checkIdx)
			var unique bool

			{
				var cmpResult int

				{
					cmpResult = bytes.Compare(v, lastVal)
				}

				unique = cmpResult != 0
			}

			b[outputIdx] = b[outputIdx] || unique
			lastVal = v
		}
	}
}

func getVectorSorter(v interface{}, order []int, asc bool) sort.Interface {
	var r sort.Interface
	switch t := v.(type) {
	case vector.IntVector:
		r := &IntVectorSorter{
			order: order,
			v:     &t,
			asc:   asc,
		}

		if t.HasNulls() {
			r.less = r.lessNull
		} else {
			r.less = r.lessNotNull
		}

		return r
	case vector.FloatVector:
		r := &FloatVectorSorter{
			order: order,
			v:     &t,
			asc:   asc,
		}

		if t.HasNulls() {
			r.less = r.lessNull
		} else {
			r.less = r.lessNotNull
		}

	default:
		log.Error().Msg("No found.")
	}
	return r
}

// Append all slices in vectors.
func (op *SortOperator) appendSlices(src []interface{}) {
	for i, it := range src {

		switch s := it.(type) {
		case vector.IntVector:
			v := op.slicesAppended[i].(vector.IntVector)
			v.Append(s.Values())
		case vector.ByteSliceVector:
			v := op.slicesAppended[i].(vector.ByteSliceVector)
			v.AppendSlice(s.Data())
		case vector.FloatVector:
			v := op.slicesAppended[i].(vector.FloatVector)
			v.Append(s.Values())
		case vector.BoolVector:
			v := op.slicesAppended[i].(vector.BoolVector)
			v.Append(s.Values())
		default:
			log.Error().Err(fmt.Errorf("not mapped type %v", s))
		}

	}
}

func (op *SortOperator) Destroy() {
	op.child.Destroy()
}

func (op *SortOperator) Next() []interface{} {

	if op.batchProc*op.ctx.Sz() >= len(op.order) {
		return nil
	}

	res := make([]interface{}, len(op.slicesAppended))

	for i, it := range op.slicesAppended {
		switch it.(type) {
		case vector.IntVector:
			res[i] = op.createIntVector(it.(vector.IntVector))
		case vector.ByteSliceVector:
			res[i] = op.createByteSliceVector(it.(vector.ByteSliceVector))
		case vector.BoolVector:
			res[i] = op.createBoolVector(it.(vector.BoolVector))
		case vector.FloatVector:
			res[i] = op.createFloatVector(it.(vector.FloatVector))
		default:
			log.Error().Err(fmt.Errorf("type not mapped %v", it))
		}
	}
	op.batchProc++
	return res
}

func (op *SortOperator) createByteSliceVector(v vector.ByteSliceVector) vector.ByteSliceVector {

	var rv vector.ByteSliceVector
	total := 0
	if v.HasNulls() {
		rv = vector.DefaultVectorPool().GetByteSliceVector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendSlice(v.Get(op.order[i]))
			if v.IsValid(op.order[i]) {
				rv.SetValid(i)
			} else {
				rv.SetInvalid(i)
			}
			total++
		}
	} else {
		rv = vector.DefaultVectorPool().GetNotNullableByteSliceVector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendSlice(v.Get(op.order[i]))
			total++
		}
	}
	rv.SetLen(total)
	return rv
}
