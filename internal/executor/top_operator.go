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
	"github.com/psilva261/timsort"
	"meerkat/internal/storage/vector"
)

// TopOperator operator takes a list of keys
// and returns the a map with the keys and the occurrences
// of these keys in storage.Vector ordered by ascending order if asc == true
// in descending order otherwise.
//
//  Example:
//
//  input      output
//  2          [ 2 , 100 ] , [ 2, 7]
//  2
//  100
//  100
//  100
//  100
//  100
//  100
//  100
//
type TopOperator struct {
	n     int
	asc   bool
	child VectorOperator
}

// NewTopOperator creates a new vector operator.
func NewTopOperator(child VectorOperator, n int, asc bool) MultiVectorOperator {
	return &TopOperator{
		n,
		asc,
		child,
	}
}

func (op *TopOperator) Init() {
	op.child.Init()

}

func (op *TopOperator) Destroy() {
	op.child.Destroy()
}

func (op *TopOperator) Next() []vector.Vector {

	m := make(map[int]int, op.n)

	for vec := op.child.Next(); vec != nil; vec = op.child.Next() {

		keys := vec.(*vector.IntVector).Values()

		for i := 0; i < len(keys); i++ {
			m[keys[i]]++
		}
	}

	kl := make([]int, 0, len(m))
	for k, _ := range m {
		kl = append(kl, k)
	}

	var f timsort.IntLessThan
	if op.asc {
		f = func(a, b int) bool {
			return a < b
		}
	} else {
		f = func(a, b int) bool {
			return a > b
		}
	}

	timsort.Ints(kl, f)

	k := make([]int, 0, 10)
	v := make([]int, 0, 10)
	for i := 0; i < op.n; i++ {
		k = append(k, kl[i])
		v = append(v, m[kl[i]])
	}

	return nil // []storage.Vector{storage.NewIntVector(k), storage.NewIntVector(v)}

}
