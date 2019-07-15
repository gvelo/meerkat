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

package catalog

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	key   = "testkey"
	value = "testvalue"
)

var time1 = time.Now()
var time2 = time.Now().Add(10 * time.Second)

func Test2PMapSet(t *testing.T) {

	m := NewTwoPhaseMap()
	r := m.Set(key, value, time.Now())
	assert.True(t, r, "Set() should set any value on initial state")
	v, found := m.Get(key)
	assert.True(t, found, "wrong lookup result")
	assert.Equal(t, v, value)

}

func Test2PMapSetOldValue(t *testing.T) {

	m := NewTwoPhaseMap()
	r := m.Set(key, value, time2)
	assert.True(t, r, "Set() should set any value on initial state")

	r = m.Set(key, "value2", time1)
	assert.False(t, r, "Set() last write should win")

	v, found := m.Get(key)
	assert.True(t, found, "wrong lookup result")
	assert.Equal(t, value, v)

}

func Test2PMapRemove(t *testing.T) {

	m := NewTwoPhaseMap()
	m.Remove(key, time1)
	r := m.Set(key, value, time2)
	assert.False(t, r, "Set() should set any value if it was delete")

}

func Test2PMapRemoveExisting(t *testing.T) {

	m := NewTwoPhaseMap()
	r := m.Set(key, value, time2)
	assert.True(t, r, "Set() should set any value on initial state")

	m.Remove(key, time2)

	v, found := m.Get(key)
	assert.False(t, found, "wrong lookup result")
	assert.Nil(t, v)

}
