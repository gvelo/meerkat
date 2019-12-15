package exec

import (
	"meerkat/internal/query/logical"
)

type BinFilter struct {
	MultiNodeImpl
}

func NewBinFilter(f *logical.Exp) *BinFilter {
	return &BinFilter{
		MultiNodeImpl: MultiNodeImpl{
			canceled: false,
			parent:   nil,
			children: nil,
		},
	}

}

func (p *BinFilter) Execute(ctx Context) (Cursor, error) {
	panic("implement me")
}

func (p *BinFilter) String() string {
	return ""
}
