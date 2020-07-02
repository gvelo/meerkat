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
	"meerkat/internal/util/sliceutil"
)

type Int64PlainDecoder struct {
}

func NewInt64PlainDecoder() *Int64PlainDecoder {
	return &Int64PlainDecoder{}
}

func (d *Int64PlainDecoder) Decode(block []byte, buf []int64) []int64 {

	data := sliceutil.B2I64(block)

	if len(buf) < len(data) {
		panic("there isn't enough space to decode integer values")
	}

	n := copy(buf, data)

	return buf[:n]

}
