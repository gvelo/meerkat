// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: sort_operator.gen.go.tmpl

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
	"math"
	"meerkat/internal/storage/vector"
)

func (op *SortOperator) createIntVector(v vector.Int64Vector) vector.Int64Vector {
	var rv vector.Int64Vector
	total := 0
	if v.HasNulls() {
		rv = vector.DefaultVectorPool().GetInt64Vector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendInt64(v.Values()[op.order[i]])
			if v.IsValid(op.order[i]) {
				rv.SetValid(i)
			} else {
				rv.SetInvalid(i)
			}
			total++
		}
	} else {
		rv = vector.DefaultVectorPool().GetNotNullableInt64Vector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendInt64(v.Values()[op.order[i]])
			total++
		}
	}
	return rv
}

type IntVectorSorter struct {
	order []int
	v     *vector.Int64Vector
	asc   bool
	less  func(i, j int) bool
}

func (v *IntVectorSorter) Len() int { return len(v.order) }

func (v *IntVectorSorter) lessNull(i, j int) bool {
	vi := v.v.Values()[v.order[i]]
	vj := v.v.Values()[v.order[j]]

	// by default the nulls should be in the last positions.
	if !v.v.IsValid(v.order[i]) {
		vi = math.MaxInt64
	}

	if !v.v.IsValid(v.order[j]) {
		vj = math.MaxInt64
	}

	if v.asc {
		return vi < vj
	} else {
		return vi > vj
	}

}

func (v *IntVectorSorter) lessNotNull(i, j int) bool {
	if v.asc {
		return v.v.Values()[v.order[i]] < v.v.Values()[v.order[j]]
	} else {
		return v.v.Values()[v.order[i]] > v.v.Values()[v.order[j]]
	}
}

func (v *IntVectorSorter) Less(i, j int) bool {
	return v.less(i, j)
}

func (v *IntVectorSorter) Swap(i, j int) {
	v.order[i], v.order[j] = v.order[j], v.order[i]
}

func (op *SortOperator) createFloatVector(v vector.Float64Vector) vector.Float64Vector {
	var rv vector.Float64Vector
	total := 0
	if v.HasNulls() {
		rv = vector.DefaultVectorPool().GetFloat64Vector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendFloat64(v.Values()[op.order[i]])
			if v.IsValid(op.order[i]) {
				rv.SetValid(i)
			} else {
				rv.SetInvalid(i)
			}
			total++
		}
	} else {
		rv = vector.DefaultVectorPool().GetNotNullableFloat64Vector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendFloat64(v.Values()[op.order[i]])
			total++
		}
	}
	return rv
}

type FloatVectorSorter struct {
	order []int
	v     *vector.Float64Vector
	asc   bool
	less  func(i, j int) bool
}

func (v *FloatVectorSorter) Len() int { return len(v.order) }

func (v *FloatVectorSorter) lessNull(i, j int) bool {
	vi := v.v.Values()[v.order[i]]
	vj := v.v.Values()[v.order[j]]

	// by default the nulls should be in the last positions.
	if !v.v.IsValid(v.order[i]) {
		vi = math.MaxFloat64
	}

	if !v.v.IsValid(v.order[j]) {
		vj = math.MaxFloat64
	}

	if v.asc {
		return vi < vj
	} else {
		return vi > vj
	}

}

func (v *FloatVectorSorter) lessNotNull(i, j int) bool {
	if v.asc {
		return v.v.Values()[v.order[i]] < v.v.Values()[v.order[j]]
	} else {
		return v.v.Values()[v.order[i]] > v.v.Values()[v.order[j]]
	}
}

func (v *FloatVectorSorter) Less(i, j int) bool {
	return v.less(i, j)
}

func (v *FloatVectorSorter) Swap(i, j int) {
	v.order[i], v.order[j] = v.order[j], v.order[i]
}

func (op *SortOperator) createBoolVector(v vector.BoolVector) vector.BoolVector {
	var rv vector.BoolVector
	total := 0
	if v.HasNulls() {
		rv = vector.DefaultVectorPool().GetBoolVector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendBool(v.Values()[op.order[i]])
			if v.IsValid(op.order[i]) {
				rv.SetValid(i)
			} else {
				rv.SetInvalid(i)
			}
			total++
		}
	} else {
		rv = vector.DefaultVectorPool().GetNotNullableBoolVector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendBool(v.Values()[op.order[i]])
			total++
		}
	}
	return rv
}
