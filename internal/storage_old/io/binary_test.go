// Copyright 2019 The Meerkat Authors
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
	"io/ioutil"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newFileWriterReader(t *testing.T) {

	fw, err := NewBinaryWriter("/tmp/test1.bin")
	if err != nil {
		t.Error(err)
	}

	WriteString("HOLA MANOLA")
	WriteVarUint64(1)
	WriteVarUint64(100)
	WriteVarUint64(2023423423423432434)
	WriteFixedUint64(2023423423423432222)
	WriteFixedUint32(32)

	//var x uint64 =
	WriteVarUint64(math.MaxUint64)

	Close()

	dat, err := ioutil.ReadFile("/tmp/test1.bin")
	fr := NewBinaryReader(dat)
	if err != nil {
		t.Error(err)
	}
	s, err := ReadString()
	if err != nil {
		t.Error(err)
	}

	i1, err := ReadVarint64()
	if err != nil {
		t.Error(err)
	}

	i2, err := ReadVarint64()
	if err != nil {
		t.Error(err)
	}

	i3, err := ReadVarint64()
	if err != nil {
		t.Error(err)
	}

	i4, err := ReadFixed64()
	if err != nil {
		t.Error(err)
	}

	i5, err := ReadFixed32()
	if err != nil {
		t.Error(err)
	}

	i6, err := ReadVarint64()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "HOLA MANOLA", s)
	assert.Equal(t, uint64(1), i1)
	assert.Equal(t, uint64(100), i2)
	assert.Equal(t, uint64(2023423423423432434), i3)
	assert.Equal(t, uint64(2023423423423432222), i4)
	assert.Equal(t, uint64(32), i5)
	assert.True(t, math.MaxUint64 == i6)
}
