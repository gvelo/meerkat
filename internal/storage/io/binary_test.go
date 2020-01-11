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
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"os"
	"path"
	"testing"
)

const testSize = 1024 * 1024

func Test(t *testing.T) {

	type TestCase struct {
		name string
		test BinaryStreamTest
	}

	testCases := []TestCase{
		{
			name: "testBytes",
			test: &TestBytes{},
		},
		{
			name: "testVarInt",
			test: &TestVarInt{},
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {

			fileName := path.Join(os.TempDir(), "binary_writer_test.bin")
			defer os.Remove(fileName)

			writer, err := NewBinaryWriter(fileName)

			if err != nil {
				t.Error(err)
				return
			}

			err = testCase.test.TestWrite(t, writer)

			if err != nil {
				return
			}

			err = writer.Close()

			if err != nil {
				t.Error(err)
				return
			}

			mf, err := MMap(fileName)

			if err != nil {
				t.Error(err)
				return
			}

			reader := mf.NewBinaryReader()

			testCase.test.TestRead(t, reader)

		})
	}

}

type BinaryStreamTest interface {
	TestWrite(t *testing.T, writer *BinaryWriter) error
	TestRead(t *testing.T, reader *BinaryReader)
}

type TestBytes struct {
	values [][]byte
}

func (tb *TestBytes) TestWrite(t *testing.T, w *BinaryWriter) error {

	tb.values = make([][]byte, testSize)

	for i := 0; i < testSize; i++ {
		tb.values[i] = randomBytes()
		err := w.WriteBytes(tb.values[i])
		if err != nil {
			t.Error(err)
			return err
		}
	}

	return nil

}

func (tb *TestBytes) TestRead(t *testing.T, reader *BinaryReader) {

	for i := 0; i < testSize; i++ {

		b, err := reader.ReadBytes()

		if err != nil {
			t.Error(err)
			return
		}
		if !assert.True(t, bytes.Equal(tb.values[i], b), "read bytes doesn't match") {
			return
		}
	}

	_, err := reader.ReadBytes()

	assert.Equal(t, err, io.ErrUnexpectedEOF, "")

	return

}

type TestVarInt struct {
	values []int
}

func (tv *TestVarInt) TestWrite(t *testing.T, w *BinaryWriter) error {

	tv.values = make([]int, testSize)

	for i := 0; i < testSize; i++ {
		tv.values[i] = rand.Int()
		err := w.WriteVarInt(tv.values[i])
		if err != nil {
			t.Error(err)
			return err
		}
	}

	return nil

}

func (tv *TestVarInt) TestRead(t *testing.T, reader *BinaryReader) {

	for i := 0; i < testSize; i++ {

		r, err := reader.ReadVarInt()

		if err != nil {
			t.Error(err)
			return
		}
		if !assert.Equal(t, tv.values[i], r, "read int doesn't match") {
			return
		}
	}

	_, err := reader.ReadVarInt()

	assert.Equal(t, err, io.ErrUnexpectedEOF, "")

	return

}

func randomBytes() []byte {
	i := rand.Intn(512)
	b := make([]byte, i)
	rand.Read(b)
	return b
}
