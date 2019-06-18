package encoding

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"unsafe"
)

const BYTES_IN_INT = 8
const BYTES_IN_FLOAT64 = 8

func UnsafeCastIntsToBytes(ints []int) []byte {
	length := len(ints) * BYTES_IN_INT
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&ints[0])), Len: length, Cap: length}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func UnsafeCastBytesToInts(bytes []byte) []int {
	length := len(bytes) / BYTES_IN_INT
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&bytes[0])), Len: length, Cap: length}
	return *(*[]int)(unsafe.Pointer(&hdr))
}

func CastBytesToString(b []byte) []string {
	buf := &bytes.Buffer{}
	buf.Write(b)
	str := make([]string, 0)
	gob.NewDecoder(buf).Decode(&str)
	return str
}

func CastStringToBytes(string []string) []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(string)
	return buf.Bytes()
}

func UnsafeCastFloatsToBytes(floats []float64) []byte {
	length := len(floats) * BYTES_IN_FLOAT64
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&floats[0])), Len: length, Cap: length}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func UnsafeCastBytesToFloats(bytes []byte) []float64 {
	length := len(bytes) / BYTES_IN_FLOAT64
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&bytes[0])), Len: length, Cap: length}
	return *(*[]float64)(unsafe.Pointer(&hdr))
}
