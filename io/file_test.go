package io

import (
	"testing"
)

func Test_newFileWriter(t *testing.T) {

	f, err := newFileWriter("/home/gabrielvelo/tmp/gofile/test")

	if err != nil {
		t.Error(err)
	}

	f.writeHeader(StringIndexV1)
	f.Close()

}
