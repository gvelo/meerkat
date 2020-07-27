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
	"io"
	"math"
)

// uvnan is the constant returned from math.NaN().
const uvnan = 0x7FF8000000000001

// TODO: check pools.
type Float64TszDecoder struct {
}

func NewFloat64TszDecoder() *Float64TszDecoder {
	return &Float64TszDecoder{}
}

func (d *Float64TszDecoder) Decode(block []byte, buf []float64) []float64 {

	dec := &floatDecoderIter{}
	dec.SetBytes(block)

	// TODO: FIX. this should explode
	//if len(buf) < len(data) {
	//	panic("there isn't enough space to decode integer values")
	//}

	n := 0
	for dec.Error() == nil && dec.Next() {
		buf[n] = dec.Value()
		n++
	}
	dec.Next()
	buf[n] = dec.Value()
	// 0.6624402935722433
	if dec.Error() != nil && dec.Error() != io.EOF {
		panic(fmt.Sprintf("Errir %v", dec.Error()))
	}

	return buf[:n]
}

type floatDecoderIter struct {
	val uint64

	leading  uint64
	trailing uint64

	br *bitstream.BitReader
	b  []byte

	first    bool
	finished bool

	err error
}

// SetBytes initializes the decoder with b. Must call before calling Next().
func (it *floatDecoderIter) SetBytes(b []byte) error {
	var v uint64
	if len(b) == 0 {
		v = uvnan
	} else {
		it.br = bitstream.NewReader(bytes.NewReader(b))

		var err error
		v, err = it.br.ReadBits(64)
		if err != nil {
			return err
		}
	}

	// Reset all fields.
	it.val = v
	it.leading = 0
	it.trailing = 0
	it.b = b
	it.first = true
	it.finished = false
	it.err = nil

	return nil
}

// Next returns true if there are remaining values to read.
func (it *floatDecoderIter) Next() bool {
	if it.err != nil || it.finished {
		return false
	}

	if it.first {
		it.first = false

		// mark as finished if there were no values.
		if it.val == uvnan { // IsNaN
			it.finished = true
			return false
		}

		return true
	}

	// read compressed value
	var bit bool
	if v, err := it.br.ReadBit(); err != nil {
		it.err = err
		return false
	} else {
		bit = v != bitstream.Zero
	}

	if !bit {
		// it.val = it.val
	} else {
		var bit bool
		if v, err := it.br.ReadBit(); err != nil {
			it.err = err
			return false
		} else {
			bit = v != bitstream.Zero
		}

		if !bit {
			// reuse leading/trailing zero bits
			// it.leading, it.trailing = it.leading, it.trailing
		} else {
			bits, err := it.br.ReadBits(5)
			if err != nil {
				it.err = err
				return false
			}
			it.leading = bits

			bits, err = it.br.ReadBits(6)
			if err != nil {
				it.err = err
				return false
			}
			mbits := bits
			// 0 significant bits here means we overflowed and we actually need 64; see comment in encoder
			if mbits == 0 {
				mbits = 64
			}
			it.trailing = 64 - it.leading - mbits
		}

		mbits := int(64 - it.leading - it.trailing)
		bits, err := it.br.ReadBits(mbits)
		if err != nil {
			it.err = err
			return false
		}

		vbits := it.val
		vbits ^= (bits << it.trailing)

		if vbits == uvnan { // IsNaN
			it.finished = true
			return false
		}
		it.val = vbits
	}

	return true
}

// Value returns the current float64 value.
func (it *floatDecoderIter) Value() float64 {
	return math.Float64frombits(it.val)
}

// Error returns the current decoding error.
func (it *floatDecoderIter) Error() error {
	return it.err
}
