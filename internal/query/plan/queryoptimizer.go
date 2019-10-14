// Copyright 2019 The Meerkat Authors
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

package plan

import (
	"meerkat/internal/query/exec"
	"meerkat/internal/query/logical"
	"meerkat/internal/schema"
)

type Optimizer interface {
	OptimizeQuery(t []logical.Node) []exec.OpNode
}

type MeerkatOptimizer struct {
	path string
}

func NewMeerkatOptimizer() *MeerkatOptimizer {
	return &MeerkatOptimizer{}
}

// This method should take the logical steps and transform in execution steps,
// returning an optimized steps tree
func (o *MeerkatOptimizer) OptimizeQuery(s schema.Schema, l []logical.Node) exec.OpNode {

	exec := make([]exec.OpNode, 0)
	for _, n := range l {
		switch t := n.(type) {
		case *logical.Projection:
			o.transformProjection(s, t)
		case *logical.RootFilter:
			o.transformRootFilter(s, t)
		case *logical.Aggregation:
			o.transformAggregation(s, t)
		}
	}

	return nil
}

func (o *MeerkatOptimizer) transformProjection(schema schema.Schema, p *logical.Projection) []exec.OpNode {

	for i, idx := range p.Indexes {

		schema.IndexByName(idx)
		ep := exec.NewProjection()

	}

	return nil
}

func (o *MeerkatOptimizer) transformRootFilter(schema schema.Schema, r *logical.RootFilter) []exec.OpNode {

	return nil
}

func (o *MeerkatOptimizer) transformAggregation(schema schema.Schema, a *logical.Aggregation) []exec.OpNode {

	return nil
}
