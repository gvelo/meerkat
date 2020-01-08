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
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"testing"
	"time"
)

func getFieldsInfo() *segment.IndexInfo {

	ii := segment.NewIndexInfo("index")
	ii.AddField("msg", segment.FieldTypeText, true)
	ii.AddField("source", segment.FieldTypeKeyword, true)
	ii.AddField("num1", segment.FieldTypeInt, true)
	ii.AddField("num2", segment.FieldTypeFloat, true)
	return ii
}

func getEvents() []segment.Event {
	return []segment.Event{{
		"_time":  int(time.Now().Nanosecond()),
		"source": "log",
		"msg":    "test message one",
		"num1":   int(1),
		"num2":   float64(1.0),
	}, {
		"_time":  int(time.Now().Nanosecond() + 1),
		"source": "log",
		"msg":    "test message two",
		"num1":   int(2),
		"num2":   float64(2.0),
	}, {
		"_time":  int(time.Now().Nanosecond() + 2),
		"source": "other",
		"msg":    "test message three",
		"num1":   int(3),
		"num2":   float64(3.0),
	}, {
		"_time":  int(time.Now().Nanosecond() + 3),
		"source": "sother",
		"msg":    "test message four",
		"num1":   int(1),
		"num2":   float64(1.0),
	}}
}

func TestStoreWriterReaderFewEvents(t *testing.T) {

	a := assert.New(t)

	p := "/tmp/"

	indexInfo := getFieldsInfo()
	writeChan := make(chan *inmem.Segment, 100)
	segment := inmem.NewSegment(indexInfo, "123456", writeChan)

	start := time.Now()
	for _, e := range getEvents() {
		segment.Add(e)
	}
	t.Logf("Index events took %v ", time.Since(start))

	start = time.Now()
	err := WriteSegment(p, segment)
	if err != nil {
		a.Fail("Failed")
	}
	t.Logf("Write events took %v ", time.Since(start))

	start = time.Now()
	if err != nil {
		a.Fail("Failed")
	}
	t.Logf("Write events took %v ", time.Since(start))

}

func createEvents(num int) []segment.Event {
	evts := make([]segment.Event, 0)
	for i := 0; i < num; i++ {
		evt := segment.Event{
			"ts":     uint64(time.Now().Nanosecond()),
			"msg":    fmt.Sprintf("test message %d ", i),
			"source": "log",
			"num1":   uint64(i),
			"num2":   uint64(float64(i)),
		}
		evts = append(evts, evt)
	}
	return evts
}

/*
func TestStoreWriterReaderMoreEvents(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/store"

	e := testFindEvent(t, p, 2, 10000, 100)

	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 2 ", e["msg"])
	a.Equal(uint64(2), e["num1"])
	a.Equal(uint64(2.0), e["num2"])

	e = testFindEvent(t, p, 100, 10000, 100)
	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 100 ", e["msg"])
	a.Equal(uint64(100), e["num1"])
	a.Equal(uint64(100.0), e["num2"])

	e = testFindEvent(t, p, 101, 10000, 100)
	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 101 ", e["msg"])
	a.Equal(uint64(101), e["num1"])
	a.Equal(uint64(101.0), e["num2"])

	e = testFindEvent(t, p, 151, 10000, 100)
	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 151 ", e["msg"])
	a.Equal(uint64(151), e["num1"])
	a.Equal(uint64(151.0), e["num2"])

	e = testFindEvent(t, p, 50000, 1000000, 1000)
	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 50000 ", e["msg"])
	a.Equal(uint64(50000), e["num1"])
	a.Equal(uint64(50000.0), e["num2"])

	e = testFindEvent(t, p, 100000, 1000000, 1000)
	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 100000 ", e["msg"])
	a.Equal(uint64(100000), e["num1"])
	a.Equal(uint64(100000.0), e["num2"])
}

func TestStoreWriter0(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/store"

	e := testFindEvent(t, p, 0, 10, 100)

	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 0 ", e["msg"])
	a.Equal(uint64(0), e["num1"])
	a.Equal(uint64(0.0), e["num2"])
}

func TestStoreWriter1(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/store"

	e := testFindEvent(t, p, 1, 10, 100)

	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 1 ", e["msg"])
	a.Equal(uint64(1), e["num1"])
	a.Equal(uint64(1.0), e["num2"])
}

func TestStoreWriter50(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/store"

	e := testFindEvent(t, p, 50, 100, 100)

	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 50 ", e["msg"])
	a.Equal(uint64(50), e["num1"])
	a.Equal(uint64(50.0), e["num2"])
}

func TestStoreWriter1000(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/store"

	e := testFindEvent(t, p, 200, 1000, 100)

	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 200 ", e["msg"])
	a.Equal(uint64(200), e["num1"])
	a.Equal(uint64(200.0), e["num2"])
}

func testFindEvent(t *testing.T, p string, id int, create int, ixl int) segment.Event {

	start := time.Now()
	_, err := WriteStore(p, createEvents(create), getFieldsInfo(), ixl)
	if err != nil {
		t.Fail()
	}
	t.Logf("write %d events took %v ", create, time.Since(start))
	start = time.Now()
	ds, _ := readers.ReadStore(p+".idx", getFieldsInfo(), ixl)
	e, err := ds.Lookup(id) //readers.ReadEvent(p, id, getFieldsInfo().Fields, ixl)
	if err != nil {
		t.Fail()
	}
	t.Logf("read event took %v ", time.Since(start))
	return e
}

*/
