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
	"time"
)

type Index struct {
	ID         string
	Name       string
	Desc       string
	Partitions int
	Fields     []Field
}

type Field struct {
	ID      string
	Name    string
	Desc    string
	Type    FieldType
	Indexed bool
	Id      bool
}

type FieldType int

const (
	Int8 FieldType = iota
	UInt8
	Int16
	UInt16
	Int32
	UInt32
	Int64
	UInt64
	Float32
	Float64
	Text
	String
	Timestamp
)

type Catalog interface {
	CreateIndex(index Index) (Index, error)
	UpdateIndex(id string, index Index, ts time.Time) (Index, error)
	DeleteIndex(id string, ts time.Time)
	GetIndex(id string) Index
	CreateField(indexID string, field Field, ts time.Time) (Field, error)
	UpdateField(indexID string, fieldID string, field Field, ts time.Time) (Field, error)
	DeleteField(indexID string, fieldID string, ts time.Time)
	SetIndexAllocation(indexID string, alloc map[string][]uint, ts time.Time)
	SetGlobalAllocation(indexID string, alloc map[string][]uint, ts time.Time)
	GetCatalogState() CatalogState
	MergeCatalogState(state CatalogState)
	GetVersion() string
}

type CatalogState struct {
}
