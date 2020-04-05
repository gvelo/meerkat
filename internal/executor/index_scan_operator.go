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
)

// NewBinaryBitmapOperator creates a new bitmap binary operator.
func NewIndexScanOperator(ctx Context, op ComparisonOperation, value interface{}, fieldName string) BitmapOperator {
	return &IndexScanOperator{
		ctx:   ctx,
		op:    op,
		fn:    fieldName,
		value: value,
	}
}

// IndexScanOperator executes a search in a column and returns the bitmap of positions
// that meet that condition.
type IndexScanOperator struct {
	ctx   Context
	op    ComparisonOperation
	fn    string
	value interface{}
}

func (op *IndexScanOperator) Init() {
	// Nothing to do yet.
}

func (op *IndexScanOperator) Destroy() {
	// Nothing to do yet.
}

func (op *IndexScanOperator) Next() *roaring.Bitmap {

	c := op.ctx.Segment().Col(op.fn)
	switch col := c.(type) {
	case storage.StringColumn:
		switch op.op {
		case eq:
			return col.Index().Search(op.value.([]byte))
		case rex:
			return col.Index().Regex(op.value.([]byte))
		case pref:
			return col.Index().Prefix(op.value.([]byte))
		}
		panic("Operator not supported")

	case storage.TextColumn:
		switch op.op {
		case eq:
			return col.Index().Search(op.value.([]byte))
		case rex:
			return col.Index().Regex(op.value.([]byte))
		case pref:
			return col.Index().Prefix(op.value.([]byte))
		}
		panic("Operator not supported")

	case storage.IntColumn:
		switch op.op {
		case ne:
			return col.Index().Ne(op.value.(int))
		case eq:
			return col.Index().Eq(op.value.(int))
		case lt:
			return col.Index().Lt(op.value.(int))
		case le:
			return col.Index().Le(op.value.(int))
		case ge:
			return col.Index().Ge(op.value.(int))
		case gt:
			return col.Index().Gt(op.value.(int))
		}
		panic("Operator not supported")

	case storage.FloatColumn:
		switch op.op {
		case ne:
			return col.Index().Ne(op.value.(float64))
		case eq:
			return col.Index().Eq(op.value.(float64))
		case lt:
			return col.Index().Lt(op.value.(float64))
		case le:
			return col.Index().Le(op.value.(float64))
		case ge:
			return col.Index().Ge(op.value.(float64))
		case gt:
			return col.Index().Gt(op.value.(float64))
		}
		panic("Operator not supported")
	}

	panic("Column type does not exists")
}
