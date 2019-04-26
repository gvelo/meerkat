package inmem

import "eventdb/segment"

type Segment struct {
	IndexName    string
	ID           string
	FieldInfo    *[]segment.FieldInfo
	eventID      int
	Idx          map[string]interface{}
	FieldStorage *FieldStorage
	PostingStore *PostingStore
	MinTS        int64
	MaxTS        int64
}
