package exec

import (
	"meerkat/internal/query/logical"
)

type Filter struct {
	MultiNodeImpl
	Op *logical.Filter
}

func NewMFilter(op *logical.Filter) *Filter {
	return &Filter{
		MultiNodeImpl: MultiNodeImpl{
			parent:   nil,
			children: make([]ExNode, 0),
		},
		Op: op,
	}

}

func NewSFilter(op *logical.Filter) *Filter {
	return &Filter{
		MultiNodeImpl: MultiNodeImpl{
			parent:   nil,
			children: make([]ExNode, 0),
		},
		Op: op,
	}

}

func (p *Filter) Execute(ctx Context) (Cursor, error) {

	return nil, nil
}

func (p *Filter) String() string {
	return ""
}
