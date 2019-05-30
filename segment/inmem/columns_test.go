package inmem

import (
	"eventdb/segment"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColumnTimeStamp_Sort(t *testing.T) {

	a := assert.New(t)
	c := &ColumnTimeStamp{
		data:    make([]int, 0),
		prev:    0,
		sorted:  false,
		fInfo:   &segment.FieldInfo{0, "_time", segment.FieldTypeTimestamp, true},
		sortMap: make([]int, 0)}

	c.Add(123)
	c.Add(455)
	c.Add(34)
	c.Add(500)

	sm := c.SortMap()
	a.Equal(1, sm[0])
	a.Equal(2, sm[1])
	a.Equal(3, sm[2])
	a.Equal(4, sm[3])

	c.Sort()

	a.Equal(3, sm[0])
	a.Equal(1, sm[1])
	a.Equal(2, sm[2])
	a.Equal(4, sm[3])

}
