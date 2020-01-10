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
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"
)

const test_size = 1024 * 1024

func Test_Binary(t *testing.T) {

	fileName := path.Join(os.TempDir(), "binary_file_test")
	//defer os.Remove(fileName)
	fmt.Println(fileName)

	w, err := NewBinaryWriter(fileName)

	if err != nil {
		t.Error(err)
		return
	}

	values := make([][]byte, test_size)
	fmt.Println("escribiendo")
	values[0] = randomBytes()
	start := time.Now()
	for i := 0; i < test_size; i++ {

		err := w.WriteBytes(values[0])
		if err != nil {
			t.Error(err)
			return
		}
	}
	duration := time.Since(start)
	fmt.Printf("fin escribiendo %v", duration)
	err = w.Close()

	if err != nil {
		t.Error(err)
	}
	fmt.Println("leyendo")
	mf, err := MMap(fileName)

	if err != nil {
		t.Error(err)
		return
	}

	r := mf.NewBinaryReader()

	for i := 0; i < test_size; i++ {

		_, err := r.ReadBytes()

		if err != nil {
			t.Error(err)
			return
		}

		//assert.True(t, bytes.Equal(values[i], b), "read bytes doesn't match")
	}

	_, err = r.ReadBytes()

	assert.Equal(t, err, io.ErrUnexpectedEOF)

}

func randomBytes() []byte {
	i := rand.Intn(789)
	fmt.Println(i)
	b := make([]byte, i)
	rand.Read(b)
	return b
}
