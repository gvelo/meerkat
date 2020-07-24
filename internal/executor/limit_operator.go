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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage/vector"
)

// NewLimitOperator creates a LimitOperator
// The Sort + Limit = Top Operator.
func NewLimitOperator(ctx Context, child MultiVectorOperator, limit int) MultiVectorOperator {
	return &LimitOperator{
		ctx:   ctx,
		child: child,
		limit: limit,
		log:   log.With().Str("src", "LimitOperator").Logger(),
	}
}

// LimitOperator scans a non indexed column, search for the []pos
// registers and returns the bitmap that meets that condition.
//
type LimitOperator struct {
	ctx   Context
	child MultiVectorOperator
	limit int
	total int
	sz    int
	done  bool
	log   zerolog.Logger
}

func (op *LimitOperator) Init() {
	op.child.Init()
	op.total = 0
	op.done = false
}

func (op *LimitOperator) Destroy() {
	op.child.Destroy()
}

func (op *LimitOperator) Next() []interface{} {

	if op.done {
		return nil
	}

	n := op.child.Next()

	if n != nil {

		l := getLen(n[0])
		if op.total+l > op.limit {
			op.done = true
			return op.cutVectors(n, op.limit-op.total)
		}
		op.total += l

	}

	return n
}

func (op *LimitOperator) cutVectors(n []interface{}, limit int) []interface{} {
	for i, _ := range n {
		switch n[i].(type) {
		case vector.Int64Vector:
			vv := n[i].(vector.Int64Vector)
			vv.SetLen(limit)
			n[i] = vv
		case vector.BoolVector:
			vv := n[i].(vector.BoolVector)
			vv.SetLen(limit)
			n[i] = vv
		case vector.Float64Vector:
			vv := n[i].(vector.Float64Vector)
			vv.SetLen(limit)
			n[i] = vv
		case vector.ByteSliceVector:
			vv := n[i].(vector.ByteSliceVector)
			vv.SetLen(limit)
			n[i] = vv
		default:
			log.Error().Msgf("Type not found %v", n[i])
		}
	}
	return n
}
