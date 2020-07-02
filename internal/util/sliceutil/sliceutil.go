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

package sliceutil

import (
	"reflect"
	"unsafe"
)

const (
	Int64SizeBytes = int(unsafe.Sizeof(int64(0)))
	Int32SizeBytes = int(unsafe.Sizeof(int32(0)))
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

func I2B(s []int) []byte {
	return asByteSlice(Int64SizeBytes, unsafe.Pointer(&s))
}

func BS2S(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func I642B(s []int64) []byte {
	return asByteSlice(Int64SizeBytes, unsafe.Pointer(&s))
}

func U642B(s []uint64) []byte {
	return asByteSlice(Int64SizeBytes, unsafe.Pointer(&s))
}

func U32B(s []uint32) []byte {
	return asByteSlice(Int32SizeBytes, unsafe.Pointer(&s))
}

func F2B(s []float64) []byte {
	return asByteSlice(Int64SizeBytes, unsafe.Pointer(&s))
}

func B2I(b []byte) []int {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	var res []int
	s := (*reflect.SliceHeader)(unsafe.Pointer(&res))
	s.Data = h.Data
	s.Len = h.Len / Int64SizeBytes
	s.Cap = h.Cap / Int64SizeBytes
	return res
}

func B2I64(b []byte) []int64 {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	var res []int64
	s := (*reflect.SliceHeader)(unsafe.Pointer(&res))
	s.Data = h.Data
	s.Len = h.Len / Int64SizeBytes
	s.Cap = h.Cap / Int64SizeBytes
	return res
}


func B2U32(b []byte) []uint32 {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	var res []uint32
	s := (*reflect.SliceHeader)(unsafe.Pointer(&res))
	s.Data = h.Data
	s.Len = h.Len / Int32SizeBytes
	s.Cap = h.Cap / Int32SizeBytes
	return res
}

func B2S(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func S2B(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
