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

type Int64Vector struct {
	valid []uint64
	buf   []int64
	l     int
}

func (v *Int64Vector) Len() int {
	return v.l
}

func (v *Int64Vector) Cap() int {
	return len(v.buf)
}

func (v *Int64Vector) RemainingLen() int {
	return len(v.buf) - v.l
}

func (v *Int64Vector) Remaining() []int64 {
	return v.buf[v.l:]
}

func (v *Int64Vector) SetLen(l int) {
	v.l = l
}

func (v *Int64Vector) Values() []int64 {
	return v.buf[:v.l]
}

func (v *Int64Vector) Buf() []int64 {
	return v.buf
}

func (v *Int64Vector) IsValid(i int) bool {
	return v.valid[uint(i)>>log2WordSize]&(1<<(uint(i)&(wordSize-1))) != 0
}

func (v *Int64Vector) SetValid(i int) {
	v.valid[uint(i)>>log2WordSize] |= 1 << (uint(i) & (wordSize - 1))
}

func (v *Int64Vector) SetInvalid(i int) {
	v.valid[i>>log2WordSize] &^= 1 << (uint(i) & (wordSize - 1))
}

func NewInt64Vector(buf []int64, valid []uint64) Int64Vector {
	return Int64Vector{
		buf:   buf,
		valid: valid,
	}
}

type Int32Vector struct {
	valid []uint64
	buf   []int32
	l     int
}

func (v *Int32Vector) Len() int {
	return v.l
}

func (v *Int32Vector) Cap() int {
	return len(v.buf)
}

func (v *Int32Vector) RemainingLen() int {
	return len(v.buf) - v.l
}

func (v *Int32Vector) Remaining() []int32 {
	return v.buf[v.l:]
}

func (v *Int32Vector) SetLen(l int) {
	v.l = l
}

func (v *Int32Vector) Values() []int32 {
	return v.buf[:v.l]
}

func (v *Int32Vector) Buf() []int32 {
	return v.buf
}

func (v *Int32Vector) IsValid(i int) bool {
	return v.valid[uint(i)>>log2WordSize]&(1<<(uint(i)&(wordSize-1))) != 0
}

func (v *Int32Vector) SetValid(i int) {
	v.valid[uint(i)>>log2WordSize] |= 1 << (uint(i) & (wordSize - 1))
}

func (v *Int32Vector) SetInvalid(i int) {
	v.valid[i>>log2WordSize] &^= 1 << (uint(i) & (wordSize - 1))
}

func NewInt32Vector(buf []int32, valid []uint64) Int32Vector {
	return Int32Vector{
		buf:   buf,
		valid: valid,
	}
}

type Float64Vector struct {
	valid []uint64
	buf   []float64
	l     int
}

func (v *Float64Vector) Len() int {
	return v.l
}

func (v *Float64Vector) Cap() int {
	return len(v.buf)
}

func (v *Float64Vector) RemainingLen() int {
	return len(v.buf) - v.l
}

func (v *Float64Vector) Remaining() []float64 {
	return v.buf[v.l:]
}

func (v *Float64Vector) SetLen(l int) {
	v.l = l
}

func (v *Float64Vector) Values() []float64 {
	return v.buf[:v.l]
}

func (v *Float64Vector) Buf() []float64 {
	return v.buf
}

func (v *Float64Vector) IsValid(i int) bool {
	return v.valid[uint(i)>>log2WordSize]&(1<<(uint(i)&(wordSize-1))) != 0
}

func (v *Float64Vector) SetValid(i int) {
	v.valid[uint(i)>>log2WordSize] |= 1 << (uint(i) & (wordSize - 1))
}

func (v *Float64Vector) SetInvalid(i int) {
	v.valid[i>>log2WordSize] &^= 1 << (uint(i) & (wordSize - 1))
}

func NewFloat64Vector(buf []float64, valid []uint64) Float64Vector {
	return Float64Vector{
		buf:   buf,
		valid: valid,
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
