// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: col_val.gen.go.tmpl

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

package colval

type IntColValues struct {
	values []int
	rid    []uint32
}

func (v *IntColValues) Len() int {
	return len(v.values)
}

func (v *IntColValues) Rid() []uint32 {
	return v.rid
}

func (v *IntColValues) Values() []int {
	return v.values
}

func NewIntColValues(values []int, rid []uint32) IntColValues {
	return IntColValues{
		values: values,
		rid:    rid,
	}
}

type UintColValues struct {
	values []uint
	rid    []uint32
}

func (v *UintColValues) Len() int {
	return len(v.values)
}

func (v *UintColValues) Rid() []uint32 {
	return v.rid
}

func (v *UintColValues) Values() []uint {
	return v.values
}

func NewUintColValues(values []uint, rid []uint32) UintColValues {
	return UintColValues{
		values: values,
		rid:    rid,
	}
}

type FloatColValues struct {
	values []float64
	rid    []uint32
}

func (v *FloatColValues) Len() int {
	return len(v.values)
}

func (v *FloatColValues) Rid() []uint32 {
	return v.rid
}

func (v *FloatColValues) Values() []float64 {
	return v.values
}

func NewFloatColValues(values []float64, rid []uint32) FloatColValues {
	return FloatColValues{
		values: values,
		rid:    rid,
	}
}

type BoolColValues struct {
	values []bool
	rid    []uint32
}

func (v *BoolColValues) Len() int {
	return len(v.values)
}

func (v *BoolColValues) Rid() []uint32 {
	return v.rid
}

func (v *BoolColValues) Values() []bool {
	return v.values
}

func NewBoolColValues(values []bool, rid []uint32) BoolColValues {
	return BoolColValues{
		values: values,
		rid:    rid,
	}
}

type IntColSource interface {
	ColSource
	Next() IntColValues
}

type UintColSource interface {
	ColSource
	Next() UintColValues
}

type FloatColSource interface {
	ColSource
	Next() FloatColValues
}

type BoolColSource interface {
	ColSource
	Next() BoolColValues
}