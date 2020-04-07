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

package executor

import (
	"meerkat/internal/storage/vector"
)

func NewSortOperator(child MultiVectorOperator, colIds []int) *SortOperator {
	return &SortOperator{child: child, colIds: colIds}
}

// SortOperator
type SortOperator struct {
	child  MultiVectorOperator // (Positions to review)
	colIds []int
	sz     int
}

func (op *SortOperator) Init() {
	op.child.Init()
}

func (op *SortOperator) Destroy() {
	op.child.Destroy()
}

func (op *SortOperator) Next() []vector.Vector {
	n := op.child.Next()

	return n
}
