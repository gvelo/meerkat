// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: vector.gen.go.tmpl

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

package vector

type IntVector struct {
	valid []uint64
	buf   []int
	l     int
}

func (v *IntVector) Len() int {
	return v.l
}

func (v *IntVector) Cap() int {
	return len(v.buf)
}

func (v *IntVector) RemainingLen() int {
	return len(v.buf) - v.l
}

func (v *IntVector) Remaining() []int {
	return v.buf[v.l:]
}

func (v *IntVector) SetLen(l int) {
	v.l = l
}

func (v *IntVector) Values() []int {
	return v.buf[:v.l]
}

func (v *IntVector) Buf() []int {
	return v.buf
}

func (v *IntVector) IsValid(i int) bool {
	return v.valid[uint(i)>>log2WordSize]&(1<<(uint(i)&(wordSize-1))) != 0
}

func (v *IntVector) SetValid(i int) {
	v.valid[uint(i)>>log2WordSize] |= 1 << (uint(i) & (wordSize - 1))
}

func (v *IntVector) SetInvalid(i int) {
	v.valid[i>>log2WordSize] &^= 1 << (uint(i) & (wordSize - 1))
}

func NewIntVector(buf []int, valid []uint64) IntVector {
	return IntVector{
		buf:   buf,
		valid: valid,
	}
}

type UintVector struct {
	valid []uint64
	buf   []uint
	l     int
}

func (v *UintVector) Len() int {
	return v.l
}

func (v *UintVector) Cap() int {
	return len(v.buf)
}

func (v *UintVector) RemainingLen() int {
	return len(v.buf) - v.l
}

func (v *UintVector) Remaining() []uint {
	return v.buf[v.l:]
}

func (v *UintVector) SetLen(l int) {
	v.l = l
}

func (v *UintVector) Values() []uint {
	return v.buf[:v.l]
}

func (v *UintVector) Buf() []uint {
	return v.buf
}

func (v *UintVector) IsValid(i int) bool {
	return v.valid[uint(i)>>log2WordSize]&(1<<(uint(i)&(wordSize-1))) != 0
}

func (v *UintVector) SetValid(i int) {
	v.valid[uint(i)>>log2WordSize] |= 1 << (uint(i) & (wordSize - 1))
}

func (v *UintVector) SetInvalid(i int) {
	v.valid[i>>log2WordSize] &^= 1 << (uint(i) & (wordSize - 1))
}

func NewUintVector(buf []uint, valid []uint64) UintVector {
	return UintVector{
		buf:   buf,
		valid: valid,
	}
}

type FloatVector struct {
	valid []uint64
	buf   []float64
	l     int
}

func (v *FloatVector) Len() int {
	return v.l
}

func (v *FloatVector) Cap() int {
	return len(v.buf)
}

func (v *FloatVector) RemainingLen() int {
	return len(v.buf) - v.l
}

func (v *FloatVector) Remaining() []float64 {
	return v.buf[v.l:]
}

func (v *FloatVector) SetLen(l int) {
	v.l = l
}

func (v *FloatVector) Values() []float64 {
	return v.buf[:v.l]
}

func (v *FloatVector) Buf() []float64 {
	return v.buf
}

func (v *FloatVector) IsValid(i int) bool {
	return v.valid[uint(i)>>log2WordSize]&(1<<(uint(i)&(wordSize-1))) != 0
}

func (v *FloatVector) SetValid(i int) {
	v.valid[uint(i)>>log2WordSize] |= 1 << (uint(i) & (wordSize - 1))
}

func (v *FloatVector) SetInvalid(i int) {
	v.valid[i>>log2WordSize] &^= 1 << (uint(i) & (wordSize - 1))
}

func NewFloatVector(buf []float64, valid []uint64) FloatVector {
	return FloatVector{
		buf:   buf,
		valid: valid,
		l:     len(buf),
	}
}

type BoolVector struct {
	valid []uint64
	buf   []bool
	l     int
}

func (v *BoolVector) Len() int {
	return v.l
}

func (v *BoolVector) Cap() int {
	return len(v.buf)
}

func (v *BoolVector) RemainingLen() int {
	return len(v.buf) - v.l
}

func (v *BoolVector) Remaining() []bool {
	return v.buf[v.l:]
}

func (v *BoolVector) SetLen(l int) {
	v.l = l
}

func (v *BoolVector) Values() []bool {
	return v.buf[:v.l]
}

func (v *BoolVector) Buf() []bool {
	return v.buf
}

func (v *BoolVector) IsValid(i int) bool {
	return v.valid[uint(i)>>log2WordSize]&(1<<(uint(i)&(wordSize-1))) != 0
}

func (v *BoolVector) SetValid(i int) {
	v.valid[uint(i)>>log2WordSize] |= 1 << (uint(i) & (wordSize - 1))
}

func (v *BoolVector) SetInvalid(i int) {
	v.valid[i>>log2WordSize] &^= 1 << (uint(i) & (wordSize - 1))
}

func NewBoolVector(buf []bool, valid []uint64) BoolVector {
	return BoolVector{
		buf:   buf,
		valid: valid,
	}
}
