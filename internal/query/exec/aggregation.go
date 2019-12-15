package exec

type Aggregation struct {
	SingleNodeImpl
}

func NewAggregation() *Aggregation {
	return &Aggregation{
		SingleNodeImpl: SingleNodeImpl{
			parent: nil,
			child:  nil,
		},
	}
}

func (p *Aggregation) Execute(ctx Context) ([][]interface{}, error) {
	panic("implement me")
}

func (p *Aggregation) String() string {
	return ""
}
