// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: bufcolval.gen.go.tmpl

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

package bufcolval

import (
	"meerkat/internal/buffer"
	"meerkat/internal/storage/colval"
)

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
	colval.ColSource
	Next() IntColValues
}

type UintColSource interface {
	colval.ColSource
	Next() UintColValues
}

type FloatColSource interface {
	colval.ColSource
	Next() FloatColValues
}

type BoolColSource interface {
	colval.ColSource
	Next() BoolColValues
}

type IntBufColSource struct {
	srcBuf   []int
	dstBuf   []int
	nulls    []bool
	rid      []uint32
	permMap  []int
	pos      int
	hasNulls bool
}

func (cs *IntBufColSource) HasNext() bool {
	return cs.pos < len(cs.srcBuf)
}

func (cs *IntBufColSource) HasNulls() bool {
	return cs.hasNulls
}

// The underlying array point to an internal buffer that will be
// overwritten by a subsequent call to Next().
func (cs *IntBufColSource) Next() colval.IntColValues {

	var i int

	for i = 0; i < len(cs.dstBuf) && cs.pos < len(cs.srcBuf); i++ {

		j := cs.permMap[cs.pos]

		if cs.hasNulls && cs.nulls[j] {
			i--
		} else {
			cs.dstBuf[i] = cs.srcBuf[j]
			cs.rid[i] = uint32(cs.pos)
		}

		cs.pos++

	}

	return colval.NewIntColValues(cs.dstBuf[:i], cs.rid[:i])

}

func NewIntBufColSource(buf *buffer.IntBuffer, dstLen int, permMap []int) *IntBufColSource {

	return &IntBufColSource{
		srcBuf:   buf.Values(),
		dstBuf:   make([]int, dstLen),
		nulls:    buf.Nulls(),
		rid:      make([]uint32, dstLen),
		permMap:  permMap,
		hasNulls: buf.Nullable(),
	}

}

type UintBufColSource struct {
	srcBuf   []uint
	dstBuf   []uint
	nulls    []bool
	rid      []uint32
	permMap  []int
	pos      int
	hasNulls bool
}

func (cs *UintBufColSource) HasNext() bool {
	return cs.pos < len(cs.srcBuf)
}

func (cs *UintBufColSource) HasNulls() bool {
	return cs.hasNulls
}

// The underlying array point to an internal buffer that will be
// overwritten by a subsequent call to Next().
func (cs *UintBufColSource) Next() colval.UintColValues {

	var i int

	for i = 0; i < len(cs.dstBuf) && cs.pos < len(cs.srcBuf); i++ {

		j := cs.permMap[cs.pos]

		if cs.hasNulls && cs.nulls[j] {
			i--
		} else {
			cs.dstBuf[i] = cs.srcBuf[j]
			cs.rid[i] = uint32(cs.pos)
		}

		cs.pos++

	}

	return colval.NewUintColValues(cs.dstBuf[:i], cs.rid[:i])

}

func NewUintBufColSource(buf *buffer.UintBuffer, dstSize int, permMap []int) *UintBufColSource {

	return &UintBufColSource{
		srcBuf:   buf.Values(),
		dstBuf:   make([]uint, dstSize),
		nulls:    buf.Nulls(),
		rid:      make([]uint32, dstSize),
		permMap:  permMap,
		hasNulls: buf.Nullable(),
	}

}

type FloatBufColSource struct {
	srcBuf   []float64
	dstBuf   []float64
	nulls    []bool
	rid      []uint32
	permMap  []int
	pos      int
	hasNulls bool
}

func (cs *FloatBufColSource) HasNext() bool {
	return cs.pos < len(cs.srcBuf)
}

func (cs *FloatBufColSource) HasNulls() bool {
	return cs.hasNulls
}

// The underlying array point to an internal buffer that will be
// overwritten by a subsequent call to Next().
func (cs *FloatBufColSource) Next() colval.FloatColValues {

	var i int

	for i = 0; i < len(cs.dstBuf) && cs.pos < len(cs.srcBuf); i++ {

		j := cs.permMap[cs.pos]

		if cs.hasNulls && cs.nulls[j] {
			i--
		} else {
			cs.dstBuf[i] = cs.srcBuf[j]
			cs.rid[i] = uint32(cs.pos)
		}

		cs.pos++

	}

	return colval.NewFloatColValues(cs.dstBuf[:i], cs.rid[:i])

}

func NewFloatBufColSource(buf *buffer.FloatBuffer, dstSize int, permMap []int) *FloatBufColSource {

	return &FloatBufColSource{
		srcBuf:   buf.Values(),
		dstBuf:   make([]float64, dstSize),
		nulls:    buf.Nulls(),
		rid:      make([]uint32, dstSize),
		permMap:  permMap,
		hasNulls: buf.Nullable(),
	}

}

type BoolBufColSource struct {
	srcBuf   []bool
	dstBuf   []bool
	nulls    []bool
	rid      []uint32
	permMap  []int
	pos      int
	hasNulls bool
}

func (cs *BoolBufColSource) HasNext() bool {
	return cs.pos < len(cs.srcBuf)
}

func (cs *BoolBufColSource) HasNulls() bool {
	return cs.hasNulls
}

// The underlying array point to an internal buffer that will be
// overwritten by a subsequent call to Next().
func (cs *BoolBufColSource) Next() colval.BoolColValues {

	var i int

	for i = 0; i < len(cs.dstBuf) && cs.pos < len(cs.srcBuf); i++ {

		j := cs.permMap[cs.pos]

		if cs.hasNulls && cs.nulls[j] {
			i--
		} else {
			cs.dstBuf[i] = cs.srcBuf[j]
			cs.rid[i] = uint32(cs.pos)
		}

		cs.pos++

	}

	return colval.NewBoolColValues(cs.dstBuf[:i], cs.rid[:i])

}

func NewBoolBufColSource(buf *buffer.BoolBuffer, dstSize int, permMap []int) *BoolBufColSource {

	return &BoolBufColSource{
		srcBuf:   buf.Values(),
		dstBuf:   make([]bool, dstSize),
		nulls:    buf.Nulls(),
		rid:      make([]uint32, dstSize),
		permMap:  permMap,
		hasNulls: buf.Nullable(),
	}

}
