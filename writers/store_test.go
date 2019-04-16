package writers

import (
	"eventdb/segment"
	"fmt"
	"testing"
	"time"
)

func getFieldsInfo() []segment.FieldInfo {

	fieldInfo := make([]segment.FieldInfo, 4)

	fieldInfo[0].FieldType = segment.FieldTypeText
	fieldInfo[0].FieldName = "msg"
	fieldInfo[1].FieldType = segment.FieldTypeKeyword
	fieldInfo[1].FieldName = "source"
	fieldInfo[2].FieldType = segment.FieldTypeInt
	fieldInfo[2].FieldName = "num1"
	fieldInfo[3].FieldType = segment.FieldTypeInt
	fieldInfo[3].FieldName = "num2"

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

func TestStoreWriterReaderEvents(t *testing.T) {
	p := "/tmp/store.bin"
	offsets,err:=WriteEvents(p, getEvents(), getFieldsInfo())
	if err != nil {
		t.Fail()
	}
	fmt.Printf("%v", offsets)

}
