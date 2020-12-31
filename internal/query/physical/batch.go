package physical

import "meerkat/internal/storage/vector"

// Col represent a column inside a Batch.
// In order to keep the column batch ordered as it was speficied in the query
// column are aranged into groups ( ie mv-expand create a new group ) and each
// group maintain an internal order
type Col struct {
	// Group represent the group a column belong to
	Group int
	// Order represent the column order inside a column group
	Order int
	// Vec is the column vector
	Vec vector.Vector
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
