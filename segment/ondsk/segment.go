package ondsk

import "eventdb/segment"

type Segment struct {
	IndexInfo *segment.IndexInfo
	ID        string
	Idx       []interface{}
	MinTS     int64
	MaxTS     int64
	Monotonic bool
}
