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
	"meerkat/internal/storage/io"
)

type Int64DdDecoder struct {
}

func NewInt64DdDecoder() *Int64DdDecoder {
	return &Int64DdDecoder{}
}

func (d *Int64DdDecoder) Decode(block []byte, buf []int64) []int64 {

	// TODO check buf for enough room

	b := io.NewBinaryReader(block)
	buf[0] = b.ReadVarint64()
	buf[1] = b.ReadVarint64() + buf[0]
	for i := 2; i < len(buf); i++ {
		buf[i] = buf[i-1] + (buf[i-1] - buf[i-2]) + b.ReadVarint64()
	}

	return buf
}
