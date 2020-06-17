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

package encoding

import (
	"meerkat/internal/storage/colval"
	"meerkat/internal/storage/io"
)

type EncodingType int

const (
	Plain EncodingType = iota
	Dict
	DictRleBitPacked
	DeltaBitPacked
	Snappy
)

type BlockWriter interface {
	WriteBlock(block []byte, baseRid uint32)
}

type Encoder interface {
	Flush()
	// TODO(gvelo) remove FlushBlocks, the encode api
	//  is block oriented and do not buffer blocks.
	//  roaringbitmap ?????
	FlushBlocks()
	Type() EncodingType
}

type Int64Encoder interface {
	Encoder
	Encode(v colval.Int64ColValues)
}

type Int32Encoder interface {
	Encoder
	Encode(v colval.Int32ColValues)
}

type Int64Decoder interface {
	Decode(block []byte, buf []int64) []int64
}

type Int32Encoder interface {
	Encoder
	Encode(v colval.Int32ColValues)
}

type Int32Decoder interface {
	Decode(block []byte, buf []int32) []int32
}

type Float64Encoder interface {
	Encoder
	Encode(v colval.Float64ColValues)
}

type Float64Decoder interface {
	Decode(block []byte, buf []float64) []float64
}

type ByteSliceEncoder interface {
	Encoder
	Encode(v colval.ByteSliceColValues)
}

type ByteSliceDecoder interface {
	Decode(block []byte) ([]byte, []int)
}

type BoolEncoder interface {
	Encoder
	Encode(v colval.BoolColValues)
}

type BoolDecoder interface {
	Decode(block []byte, buf []bool) []bool
}

func DeltaEncode(src []int, dst []int) {

	dst[0] = src[0]

	for i := 1; i < len(src); i++ {
		dst[i] = src[i] - src[i-1]
	}

}

func DeltaDecode(data []int) {

	for i := 1; i < len(data); i++ {
		data[i] = data[i-1] + data[i]
	}

}

func GetIntDecoder(d EncodingType, b []byte, bounds io.Bounds, blockLen int) (Int64Decoder, BlockReader) {

	var dec Int64Decoder
	var br BlockReader

	switch d {
	case Plain:
		dec = NewInt64PlainDecoder()
		br = NewScalarPlainBlockReader(b, bounds, blockLen)
	default:
		panic("unknown encoding type")
	}

	return dec, br
}

func GetBinaryDecoder(d EncodingType, b []byte, bounds io.Bounds) (ByteSliceDecoder, BlockReader) {

	var dec ByteSliceDecoder
	var br BlockReader

	switch d {
	case Plain:
		dec = NewByteSlicePlainDecoder()
		br = NewByteSliceBlockReader(b, bounds)
	case Snappy:
		dec = NewByteSliceSnappyDecoder()
		br = NewByteSliceBlockReader(b, bounds)
	default:
		panic("unknown encoding type")
	}

	return dec, br
}
