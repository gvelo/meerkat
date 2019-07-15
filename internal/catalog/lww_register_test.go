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

func TestLWWRegisterSetOnInitialState(t *testing.T) {

	reg := &LWWRegister{}
	r := reg.Set(1, time.Now())
	assert.True(t, r, "Set() should set any value on initial state")
	assert.Equal(t, 1, reg.Get())

}

func TestLWWRegisterSetOldValue(t *testing.T) {

	reg := &LWWRegister{
		value: 1,
		time:  time.Now(),
	}

	r := reg.Set(0, time.Now().Add(-10*time.Second))
	assert.False(t, r, "Set() should not set the value if the new TS is before the curren TS")
	assert.Equal(t, 1, reg.Get())

}


func TestLWWRegisterSetValue(t *testing.T) {

	reg := &LWWRegister{
		value: 1,
		time:  time.Now(),
	}

	r := reg.Set(2, time.Now().Add(10*time.Second))
	assert.True(t, r, "Set() should return true when the value was successfully set")
	assert.Equal(t, 2, reg.Get())

}
