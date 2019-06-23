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
