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

package executor

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage"
	"testing"
)

type BitmapIndexOperator struct {
	i     storage.IntIndex
	op    Operation
	value interface{}
}

func (b BitmapIndexOperator) Init() {
	panic("implement me")
}

func (b BitmapIndexOperator) Destroy() {
	panic("implement me")
}

func (b BitmapIndexOperator) Next() *roaring.Bitmap {
	panic("implement me")
}

func TestBinaryBitmapOperator(t *testing.T) {
	a := assert.New(t)

	op := NewBinaryBitmapOperator(list, 5, true)

	r := op.Next()
	k := r[0].(storage.IntVector).ValuesAsInt()
	v := r[1].(storage.IntVector).ValuesAsInt()

	a.Len(k, 5, "length is wrong ")
	a.Len(v, 5, "length is wrong ")

	c := 10
	for i := 0; i < len(k); i++ {
		a.Equal(k[i], c)
		a.Equal(v[i], c)
		c++
	}

}
