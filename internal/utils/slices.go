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

package utils

import (
	"reflect"
	"unsafe"
)

const (
	// Int64SizeBytes specifies the number of bytes required to store a single int64 in memory
	Int64SizeBytes = int(unsafe.Sizeof(int64(0)))
)

func asByteSlice(size int, p unsafe.Pointer) []byte {
	h := (*reflect.SliceHeader)(p)
	var res []byte
	s := (*reflect.SliceHeader)(unsafe.Pointer(&res))
	s.Data = h.Data
	s.Len = h.Len * size
	s.Cap = h.Cap * size
	return res
}

func IntAsByte(s []int) []byte {
	return asByteSlice(Int64SizeBytes, unsafe.Pointer(&s))
}

func Float64AsByte(s []float64) []byte {
	return asByteSlice(Int64SizeBytes, unsafe.Pointer(&s))
}

func BytesAsInt(b []byte) []int {

	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	var res []int
	s := (*reflect.SliceHeader)(unsafe.Pointer(&res))
	s.Data = h.Data
	s.Len = h.Len / Int64SizeBytes
	s.Cap = h.Cap / Int64SizeBytes
	return res
}
