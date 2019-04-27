package inmem

import (
	"eventdb/segment"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getFieldsInfo() []segment.FieldInfo {
	fieldInfo := make([]segment.FieldInfo, 4)

	fieldInfo[0].FieldID = 0
	fieldInfo[0].FieldType = segment.FieldTypeText
	fieldInfo[0].FieldName = "msg"
	fieldInfo[1].FieldID = 1
	fieldInfo[1].FieldType = segment.FieldTypeKeyword
	fieldInfo[1].FieldName = "source"

	return fieldInfo
}

func createEvents() []map[string]interface{} {

	events := make([]map[string]interface{}, 0)

	for i := 0; i < 1000; i++ {
		event := make(map[string]interface{})
		event["msg"] = fmt.Sprintf("event %v", i)
		event["src"] = "log"
		events = append(events, event)
	}

	return events

}

func TestAddEvent(t *testing.T) {

	fieldInfo := getFieldsInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment("testindex", "testid", fieldInfo, writeChan)

	for _, event := range events {
		segment.Add(event)
	}

	segment.Write()

}

func TestAddEventOnvalidState(t *testing.T) {

	fieldInfo := getFieldsInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment("testindex", "testid", fieldInfo, writeChan)

	segment.Add(events[0])
	segment.Write()

	assert.Panics(t, func() { segment.Add(events[0]) }, "add event in invalid state should panic")

}

func TestWriteOnInvalidState(t *testing.T) {

	fieldInfo := getFieldsInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment("testindex", "testid", fieldInfo, writeChan)

	segment.Add(events[0])
	segment.Write()

	assert.Panics(t, segment.Write, "write segment on invalid state should panic")

}

func TestCloseOnInvalidState(t *testing.T) {

	fieldInfo := getFieldsInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment("testindex", "testid", fieldInfo, writeChan)

	segment.Add(events[0])

	assert.Panics(t, segment.Close, "close segment on invalid state should panic")

}
