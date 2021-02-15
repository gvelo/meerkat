package physical

import (
	"fmt"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
)

func NewColumnReaderOp(iterator storage.Iterator) *ColumnReaderOp {
	return &ColumnReaderOp{
		iterator: iterator,
	}
}

type ColumnReaderOp struct {
	iterator storage.Iterator
}

func (c *ColumnReaderOp) Init()  {}
func (c *ColumnReaderOp) Close() {}

func (c *ColumnReaderOp) Next() vector.Vector {

	fmt.Println("ColumnReaderOp Next()")

	if c.iterator.HasNext() {
		switch i := c.iterator.(type) {
		case storage.Int64Iterator:
			v := i.Next()
			return &v
		case storage.ByteSliceIterator:
			v := i.Next()
			return &v
		default:
			panic("unknown iterator")
		}
	}

	// if we are at EOF return zero length vector.
	// TODO(gvelo): the iterator semantic should be consistent with
	// operators. Iterator.Next() should return a zero length vector
	// to signal EOF

	return &vector.ZeroVector{}

}

func (c *ColumnReaderOp) Accept(v Visitor) {}
