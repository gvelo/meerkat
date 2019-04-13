package io

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newFileWriterReader(t *testing.T) {

	fw, err := newFileWriter("/tmp/test1.bin")
	if err != nil {
		t.Error(err)
	}

	fw.WriteEncodedStringBytes("HOLA MANOLA")
	fw.WriteEncodedVarint(1)
	fw.WriteEncodedVarint(100)
	fw.WriteEncodedVarint(2023423423423432434)
	fw.WriteEncodedFixed64(2023423423423432222)
	fw.WriteEncodedFixed32(32)

	fw.Close()

	fr, err := newFileReader("/tmp/test1.bin")
	if err != nil {
		t.Error(err)
	}
	s, err := fr.DecodeStringBytes()
	if err != nil {
		t.Error(err)
	}

	i1, err := fr.DecodeVarint()
	if err != nil {
		t.Error(err)
	}

	i2, err := fr.DecodeVarint()
	if err != nil {
		t.Error(err)
	}

	i3, err := fr.DecodeVarint()
	if err != nil {
		t.Error(err)
	}

	i4, err := fr.DecodeFixed64()
	if err != nil {
		t.Error(err)
	}

	i5, err := fr.DecodeFixed32()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "HOLA MANOLA", s)
	assert.Equal(t, uint64(1), i1)
	assert.Equal(t, uint64(100), i2)
	assert.Equal(t, uint64(2023423423423432434), i3)
	assert.Equal(t, uint64(2023423423423432222), i4)
	assert.Equal(t, uint64(32), i5)
}
