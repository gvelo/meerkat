package exec

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/spf13/cast"
)

type Cursor interface {
	Size() int
	Bitmap() *roaring.Bitmap
	Record() array.Record
}

type CursorImpl struct {
	bitmap *roaring.Bitmap
	record array.Record
}

func (c *CursorImpl) Size() int64 {

	if c.record != nil {
		return c.record.NumRows()
	}

	if c.bitmap != nil {
		i, err := cast.ToInt64E(c.bitmap.GetSizeInBytes())
		if err != nil {
			panic(err) // uff!
		}
		return i
	}

	return 0
}

func (c *CursorImpl) Bitmap() *roaring.Bitmap {
	return c.bitmap
}

func (c *CursorImpl) Record() array.Record {
	return c.record
}
