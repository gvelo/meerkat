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

package storage

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"meerkat/internal/buffer"
	"meerkat/internal/schema"
	"meerkat/internal/storage/vector"
	"meerkat/internal/util/testutil"
	"path/filepath"
	"testing"
	"time"
)

const (
	// 16k vector x 3 plus random
	testLen = 1024*8*2*3 + 432
)

func TestSegmentWriter(t *testing.T) {

	indexInfo := createIndexInfo()
	buf := createBuffers(indexInfo)

	filePath := "/Users/sebad/meerkat/segments"

	//fmt.Println("path ", filePath)

	sid := uuid.New()

	sw := NewSegmentWriter(filePath, sid, buf)

	err := sw.Write()

	if err != nil {
		t.Fatal(err)
	}

	filePath = filepath.Join(filePath, sid.String())

	seg, err := ReadSegment(filePath)

	if err != nil {
		t.Fatal(err)
	}

	fid := make(map[string]schema.Field)

	for _, f := range indexInfo.Fields {
		fid[f.Id] = f
	}

	for id, col := range seg.columns {

		f, ok := fid[id]

		if !ok {
			t.Fatal("cannot find field")
		}

		b, ok := buf.Col(id)

		if !ok {
			t.Fatal("cannot find buffer")
		}

		testCol(t, f, col, b)

	}

}

func createBuffers(indexInfo schema.IndexInfo) *buffer.Table {

	now := int(time.Now().UnixNano())
	table := buffer.NewTable(indexInfo)

	for i := 0; i < testLen; i++ {
		r := buffer.NewRow(len(indexInfo.Fields))
		for _, f := range indexInfo.Fields {
			switch f.FieldType {
			case schema.FieldType_TIMESTAMP:
				now += rand.Intn(2000)
				r.AddCol(f.Id, now)
			case schema.FieldType_INT:
				if f.Nullable {
					if rand.Intn(3) == 2 {
						r.AddCol(f.Id, rand.Int())
					}
				} else {
					r.AddCol(f.Id, rand.Int())
				}
			case schema.FieldType_STRING:
				if f.Nullable {
					if rand.Intn(3) == 2 {
						r.AddCol(f.Id, testutil.RandomString(25))
					}
				} else {
					r.AddCol(f.Id, testutil.RandomString(25))
				}
			}

		}

		table.AppendRow(r)

	}

	return table

}

func createIndexInfo() schema.IndexInfo {

	return schema.IndexInfo{
		Id:             "test-index",
		Name:           "test-index",
		Desc:           "test-index",
		Created:        time.Time{},
		Updated:        time.Time{},
		PartitionAlloc: schema.PartitionAlloc{},
		Fields: []schema.Field{
			{
				Id:        "_ts",
				Name:      "_ts",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_TIMESTAMP,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "intFieldId",
				Name:      "intField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_INT,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "intNullable",
				Name:      "intNullable",
				Desc:      "intNullable",
				IndexId:   "test-index",
				FieldType: schema.FieldType_INT,
				Nullable:  true,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "stringFieldId",
				Name:      "stringField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_STRING,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "stringNullFieldId",
				Name:      "stringNullField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_STRING,
				Nullable:  true,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
		},
	}

}

func testCol(t *testing.T, field schema.Field, col interface{}, buf buffer.Buffer) {

	t.Log("Testing field ", field)

	switch field.FieldType {
	case schema.FieldType_TIMESTAMP:
		testIterINTField(t, field, col.(*intColumn), buf.(*buffer.IntBuffer))
		testReadINTField(t, field, col.(*intColumn), buf.(*buffer.IntBuffer))
	case schema.FieldType_INT:
		testIterINTField(t, field, col.(*intColumn), buf.(*buffer.IntBuffer))
		testReadINTField(t, field, col.(*intColumn), buf.(*buffer.IntBuffer))
	case schema.FieldType_STRING:
		testIterStringField(t, field, col.(*binaryColumn), buf.(*buffer.ByteSliceBuffer))
		testReadStringField(t, field, col.(*binaryColumn), buf.(*buffer.ByteSliceBuffer))
	default:
		t.Fatal("unknown column type")
	}

}

func testIterINTField(t *testing.T, f schema.Field, col *intColumn, buf *buffer.IntBuffer) {

	var values []int
	var nulls []bool
	var valids int

	iter := col.Iterator()

	for iter.HasNext() {

		v := iter.Next()

		values = append(values, v.Values()...)

		if f.Nullable {
			for n := 0; n < v.Len(); n++ {
				if v.IsValid(n) {
					valids++
					nulls = append(nulls, false)
				} else {
					nulls = append(nulls, true)
				}
			}
		}

	}

	assert.Equal(t, buf.Values(), values)

	if f.Nullable {
		assert.Equal(t, buf.Nulls(), nulls)
	}

}

func testReadINTField(t *testing.T, f schema.Field, col *intColumn, buf *buffer.IntBuffer) {

	v := vector.DefaultVectorPool().GetIntVector()
	l := v.Cap()

	var rids []uint32

	for i := 0; i < testLen && len(rids) < l; i++ {

		if rand.Intn(20) == 0 {
			rids = append(rids, uint32(i))
		}

	}

	r := col.Reader()

	vec := r.Read(rids)

	for i, rid := range rids {

		if f.Nullable {
			assert.Equal(t, !buf.Nulls()[rid], vec.IsValid(i))
			if !buf.Nulls()[rid] {
				assert.Equal(t, buf.Values()[rid], vec.Values()[i])
			}
		} else {
			assert.Equal(t, buf.Values()[rid], vec.Values()[i])
		}

	}

}

func testIterStringField(t *testing.T, f schema.Field, col *binaryColumn, buf *buffer.ByteSliceBuffer) {

	var values [][]byte
	var nulls []bool
	var valids int

	iter := col.Iterator()

	for iter.HasNext() {

		v := iter.Next()

		for i := 0; i < v.Len(); i++ {
			values = append(values, v.Get(i))
		}

		if f.Nullable {
			for n := 0; n < v.Len(); n++ {
				if v.IsValid(n) {
					valids++
					nulls = append(nulls, false)
				} else {
					nulls = append(nulls, true)
				}
			}
		}

	}

	assert.Equal(t, buf.Len(), len(values))

	for i := 0; i < buf.Len(); i++ {
		assert.Equal(t, buf.Get(i), values[i])
	}

	if f.Nullable {
		assert.Equal(t, buf.Nulls(), nulls)
	}

}

func testReadStringField(t *testing.T, f schema.Field, col *binaryColumn, buf *buffer.ByteSliceBuffer) {

	v := vector.DefaultVectorPool().GetByteSliceVector()
	l := v.Cap()

	var rids []uint32

	for i := 0; i < testLen && len(rids) < l; i++ {

		if rand.Intn(20) == 0 {
			rids = append(rids, uint32(i))
		}

	}

	r := col.Reader()

	vec := r.Read(rids)

	for i, rid := range rids {

		if f.Nullable {
			assert.Equal(t, !buf.Nulls()[rid], vec.IsValid(i))
			if !buf.Nulls()[rid] {
				assert.Equal(t, buf.Get(int(rid)), vec.Get(i))
			}
		} else {
			assert.Equal(t, buf.Get(int(rid)), vec.Get(i))
		}

	}

}
