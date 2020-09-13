package exec

import (
	"context"
	"meerkat/internal/cluster"
	"meerkat/internal/query/logical"
	"meerkat/internal/query/physical"
	"meerkat/internal/storage"
)

type ExecutableGraph interface {
	Outputs() //[]
	Start()
	Cancel()
	Context() context.Context
}

type ExecutableGraphBuilder interface {
	Build(root *logical.Fragment) ExecutableGraph
}

// necesitamos el query id para loguo o viene el el plan ?
func NewExecutableGraphBuilder(connReg cluster.ConnRegistry,
	segReg storage.SegmentRegistry,
	streamReg StreamRegistry) ExecutableGraphBuilder {

}

type executableGraphBuilder struct {
}

func (b *executableGraphBuilder) build(root *logical.Fragment) *physical.OutputOp {

	// tener en cuenta que los miembros locales nose deben manejar mediante ExchaneIN/OUT

}
