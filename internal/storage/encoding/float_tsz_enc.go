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
	"bytes"
	"fmt"
	"github.com/dgryski/go-bitstream"
	"math"
	"math/bits"
	"meerkat/internal/storage/colval"
)

type Float64TszEncoder struct {
	bw  BlockWriter
	val float64
	err error

	leading  uint64
	trailing uint64

	buf       bytes.Buffer
	binWriter *bitstream.BitWriter

	first    bool
	finished bool
}

func NewFloat64TszEncoder(bw BlockWriter) *Float64TszEncoder {
	s := Float64TszEncoder{
		first:   true,
		leading: ^uint64(0),
		bw:      bw,
	}

	s.binWriter = bitstream.NewWriter(&s.buf)

	return &s
}

func (e *Float64TszEncoder) Flush() {
	// Flush indicates there are no more values to encode.
	if !e.finished {
		// write an end-of-stream record
		e.finished = true
		e.Write(math.NaN())
		e.binWriter.Flush(bitstream.Zero)
	}
}

func (e *Float64TszEncoder) FlushBlocks() {
}

func (e *Float64TszEncoder) Type() EncodingType {
	return Plain
}

func (e *Float64TszEncoder) Encode(v colval.Float64ColValues) {

	for _, f := range v.Values() {
		e.Write(f)
	}
	e.Flush()
	e.bw.WriteBlock(e.buf.Bytes(), v.Rid()[0])
}

// Write encodes v to the underlying buffer.
func (e *Float64TszEncoder) Write(v float64) {
	// Only allow NaN as a sentinel value
	if math.IsNaN(v) && !e.finished {
		e.err = fmt.Errorf("unsupported value: NaN")
		return
	}
	if e.first {
		// first point
		e.val = v
		e.first = false
		e.binWriter.WriteBits(math.Float64bits(v), 64)
		return
	}

	vDelta := math.Float64bits(v) ^ math.Float64bits(e.val)

	if vDelta == 0 {
		e.binWriter.WriteBit(bitstream.Zero)
	} else {
		e.binWriter.WriteBit(bitstream.One)

		leading := uint64(bits.LeadingZeros64(vDelta))
		trailing := uint64(bits.TrailingZeros64(vDelta))

		// Clamp number of leading zeros to avoid overflow when encoding
		leading &= 0x1F
		if leading >= 32 {
			leading = 31
		}

		// TODO(dgryski): check if it's 'cheaper' to reset the leading/trailing bits instead
		if e.leading != ^uint64(0) && leading >= e.leading && trailing >= e.trailing {
			e.binWriter.WriteBit(bitstream.Zero)
			e.binWriter.WriteBits(vDelta>>e.trailing, 64-int(e.leading)-int(e.trailing))
		} else {
			e.leading, e.trailing = leading, trailing

			e.binWriter.WriteBit(bitstream.One)
			e.binWriter.WriteBits(leading, 5)

			// Note that if leading == trailing == 0, then sigbits == 64.  But that
			// value doesn't actually fit into the 6 bits we have.
			// Luckily, we never need to encode 0 significant bits, since that would
			// put us in the other case (vdelta == 0).  So instead we write out a 0 and
			// adjust it back to 64 on unpacking.
			sigbits := 64 - leading - trailing
			e.binWriter.WriteBits(sigbits, 6)
			e.binWriter.WriteBits(vDelta>>trailing, int(sigbits))
		}
	}

	e.val = v
}
