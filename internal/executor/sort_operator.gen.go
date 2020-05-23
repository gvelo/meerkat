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
	"meerkat/internal/storage/vector"
)

func (op *SortOperator) createIntVector(v vector.IntVector) vector.IntVector {
	var rv vector.IntVector
	total := 0
	if v.HasNulls() {
		rv = vector.DefaultVectorPool().GetIntVector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendInt(v.Values()[op.order[i]])
			if v.IsValid(op.order[i]) {
				rv.SetValid(i)
			} else {
				rv.SetInvalid(i)
			}
			total++
		}
	} else {
		rv = vector.DefaultVectorPool().GetNotNullableIntVector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendInt(v.Values()[op.order[i]])
			total++
		}
	}
	return rv
}

func (op *SortOperator) createUintVector(v vector.UintVector) vector.UintVector {
	var rv vector.UintVector
	total := 0
	if v.HasNulls() {
		rv = vector.DefaultVectorPool().GetUintVector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendUint(v.Values()[op.order[i]])
			if v.IsValid(op.order[i]) {
				rv.SetValid(i)
			} else {
				rv.SetInvalid(i)
			}
			total++
		}
	} else {
		rv = vector.DefaultVectorPool().GetNotNullableUintVector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendUint(v.Values()[op.order[i]])
			total++
		}
	}
	return rv
}

func (op *SortOperator) createFloatVector(v vector.FloatVector) vector.FloatVector {
	var rv vector.FloatVector
	total := 0
	if v.HasNulls() {
		rv = vector.DefaultVectorPool().GetFloatVector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendFloat(v.Values()[op.order[i]])
			if v.IsValid(op.order[i]) {
				rv.SetValid(i)
			} else {
				rv.SetInvalid(i)
			}
			total++
		}
	} else {
		rv = vector.DefaultVectorPool().GetNotNullableFloatVector()
		for i := op.batchProc * op.ctx.Sz(); i < len(op.order); i++ {
			rv.AppendFloat(v.Values()[op.order[i]])
			total++
		}
	}
	return rv
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
