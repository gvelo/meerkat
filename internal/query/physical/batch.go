package physical

import "meerkat/internal/storage/vector"

type Col struct {
	Group int
	Order int
	Vec   vector.Vector
}

type Batch struct {
	Len     int
	Columns map[string]Col
}
