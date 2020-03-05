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
	IsNull(i int) bool
	SetNull(i int)
}

type ByteSliceVector struct {
	nulls   []uint64
	data    []byte
	offsets []int
	l       int
}

func (v *ByteSliceVector) Len() int {
	return len(v.offsets)
}

func (v *ByteSliceVector) SetLen(l int) {
	v.l = l
}

func (v *ByteSliceVector) Data() []byte {
	return v.data[:v.l]
}

func (v *ByteSliceVector) Offsets() []int {
	return v.offsets[:v.l]
}

func (v *ByteSliceVector) IsNull(i int) bool {
	return v.nulls[uint(i)>>log2WordSize]&(1<<(uint(i)&(wordSize-1))) != 0
}

func (v *ByteSliceVector) SetNull(i int) {
	v.nulls[uint(i)>>log2WordSize] |= 1 << (uint(i) & (wordSize - 1))
}

func (v *ByteSliceVector) Get(i int) []byte {

	var start int

	if i > 0 {
		start = v.offsets[i-1]
	}

	return v.data[start:v.offsets[i]]

}

func NewByteSliceVector(data []byte, nulls []uint64, offsets []int) ByteSliceVector {
	return ByteSliceVector{
		offsets: offsets,
		data:    data,
		nulls:   nulls,
	}
}
