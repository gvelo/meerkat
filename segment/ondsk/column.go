package ondsk

import "github.com/RoaringBitmap/roaring"

type ColumnInt struct {
}

func (c *ColumnInt) Scan() *ColumnIterator {

}

func (c *ColumnInt) Filter(bitmap roaring.Bitmap) *ColumnIterator {

}

func (c ColumnInt) Close() error {

}

type ColumnIterator struct {
}

func (ci *ColumnIterator) Next(buf []int) (int, error) {

}

// TODO add idx lookup methods.
