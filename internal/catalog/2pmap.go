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

import "time"

type entry struct {
	time        time.Time
	value       interface{}
	deleted     bool
	deletedTime time.Time
}

// TwoPhaseMap is an 2P Set implementation using a map.
type TwoPhaseMap struct {
	time   time.Time
	values map[string]*entry
}

func NewTwoPhaseMap() *TwoPhaseMap {
	return &TwoPhaseMap{
		values: make(map[string]*entry),
	}
}

func (m *TwoPhaseMap) Set(id string, value interface{}, time time.Time) bool {

	e, found := m.values[id]

	if found {

		if e.deleted {
			return false
		}

		if e.time.After(time) {
			return false
		}

	}

	m.values[id] = &entry{time: time, value: value}

	return true

}

func (m *TwoPhaseMap) Get(id string) (interface{}, bool) {

	e, found := m.values[id]

	if found && !e.deleted {
		return e.value, true
	}

	return nil, false
}

func (m *TwoPhaseMap) Remove(id string, time time.Time) {

	e := &entry{
		deletedTime: time,
		deleted:     true,
	}

	m.values[id] = e
}
