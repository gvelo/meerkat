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

// FIX: did not find a better solution to prevent circular import with storage
func mapColumnTypeToVector(t storage.ColumnType) vector.Vector {
	switch t {
	case storage.ColumnType_BOOL:
		return &vector.BoolVector{}
	case storage.ColumnType_DATETIME, storage.ColumnType_INT64, storage.ColumnType_TIMESTAMP:
		return &vector.Int64Vector{}
	case storage.ColumnType_INT32:
		return &vector.Int32Vector{}
	case storage.ColumnType_FLOAT64:
		return &vector.Float64Vector{}
	case storage.ColumnType_STRING:
		return &vector.ByteSliceVector{}
	default:
		panic("No vector found.")
	}
}

// Util function to clone the column.
func (c Col) Clone() Col {
	v := vector.DefaultVectorPool().GetVector(mapColumnTypeToVector(c.ColumnType), c.Vec.HasNulls())
	res := Col{
		Group:      c.Group,
		Order:      c.Order,
		Vec:        v,
		ColumnType: c.ColumnType,
	}
	return res
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

// Util function to clone the batch.
func (b Batch) Clone() Batch {
	res := NewBatch()
	res.Len = b.Len
	res.Columns = make(map[string]Col)

	for key, _ := range b.Columns {
		res.Columns[key] = b.Columns[key].Clone()
	}

	return res
}
