package inmem

// Type represent the type of a field.
type Encoding uint

const (
	RLE Encoding = iota
	Simple8B
	DoubleDelta
	Raw
	ZigZag
	Snappy
)

type Page struct {
	Enc         Encoding
	StartID     int
	PayloadSize int
	Total       int
}

type PageDescriptor struct {
	StartID int
	Offset  int
}
