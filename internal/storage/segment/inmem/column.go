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
	"github.com/psilva261/timsort"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage/segment"
)

type Column interface {
	Size() int
	FieldInfo() *segment.FieldInfo
	Add(value interface{})
	SetSortMap(sMap []int)
	Get(idx int) interface{}
	Data() interface{}
}

type ColumnInt struct {
	data  []int // TODO: use a buffer pool.
	fInfo *segment.FieldInfo
	sMap  []int
}

func (c *ColumnInt) Add(value interface{}) {
	c.data = append(c.data, value.(int))
}

func (c *ColumnInt) Size() int {
	return len(c.data)
}

func (c *ColumnInt) Data() interface{} {
	return c.data
}

func (c *ColumnInt) FieldInfo() *segment.FieldInfo {
	return c.fInfo
}

func (c *ColumnInt) SetSortMap(sMap []int) {
	c.sMap = sMap
}

func (c *ColumnInt) Get(idx int) interface{} {
	if c.sMap != nil {
		return c.data[c.sMap[idx]]
	} else {
		return c.data[idx]
	}
}

type ColumnTimeStamp struct {
	data    []int // TODO: use a buffer pool.
	prev    int
	sorted  bool
	fInfo   *segment.FieldInfo
	sortMap []int
}

func (c *ColumnTimeStamp) Add(value interface{}) {

	var v int
	switch value.(type) {
	case uint64:
		v = int(value.(uint64))
	default:
		v = value.(int)
	}

	if c.sorted && c.prev > v {
		c.sorted = false
	}

	c.prev = v

	c.data = append(c.data, v)

}

func (c *ColumnTimeStamp) Get(idx int) interface{} {
	return c.data[idx]
}

func (c *ColumnTimeStamp) Data() interface{} {
	return c.data
}

func (c *ColumnTimeStamp) Size() int {
	return len(c.data)
}

func (c *ColumnTimeStamp) Sorted() bool {
	return c.sorted
}

func (c *ColumnTimeStamp) Len() int {
	return len(c.data)
}

func (c *ColumnTimeStamp) Less(i int, j int) bool {
	return c.data[i] < c.data[j]
}

func (c *ColumnTimeStamp) Swap(i int, j int) {
	c.data[i], c.data[j] = c.data[j], c.data[i]
	c.sortMap[i], c.sortMap[j] = c.sortMap[j], c.sortMap[i]
}

func (c *ColumnTimeStamp) Sort() {
	c.sortMap = make([]int, len(c.data))
	for i := 0; i < len(c.data); i++ {
		c.sortMap[i] = i
	}
	timsort.TimSort(c)
	c.sorted = true
}

func (c *ColumnTimeStamp) SortMap() []int {
	return c.sortMap
}
func (c *ColumnTimeStamp) SetSortMap(sMap []int) {
	// do nothing, ts columns are already sorted
}

func (c *ColumnTimeStamp) First() int {
	return c.data[0]
}

func (c *ColumnTimeStamp) Last() int {
	return c.data[len(c.data)]
}

func (c *ColumnTimeStamp) FieldInfo() *segment.FieldInfo {
	return c.fInfo
}

type ColumnStr struct {
	data  []string // TODO: use a buffer pool.
	fInfo *segment.FieldInfo
	sMap  []int
}

func (c *ColumnStr) Add(value interface{}) {
	c.data = append(c.data, value.(string))
}

func (c *ColumnStr) Get(idx int) interface{} {
	if c.sMap != nil {
		return c.data[c.sMap[idx]]
	} else {
		return c.data[idx]
	}
}

func (c *ColumnStr) Data() interface{} {
	return c.data
}

func (c *ColumnStr) Size() int {
	return len(c.data)
}

func (c *ColumnStr) FieldInfo() *segment.FieldInfo {
	return c.fInfo
}

func (c *ColumnStr) SetSortMap(sMap []int) {
	c.sMap = sMap
}

type ColumnFloat struct {
	data  []float64 // TODO: use a buffer pool.
	fInfo *segment.FieldInfo
	sMap  []int
}

func (c *ColumnFloat) Add(value interface{}) {
	c.data = append(c.data, value.(float64))
}

func (c *ColumnFloat) Get(idx int) interface{} {
	if c.sMap != nil {
		return c.data[c.sMap[idx]]
	} else {
		return c.data[idx]
	}
}

func (c *ColumnFloat) Data() interface{} {
	return c.data
}

func (c *ColumnFloat) Size() int {
	return len(c.data)
}

func (c *ColumnFloat) FieldInfo() *segment.FieldInfo {
	return c.fInfo
}

func (c *ColumnFloat) SetSortMap(sMap []int) {
	c.sMap = sMap
}

func NewColumnt(fInfo *segment.FieldInfo) Column {
	switch fInfo.Type {
	case segment.FieldTypeTimestamp:
		return &ColumnTimeStamp{
			data:   make([]int, 0),
			prev:   0,
			sorted: true,
			fInfo:  fInfo,
		}
	case segment.FieldTypeInt:
		return &ColumnInt{
			data:  make([]int, 0),
			fInfo: fInfo,
		}
	case segment.FieldTypeKeyword:
		return &ColumnStr{
			data:  make([]string, 0),
			fInfo: fInfo,
		}
	case segment.FieldTypeText:
		return &ColumnStr{
			data:  make([]string, 0),
			fInfo: fInfo,
		}
	case segment.FieldTypeFloat:
		return &ColumnFloat{
			data:  make([]float64, 0),
			fInfo: fInfo,
		}
	default:
		log.Panic().Str("component", "column").Msg("unknown field type")
		return nil
	}
}
