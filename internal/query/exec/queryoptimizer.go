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

package exec

import (
	"log"
	"meerkat/internal/query/logical"
)

type ExecutionError struct {
	Err string
}

func (e ExecutionError) Error() string {
	return e.Err
}

type Optimizer interface {
	OptimizeQuery(t []logical.Node) (ExNode, error)
}

type MeerkatOptimizer struct {
	path string
}

func NewMeerkatOptimizer() *MeerkatOptimizer {
	return &MeerkatOptimizer{}
}

// This method should take the logical steps and transform in execution steps,
// returning an optimized steps tree
func (o *MeerkatOptimizer) OptimizeQuery(ctx Context, l []logical.Node) (ExNode, error) {

	//get all the metadata to process the query

	// optimize the query with DP

	// Return the lower cost plan.

	// me arden los ojos!
	var bn ExNode
	var root ExNode
	for _, n := range l {

		var err error
		var an ExNode

		switch t := n.(type) {
		case *logical.Projection:
			an, err = o.transformProjection(ctx, t)
			root = an
		case *logical.RootFilter:
			an, err = o.transformRootFilter(ctx, t)
		case *logical.Aggregation:
			an, err = o.transformAggregation(ctx, t)
		}

		if err != nil {
			return nil, err
		} else {
			if bn != nil {
				bn.AddChild(an)
			}
			bn = an
		}

	}

	return root, nil
}

// TODO(sebad): make multi_index
func (o *MeerkatOptimizer) transformProjection(ctx Context, p *logical.Projection) (ExNode, error) {

	ep, err := NewProjection(p)
	if err != nil {
		return nil, err
	}
	return ep, nil
}

func (o *MeerkatOptimizer) transformRootFilter(ctx Context, r *logical.RootFilter) (ExNode, error) {
	return buildFilterOps(nil, r.RootFilter)
}

func buildFilterOps(parent ExNode, r interface{}) (ExNode, error) {

	switch t := r.(type) {

	case *logical.Filter:
		/*
			f := NewFilter(t.Op)
			f.parent = parent

			n, err := buildFilterOps(f, t.Left)
			if err != nil {
				f.children = append(f.children, n)
			}
			n, err = buildFilterOps(f, t.Right)
			if err != nil {
				f.children = append(f.children, n)
			}
		*/
		log.Fatal(t)
		return nil, nil

	case logical.Expression:
		/*
			f := NewExpressionFilter()
			f.parent = parent

			n, err := buildFilterOps(f, t.Left)
			if err != nil {
				f.children = append(f.children, n)
			}
			n, err = buildFilterOps(f, t.Right)
			if err != nil {
				f.children = append(f.children, n)
			}
		*/
		return nil, nil

	default:
		return nil, &ExecutionError{
			Err: "Error building filters",
		}
	}

}

func (o *MeerkatOptimizer) transformAggregation(ctx Context, a *logical.Aggregation) (ExNode, error) {

	return nil, nil
}
