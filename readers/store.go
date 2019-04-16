package readers

import (
	"eventdb/collection"
	"eventdb/segment"
)

func ReadEvents(name string, ids []uint64, info []segment.FieldInfo) ([]segment.Event, error) {
	return nil, nil
}

func ReadEventsIndex(name string) (*collection.SkipList, error) {
	return nil, nil
}
