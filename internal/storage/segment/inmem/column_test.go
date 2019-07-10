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

package inmem

import (
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage/segment"
	"testing"
)

func TestColumnTimeStamp_Sort(t *testing.T) {

	a := assert.New(t)
	c := &ColumnTimeStamp{
		data:   make([]int, 0),
		prev:   0,
		sorted: false,
		fInfo:  &segment.FieldInfo{0, "_time", segment.FieldTypeTimestamp, true}}

	c.Add(123)
	c.Add(455)
	c.Add(34)
	c.Add(500)

	a.Nil(c.SortMap())

	c.Sort()
	sm := c.SortMap()
	a.Equal(2, sm[0])
	a.Equal(0, sm[1])
	a.Equal(1, sm[2])
	a.Equal(3, sm[3])

}
