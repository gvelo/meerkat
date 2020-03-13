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
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"meerkat/internal/buffer"
	"meerkat/internal/schema"
	"os"
	"path"
	"testing"
	"time"
)

const (
	testLen = 126433
)

func TestSegmentWriter_Write(t *testing.T) {

	index := schema.IndexInfo{
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
			//{
			//	Id:        "stringFieldId",
			//	Name:      "stringField",
			//	Desc:      "",
			//	IndexId:   "test-index",
			//	FieldType: schema.FieldType_STRING,
			//	Nullable:  false,
			//	Created:   time.Time{},
			//	Updated:   time.Time{},
			//},
		},
	}

	table := buffer.NewTable(index)

	for i := 0; i < testLen; i++ {
		r := buffer.NewRow(4)
		r.AddCol("_ts", int(time.Now().UnixNano()))
		//r.AddCol("_id", uuid.New())
		r.AddCol("intFieldId", rand.Int())
		if rand.Intn(3) == 2 {
			r.AddCol("intNullable", i)
		}
		///r.AddCol("stringFieldId", fmt.Sprintf("row number %v", i))
		table.AppendRow(r)
	}

	path := "/tmp/segment"

	sid := uuid.New()

	fmt.Println(sid)

	sw := NewSegmentWriter("/tmp/segment", sid, table)

	err := sw.Write()

	if err != nil {
		t.Error(err)
		return
	}

	seg, err := ReadSegment(path)

	if err != nil {
		t.Error(err)
		return
	}

	col := seg.columns["_ts"].(*intColumn)

	var buf []int

	iter := col.Iterator()

	for iter.HasNext() {
		v := iter.Next()
		buf = append(buf, v.Values()...)
	}

	assert.Len(t, buf, testLen, "wrong column len")

	colbuf, _ := table.Col("_ts")

	intbuf := colbuf.(*buffer.IntBuffer)

	assert.Equal(t, intbuf.Values(), buf, "values doesnt match")

	nulCol, _ := table.Col("intNullable")

	nulintcol := nulCol.(*buffer.IntBuffer)

	col = seg.columns["intNullable"].(*intColumn)
	iter = col.Iterator()
	buf = make([]int, 0)
	var nulls []bool

	for iter.HasNext() {
		fmt.Println("parace que hay next")
		v := iter.Next()
		buf = append(buf, v.Values()...)
		for n := 0; n < v.Len(); n++ {
			if v.IsValid(n) {
				nulls = append(nulls, false)
			} else {
				nulls = append(nulls, true)
			}
		}
	}

	fmt.Println(buf[:200])
	fmt.Println(nulintcol.Values()[:200])

}

func TestSegmentWriter(t *testing.T) {

	indexInfo := createIndexInfo()
	buf := createBuffers(indexInfo)

	filePath := path.Join(os.TempDir(), "segment_test")

	fmt.Println("path ", filePath)

	sid := uuid.New()

	sw := NewSegmentWriter(filePath, sid, buf)

	err := sw.Write()

	if err != nil {
		t.Fatal(err)
	}

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

	table := buffer.NewTable(indexInfo)

	for i := 0; i < testLen; i++ {
		r := buffer.NewRow(len(indexInfo.Fields))
		for _, f := range indexInfo.Fields {
			switch f.FieldType {
			case schema.FieldType_TIMESTAMP:
				r.AddCol(f.Id, int(time.Now().UnixNano()))
			case schema.FieldType_INT:
				if f.Nullable {
					if rand.Intn(3) == 2 {
						r.AddCol(f.Id, rand.Int())
					}
				} else {
					r.AddCol(f.Id, rand.Int())
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
			//{
			//	Id:        "stringFieldId",
			//	Name:      "stringField",
			//	Desc:      "",
			//	IndexId:   "test-index",
			//	FieldType: schema.FieldType_STRING,
			//	Nullable:  false,
			//	Created:   time.Time{},
			//	Updated:   time.Time{},
			//},
		},
	}

}

func testCol(t *testing.T, field schema.Field, col interface{}, buf buffer.Buffer) {

	fmt.Println("===========================")
	fmt.Println("testing col", field.Id)
	fmt.Println("===========================")
	switch field.FieldType {
	case schema.FieldType_TIMESTAMP:
		testINTField(t, field, col.(*intColumn), buf.(*buffer.IntBuffer))
	case schema.FieldType_INT:
		testINTField(t, field, col.(*intColumn), buf.(*buffer.IntBuffer))
	default:
		t.Fatal("unknown column type")
	}

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()

}

func testINTField(t *testing.T, f schema.Field, col *intColumn, buf *buffer.IntBuffer) {

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
