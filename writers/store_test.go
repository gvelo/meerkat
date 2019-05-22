package writers

import (
	"eventdb/readers"
	"eventdb/segment"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func getFieldsInfo() *segment.IndexInfo {

	fieldInfo := make([]*segment.FieldInfo, 5)

	fieldInfo[0] = &segment.FieldInfo{Type: segment.FieldTypeTimestamp, Name: "ts"}
	fieldInfo[1] = &segment.FieldInfo{Type: segment.FieldTypeText, Name: "msg"}
	fieldInfo[2] = &segment.FieldInfo{Type: segment.FieldTypeKeyword, Name: "source"}
	fieldInfo[3] = &segment.FieldInfo{Type: segment.FieldTypeInt, Name: "num1"}
	fieldInfo[4] = &segment.FieldInfo{Type: segment.FieldTypeInt, Name: "num2"}

	return &segment.IndexInfo{Fields: fieldInfo}
}

func getEvents() []segment.Event {
	return []segment.Event{{
		"ts":     uint64(time.Now().Nanosecond()),
		"source": "log",
		"msg":    "test message one",
		"num1":   uint64(1),
		"num2":   uint64(1.0),
	}, {
		"ts":     uint64(time.Now().Nanosecond() + 1),
		"source": "log",
		"msg":    "test message two",
		"num1":   uint64(2),
		"num2":   uint64(2.0),
	}, {
		"ts":     uint64(time.Now().Nanosecond() + 2),
		"source": "other",
		"msg":    "test message three",
		"num1":   uint64(3),
		"num2":   uint64(3.0),
	}, {
		"ts":     uint64(time.Now().Nanosecond() + 3),
		"source": "sother",
		"msg":    "test message four",
		"num1":   uint64(1),
		"num2":   uint64(1.0),
	}}
}

func TestStoreWriterReaderFewEvents(t *testing.T) {

	a := assert.New(t)

	p := "/tmp/store"

	start := time.Now()
	offsets, err := WriteStore(p, getEvents(), getFieldsInfo(), 100)
	log.Println(fmt.Sprintf("offsets W %v", offsets))
	if err != nil {
		t.Fail()
	}

	t.Logf("write events took %v ", time.Since(start))
	a.Nil(err)
	a.NotNil(offsets)

	start = time.Now()

	ds, _ := readers.ReadStore(p+".idx", getFieldsInfo(), 100)
	e, err := ds.Lookup(2)
	if err != nil {
		t.Fail()
	}
	t.Logf("find event took idx %v ", time.Since(start))
	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message three", e["msg"])

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
