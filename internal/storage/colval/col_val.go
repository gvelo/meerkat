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

//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@../scalar_types.tmpldata col_val.gen.go.tmpl

package colval

type ByteSliceColValues struct {
	rid     []uint32
	data    []byte
	offsets []int
}

func NewByteSliceColValues(data []byte, rid []uint32, offsets []int) ByteSliceColValues {
	return ByteSliceColValues{
		offsets: offsets,
		data:    data,
		rid:     rid,
	}
}

func (v *ByteSliceColValues) Len() int {
	return len(v.offsets)
}

func (v *ByteSliceColValues) Rid() []uint32 {
	return v.rid
}

func (v *ByteSliceColValues) Data() []byte {
	return v.data
}

func (v *ByteSliceColValues) Offsets() []int {
	return v.offsets
}

func (v *ByteSliceColValues) Get(i int) []byte {

	var start int

	if i > 0 {
		start = v.offsets[i-1]
	}

	return v.data[start:v.offsets[i]]

}
