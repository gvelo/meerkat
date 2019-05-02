package inmem

import (
	"eventdb/segment"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getIndexInfo() *segment.IndexInfo {

	indexInfo := segment.NewIndexInfo("test_index")
	indexInfo.AddField("msg", segment.FieldTypeText, true)
	indexInfo.AddField("source", segment.FieldTypeKeyword, true)

	return indexInfo
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

	indexInfo := getIndexInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment(indexInfo, "testid", writeChan)

	for _, event := range events {
		segment.Add(event)
	}

	segment.Write()

}

func TestAddEventOnvalidState(t *testing.T) {

	indexInfo := getIndexInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment(indexInfo, "testid", writeChan)

	segment.Add(events[0])
	segment.Write()

	assert.Panics(t, func() { segment.Add(events[0]) }, "add event in invalid state should panic")

}

func TestWriteOnInvalidState(t *testing.T) {

	indexInfo := getIndexInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment(indexInfo, "testid", writeChan)

	segment.Add(events[0])
	segment.Write()

	assert.Panics(t, segment.Write, "write segment on invalid state should panic")

}

func TestCloseOnInvalidState(t *testing.T) {

	indexInfo := getIndexInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment(indexInfo, "testid", writeChan)

	segment.Add(events[0])

	assert.Panics(t, segment.Close, "close segment on invalid state should panic")

}
