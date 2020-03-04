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

import "meerkat/internal/utils"

type Vector interface {
	Len() int
	Rid() []uint32
	Data() []byte
}

type IntVector interface {
	Vector
	Values() []int
}

type FloatVector interface {
	Vector
	Values() []float64
}

type ByteSliceVector interface {
	Vector
	Offsets() []int
	Get(i int) []byte
}

type UUIDVector interface {
	Vector
	Get(i int) []byte
}

type intVector struct {
	data []int
	rid  []uint32
}

func (v intVector) Len() int {
	return len(v.data)
}

func (v intVector) Rid() []uint32 {
	return v.rid
}

func (v intVector) Values() []int {
	return v.data
}

func (v intVector) Data() []byte {
	return utils.I2B(v.data)
}

func NewIntVector(data []int, rid []uint32) IntVector {
	return &intVector{
		data: data,
		rid:  rid,
	}
}

type floatVector struct {
	data []float64
	rid  []uint32
}

func (v floatVector) Len() int {
	return len(v.data)
}

func (v floatVector) Rid() []uint32 {
	return v.rid
}

func (v floatVector) Values() []float64 {
	return v.data
}

func (v floatVector) Data() []byte {
	return utils.F2B(v.data)
}

func NewFloatVector(data []float64, rid []uint32) FloatVector {
	return floatVector{
		data: data,
		rid:  rid,
	}
}

type byteSliceVector struct {
	rid     []uint32
	data    []byte
	offsets []int
}

func NewByteSliceVector(data []byte, rid []uint32, offsets []int) ByteSliceVector {
	return byteSliceVector{
		offsets: offsets,
		data:    data,
		rid:     rid,
	}
}

func (v byteSliceVector) Len() int {
	return len(v.offsets)
}

func (v byteSliceVector) Rid() []uint32 {
	return v.rid
}

func (v byteSliceVector) Data() []byte {
	return v.data
}

func (v byteSliceVector) Offsets() []int {
	return v.offsets
}

func (v byteSliceVector) Get(i int) []byte {

	var start int

	if i > 0 {
		start = v.offsets[i-1]
	}

	return v.data[start:v.offsets[i]]

}
