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

	fieldInfo[0].FieldType = segment.FieldTypeTimestamp
	fieldInfo[0].FieldName = "ts"
	fieldInfo[1].FieldType = segment.FieldTypeText
	fieldInfo[1].FieldName = "msg"
	fieldInfo[2].FieldType = segment.FieldTypeKeyword
	fieldInfo[2].FieldName = "source"
	fieldInfo[3].FieldType = segment.FieldTypeInt
	fieldInfo[3].FieldName = "num1"
	fieldInfo[4].FieldType = segment.FieldTypeInt
	fieldInfo[4].FieldName = "num2"

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

	ass := assert.New(t)

	p := "/tmp/store.bin"
	offsets, err := WriteEvents(p, getEvents(), getFieldsInfo())
	if err != nil {
		t.Fail()
	}

	err = WriteStoreIdx(p, offsets)
	if err != nil {
		t.Fail()
	}

	e, err := readers.ReadEvent(p, 3, getFieldsInfo())
	if err != nil {
		t.Fail()
	}
	ass.NotNil(e)
	ass.True(len(e) == 5)
	assert.Equal(t, e["msg"], "test message three")
	assert.Equal(t, e["num1"], uint64(3))
	assert.Equal(t, e["num2"], uint64(3.0))

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
	ass := assert.New(t)
	p := "/tmp/store2.bin"

	e := testFindEvent(t, p, 3)

	ass.NotNil(e)
	ass.True(len(e) == 5)
	assert.Equal(t, e["msg"], "test message 2 ")
	assert.Equal(t, e["num1"], uint64(2))
	assert.Equal(t, e["num2"], uint64(2.0))

	e = testFindEvent(t, p, 100)
	ass.NotNil(e)
	ass.True(len(e) == 5)
	assert.Equal(t, e["msg"], "test message 99 ")
	assert.Equal(t, e["num1"], uint64(99))
	assert.Equal(t, e["num2"], uint64(99.0))

	e = testFindEvent(t, p, 101)
	ass.NotNil(e)
	ass.True(len(e) == 5)
	assert.Equal(t, e["msg"], "test message 100 ")
	assert.Equal(t, e["num1"], uint64(100))
	assert.Equal(t, e["num2"], uint64(100.0))

	e = testFindEvent(t, p, 151)
	ass.NotNil(e)
	ass.True(len(e) == 5)
	assert.Equal(t, e["msg"], "test message 150 ")
	assert.Equal(t, e["num1"], uint64(150))
	assert.Equal(t, e["num2"], uint64(150.0))

	e = testFindEvent(t, p, 50000)
	ass.NotNil(e)
	ass.True(len(e) == 5)
	assert.Equal(t, e["msg"], "test message 49999 ")
	assert.Equal(t, e["num1"], uint64(49999))
	assert.Equal(t, e["num2"], uint64(49999.0))
}

func TestStoreWriter100000(t *testing.T) {
	ass := assert.New(t)
	p := "/tmp/store2.bin"

	e := testFindEvent(t, p, 100000)

	ass.NotNil(e)
	ass.True(len(e) == 5)
	assert.Equal(t, e["msg"], "test message 99999 ")
	assert.Equal(t, e["num1"], uint64(99999))
	assert.Equal(t, e["num2"], uint64(99999.0))
}

func testFindEvent(t *testing.T, p string, id uint64) segment.Event {

	offsets, err := WriteEvents(p, createEvents(100000), getFieldsInfo())
	if err != nil {
		t.Fail()
	}

	err = WriteStoreIdx(p, offsets)
	if err != nil {
		t.Fail()
	}

	e, err := readers.ReadEvent(p, id, getFieldsInfo())
	if err != nil {
		t.Fail()
	}
	return e
}
