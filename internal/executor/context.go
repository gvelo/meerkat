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
	"meerkat/internal/storage"
)

// save a slice that maps the vector index to column ids in segment.
const ColumnIndexToColumnName = "COLUMN_INDEX_TO_COLUMN_NAME"

type Context interface {
	Value(key string, value interface{})
	Get(key string) (interface{}, bool)
	Segment() storage.Segment
}

type ctx struct {
	segment storage.Segment
	m       map[string]interface{}
}

func NewContext(s storage.Segment) Context {
	return &ctx{segment: s, m: make(map[string]interface{})}
}

func (c ctx) Value(key string, value interface{}) {
	c.m[key] = value
}

func (c ctx) Get(key string) (interface{}, bool) {
	i, ok := c.m[key]
	return i, ok
}

func (c ctx) Segment() storage.Segment {
	return c.segment
}
