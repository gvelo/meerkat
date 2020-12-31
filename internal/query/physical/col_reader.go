package physical

import (
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
)

func NewInt64ColumnReader(iterator storage.Int64Iterator) *Int64ColumnReader {
	return &Int64ColumnReader{
		iterator: iterator,
	}
}

type Int64ColumnReader struct {
	iterator storage.Int64Iterator
}

func (c *Int64ColumnReader) Init()  {}
func (c *Int64ColumnReader) Close() {}

func (c *Int64ColumnReader) Next() vector.Vector {
	v := c.iterator.Next()
	return &v
}
