package exec

import "meerkat/internal/query/logical"

type Filter struct {
	NodeImp
}

func NewFilter(f *logical.Exp) *Filter {

	return &Filter{
		NodeImp: NodeImp{
			parent:   nil,
			children: make([]OpNode, 0),
		},
	}

}

func (p *Filter) Execute(ctx Context) ([][]interface{}, error) {
	panic("implement me")
}

func (p *Filter) String() string {
	return ""
}
