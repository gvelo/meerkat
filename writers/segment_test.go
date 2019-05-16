package writers

import (
	"eventdb/readers"
	"eventdb/segment"
	"eventdb/segment/inmem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSegmentWriterReader(t *testing.T) {

	assert := assert.New(t)

	indexInfo := segment.NewIndexInfo("test-index")
	indexInfo.AddField("testfield", segment.FieldTypeKeyword, true)
	indexInfo.AddField("ts", segment.FieldTypeTimestamp, false)

	s := inmem.NewSegment(indexInfo, "3dfa542d", nil)

	for i := 0; i < 100; i++ {
		e := make(map[string]interface{})
		e["testfield"] = fmt.Sprintf("test %v", i)
		e["ts"] = uint64(time.Now().Add(time.Duration(i)).Nanosecond())
		s.Add(e)
	}

	sw := NewSegmentWriter("/tmp", s)

	err := sw.Write()

	if !assert.NoErrorf(err, "an error occurred while writing the segment: %v", err) {
		return
	}

	sr := readers.NewSegmentReader("/tmp")

	odSegment, err := sr.Read()

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
