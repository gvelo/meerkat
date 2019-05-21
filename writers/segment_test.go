package writers

import (
	"eventdb/readers"
	"eventdb/segment"
	"eventdb/segment/inmem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestSegmentWriterReader(t *testing.T) {

	assert := assert.New(t)

	indexInfo := segment.NewIndexInfo("test-index")
	indexInfo.AddField("testfield", segment.FieldTypeKeyword, true)
	indexInfo.AddField("ts", segment.FieldTypeTimestamp, false)
	indexInfo.AddField("mun1", segment.FieldTypeInt, true)
	indexInfo.AddField("float", segment.FieldTypeFloat, true)

	s := inmem.NewSegment(indexInfo, "3dfa542d", nil)

	for i := 0; i < 100; i++ {
		e := make(map[string]interface{})
		e["testfield"] = fmt.Sprintf("test %v", i)
		e["ts"] = uint64(time.Now().Add(time.Duration(i)).Nanosecond())
		e["mun1"] = uint64(i)
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

func TestSegmentSorted(t *testing.T) {

	assert := assert.New(t)

	indexInfo := segment.NewIndexInfo("test-index")
	indexInfo.AddField("testfield", segment.FieldTypeKeyword, true)
	indexInfo.AddField("ts", segment.FieldTypeTimestamp, false)
	indexInfo.AddField("mun1", segment.FieldTypeInt, true)
	indexInfo.AddField("float", segment.FieldTypeFloat, true)

	s := inmem.NewSegment(indexInfo, "3dfa542d", nil)

	for i := 0; i < 100; i++ {
		e := make(map[string]interface{})
		e["testfield"] = fmt.Sprintf("test %v", i)
		e["ts"] = uint64(time.Now().Add(time.Duration(i + rand.Intn(100000))).Nanosecond())
		e["mun1"] = uint64(i)
		e["float"] = float64(i)
		s.Add(e)
	}

	sw := NewSegmentWriter("/tmp", s)

	assert.False(isSortedByTs(s.FieldStorage))

	err := sw.Write()

	assert.True(isSortedByTs(s.FieldStorage))

	if !assert.NoErrorf(err, "an error occurred while writing the segment: %v", err) {
		return
	}
}

func isSortedByTs(events []interface{}) bool {
	var ant uint64 = 0
	for _, x := range events {
		if x.(map[string]interface{})["ts"].(uint64) < ant {
			return false
		} else {
			ant = x.(map[string]interface{})["ts"].(uint64)
		}
	}
	return true
}
