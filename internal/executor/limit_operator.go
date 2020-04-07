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

import "meerkat/internal/storage/vector"

// NewLimitOperator creates a ColumnScanOperator
func NewLimitOperator(ctx Context, child MultiVectorOperator, limit int) MultiVectorOperator {
	return &LimitOperator{
		ctx:   ctx,
		child: child,
		limit: limit,
	}
}

// LimitOperator scans a non indexed column, search for the []pos
// registers and returns the bitmap that meets that condition.
//
type LimitOperator struct {
	ctx   Context
	child MultiVectorOperator
	limit int
}

func (op *LimitOperator) Init() {
	op.child.Init()
}

func (op *LimitOperator) Destroy() {
	op.child.Destroy()
}

func (op *LimitOperator) Next() []vector.Vector {
	n := op.child.Next()

	return n
}
