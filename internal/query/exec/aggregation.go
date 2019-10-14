package exec

type Aggregation struct {
	NodeImp
}

func NewAggregation() *Aggregation {
	return &Aggregation{
		NodeImp: NodeImp{
			parent:   nil,
			children: make([]OpNode, 0),
		},
	}
}

func (p *Aggregation) Execute(ctx Context) ([][]interface{}, error) {
	panic("implement me")
}

func (p *Aggregation) String() string {
	return ""
}
