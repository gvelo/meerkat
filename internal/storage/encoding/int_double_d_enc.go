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

type Int64DdEncoder struct {
	bw BlockWriter
}

func NewInt64DdEncoder(bw BlockWriter) *Int64DdEncoder {
	return &Int64DdEncoder{
		bw: bw,
	}
}

//TODO check this.
func (e *Int64DdEncoder) Flush() {

}

func (e *Int64DdEncoder) FlushBlocks() {

}

func (e *Int64DdEncoder) Type() Type {
	return DoubleDelta
}

func (e *Int64DdEncoder) Encode(v colval.Int64ColValues) {

	b := io.NewBuffer(64 * 1024) //TODO calculate this.
	b.WriteVarInt64(v.Values()[0])

	b.WriteVarInt64(v.Values()[1] - v.Values()[0])
	for i := 2; i < v.Len(); i++ {
		b.WriteVarInt64((v.Values()[i] - v.Values()[i-1]) - (v.Values()[i-1] - v.Values()[i-2]))
	}
	e.bw.WriteBlock(b.Data(), v.Rid()[0])
}
