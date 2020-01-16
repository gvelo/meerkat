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

package buffer

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/google/uuid"
	"log"
	"meerkat/internal/schema"
	"reflect"
	"unsafe"
)

//TODO(gvelo):generate

const (
	// Int64SizeBytes specifies the number of bytes required to store a single int64 in memory
	Int64SizeBytes = int(unsafe.Sizeof(int64(0)))
)

func CastToBytes(size int, p unsafe.Pointer) []byte {
	h := (*reflect.SliceHeader)(p)
	var res []byte
	s := (*reflect.SliceHeader)(unsafe.Pointer(&res))
	s.Data = h.Data
	s.Len = h.Len * size
	s.Cap = h.Cap * size
	return res
}

type Buffer interface {
	Len() int
	Size() int
	Nulls() *roaring.Bitmap
	AppendNull()
	Append(interface{})
	AppendBuffer(interface{})
	Nullable() bool
	Bytes() []byte
}

type IntBuffer struct {
	nulls *roaring.Bitmap
	buf   []int
}

func NewIntBuffer(nullable bool, capacity int) *IntBuffer {
	b := &IntBuffer{
		buf: make([]int, 0, capacity),
	}
	if nullable {
		b.nulls = roaring.NewBitmap()
	}
	return b
}

func (b *IntBuffer) Len() int {
	return len(b.buf)
}

func (b *IntBuffer) Size() int {
	return len(b.buf) * Int64SizeBytes
}

func (b *IntBuffer) Nulls() *roaring.Bitmap {
	return b.nulls
}

func (b *IntBuffer) AppendNull() {
	b.nulls.AddInt(len(b.buf))
	b.buf = append(b.buf, 0)
}

func (b *IntBuffer) Nullable() bool {
	return b.nulls != nil
}

func (b *IntBuffer) Bytes() []byte {
	return CastToBytes(Int64SizeBytes, unsafe.Pointer(&b.buf))
}

func (b *IntBuffer) Int() []int {
	return b.buf
}

func (b *IntBuffer) AppendInt(i int) {
	b.buf = append(b.buf, i)
}

func (b *IntBuffer) Append(i interface{}) {
	b.AppendInt(i.(int))
}

func (b *IntBuffer) AppendIntBuffer(s *IntBuffer) {
	b.buf = append(b.buf, s.buf...)
	if s.nulls != nil && b.nulls != nil {
		bshift := roaring.AddOffset(s.nulls, uint32(len(b.buf)))
		b.nulls.Or(bshift)
	}
}

func (b *IntBuffer) AppendBuffer(buf interface{}) {
	b.AppendIntBuffer(buf.(*IntBuffer))
}

// UINT

type UintBuffer struct {
	nulls *roaring.Bitmap
	buf   []uint
}

func NewUintBuffer(nullable bool, capacity int) *UintBuffer {
	b := &UintBuffer{
		buf: make([]uint, 0, capacity),
	}
	if nullable {
		b.nulls = roaring.NewBitmap()
	}
	return b
}

func (b *UintBuffer) Len() int {
	return len(b.buf)
}

func (b *UintBuffer) Size() int {
	return len(b.buf) * Int64SizeBytes
}

func (b *UintBuffer) Nulls() *roaring.Bitmap {
	return b.nulls
}

func (b *UintBuffer) AppendNull() {
	b.nulls.AddInt(len(b.buf))
	b.buf = append(b.buf, 0)
}

func (b *UintBuffer) Nullable() bool {
	return b.nulls != nil
}

func (b *UintBuffer) Bytes() []byte {
	return CastToBytes(Int64SizeBytes, unsafe.Pointer(&b.buf))
}

func (b *UintBuffer) AppendUint(i uint) {
	b.buf = append(b.buf, i)
}

func (b *UintBuffer) Append(i interface{}) {
	b.AppendUint(i.(uint))
}

func (b *UintBuffer) AppendUintBuffer(s *UintBuffer) {
	b.buf = append(b.buf, s.buf...)
	if s.nulls != nil && b.nulls != nil {
		bshift := roaring.AddOffset(s.nulls, uint32(len(b.buf)))
		b.nulls.Or(bshift)
	}
}

func (b *UintBuffer) AppendBuffer(buf interface{}) {
	b.AppendUintBuffer(buf.(*UintBuffer))
}

//Float

type Float64Buffer struct {
	nulls *roaring.Bitmap
	buf   []float64
}

func NewFloat64Buffer(nullable bool, capacity int) *Float64Buffer {
	b := &Float64Buffer{
		buf: make([]float64, 0, capacity),
	}
	if nullable {
		b.nulls = roaring.NewBitmap()
	}
	return b
}

func (b *Float64Buffer) Len() int {
	return len(b.buf)
}

func (b *Float64Buffer) Size() int {
	return len(b.buf) * Int64SizeBytes
}

func (b *Float64Buffer) Nulls() *roaring.Bitmap {
	return b.nulls
}

func (b *Float64Buffer) AppendNull() {
	b.nulls.AddInt(len(b.buf))
	b.buf = append(b.buf, 0)
}

func (b *Float64Buffer) Nullable() bool {
	return b.nulls != nil
}

func (b *Float64Buffer) Bytes() []byte {
	return CastToBytes(Int64SizeBytes, unsafe.Pointer(&b.buf))
}

func (b *Float64Buffer) AppendFloat(f float64) {
	b.buf = append(b.buf, f)
}

func (b *Float64Buffer) Append(i interface{}) {
	b.AppendFloat(i.(float64))
}

func (b *Float64Buffer) AppendFloat64Buffer(s *Float64Buffer) {
	b.buf = append(b.buf, s.buf...)
	if s.nulls != nil && b.nulls != nil {
		bshift := roaring.AddOffset(s.nulls, uint32(len(b.buf)))
		b.nulls.Or(bshift)
	}
}

func (b *Float64Buffer) AppendBuffer(buf interface{}) {
	b.AppendFloat64Buffer(buf.(*Float64Buffer))
}

// slice

type SliceBuffer struct {
	nulls   *roaring.Bitmap
	buf     []byte
	offsets []int
}

func NewSliceBuffer(nullable bool, capacity int) *SliceBuffer {
	b := &SliceBuffer{
		buf: make([]byte, 0, capacity),
	}
	if nullable {
		b.nulls = roaring.NewBitmap()
	}
	return b
}

func (b *SliceBuffer) Len() int {
	return len(b.offsets)
}

func (b *SliceBuffer) Size() int {
	return len(b.buf) + len(b.offsets)*Int64SizeBytes
}

func (b *SliceBuffer) Nulls() *roaring.Bitmap {
	return b.nulls
}

func (b *SliceBuffer) AppendNull() {
	b.nulls.AddInt(len(b.offsets))
	b.offsets = append(b.offsets, 0)
}

func (b *SliceBuffer) Nullable() bool {
	return b.nulls != nil
}

func (b *SliceBuffer) Bytes() []byte {
	return b.buf
}

func (b *SliceBuffer) OffsetBytes() []byte {
	return CastToBytes(Int64SizeBytes, unsafe.Pointer(&b.offsets))
}

func (b *SliceBuffer) AppendSlice(s []byte) {
	b.offsets = append(b.offsets, len(b.buf))
	b.buf = append(b.buf, s...)
}

func (b *SliceBuffer) AppendString(s string) {
	b.offsets = append(b.offsets, len(b.buf))
	b.buf = append(b.buf, s...)
}

func (b *SliceBuffer) Append(i interface{}) {
	b.AppendString(i.(string))
}

func (b *SliceBuffer) AppendSliceBuffer(s *SliceBuffer) {
	l := len(b.offsets)
	o := len(b.buf)
	b.buf = append(b.buf, s.buf...)
	b.offsets = append(b.offsets, s.offsets...)
	// shift the offsets
	for i := l; i < len(b.offsets); i++ {
		b.offsets[i] = b.offsets[i] + o
	}
	if s.nulls != nil && b.nulls != nil {
		bshift := roaring.AddOffset(s.nulls, uint32(len(b.buf)))
		b.nulls.Or(bshift)
	}
}

func (b *SliceBuffer) AppendBuffer(buf interface{}) {
	b.AppendSliceBuffer(buf.(*SliceBuffer))
}

func (b *SliceBuffer) Each(f func(int, []byte) bool) {

	for i, start := range b.offsets {

		var end int

		if i == (len(b.offsets) - 1) {
			end = len(b.buf)
		} else {
			end = b.offsets[i+1]
		}

		if !f(i, b.buf[start:end]) {
			return
		}

	}

}

// UUID

type UUIDBuffer struct {
	nulls    *roaring.Bitmap
	buf      []byte
	nullUUID []byte
}

func NewUUIDBuffer(nullable bool, capacity int) *UUIDBuffer {
	b := &UUIDBuffer{
		buf: make([]byte, 0, capacity),
	}
	if nullable {
		b.nulls = roaring.NewBitmap()
		b.nullUUID = make([]byte, 16)
	}
	return b
}

func (b *UUIDBuffer) Len() int {
	return len(b.buf) >> 4
}

func (b *UUIDBuffer) Size() int {
	return len(b.buf)
}

func (b *UUIDBuffer) Nulls() *roaring.Bitmap {
	return b.nulls
}

func (b *UUIDBuffer) AppendNull() {
	b.nulls.AddInt(b.Len())
	b.buf = append(b.buf, b.nullUUID...)
}

func (b *UUIDBuffer) Nullable() bool {
	return b.nulls != nil
}

func (b *UUIDBuffer) Bytes() []byte {
	return b.buf
}

func (b *UUIDBuffer) AppendUUID(uuid uuid.UUID) {
	bytes := [16]byte(uuid)
	b.buf = append(b.buf, bytes[:]...)
}

func (b *UUIDBuffer) Append(i interface{}) {
	b.AppendUUID(i.(uuid.UUID))
}

func (b *UUIDBuffer) AppendUUIDBuffer(s *UUIDBuffer) {

	b.buf = append(b.buf, s.buf...)

	if s.nulls != nil && b.nulls != nil {
		bshift := roaring.AddOffset(s.nulls, uint32(b.Len()))
		b.nulls.Or(bshift)
	}

}

func (b *UUIDBuffer) AppendBuffer(buf interface{}) {
	b.AppendUUIDBuffer(buf.(*UUIDBuffer))
}

func (b *UUIDBuffer) Each(f func(int, uuid.UUID) bool) {

	for i := 0; i < b.Len(); i++ {
		start := i << 4
		end := start + 16
		var uid [16]byte
		copy(uid[:], b.buf[start:end])
		if !f(i, uid) {
			return
		}
	}

}

// bool

type BoolBuffer struct {
	nulls *roaring.Bitmap
	buf   *roaring.Bitmap
	len   int
}

func (b BoolBuffer) Len() int {
	return b.len
}

func (b BoolBuffer) Size() int {
	return int(b.buf.GetSizeInBytes())
}

func (b BoolBuffer) Nulls() *roaring.Bitmap {
	return b.nulls
}

func (b BoolBuffer) AppendNull() {
	b.nulls.AddInt(b.len)
	b.len++
}

func (b BoolBuffer) Nullable() bool {
	return b.nulls != nil
}

func (b BoolBuffer) Bytes() []byte {
	bytes, err := b.buf.ToBytes()
	if err != nil {
		panic(err)
	}
	return bytes
}

func (b *BoolBuffer) AppendBool(v bool) {
	if v {
		b.buf.AddInt(b.len)
	}
	b.len++
}

func (b *BoolBuffer) Append(i interface{}) {
	b.AppendBool(i.(bool))
}

func (b *BoolBuffer) Values() *roaring.Bitmap {
	return b.buf
}

func (b *BoolBuffer) AppendBoolBuffer(s *BoolBuffer) {

	vshift := roaring.AddOffset(s.nulls, uint32(b.Len()))
	b.nulls.Or(vshift)

	if s.nulls != nil && b.nulls != nil {
		nshift := roaring.AddOffset(s.nulls, uint32(b.Len()))
		b.nulls.Or(nshift)
	}

}

func (b *BoolBuffer) AppendBuffer(buf interface{}) {
	b.AppendBoolBuffer(buf.(*BoolBuffer))
}

type Row struct {
	cols     map[string]interface{}
	colCount int
}

func NewRow(colCount int) *Row {
	return &Row{
		cols:     make(map[string]interface{}, colCount),
		colCount: colCount,
	}
}

func (r *Row) AddCol(id string, value interface{}) {
	r.cols[id] = value
}

func (r *Row) GetCol(id string) (interface{}, bool) {
	v, found := r.cols[id]
	return v, found
}

func (r *Row) Reset() {
	r.cols = make(map[string]interface{}, r.colCount)
}

type Table struct {
	cols  map[string]Buffer
	index schema.IndexInfo
	len   int
}

func (t *Table) AppendRow(r *Row) {

	t.len++

	for _, f := range t.index.Fields {

		v, ok := r.GetCol(f.Id)

		if !ok {
			t.cols[f.Id].AppendNull()
			continue
		}

		t.cols[f.Id].Append(v)

	}
}

func (t *Table) AppendTable(b *Table) {

	for id, col := range t.cols {

		src, ok := b.Col(id)

		if !ok {
			// TODO(gvelo): handle schema mutation properly.
			panic("column not found")
		}

		col.AppendBuffer(src)

	}

	t.len = t.len + b.len

}

func (t *Table) Col(id string) (Buffer, bool) {
	b, ok := t.cols[id]
	return b, ok
}

func (t *Table) Cols() map[string]Buffer {
	return t.cols
}

func (t *Table) Len() int {
	return t.len
}

func (t *Table) Index() schema.IndexInfo {
	return t.index
}

func NewTable(idx schema.IndexInfo) *Table {

	t := &Table{
		cols:  make(map[string]Buffer, len(idx.Fields)),
		index: idx,
	}

	for _, f := range idx.Fields {
		switch f.FieldType {
		case schema.FieldType_FLOAT:
			t.cols[f.Id] = NewFloat64Buffer(f.Nullable, 0)
		case schema.FieldType_INT:
			t.cols[f.Id] = NewIntBuffer(f.Nullable, 0)
		case schema.FieldType_STRING:
			t.cols[f.Id] = NewSliceBuffer(f.Nullable, 0)
		case schema.FieldType_TEXT:
			t.cols[f.Id] = NewSliceBuffer(f.Nullable, 0)
		case schema.FieldType_TIMESTAMP:
			t.cols[f.Id] = NewIntBuffer(f.Nullable, 0)
		case schema.FieldType_BOOLEAN:
			t.cols[f.Id] = NewIntBuffer(f.Nullable, 0)
		case schema.FieldType_UINT:
			t.cols[f.Id] = NewUintBuffer(f.Nullable, 0)
		case schema.FieldType_UUID:
			t.cols[f.Id] = NewUUIDBuffer(f.Nullable, 0)
		default:
			log.Panicf("invalid field type %v", f.FieldType)
		}
	}

	return t

}
