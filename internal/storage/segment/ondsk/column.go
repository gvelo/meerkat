package ondsk

import (
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage/segment/inmem"
)

type ColumnInt struct {
	value  *OnDiskSkipList
	posIdx *OnDiskColumn
}

func (c *ColumnInt) Scan() *ColumnIterator {
	return nil
}

func (c *ColumnInt) Filter(bitmap roaring.Bitmap) *ColumnIterator {
	// levantar el indice y ahsta que se me acaben los valores del bipmap
	// hacer el par devolver la pagina sin descomprimir.
	return nil
}

func (c ColumnInt) Close() error {
	return nil
}

type ColumnIterator struct {
}

func (ci *ColumnIterator) Next() *inmem.Page {
	return nil
}

func (ci *ColumnIterator) HasNext() bool {
	return true
}

// TODO add idx lookup methods.
