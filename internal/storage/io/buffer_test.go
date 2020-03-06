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

package io

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

const (
	numOfValues = 1024 * 2
	bufCap      = binary.MaxVarintLen64 * numOfValues
)

func TestDecoderBuffer_Uvarint(t *testing.T) {

	values := RandomUVarInt(numOfValues)

	e := NewEncoderBuffer(bufCap)

	for _, v := range values {
		e.WriteUvarint(v)
	}

	d := NewDecoderBuffer()

	d.SetBytes(e.Bytes())

	for i := 0; i < numOfValues; i++ {
		n := d.ReadUvarint()

		assert.Equal(t, values[i], n, "wrong int values")

	}

}

func RandomUVarInt(size int) []int {

	var r []int

	for j := 0; j < size; j++ {
		n := uint64(0)
		s := rand.Intn(8) + 1
		for i := 0; i < s; i++ {
			n = n | uint64(byte(rand.Int()))<<uint(8*i)
		}
		r = append(r, int(n))
	}

	return r

}
