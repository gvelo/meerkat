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
	"fmt"
	"meerkat/internal/cluster"
	"meerkat/internal/query/logical"
	"meerkat/internal/query/physical"
	"meerkat/internal/storage"
	"sync"
)

type ExecutableGraphBuilder interface {
	BuildCoordinatorGraph(writer QueryOutputWriter, outputs []*logical.Fragment) (physical.OutputOp, error)
	BuildNodeGraph(outputs []*logical.Fragment, outputWg *sync.WaitGroup) ([]physical.OutputOp, error)
}

// necesitamos el query id para loguo o viene el el plan ?
func NewExecutableGraphBuilder(connReg cluster.ConnRegistry,
	segReg storage.SegmentRegistry,
	streamReg StreamRegistry) ExecutableGraphBuilder {

	return &executableGraphBuilder{
		segReg:    segReg,
		streamReg: streamReg,
	}

}

type executableGraphBuilder struct {
	segReg    storage.SegmentRegistry
	streamReg StreamRegistry
}

func (e *executableGraphBuilder) BuildCoordinatorGraph(writer QueryOutputWriter, outputs []*logical.Fragment) (physical.OutputOp, error) {
	fmt.Println("BuildCoordinatorGraph : ", outputs)
	return nil, nil
}

func (e *executableGraphBuilder) BuildNodeGraph(outputs []*logical.Fragment, outputWg *sync.WaitGroup) ([]physical.OutputOp, error) {
	fmt.Println("BuildNodeGraph :", outputs)
	return nil, nil
}

type GraphBuilder struct {
	childs []physical.BatchOperator
}

func (g *GraphBuilder) VisitPre(n logical.Node) logical.Node {
	panic("VisitPre implement me")
}

func (g *GraphBuilder) VisitPost(n logical.Node) logical.Node {
	panic("VisitPre implement me")
}
