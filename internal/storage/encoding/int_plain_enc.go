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
	"meerkat/internal/utils"
)

type IntPlainEncoder struct {
	bw BlockWriter
}

func NewIntPlainEncoder(bw BlockWriter) *IntPlainEncoder {
	return &IntPlainEncoder{
		bw: bw,
	}
}

func (e *IntPlainEncoder) Flush() {
}

func (e *IntPlainEncoder) FlushBlocks() {
}

func (e *IntPlainEncoder) Type() EncodingType {
	return Plain
}

func (e *IntPlainEncoder) Encode(v colval.IntColValues) {
	b := utils.I2B(v.Values())
	e.bw.WriteBlock(b, v.Rid()[0])
}
