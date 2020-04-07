// Copyright 2020 The Meerkat Authors
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

package executor

import (
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
)

func NewReaderOperator(ctx Context, child BitmapOperator, colName string) VectorOperator {
	return &ReaderOperator{
		ctx:     ctx,
		child:   child,
		colName: colName,
	}
}

// ReaderOperator reads all positions in the bitmap
type ReaderOperator struct {
	ctx     Context
	child   BitmapOperator
	colName string
	it      roaring.ManyIntIterable
	sz      int
}

func (r *ReaderOperator) Init() {
	r.child.Init()
	n := r.child.Next()
	if n != nil {
		r.it = n.ManyIterator()
	}
}

func (r *ReaderOperator) GetName() string {
	return r.colName
}

func (r *ReaderOperator) Destroy() {
	r.child.Destroy()
}

func (r *ReaderOperator) Next() vector.Vector {

	if r.it == nil {
		return nil
	}

	buff := make([]uint32, 0, r.sz)
	s := r.it.NextMany(buff)
	if s != 0 {
		c := r.ctx.Segment().Col(r.colName)
		v := c.(storage.IntColumn).Reader().Read(buff) // Check error? TODO(sebad): hacer un operator por tipo
		return &v
	} else {
		return nil
	}
}
