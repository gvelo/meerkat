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
	"meerkat/internal/query/logical"
	"meerkat/internal/storage/readers"
	"meerkat/internal/storage/segment/ondsk"
	"meerkat/internal/tools"
	"path/filepath"
)

type Optimizer interface {
	// TODO:(sebad) execution
	OptimizeQuery(t *logical.Projection) *logical.Projection
}

type MeerkatOptimizer struct {
	path string
}

func NewMeerkatOptimizer() *MeerkatOptimizer {
	return &MeerkatOptimizer{}
}

func (o *MeerkatOptimizer) OptimizeQuery(t *logical.Projection) *logical.Projection {

	return t
}

//TODO(sebad): DO OPTIMIZE should send indexed Filters first bottom left ?
func (o *MeerkatOptimizer) optimizeFilters(f interface{}) interface{} {
	tools.Log("Optimize filters!")
	return nil
}

// TODO(sebad): get data from Schema not Segment.
func (o *MeerkatOptimizer) getMetadata(i *logical.Projection) *ondsk.Segment {

	var file = filepath.Join(o.path, i.String())

	s, err := readers.ReadSegment(file)
	if err != nil {
		// panic(fmt.Sprintf(" %v does not exist ", i .GetIndexName()))
	}

	return s
}
