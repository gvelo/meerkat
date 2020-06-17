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
	"meerkat/internal/util/sliceutil"
)

type Int64PlainEncoder struct {
	bw BlockWriter
}

func NewInt64PlainEncoder(bw BlockWriter) *Int64PlainEncoder {
	return &Int64PlainEncoder{
		bw: bw,
	}
}

func (e *Int64PlainEncoder) Flush() {
}

func (e *Int64PlainEncoder) FlushBlocks() {
}

func (e *Int64PlainEncoder) Type() EncodingType {
	return Plain
}

func (e *Int64PlainEncoder) Encode(v colval.Int64ColValues) {
	b := sliceutil.I642B(v.Values())
	e.bw.WriteBlock(b, v.Rid()[0])
}
