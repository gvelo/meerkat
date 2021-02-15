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

//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@../scalar_types.tmpldata vector.gen.go.tmpl

package vector

import "meerkat/internal/util/sliceutil"

const (
	log2WordSize = uint(6)
	wordSize     = uint(64)
)

type Vector interface {
	Len() int
	Cap() int
	HasNulls() bool
	AsBytes() []byte
	ValidityAsBytes() []byte
}

type Pool interface {
	GetInt64Vector() Int64Vector
	GetInt32Vector() Int32Vector
	GetFloat64Vector() Float64Vector
	GetByteSliceVector() ByteSliceVector
	GetBoolVector() BoolVector

	GetNotNullableBoolVector() BoolVector
	GetNotNullableFloat64Vector() Float64Vector
	GetNotNullableByteSliceVector() ByteSliceVector
	GetNotNullableInt64Vector() Int64Vector
	GetNotNullableInt32Vector() Int32Vector

	PutInt64Vector(vector Int64Vector)
}

type ByteSliceVector struct {
	valid   []uint64
	buf     []byte
	offsets []int
	l       int
	c       int
}

func (v *ByteSliceVector) HasNulls() bool {
	return v.valid != nil
}

func (v *ByteSliceVector) Len() int {
	return v.l
}

func (v *ByteSliceVector) Cap() int {
	return v.c
}

func (v *ByteSliceVector) SetLen(l int) {
	v.l = l
}

func (v *ByteSliceVector) Buffer() []byte {
	return v.buf[:v.l]
}

func (v *ByteSliceVector) Offsets() []int {
	return v.offsets[:v.l]
}

func (v *ByteSliceVector) AsBytes() []byte {
	return v.buf[:v.l]
}

func (v *ByteSliceVector) OffsetsAsBytes() []byte {
	return sliceutil.I2B(v.offsets[:v.l])
}

func (v *ByteSliceVector) ValidityAsBytes() []byte {
	return sliceutil.U642B(v.valid[:v.l/8])
}

func (v *ByteSliceVector) IsValid(i int) bool {
	return v.valid[uint(i)>>log2WordSize]&(1<<(uint(i)&(wordSize-1))) != 0
}

func (v *ByteSliceVector) SetValid(i int) {
	v.valid[uint(i)>>log2WordSize] |= 1 << (uint(i) & (wordSize - 1))
}

func (v *ByteSliceVector) SetInvalid(i int) {
	v.valid[i>>log2WordSize] &^= 1 << (uint(i) & (wordSize - 1))
}

func (v *ByteSliceVector) AppendSlice(slice []byte) {

	// Remove when we use fixed vectors
	if v.l == v.c {
		panic("vector out of bounds")
	}

	v.buf = append(v.buf, slice...)
	v.offsets = append(v.offsets, len(v.buf))
	if v.HasNulls() {
		v.SetValid(v.l)
	}

	v.l++

}

func (v *ByteSliceVector) AppendNull() {

	if v.l == v.c {
		panic("vector out of bounds")
	}

	// v.offsets[v.l] = len(v.buf)
	v.offsets = append(v.offsets, len(v.buf))

	// TODO(gvelo) memset
	v.SetInvalid(v.l)
	v.l++
}

func (v *ByteSliceVector) Get(i int) []byte {

	var start int

	if i > 0 {
		start = v.offsets[i-1]
	}

	return v.buf[start:v.offsets[i]]

}

func (v *ByteSliceVector) Remaining() int {
	return v.c - v.l
}

func NewByteSliceVector(data []byte, offsets []int, valid []uint64) ByteSliceVector {
	return ByteSliceVector{
		offsets: offsets,
		buf:     data,
		valid:   valid,
	}
}

func NewByteSliceVectorFromByteArray(data [][]byte, valid []uint64) ByteSliceVector {

	buff := make([]byte, 0)
	offsets := make([]int, 0)

	offset := 0
	for _, it := range data {
		offset = offset + len(it)
		offsets = append(offsets, offset)
		buff = append(buff, it...)
	}
	return NewByteSliceVector(buff, offsets, valid)

}

func DefaultVectorPool() Pool {
	return &defaultPool{}
}

type defaultPool struct {
}

func (p *defaultPool) GetNotNullableInt32Vector() Int32Vector {
	return Int32Vector{
		valid: make([]uint64, 8192*2),
		buf:   make([]int32, 8192*2),
	}
}

func (p *defaultPool) GetNotNullableBoolVector() BoolVector {
	return BoolVector{
		valid: nil,
		buf:   make([]bool, 8192*2),
	}
}

func (p *defaultPool) GetNotNullableFloat64Vector() Float64Vector {
	return Float64Vector{
		valid: nil,
		buf:   make([]float64, 8192*2),
	}
}

func (p *defaultPool) GetNotNullableByteSliceVector() ByteSliceVector {
	// TODO: Check this, we should use a fixed segment to make faster the loading from disk,
	// the issue with this approach would be, that we need to manage different batch sizes
	// in numeric and byteslices vectors.
	return ByteSliceVector{
		valid:   nil,
		buf:     make([]byte, 0), // we do the appends
		offsets: make([]int, 0),  // we do the appends
		c:       0,
	}
}

func (p *defaultPool) GetNotNullableInt64Vector() Int64Vector {
	// TODO: parametrize vector capacity.
	return Int64Vector{
		valid: nil,
		buf:   make([]int64, 8192*2),
	}
}

func (p *defaultPool) GetBoolVector() BoolVector {
	// TODO: parametrize vector capacity.
	return BoolVector{
		valid: make([]uint64, 8192*2),
		buf:   make([]bool, 8192*2),
	}
}

func (p *defaultPool) GetFloat64Vector() Float64Vector {
	// TODO: parametrize vector capacity.
	return Float64Vector{
		valid: make([]uint64, 8192*2),
		buf:   make([]float64, 8192*2),
	}
}

func (*defaultPool) GetInt64Vector() Int64Vector {

	// TODO: parametrize vector capacity.
	return Int64Vector{
		valid: make([]uint64, 8192*2),
		buf:   make([]int64, 8192*2),
	}
}

func (*defaultPool) GetInt32Vector() Int32Vector {

	// TODO: parametrize vector capacity.
	return Int32Vector{
		valid: make([]uint64, 8192*2),
		buf:   make([]int32, 8192*4),
	}
}
func (*defaultPool) GetByteSliceVector() ByteSliceVector {

	// TODO: parametrize vector capacity.
	return ByteSliceVector{
		valid: make([]uint64, 8192*2),
		c:     8192 * 2,
	}

}

func (*defaultPool) PutInt64Vector(vector Int64Vector) {
	panic("implement me")
}

type ZeroVector struct {
}

func (v *ZeroVector) Len() int {
	return 0
}

func (v *ZeroVector) Cap() int {
	return 0
}

func (v *ZeroVector) HasNulls() bool {
	return false
}

func (v *ZeroVector) AsBytes() []byte {
	panic("not implemented")
}

func (v *ZeroVector) ValidityAsBytes() []byte {
	panic("not implemented")
}
