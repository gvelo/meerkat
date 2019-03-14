package index

import (
	"testing"
)

func TestAddEvent(T *testing.T) {

	fieldInfo := make(map[string]FieldType)

	fieldInfo["msg"] = FieldTypeText
	fieldInfo["src"] = FieldTypeKeyword

	index := newIndex("test", fieldInfo)

	event := &Event{
		Timestamp: 0,
		Fields:    make(map[string]interface{}),
	}

	event.Fields["msg"] = "testmsg a new mesage"
	event.Fields["src"] = "log"

	index.addEvent(event)

}
