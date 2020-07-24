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
	"time"
)

//TODO(sebad): cambiar interface cuando corresponda.
//TODO(sebad): pasar todo a int!

// Bucket operator takes a list of timestamps in nanoseconds and a @time.Duration
// and returns the same list rounded by that duration
//
// Example:
//  Duration = 1 min
//
//  int64 nanosecond count
//  input         output
//
// 12:00:02      12:00:00
// 12:01:22      12:01:00
// 12:01:25      12:01:00
// 12:01:54      12:01:00
// 12:07:43      12:07:00
//

//
//  Bucket operator takes a list of int and a span,
//  and returns the same list rounded by that span.
//
// Example Int: span = 5
//
//  |----|----|----|----|----|----|----|----|----|
// 100  102  103  104  105  106  107  108  109  110
//  |---------|-------------------|--------------|
//  5  				   5                         5
// Example Int: span = 5
//
//  |----|----|----|----|----|----|----|----|----|
// 100  102  103  104  105  106  107  108  109  110
//  |--------------------|------------------------|
//            5  	                5
//
// input    output
//  100      100
//  102      100
//  103      105
//  107      105
//
type BucketOperator struct {
	ctx   Context
	span  int64
	tspan time.Duration
	child MultiVectorOperator
	log   zerolog.Logger
}

func (op *BucketOperator) Init() {
	op.child.Init()
}

// NewBucketOperator creates a new BucketOperator.
func NewBucketOperator(ctx Context, child MultiVectorOperator, span int64) MultiVectorOperator {
	return &BucketOperator{
		ctx,
		span,
		0,
		child,
		log.With().Str("src", "BucketOperator").Logger(),
	}
}

// NewTimeBucketOperator creates a new  TimeBucketOperator .
func NewTimeBucketOperator(ctx Context, child MultiVectorOperator, span string) MultiVectorOperator {
	d, err := time.ParseDuration(span)
	if err != nil {
		log.Error().Err(err)
	}
	return &BucketOperator{
		ctx,
		0,
		d,
		child,
		log.With().Str("src", "BucketOperator").Logger(),
	}
}

func (op *BucketOperator) Destroy() {
	op.child.Destroy()
}

func (op *BucketOperator) Next() []interface{} {

	n := op.child.Next()

	if n != nil {

		// TODO: remove this Magic Number
		vv := n[0].(vector.Int64Vector)
		vp := &vv
		ts := vp.Values()

		if len(ts) == 0 {
			return nil
		}
		init := getNextSpan(ts[0], op.span, op.tspan)

		for i, _ := range ts {
			if ts[i] < init+op.span {
				ts[i] = init
			} else {
				init := getNextSpan(ts[0], op.span, op.tspan)
				ts[i] = init
			}
		}

		v := vector.NewInt64Vector(ts, []uint64{})
		v.SetLen(len(ts))

		n[0] = v
		return n
	}

	return nil
}

func getNextSpan(t, s int64, d time.Duration) int64 {
	if s != 0 {

		if s <= 0 {
			return t
		}

		r := t % s

		if lessThanHalf(r, s) {
			return t - r
		}
		return t + (s - r)

	} else {
		return time.Unix(0, int64(t)).Round(d).UnixNano()
	}
}

func lessThanHalf(x, y int64) bool {
	return uint(x)+uint(x) < uint(y)
}
