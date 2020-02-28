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
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestDeltaEnc(t *testing.T) {

	s := 1024

	src := make([]int, s)
	dst := make([]int, s)

	r := make([]int, s)

	src[0] = 10

	for i := 1; i < s; i++ {
		src[i] = src[i-1] + rand.Intn(20)
	}

	DeltaEncode(src, dst)
	DeltaDecode(dst, r)

	assert.Equal(t, src, r, "slice doesn't match")

}
