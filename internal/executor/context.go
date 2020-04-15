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
	Segment() storage.ColumnFinder
	IndexInfo() *schema.IndexInfo
	GetFieldProcessed() *ProcessedField
	SetFieldProcessed(fields []schema.Field)
}

type ctx struct {
	cf storage.ColumnFinder
	m  map[string]interface{}
	ii *schema.IndexInfo
	fp *ProcessedField
}

func NewContext(s storage.ColumnFinder, ii *schema.IndexInfo) Context {
	return &ctx{
		cf: s,
		m:  make(map[string]interface{}),
		ii: ii,
	}
}

func (c *ctx) Value(key string, value interface{}) {
	c.m[key] = value
}

func (c *ctx) Get(key string) (interface{}, bool) {
	i, ok := c.m[key]
	return i, ok
}

func (c *ctx) Segment() storage.ColumnFinder {
	return c.cf
}

func (c *ctx) IndexInfo() *schema.IndexInfo {
	return c.ii
}

func (c *ctx) SetFieldProcessed(fields []schema.Field) {
	c.fp = &ProcessedField{
		Fields: fields,
	}
}

func (c *ctx) GetFieldProcessed() *ProcessedField {
	return c.fp
}
