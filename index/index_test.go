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

	fieldInfo := make([]FieldInfo, 2)

	fieldInfo[0].fieldType = FieldTypeText
	fieldInfo[0].fieldName = "msg"
	fieldInfo[1].fieldType = FieldTypeKeyword
	fieldInfo[1].fieldName = "source"

	index := newInMemoryIndex(fieldInfo)

	events := createEvents()

	for _, e := range events {
		index.addEvent(e)
	}

	bitmap := index.lookup(fieldInfo[1], "log")
	cardinallity := bitmap.GetCardinality()

	if cardinallity != 2 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap = index.lookup(fieldInfo[1], "other")
	cardinallity = bitmap.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

}

func TestCardinalityText(T *testing.T) {
	fieldInfo := make([]FieldInfo, 2)

	fieldInfo[0].fieldType = FieldTypeText
	fieldInfo[0].fieldName = "msg"
	fieldInfo[1].fieldType = FieldTypeKeyword
	fieldInfo[1].fieldName = "source"

	index := newInMemoryIndex(fieldInfo)

	events := createEvents()

	for _, e := range events {
		index.addEvent(e)
	}

	bitmap := index.lookup(fieldInfo[0], "test")
	cardinallity := bitmap.GetCardinality()

	if cardinallity != 3 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap = index.lookup(fieldInfo[0], "three")
	cardinallity = bitmap.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

}

func TestBitmapOr(T *testing.T) {

	fieldInfo := make([]FieldInfo, 2)

	fieldInfo[0].fieldType = FieldTypeText
	fieldInfo[0].fieldName = "msg"
	fieldInfo[1].fieldType = FieldTypeKeyword
	fieldInfo[1].fieldName = "source"

	index := newInMemoryIndex(fieldInfo)

	events := createEvents()

	for _, e := range events {
		index.addEvent(e)
	}

	bitmap1 := index.lookup(fieldInfo[0], "one")
	cardinallity := bitmap1.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap2 := index.lookup(fieldInfo[0], "two")
	cardinallity = bitmap2.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap3 := index.lookup(fieldInfo[0], "three")
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
