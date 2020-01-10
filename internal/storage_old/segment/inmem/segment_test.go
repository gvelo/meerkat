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

package inmem

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"meerkat/internal/storage/segment"
	"testing"
	"time"
)

func getIndexInfo() *segment.IndexInfo {

	indexInfo := segment.NewIndexInfo("test_index")
	indexInfo.AddField("msg", segment.FieldTypeText, true)
	indexInfo.AddField("src", segment.FieldTypeKeyword, true)
	indexInfo.AddField("number", segment.FieldTypeInt, true)
	indexInfo.AddField("float", segment.FieldTypeFloat, true)

	return indexInfo
}

func createEvents() []segment.Event {

	events := make([]segment.Event, 0)

	for i := 0; i < 1000; i++ {
		event := make(segment.Event)
		event["msg"] = fmt.Sprintf("event %v", i)
		event["src"] = "log"
		event["number"] = uint64(1)
		event["float"] = math.Float64bits(123.12)
		event["ts"] = uint64(time.Now().Add(time.Duration(i)).Nanosecond())
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
		Add(event)
	}

	Write()

}

func TestAddEventOnvalidState(t *testing.T) {

	indexInfo := getIndexInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment(indexInfo, "testid", writeChan)

	Add(events[0])
	Write()

	assert.Panics(t, func() { Add(events[0]) }, "add event in invalid state should panic")

}

func TestWriteOnInvalidState(t *testing.T) {

	indexInfo := getIndexInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment(indexInfo, "testid", writeChan)

	Add(events[0])
	Write()

	assert.Panics(t, Write, "write segment on invalid state should panic")

}

func TestCloseOnInvalidState(t *testing.T) {

	indexInfo := getIndexInfo()
	events := createEvents()
	writeChan := make(chan *Segment, 100)

	segment := NewSegment(indexInfo, "testid", writeChan)

	Add(events[0])

	assert.Panics(t, Close, "close segment on invalid state should panic")

}
