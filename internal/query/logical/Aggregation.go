package logical

// logical aggregation be implemented by all functions
type Aggregation struct {
	Function string   // function to apply to fields
	Fields   []string // fields to apply

	parent   Node
	children []Node
}

func NewAggregation(f string, fields []string) *Aggregation {
	return &Aggregation{
		Function: f,
		Fields:   fields,
	}
}

func (p *Aggregation) String() string {
	return "Aggregation"
}
