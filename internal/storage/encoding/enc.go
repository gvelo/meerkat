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

import "meerkat/internal/storage/vector"

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
	FlushBlocks()
	Type() EncodingType
}

type IntEncoder interface {
	Encoder
	Encode(vec vector.IntVector)
}

type IntDecoder interface {
	Decode(block []byte, buf []int) []int
}

type UintEncoder interface {
	Encoder
	Encode(vec vector.IntVector)
}

type FloatEncoder interface {
	Encoder
	Encode(vec vector.FloatVector)
}

type ByteSliceEncoder interface {
	Encoder
	Encode(vec vector.ByteSliceVector)
}

type UUIDEncoder interface {
	Encoder
	Encode(vec vector.UUIDVector)
}

type ByteSliceDecoder interface {
	Decode(block []byte, data []byte, offsets []int) ([]byte, []int)
}

type UUIDDecoder interface {
	Decode(block []byte, data []byte) []byte
}

func DeltaEncode(src []int, dst []int) {

	dst[0] = src[0]

	for i := 1; i < len(src); i++ {
		dst[i] = src[i] - src[i-1]
	}

}

func DeltaDecode(src []int, dst []int) {

	dst[0] = src[0]

	for i := 1; i < len(src); i++ {
		dst[i] = dst[i-1] + src[i]
	}

}
