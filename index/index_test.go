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

func getFieldsInfo() []FieldInfo {
	fieldInfo := make([]FieldInfo, 4)

	fieldInfo[0].fieldType = FieldTypeText
	fieldInfo[0].fieldName = "msg"
	fieldInfo[1].fieldType = FieldTypeKeyword
	fieldInfo[1].fieldName = "source"
	fieldInfo[2].fieldType = FieldTypeInt
	fieldInfo[2].fieldName = "num1"
	fieldInfo[3].fieldType = FieldTypeInt
	fieldInfo[3].fieldName = "num2"

	return fieldInfo
}

func createEvents() []*Event {
	return []*Event{{
		Timestamp: 0,
		Fields: map[string]interface{}{
			"source": "log",
			"msg":    "test message one",
			"num1":   1,
			"num2":   1.0,
		},
	}, {
		Timestamp: 0,
		Fields: map[string]interface{}{
			"source": "log",
			"msg":    "test message two",
			"num1":   2,
			"num2":   2.0,
		},
	}, {
		Timestamp: 0,
		Fields: map[string]interface{}{
			"source": "other",
			"msg":    "test message three",
			"num1":   3,
			"num2":   3.0,
		},
	}}
}

func TestCardinalityKeyword(T *testing.T) {

	fieldInfo := getFieldsInfo()

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

	index := newInMemoryIndex(getFieldsInfo())
	events := createEvents()

	for _, e := range events {
		index.addEvent(e)
	}

	bitmap := index.lookup(getFieldsInfo()[0], "test")
	cardinallity := bitmap.GetCardinality()

	if cardinallity != 3 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap = index.lookup(getFieldsInfo()[0], "three")
	cardinallity = bitmap.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

}

func TestBitmapOr(T *testing.T) {

	index := newInMemoryIndex(getFieldsInfo())
	events := createEvents()

	for _, e := range events {
		index.addEvent(e)
	}

	bitmap1 := index.lookup(getFieldsInfo()[0], "one")
	cardinallity := bitmap1.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap2 := index.lookup(getFieldsInfo()[0], "two")
	cardinallity = bitmap2.GetCardinality()

	if cardinallity != 1 {
		T.Errorf("wrong cardinallity expect %v got %v", 1, cardinallity)
	}

	bitmap3 := index.lookup(getFieldsInfo()[0], "three")
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

func TestInsertInFieldTypeInt(T *testing.T) {
	/*
		index := newInMemoryIndex(getFieldsInfo())
		events := createEvents()

		for _, e := range events {
			index.addEvent(e)
		}
	*/
	//TODO: terminar el test.

}
