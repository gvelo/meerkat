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
	"meerkat/internal/storage/segment/ondsk"
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
// returning an optimized execution steps list
func (o *MeerkatOptimizer) OptimizeQuery(l []logical.Node) []exec.OpNode {

	exec := make([]exec.OpNode, 0)
	for _, n := range l {
		switch t := n.(type) {
		case *logical.Projection:
			o.transformProjection(t)
		case *logical.RootFilter:
			o.transformRootFilter(t)
		case *logical.Aggregation:
			o.transformAggregation(t)
		}
	}

	return exec
}

func (o *MeerkatOptimizer) transformProjection(p *logical.Projection) []exec.OpNode {

	return nil
}

func (o *MeerkatOptimizer) transformRootFilter(r *logical.RootFilter) []exec.OpNode {

	return nil
}

func (o *MeerkatOptimizer) transformAggregation(a *logical.Aggregation) []exec.OpNode {

	return nil
}

// TODO(sebad): get data from Schema not Segment.
func (o *MeerkatOptimizer) getMetadata(i *logical.Projection) *ondsk.Segment {

	return nil
}
