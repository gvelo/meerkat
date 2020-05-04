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

package schema2

import (
	"github.com/google/uuid"
	"time"
)

//go:generate protoc  -I . -I ../../build/proto/   --plugin ../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc:. ./schemapb/schema.proto

type Databases struct {
	byName map[string]*DatabaseDesc
	byId   map[uuid.UUID]*DatabaseDesc
}

func (d *Databases) ByName(name string) *DatabaseDesc {
	return d.byName[name]
}

func (d *Databases) ById(id uuid.UUID) *DatabaseDesc {
	return d.byId[id]
}

func (d *Databases) Each(f func(*DatabaseDesc)) {
	for _, db := range d.byId {
		f(db)
	}
}

func (d *Databases) copy() *Databases {
	c := *d
	c.byName = make(map[string]*DatabaseDesc)
	c.byId = make(map[uuid.UUID]*DatabaseDesc)
	for _, desc := range d.byId {
		dbCopy := desc.copy()
		c.byName[desc.name] = dbCopy
		c.byId[desc.id] = dbCopy
	}
	return &c
}

type DatabaseDesc struct {
	id      uuid.UUID
	name    string
	desc    string
	tables  *Tables
	created time.Time
	updated time.Time
}

func (db *DatabaseDesc) Id() uuid.UUID      { return db.id }
func (db *DatabaseDesc) Name() string       { return db.name }
func (db *DatabaseDesc) Desc() string       { return db.desc }
func (db *DatabaseDesc) Updated() time.Time { return db.updated }
func (db *DatabaseDesc) Created() time.Time { return db.created }
func (db *DatabaseDesc) Tables() *Tables    { return db.tables }

func (db *DatabaseDesc) copy() *DatabaseDesc {
	c := *db
	c.tables = db.tables.copy()
	return &c
}

type Tables struct {
	byName map[string]*TableDesc
	byId   map[uuid.UUID]*TableDesc
}

func (t *Tables) ByName(name string) *TableDesc { return t.byName[name] }
func (t *Tables) ById(id uuid.UUID) *TableDesc  { return t.byId[id] }

func (t *Tables) copy() *Tables {
	c := *t
	c.byId = make(map[uuid.UUID]*Tables)
	c.byName = make(map[string]*Tables)
	for _, desc := range t.byId {
		tableCopy := desc.copy()
		c.byId[tableCopy.id] = tableCopy
		c.byName[tableCopy.name] = tableCopy
	}
	return &c
}

func (t *Tables) Each(f func(*TableDesc)) {
	for _, tableDesc := range t.byId {
		f(tableDesc)
	}
}

type TableDesc struct {
	id              uuid.UUID
	databaseId      uuid.UUID
	name            string
	numOfPartitions uint32
	desc            string
	columns         *Columns
	pAlloc          *PartitionAlloc
	created         time.Time
	updated         time.Time
}

func (t *TableDesc) Id() uuid.UUID                   { return t.id }
func (t *TableDesc) DatabaseId() uuid.UUID           { return t.databaseId }
func (t *TableDesc) Name() string                    { return t.name }
func (t *TableDesc) NumOfPartitions() uint32         { return t.numOfPartitions }
func (t *TableDesc) Desc() string                    { return t.desc }
func (t *TableDesc) Columns() *Columns               { return t.columns }
func (t *TableDesc) PartitionAlloc() *PartitionAlloc { return t.pAlloc }
func (t *TableDesc) Created() time.Time              { return t.created }
func (t *TableDesc) Updated() time.Time              { return t.updated }

func (t *TableDesc) copy() *TableDesc {
	tCopy := *t
	tCopy.columns = t.columns.copy()
	return &tCopy
}

type Columns struct {
	byName map[string]*ColumnDesc
	byId   map[uuid.UUID]*ColumnDesc
}

func (c *Columns) ByName(name string) *ColumnDesc { return c.byName[name] }
func (c *Columns) ById(id uuid.UUID) *ColumnDesc  { return c.byId[id] }

func (c *Columns) copy() *Columns {
	copy := *c
	copy.byName = make(map[string]*ColumnDesc)
	copy.byId = make(map[uuid.UUID]*ColumnDesc)
	for _, desc := range c.byId {
		colCopy := desc.copy()
		copy.byId[desc.id] = colCopy
		copy.byName[desc.name] = colCopy
	}
	return &copy
}

func (c *Columns) Each(f func(desc *ColumnDesc)) {
	for _, col := range c.byId {
		f(col)
	}
}

type ColumnDesc struct {
	id       uuid.UUID
	tableId  uuid.UUID
	name     string
	desc     string
	dataType DataType
	nullable bool
	indexed  bool
	fullText bool
	created  time.Time
	updated  time.Time
}

func (c *ColumnDesc) Id() uuid.UUID      { return c.id }
func (c *ColumnDesc) TableId() uuid.UUID { return c.tableId }
func (c *ColumnDesc) Name() string       { return c.name }
func (c *ColumnDesc) Desc() string       { return c.desc }
func (c *ColumnDesc) DataType() DataType { return c.dataType }
func (c *ColumnDesc) Nullable() bool     { return c.nullable }
func (c *ColumnDesc) Indexed() bool      { return c.indexed }
func (c *ColumnDesc) Created() time.Time { return c.created }
func (c *ColumnDesc) Updated() time.Time { return c.updated }

func (c *ColumnDesc) copy() *ColumnDesc {
	copy := *c
	return &copy
}

type PartitionAlloc struct {
	alloc [][]string
}

func (p *PartitionAlloc) Each(partition int, f func(string)) {
	for _, member := range p.alloc[partition] {
		f(member)
	}
}

type OpType int

const (
	OpCreate = iota
	OpUpdate
	OpDelete
)

type EntityType int

const (
	Database = iota
	Table
	Column
	PartitionMap
)

type DatabaseInfo struct {
	Name string
	Desc string
}

type TableInfo struct {
	Name            string
	Desc            string
	NumOfPartitions int
}

type Op struct {
	OpType     OpType
	EntityType EntityType
	EntityId   uuid.UUID
	EntityName string
	OpData     interface{}
}

type Schema struct {
}

func (s *Schema) Databases() databases {

}

func (s *Schema) Update(operations []Op) {

}

//////////////////

type Root struct {
	db          *Databases
	tablesById  map[uuid.UUID]*TableDesc
	columnsById map[uuid.UUID]*ColumnDesc
	pAllocById  map[uuid.UUID]*PartitionAlloc
}

func (r *Root) copy() *Root {

	root := &Root{
		db: r.db.copy(),
	}

	root.index()

	return root

}

func (r *Root) index() {
	r.db.Each(func(desc *DatabaseDesc) {
		desc.tables.Each(func(desc *TableDesc) {
			r.tablesById[desc.id] = desc
			r.pAllocById[desc.id] = desc.pAlloc
			desc.columns.Each(func(desc *ColumnDesc) {
				r.columnsById[desc.id] = desc
			})
		})
	})
}

func (s *Schema) CreateDatabase(r *Root, db *DatabaseDesc) {

}

func (s *Schema) UpdateDatabase(tx Root, db *DatabaseDesc) {

}
