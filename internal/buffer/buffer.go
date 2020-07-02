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

//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@buffer.gen.go.tmpldata buffer.gen.go.tmpl

package buffer

type Buffer interface {
	Len() int
	Size() int
	Nulls() []bool
	AppendNull()
	Append(interface{})
	AppendBuffer(interface{})
	Nullable() bool
}

type ByteSliceBuffer struct {
	nulls    []bool
	buf      []byte
	offsets  []int
	nullable bool
}

func NewByteSliceBuffer(nullable bool, capacity int) *ByteSliceBuffer {

	b := &ByteSliceBuffer{
		buf:      make([]byte, 0, capacity),
		nullable: nullable,
	}

	if nullable {
		b.nulls = make([]bool, 0, capacity)
	}

	return b

}

func (b *ByteSliceBuffer) Len() int {
	return len(b.offsets)
}

func (b *ByteSliceBuffer) Size() int {
	return (len(b.buf) + len(b.offsets)) * 8
}

func (b *ByteSliceBuffer) Nulls() []bool {
	return b.nulls
}

func (b *ByteSliceBuffer) AppendNull() {
	b.nulls = append(b.nulls, true)
	b.offsets = append(b.offsets, len(b.buf))
}

func (b *ByteSliceBuffer) Nullable() bool {
	return b.nulls != nil
}

func (b *ByteSliceBuffer) AppendSlice(s []byte) {

	b.buf = append(b.buf, s...)
	b.offsets = append(b.offsets, len(b.buf))

	if b.nullable {
		b.nulls = append(b.nulls, false)
	}

}

func (b *ByteSliceBuffer) AppendString(s string) {

	b.buf = append(b.buf, s...)
	b.offsets = append(b.offsets, len(b.buf))

	if b.nullable {
		b.nulls = append(b.nulls, false)
	}

}

func (b *ByteSliceBuffer) Append(i interface{}) {
	b.AppendString(i.(string))
}

func (b *ByteSliceBuffer) AppendSliceBuffer(s *ByteSliceBuffer) {

	if b.nullable != s.nullable {
		panic("schema mutation on ingestion not supported yet")
	}

	l := len(b.offsets)
	o := len(b.buf)

	b.buf = append(b.buf, s.buf...)
	b.offsets = append(b.offsets, s.offsets...)

	// shift the offsets

	for i := l; i < len(b.offsets); i++ {
		b.offsets[i] = b.offsets[i] + o
	}

	if b.nullable {
		b.nulls = append(b.nulls, s.nulls...)
	}

}

func (b *ByteSliceBuffer) AppendBuffer(buf interface{}) {
	b.AppendSliceBuffer(buf.(*ByteSliceBuffer))
}

func (b *ByteSliceBuffer) Get(i int) []byte {

	var start int

	if i > 0 {
		start = b.offsets[i-1]
	}

	return b.buf[start:b.offsets[i]]

}

func (b *ByteSliceBuffer) Each(f func(int, []byte) bool) {

	for i, end := range b.offsets {

		var start int

		if i == 0 {
			start = 0
		} else {
			start = b.offsets[i-1]
		}

		if !f(i, b.buf[start:end]) {
			return
		}

	}

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
	cols map[string]Buffer
	//index schema.IndexInfo
	len int
}

//func (t *Table) AppendRow(r *Row) {
//
//	t.len++
//
//	for _, f := range t.index.Fields {
//
//		v, ok := r.GetCol(f.Id)
//
//		if !ok {
//			t.cols[f.Id].AppendNull()
//			continue
//		}
//
//		t.cols[f.Id].Append(v)
//
//	}
//}

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

//func (t *Table) Index() schema.IndexInfo {
//	return t.index
//}

//func NewTable(idx schema.IndexInfo) *Table {
//
//	t := &Table{
//		cols:  make(map[string]Buffer, len(idx.Fields)),
//		index: idx,
//	}
//
//	for _, f := range idx.Fields {
//		switch f.FieldType {
//		case schema.FieldType_FLOAT:
//			t.cols[f.Id] = NewFloatBuffer(f.Nullable, 0)
//		case schema.FieldType_INT:
//			t.cols[f.Id] = NewIntBuffer(f.Nullable, 0)
//		case schema.FieldType_STRING:
//			t.cols[f.Id] = NewByteSliceBuffer(f.Nullable, 0)
//		case schema.FieldType_TEXT:
//			t.cols[f.Id] = NewByteSliceBuffer(f.Nullable, 0)
//		case schema.FieldType_TIMESTAMP:
//			t.cols[f.Id] = NewIntBuffer(f.Nullable, 0)
//		case schema.FieldType_BOOLEAN:
//			t.cols[f.Id] = NewBoolBuffer(f.Nullable, 0)
//		case schema.FieldType_UINT:
//			t.cols[f.Id] = NewUintBuffer(f.Nullable, 0)
//		default:
//			log.Panicf("invalid field type %v", f.FieldType)
//		}
//	}
//
//	return t
//
//}
