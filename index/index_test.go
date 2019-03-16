package index

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
)

func newEvent() *Event {
	return &Event{
		Timestamp: 0,
		Fields:    make(map[string]interface{}),
	}
}

func createEvents() []*Event {
	return []*Event{{
		Timestamp: 0,
		Fields: map[string]interface{}{
			"source": "log",
			"msg":    "test message one",
		},
	}, {
		Timestamp: 0,
		Fields: map[string]interface{}{
			"source": "log",
			"msg":    "test message two",
		},
	}, {
		Timestamp: 0,
		Fields: map[string]interface{}{
			"source": "other",
			"msg":    "test message three",
		},
	}}
}

func TestCardinalityKeyword(T *testing.T) {

	fieldInfo := make(map[string]FieldType)

	fieldInfo["msg"] = FieldTypeText
	fieldInfo["source"] = FieldTypeKeyword

	index := newIndex("test", fieldInfo)

	events := createEvents()

	for _, e := range events {
		index.addEvent(e)
	}

	bitmap := index.lookup("source", "log")
	cardinallity := bitmap.GetCardinality()

	if cardinallity != 2 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap = index.lookup("source", "other")
	cardinallity = bitmap.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

}

func TestCardinalityText(T *testing.T) {

	fieldInfo := make(map[string]FieldType)

	fieldInfo["msg"] = FieldTypeText
	fieldInfo["source"] = FieldTypeKeyword

	index := newIndex("test", fieldInfo)

	events := createEvents()

	for _, e := range events {
		index.addEvent(e)
	}

	bitmap := index.lookup("msg", "test")
	cardinallity := bitmap.GetCardinality()

	if cardinallity != 3 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap = index.lookup("msg", "three")
	cardinallity = bitmap.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

}

func TestBitmapOr(T *testing.T) {

	fieldInfo := make(map[string]FieldType)

	fieldInfo["msg"] = FieldTypeText
	fieldInfo["source"] = FieldTypeKeyword

	index := newIndex("test", fieldInfo)

	events := createEvents()

	for _, e := range events {
		index.addEvent(e)
	}

	bitmap1 := index.lookup("msg", "one")
	cardinallity := bitmap1.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap2 := index.lookup("msg", "two")
	cardinallity = bitmap2.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap3 := index.lookup("msg", "three")
	cardinallity = bitmap3.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	res := roaring.FastOr(bitmap1, bitmap2, bitmap3)

	cardinallity = res.GetCardinality()

	if cardinallity != 3 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

}
