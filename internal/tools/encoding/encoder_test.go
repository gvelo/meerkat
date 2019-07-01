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

package encoding

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"meerkat/internal/tools/utils"
	"testing"
)

func TestSnappyEncoder_Encode(t *testing.T) {

	assert := assert.New(t)
	enc := NewSnappyEncoder()

	slice := &utils.SlicerString{}

	for i := 0; i < 10000; i++ {
		slice.Add(fmt.Sprintf("String numero %d", i))
	}

	s := enc.Encode(slice)
	assert.NotNil(s)

	dec := NewSnappyDecoder()
	p := s[0]
	r := dec.Decode(p.Data)
	print("s %v", r.(string))

}

func TestRLEIntegerEncoder_Encode(t *testing.T) {

	assert := assert.New(t)
	enc := NewRLEEncoder()

	slice := &utils.SlicerInt{}
	x := 0
	for i := 0; i < 10000; i++ {
		if i%1000 == 0 {
			x = x + 1
		}
		slice.Add(x)
	}

	s := enc.Encode(slice)

	assert.NotNil(s)
	assert.Equal(10, len(s))

	dec := NewRLEDecoder()
	p := s[0]
	r := dec.Decode(p.Data)
	assert.Equal(1000, len(r.([]int)))

}
