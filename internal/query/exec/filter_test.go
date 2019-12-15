package exec

import (
	"meerkat/internal/query/logical"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"meerkat/internal/storage/writers"
	"testing"
	"time"
)

func getFieldsInfo() *segment.IndexInfo {

	ii := segment.NewIndexInfo("index")
	ii.AddField("msg", segment.FieldTypeText, true)
	ii.AddField("source", segment.FieldTypeKeyword, true)
	ii.AddField("num1", segment.FieldTypeInt, true)
	ii.AddField("num2", segment.FieldTypeFloat, true)
	return ii
}

func getEvents() []segment.Event {
	return []segment.Event{{
		"_time":  int(time.Now().Nanosecond()),
		"source": "log",
		"msg":    "test message one",
		"num1":   int(1),
		"num2":   float64(1.0),
	}, {
		"_time":  int(time.Now().Nanosecond() + 1),
		"source": "log",
		"msg":    "test message two",
		"num1":   int(2),
		"num2":   float64(2.0),
	}, {
		"_time":  int(time.Now().Nanosecond() + 2),
		"source": "other",
		"msg":    "test message three",
		"num1":   int(3),
		"num2":   float64(3.0),
	}, {
		"_time":  int(time.Now().Nanosecond() + 3),
		"source": "sother",
		"msg":    "test message four",
		"num1":   int(1),
		"num2":   float64(1.0),
	}}
}

func Test_Filter(t *testing.T) {

	p := "/tmp/"

	indexInfo := getFieldsInfo()
	writeChan := make(chan *inmem.Segment, 100)
	segment := inmem.NewSegment(indexInfo, "123456", writeChan)

	for _, e := range getEvents() {
		segment.Add(e)
	}

	writers.WriteSegment(p, segment)

	ex1 := logical.NewExp(logical.FLOAT, "12.5")
	ex2 := logical.NewExp(logical.IDENTIFIER, "col")
	lf := logical.NewFilter(ex2, logical.EQ, ex1)

	f := NewSFilter(lf)

	f.Execute(nil)

}
