package exec

type Cursor interface {
	Size() int
	Data() []interface{}
	Labels() []string
}

type Int64Cursor struct {
	data []int64
}

func (c *Int64Cursor) Size() int {
	return len(c.data)
}

func (c *Int64Cursor) Data() interface{} {
	return c.data
}

type MultiFieldCursor interface {
	Size() int
	Cursors() []Cursor
	FieldNames() []string
}
