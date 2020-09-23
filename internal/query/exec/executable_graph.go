// Copyright 2020 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


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
