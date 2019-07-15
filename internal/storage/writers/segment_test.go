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

package writers

import (
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"meerkat/internal/storage/readers"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"testing"
	"time"
)

func TestSegmentWriterReader(t *testing.T) {

	assert := assert.New(t)

	indexInfo := segment.NewIndexInfo("test-index")
	indexInfo.AddField("testfield", segment.FieldTypeKeyword, true)
	indexInfo.AddField("mun1", segment.FieldTypeInt, true)
	indexInfo.AddField("float", segment.FieldTypeFloat, true)

	s := inmem.NewSegment(indexInfo, "3dfa542d", nil)

	for i := 0; i < 100; i++ {
		e := make(segment.Event)
		e["testfield"] = fmt.Sprintf("test %v", i)
		e["ts"] = uint64(time.Now().Add(time.Duration(i)).Nanosecond())
		e["mun1"] = i
		e["float"] = float64(i)
		s.Add(e)
	}

	sw := NewSegmentWriter("/tmp", s)

	err := sw.Write()

	if !assert.NoErrorf(err, "an error occurred while writing the segment: %v", err) {
		return
	}

	odSegment, err := readers.ReadSegment("/tmp")

	if !assert.NoErrorf(err, "an error occurred while reading the segment: %v", err) {
		return
	}

	assert.Equal(indexInfo.Name, odSegment.IndexInfo.Name, "wrong index name")

	odFields := odSegment.IndexInfo.Fields

	for i, field := range indexInfo.Fields {
		assert.Equal(field.Name, odFields[i].Name, "field name doesn't match")
		assert.Equal(field.Type, odFields[i].Type, "field type doesn't match")
	}

}

func findColumnByName(columns []inmem.Column, name string) inmem.Column {
	for _, c := range columns {
		if c.FieldInfo().Name == name {
			return c
		}
	}
	return nil
}

func TestSegmentSorted(t *testing.T) {

	assert := assert.New(t)

	indexInfo := segment.NewIndexInfo("test-index")
	indexInfo.AddField("testfield", segment.FieldTypeKeyword, true)
	indexInfo.AddField("mun1", segment.FieldTypeInt, true)
	indexInfo.AddField("float", segment.FieldTypeFloat, true)

	s := inmem.NewSegment(indexInfo, "3dfa542d", nil)

	for i := 0; i < 100; i++ {
		e := make(segment.Event)
		e["testfield"] = fmt.Sprintf("test %v", i)
		e["ts"] = uint64(time.Now().Add(time.Duration(i + rand.Intn(100000))).Nanosecond())
		e["mun1"] = i
		e["float"] = float64(i)
		s.Add(e)
	}

	sw := NewSegmentWriter("/tmp", s)

	findColumnByName(sw.segment.Columns, "_time")
	//assert.False(isSortedByTs(c.))

	err := sw.Write()

	//assert.True(isSortedByTs())

	if !assert.NoErrorf(err, "an error occurred while writing the segment: %v", err) {
		return
	}
}

func TestColumnWrite_Num(t *testing.T) {

	assert := assert.New(t)

	indexInfo := segment.NewIndexInfo("test-index")
	indexInfo.AddField("mun1", segment.FieldTypeInt, true)

	s := inmem.NewSegment(indexInfo, "123456", nil)

	for i := 0; i < 10000; i++ {
		e := make(segment.Event)
		e["_time"] = uint64(time.Now().Add(time.Duration(i + rand.Intn(100000))).Nanosecond())
		e["mun1"] = i
		s.Add(e)
	}

	sw := NewSegmentWriter("/Users/sdominguez/desa/event_db_data", s)

	sw.Write()

	segment, _ := readers.ReadSegment("/Users/sdominguez/desa/event_db_data")

	log.Printf("Segmento %v", segment)

	for z := 1; z < 2; z++ {
		it := segment.Columns[z].Scan()
		t.Logf("Processing column %s", segment.IndexInfo.Fields[z].Name)
		i := 0

		for it.HasNext() {
			it.Next()
			i++
		}
		assert.Equal(6, i) /// da 6 ver ??

		b := roaring.NewBitmap()
		b.Add(1)
		b.Add(2)
		b.Add(3)
		b.Add(4)
		b.Add(990)
		segment.Columns[z].SetFilter(b)

		i = 0
		pages := make([]*inmem.Page, 0)
		for it.HasNext() {
			p := it.Next()
			pages = append(pages, p)
			i++
		}
		assert.Equal(0, pages[0].StartID)
		assert.Equal(1, i)

		b.Add(1001)
		segment.Columns[z].SetFilter(b)

		i = 0
		for it.HasNext() {
			it.Next()
			i++
		}

		assert.Equal(2, i)
	}

}

func TestColumnWrite_Float(t *testing.T) {

	assert := assert.New(t)

	indexInfo := segment.NewIndexInfo("test-index")
	indexInfo.AddField("float", segment.FieldTypeFloat, true)

	s := inmem.NewSegment(indexInfo, "123456", nil)

	for i := 0; i < 10000; i++ {
		e := make(segment.Event)
		e["_time"] = uint64(time.Now().Add(time.Duration(i + rand.Intn(100000))).Nanosecond())
		e["float"] = float64(i)
		s.Add(e)
	}

	sw := NewSegmentWriter("/Users/sdominguez/desa/event_db_data", s)

	sw.Write()

	segment, _ := readers.ReadSegment("/Users/sdominguez/desa/event_db_data")

	log.Printf("Segmento %v", segment)

	for z := 1; z < 2; z++ {
		it := segment.Columns[z].Scan()
		t.Logf("Processing column %s", segment.IndexInfo.Fields[z].Name)
		i := 0

		for it.HasNext() {
			it.Next()
			i++
		}
		assert.Equal(10, i)

		b := roaring.NewBitmap()
		b.Add(1)
		b.Add(2)
		b.Add(3)
		b.Add(4)
		b.Add(990)
		b.Add(1001)
		segment.Columns[z].SetFilter(b)

		i = 0
		pages := make([]*inmem.Page, 0)
		for it.HasNext() {
			p := it.Next()
			pages = append(pages, p)
			i++
		}
		assert.Equal(0, pages[0].StartID)
		assert.Equal(1, i)

		b.Add(1002)
		segment.Columns[z].SetFilter(b)

		i = 0
		for it.HasNext() {
			it.Next()
			i++
		}

		assert.Equal(2, i)
	}

}
