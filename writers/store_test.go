package writers

import (
	"eventdb/readers"
	"eventdb/segment"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getFieldsInfo() []segment.FieldInfo {

	fieldInfo := make([]segment.FieldInfo, 5)

	fieldInfo[0].Type = segment.FieldTypeTimestamp
	fieldInfo[0].Name = "ts"
	fieldInfo[1].Type = segment.FieldTypeText
	fieldInfo[1].Name = "msg"
	fieldInfo[2].Type = segment.FieldTypeKeyword
	fieldInfo[2].Name = "source"
	fieldInfo[3].Type = segment.FieldTypeInt
	fieldInfo[3].Name = "num1"
	fieldInfo[4].Type = segment.FieldTypeInt
	fieldInfo[4].Name = "num2"

	return fieldInfo
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

	p := "/tmp/store.bin"

	start := time.Now()
	offsets, err := WriteEvents(p, getEvents(), getFieldsInfo(), 100)
	if err != nil {
		t.Fail()
	}
	t.Logf("write events took %v ", time.Since(start))

	start = time.Now()
	err = WriteStoreIdx(p, offsets, 100)
	if err != nil {
		t.Fail()
	}
	t.Logf("write events took idx %v ", time.Since(start))

	start = time.Now()
	e, err := readers.ReadEvent(p, 2, getFieldsInfo(), 100)
	if err != nil {
		t.Fail()
	}
	t.Logf("find event took idx %v ", time.Since(start))
	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal(e["msg"], "test message three")
	a.Equal(uint64(3), e["num1"])
	a.Equal(uint64(3.0), e["num2"])

}

func createEvents(num int) []segment.Event {
	evts := make([]segment.Event, 0)
	for i := 0; i < num; i++ {
		evt := segment.Event{
			"ts":     uint64(time.Now().Nanosecond()),
			"source": "log",
			"msg":    fmt.Sprintf("test message %d ", i),
			"num1":   uint64(i),
			"num2":   uint64(float64(i)),
		}
		evts = append(evts, evt)
	}
	return evts
}

func TestStoreWriterReaderMoreEvents(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/store2.bin"

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
	p := "/tmp/store4.bin"

	e := testFindEvent(t, p, 0, 10, 100)

	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 0 ", e["msg"])
	a.Equal(uint64(0), e["num1"])
	a.Equal(uint64(0.0), e["num2"])
}

func TestStoreWriter1(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/store4.bin"

	e := testFindEvent(t, p, 1, 10, 100)

	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 1 ", e["msg"])
	a.Equal(uint64(1), e["num1"])
	a.Equal(uint64(1.0), e["num2"])
}

func TestStoreWriter50(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/store4.bin"

	e := testFindEvent(t, p, 50, 100, 100)

	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 50 ", e["msg"])
	a.Equal(uint64(50), e["num1"])
	a.Equal(uint64(50.0), e["num2"])
}

func TestStoreWriter200(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/store3.bin"

	e := testFindEvent(t, p, 200, 1000, 100)

	a.NotNil(e)
	a.True(len(e) == 5)
	a.Equal("test message 200 ", e["msg"])
	a.Equal(uint64(200), e["num1"])
	a.Equal(uint64(200.0), e["num2"])
}

func testFindEvent(t *testing.T, p string, id uint64, create int, ixl uint64) segment.Event {

	start := time.Now()
	_, err := WriteEvents(p, createEvents(create), getFieldsInfo(), ixl)
	if err != nil {
		t.Fail()
	}
	t.Logf("write %d events took %v ", create, time.Since(start))
	start = time.Now()
	e, err := readers.ReadEvent(p, id, getFieldsInfo(), ixl)
	if err != nil {
		t.Fail()
	}
	t.Logf("read event took %v ", time.Since(start))
	return e
}
