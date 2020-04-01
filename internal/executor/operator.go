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

// TODO: hacer un operador por tipo asi sacamos los switchs.

import (
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
)

// BinaryOperation represents an operation between two expressions
type BinaryOperation int

const (
	and BinaryOperation = iota
	or
	xor
)

// ComparisonOperation represents an comparison between two expressions
type ComparisonOperation int

const (
	eq ComparisonOperation = iota
	lt
	le
	gt
	ge
	ne
	isNull
	rex
	pref
)

// Operator represents an Operator of a physical plan execution.
// Operators take one or more inputs and produce an output in the form
// of vectors, bitmaps or some other type.
type Operator interface {
	// Init initializes the Operator acquiring the required resources.
	// Init will call the init method on all it's input operators.
	Init()

	// Destroy the Operator releasing all the acquired resources.
	// Destroy will cascade calling the Destroy method on all it's
	// children operators.
	Destroy()
}

// ColumnOperator represents an Operator which acts to a column
// GetName returns the column name.
type ColumnOperator interface {
	GetName() []byte
}

// BitmapOperator is an Operator that produces bitmaps as output
type BitmapOperator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() *roaring.Bitmap
}

// Uint32Operator is an Operator that []uint32 bitmaps as output
type Uint32Operator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() []uint32
}

// VectorOperator is an Operator that produces Vector as output
type VectorOperator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() vector.Vector
}

// MultiVectorOperator is an Operator that produces a Vector array as output
type MultiVectorOperator interface {
	Operator

	// Next returns the next result from this operator or nil
	// if there is no more data to process.
	// TODO(gvelo) should we destroy the operator automatically when
	//  there is no more data ?
	Next() []vector.Vector
}

// NewBinaryBitmapOperator creates a new bitmap binary operator.
func NewBinaryBitmapOperator(ctx Context, op BinaryOperation, left BitmapOperator, right BitmapOperator) *BinaryBitmapOperator {
	return &BinaryBitmapOperator{
		ctx:   ctx,
		op:    op,
		left:  left,
		right: right,
	}
}

// BinaryBitmapOperator executes a binary operation between two bitmaps
// and returns a new bitmap.
type BinaryBitmapOperator struct {
	ctx   Context
	op    BinaryOperation
	left  BitmapOperator
	right BitmapOperator
}

func (op *BinaryBitmapOperator) Init() {
	op.left.Init()
	op.right.Init()
}

func (op *BinaryBitmapOperator) Destroy() {
	op.left.Destroy()
	op.right.Destroy()
}

func (op *BinaryBitmapOperator) Next() *roaring.Bitmap {

	// parallelize
	l := op.left.Next()
	r := op.left.Next()

	switch op.op {
	case and:
		return roaring.And(l, r)
	case or:
		return roaring.Or(l, r)
	case xor:
		return roaring.Xor(l, r)
	}
	panic("Operator not supported")
}

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

// NewColumnScanOperator creates a ColumnScanOperator
func NewIntColumnScanOperator(ctx Context, op ComparisonOperation, value int, fieldName string, size int) Uint32Operator {
	var v func(x, y int) bool
	switch op {
	case eq:
		v = func(x, y int) bool {
			return x == y
		}
	case gt:
		v = func(x, y int) bool {
			return x > y
		}
	case ge:
		v = func(x, y int) bool {
			return x >= y
		}
	case le:
		v = func(x, y int) bool {
			return x <= y
		}
	case lt:
		v = func(x, y int) bool {
			return x < y
		}
	case ne:
		v = func(x, y int) bool {
			return x != y
		}
	case isNull:
		v = nil
	}

	return &IntColumnScanOperator{
		ctx:   ctx,
		opFn:  v,
		value: value,
		fn:    fieldName,
		sz:    size,
	}
}

type IntColumnScanOperator struct {
	ctx      Context
	opFn     func(x, y int) bool
	fn       string
	value    int
	sz       int
	iterator storage.IntIterator
	scanLeft []int
	lastRid  uint32
}

func (op *IntColumnScanOperator) Init() {
	c := op.ctx.Segment().Col(op.fn).(storage.IntColumn)
	op.iterator = c.Iterator()
}

func (op *IntColumnScanOperator) Destroy() {
}

func (op *IntColumnScanOperator) processVector(src []int, x int) []uint32 {

	r := make([]uint32, 0, op.sz)

	for ; x < len(src) && len(r) < op.sz; x++ {
		if op.opFn(src[x], op.value) {
			r = append(r, op.lastRid)
		}
		op.lastRid++
	}

	if len(r) == op.sz {
		op.scanLeft = src[x:]
	}

	return r
}

func (op *IntColumnScanOperator) Next() []uint32 {

	r := make([]uint32, 0, op.sz)

	if len(op.scanLeft) > 0 {

		x := 0
		for ; x < len(op.scanLeft) && len(r) < op.sz; x++ {
			if op.opFn(op.scanLeft[x], op.value) {
				r = append(r, op.lastRid)
			}
			op.lastRid++
		}

		if len(r) == op.sz {
			op.scanLeft = op.scanLeft[x:]
			return r
		}
	}

	for op.iterator.HasNext() {

		intVector := op.iterator.Next()
		values := intVector.Values()

		x := 0
		for ; x < len(values) && len(r) < op.sz; x++ {
			if op.opFn(values[x], op.value) {
				r = append(r, op.lastRid)
			}
			op.lastRid++
		}

		if len(r) == op.sz {
			op.scanLeft = values[x:]
			return r
		}

	}

	op.scanLeft = nil
	if len(r) > 0 {
		return r
	} else {
		return nil
	}

}

// NewColumnScanOperator creates a ColumnScanOperator
func NewLimitOperator(ctx Context, child MultiVectorOperator, limit int) MultiVectorOperator {
	return &LimitOperator{
		ctx:   ctx,
		child: child,
		limit: limit,
	}
}

// LimitOperator scans a non indexed column, search for the []pos
// registers and returns the bitmap that meets that condition.
//
type LimitOperator struct {
	ctx   Context
	child MultiVectorOperator
	limit int
}

func (op *LimitOperator) Init() {
	op.child.Init()
}

func (op *LimitOperator) Destroy() {
	op.child.Destroy()
}

func (op *LimitOperator) Next() []vector.Vector {
	n := op.child.Next()

	return n
}

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

	buff := make([]uint32, 0, 1000)
	s := r.it.NextMany(buff)
	if s != 0 {
		c := r.ctx.Segment().Col(r.colName)
		v := c.(storage.IntColumn).Reader().Read(buff) // Check error? TODO(sebad): hacer un operator por tipo
		return &v
	} else {
		return nil
	}
}

func NewBufferOperator(ctx Context, children []VectorOperator) MultiVectorOperator {
	return &BufferOperator{
		ctx:      ctx,
		children: children,
	}
}

// BufferOperator reads all positions in the bitmap
// and other and
type BufferOperator struct {
	ctx      Context
	vKeys    [][]byte
	children []VectorOperator
}

func (r *BufferOperator) Init() {
	r.vKeys = make([][]byte, 0)
	for i, c := range r.children {
		r.vKeys[i] = c.(ColumnOperator).GetName()
		c.Init()
	}
	r.ctx.Value(ColumnIndexToColumnName, r.vKeys)
}

func (r *BufferOperator) Destroy() {
	for _, c := range r.children {
		c.Destroy()
	}
}

func (r *BufferOperator) Next() []vector.Vector {
	op := make([]vector.Vector, 0, len(r.children))
	for i, c := range r.children {
		// Paralelize
		op[i] = c.Next()

	}
	return op
}

// Decodifica el diccionario y lo manda....
func NewMaterialize(ctx Context, child MultiVectorOperator) MaterializeOperator {
	return MaterializeOperator{
		ctx:   ctx,
		child: child,
	}
}

// MaterializeOperator operator
type MaterializeOperator struct {
	ctx   Context
	child MultiVectorOperator
}

func (r *MaterializeOperator) Init() {
	r.child.Init()
}

func (r *MaterializeOperator) Destroy() {
	r.child.Destroy()
}

func (r *MaterializeOperator) Next() []vector.Vector {
	n := r.child.Next()
	var keys [][]byte
	v, ok := r.ctx.Get(ColumnIndexToColumnName)
	if ok {
		keys = v.([][]byte)
	} else {
		panic("No ColumnIndexToColumnName")
	}

	// Aca tengo los valores de los objetos... que son null o no segun el vector lo tengo que validar
	if n != nil {
		res := make([]vector.Vector, 0, len(n))
		for i, _ := range n {

			switch vec := n[i].(type) {

			case *vector.ByteSliceVector: // No estoy seguro que pueda caer a esta algura.
				res[i] = vec
			case *vector.FloatVector:
				res[i] = vec
			case *vector.IntVector:
				fName := keys[i]
				col := r.ctx.Segment().Col(string(fName))
				_, ok := col.(storage.StringColumn) // c
				if ok {                             // Its a string col, we should dict decode here.
					// Here we sould create a vector
					// b, err := c.Dict().DecodeByteSlice()
					res[i] = vec
				} else {
					res[i] = vec
				}
			}
		}
		return res
	}
	return nil
}
