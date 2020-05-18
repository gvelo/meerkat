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

const (
	log2WordSize = uint(6)
	wordSize     = uint(64)
)

type Vector interface {
	Len() int
	Cap() int
}

type Pool interface {
	GetIntVector() IntVector
	GetByteSliceVector() ByteSliceVector
	PutIntVector(vector IntVector)
}

type ByteSliceVector struct {
	valid   []uint64
	buf     []byte
	offsets []int
	l       int
	c       int
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

func (v *ByteSliceVector) Data() []byte {
	return v.buf
}

func (v *ByteSliceVector) Buf() []byte {
	return v.buf
}

func (v *ByteSliceVector) Offsets() []int {
	return v.offsets
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

	if v.l == v.c {
		panic("vector out of bounds")
	}

	v.buf = append(v.buf, slice...)
	v.offsets = append(v.offsets, len(v.buf))
	v.SetValid(v.l)
	v.l++

}

func (v *ByteSliceVector) AppendNull() {

	if v.l == v.c {
		panic("vector out of bounds")
	}

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
		l:       len(offsets),
	}
}

func NewByteSliceVectorFromByteArray(v [][]byte) ByteSliceVector {
	buff := make([]byte, 0)
	offsets := make([]int, 0)

	offset := 0
	for _, it := range v {
		offset = offset + len(it)
		offsets = append(offsets, offset)
		buff = append(buff, it...)
	}
	return NewByteSliceVector(buff, offsets, []uint64{})
}

func DefaultVectorPool() Pool {
	return &defaultPool{}
}

type defaultPool struct {
}

func (*defaultPool) GetIntVector() IntVector {

	// TODO: parametrize vector capacity.

	return IntVector{
		valid: make([]uint64, 8192*2),
		buf:   make([]int, 8192*2),
	}
}

func (*defaultPool) GetByteSliceVector() ByteSliceVector {

	// TODO: parametrize vector capacity.

	return ByteSliceVector{
		valid: make([]uint64, 8192*2),
		c:     8192 * 2,
	}

}

func (*defaultPool) PutIntVector(vector IntVector) {
	panic("implement me")
}
