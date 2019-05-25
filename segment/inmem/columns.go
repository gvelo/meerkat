package inmem

import (
	"eventdb/segment"
	"github.com/rs/zerolog/log"
)

type Column interface {
	Size() int
	FieldInfo() *segment.FieldInfo
	Add(value interface{})
	SetSortMap(sMap []int)
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

func (c *ColumnInt) FieldInfo() *segment.FieldInfo {
	return c.fInfo
}

func (c *ColumnInt) SetSortMap(sMap []int) {
	c.sMap = sMap
}

func (c *ColumnInt) Get(idx int) int {
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

	v := value.(int)

	if c.sorted && c.prev > v {
		c.sorted = false
	}

	c.data = append(c.data, v)

}

func (c *ColumnTimeStamp) Get(idx int) int {
	return c.data[idx]
}

func (c *ColumnTimeStamp) Size() int {
	return len(c.data)
}

func (c *ColumnTimeStamp) Sorted() bool {
	return c.sorted
}

func (c *ColumnTimeStamp) Sort() {
	// TODO Timsort the columnt.
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

func (c *ColumnStr) Get(idx int) string {
	if c.sMap != nil {
		return c.data[c.sMap[idx]]
	} else {
		return c.data[idx]
	}
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

func NewColumnt(fInfo *segment.FieldInfo) Column {
	switch fInfo.Type {
	case segment.FieldTypeTimestamp:
		return &ColumnTimeStamp{
			data:   make([]int, 0),
			prev:   0,
			sorted: false,
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
	default:
		log.Panic().Str("component", "column").Msg("unknown field type")
		return nil
	}
}
