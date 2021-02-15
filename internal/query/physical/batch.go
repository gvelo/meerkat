package physical

import (
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
)

// Col represent a column inside a Batch.
// In order to keep the column batch ordered as it was speficied in the query
// column are aranged into groups ( ie mv-expand create a new group ) and each
// group maintain an internal order
type Col struct {
	// Group represent the group a column belong to
	Group int64
	// Order represent the column order inside a column group
	Order int64
	// Vec is the column vector
	Vec vector.Vector
	// ColumnType represent the column type
	ColumnType storage.ColumnType
}

type Batch struct {
	Len     int
	Columns map[string]Col
}

func NewBatch() Batch {
	return Batch{
		Len:     0,
		Columns: make(map[string]Col),
	}
}
