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
	"errors"
	"meerkat/internal/schema"
	"meerkat/internal/storage"
)

// save a slice that maps the vector index to column ids in cf.

type ProcessedField struct {
	Fields []schema.Field
}

func (pf *ProcessedField) FindField(name string) (schema.Field, int, error) {

	for i, it := range pf.Fields {
		if it.Name == name {
			return it, i, nil
		}
	}

	return schema.Field{}, 0, errors.New("Field Not Found")
}

type Context interface {
	Value(key string, value interface{})
	Get(key string) (interface{}, bool)
	GetFieldProcessed() *ProcessedField
	SetFieldProcessed(fields []schema.Field)
	Segment() *storage.Segment
	Sz() int
}

type ctx struct {
	s  *storage.Segment
	m  map[string]interface{}
	fp *ProcessedField
	sz int
}

func NewContext(s *storage.Segment, sz int) Context {
	return &ctx{
		s:  s,
		m:  make(map[string]interface{}),
		sz: sz,
	}
}

func (c *ctx) Value(key string, value interface{}) {
	c.m[key] = value
}

func (c *ctx) Sz() int {
	return c.sz
}

func (c *ctx) Get(key string) (interface{}, bool) {
	i, ok := c.m[key]
	return i, ok
}

func (c *ctx) Segment() *storage.Segment {
	return c.s
}

func (c *ctx) SetFieldProcessed(fields []schema.Field) {
	c.fp = &ProcessedField{
		Fields: fields,
	}
}

func (c *ctx) GetFieldProcessed() *ProcessedField {
	return c.fp
}
