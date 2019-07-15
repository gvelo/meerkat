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

package ondsk

import (
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage/segment/inmem"
)

type Column interface {
	Scan() ColumnIterator
	SetFilter(bitmap *roaring.Bitmap)
	Close() error
}

type ColumnImpl struct {
	IdxPag            *OnDiskColumnIdx
	Pag               *OnDiskColumn
	actualPage        *inmem.Page
	pagesFiltered     []*inmem.Page
	pageFilteredIndex int
	err               error
}

func (c *ColumnImpl) HasNext() bool {
	if c.pagesFiltered == nil {
		if c.actualPage == nil {
			c.actualPage, c.err = c.Pag.Br.ReadPageHeader()
		} else {
			n := c.Pag.Br.Offset + c.actualPage.PayloadSize
			if c.Pag.Br.Size > n {
				c.Pag.Br.Offset = c.Pag.Br.Offset + c.actualPage.PayloadSize
				c.actualPage, c.err = c.Pag.Br.ReadPageHeader()
			} else {
				return false
			}
		}
	} else {
		if c.pageFilteredIndex < len(c.pagesFiltered) {
			c.actualPage = c.pagesFiltered[c.pageFilteredIndex]
			c.pageFilteredIndex++
		} else {
			return false
		}

	}
	return c.err == nil
}

func (c *ColumnImpl) Next() *inmem.Page {
	return c.actualPage
}

func (c *ColumnImpl) Scan() ColumnIterator {
	return c
}

func (c *ColumnImpl) SetFilter(bitmap *roaring.Bitmap) {

	c.pageFilteredIndex = 0
	c.pagesFiltered = make([]*inmem.Page, 0)

	var lastPage *inmem.Page

	it := bitmap.Iterator()
	for it.HasNext() {
		i := it.Next()
		if lastPage == nil {
			lastPage = processPage(c, i)
		} else {
			if int(i) > lastPage.StartID+lastPage.Total {
				lastPage = processPage(c, i)
			}
		}
	}
	return
}

func processPage(c *ColumnImpl, i uint32) *inmem.Page {
	offset, _ := c.IdxPag.Lookup(int(i))
	c.Pag.Br.Offset = int(offset)
	p, err := c.Pag.Br.ReadPageHeader()
	if err != nil {
		panic(err)
	}
	c.pagesFiltered = append(c.pagesFiltered, p)
	return p
}

func (c *ColumnImpl) Close() error {
	return nil
}

type ColumnIterator interface {
	Next() *inmem.Page
	HasNext() bool
}
