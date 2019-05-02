package writers

import (
	"eventdb/segment"
	"eventdb/segment/inmem"
	"fmt"
	"testing"
)

func TestNewSegmentWriter(t *testing.T) {

	indexInfo := segment.NewIndexInfo("test-index")
	indexInfo.AddField("testfield", segment.FieldTypeKeyword, true)

	s := inmem.NewSegment(indexInfo, "3dfa542d", nil)

	for i := 0; i < 100; i++ {
		e := make(map[string]interface{})
		e["testfield"] = fmt.Sprintf("test %v", i)
		s.Add(e)
	}

	sw := NewSegmentWriter("/tmp/data", s)

	err := sw.Write()

	if err != nil {
		t.Fatal(err)
	}

}
