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

/*
This code is originally from: https://github.com/dgryski/go-tsz and  https://github.com/influxdata/influxdb has been modified to remove
use our encoder interface.

It implements the float compression as presented in: http://www.vldb.org/pvldb/vol8/p1816-teller.pdf.
This implementation uses a sentinel value of NaN which means that float64 NaN cannot be stored using
this version.
*/

type Float64PlainEncoder struct {
	bw BlockWriter
}

func NewFloat64PlainEncoder(bw BlockWriter) *Float64PlainEncoder {
	return &Float64PlainEncoder{
		bw: bw,
	}
}

func (e *Float64PlainEncoder) Flush() {
}

func (e *Float64PlainEncoder) FlushBlocks() {
}

func (e *Float64PlainEncoder) Type() Type {
	return Plain
}

func (e *Float64PlainEncoder) Encode(v colval.Float64ColValues) {
	b := sliceutil.F2B(v.Values())
	e.bw.WriteBlock(b, v.Rid()[0])
}
