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
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage"
	"testing"
)

type TestOperator struct {
	val [][]int
	i   int
}

func (op *TestOperator) Destroy() {
	panic("implement me")
}

func NewTestOperator(vals [][]int) VectorOperator {
	return &TestOperator{
		val: vals,
	}
}

func (op *TestOperator) Init() {
	op.i = 0
}

func (op *TestOperator) Next() storage.Vector {
	if op.i < len(op.val) {
		// sc := storage.NewIntVector(op.val[op.i])
		op.i++
		return nil
	}
	return nil
}

func setUpTop() VectorOperator {
	l := make([][]int, 0)
	l = append(l, appendNTimes(10, 10))
	l = append(l, appendNTimes(11, 11))
	l = append(l, appendNTimes(12, 12))
	l = append(l, appendNTimes(13, 13))
	l = append(l, appendNTimes(14, 14))
	l = append(l, appendNTimes(15, 15))
	l = append(l, appendNTimes(16, 16))
	l = append(l, appendNTimes(17, 17))
	l = append(l, appendNTimes(18, 18))
	l = append(l, appendNTimes(19, 19))
	l = append(l, appendNTimes(20, 20))
	l = append(l, appendNTimes(21, 21))
	return NewTestOperator(l)
}

func setUpTopN(n int) VectorOperator {
	l := make([][]int, 0)
	l = append(l, appendNTimes(n+10, n/10))
	l = append(l, appendNTimes(n+11, n/10))
	l = append(l, appendNTimes(n+12, n/10))
	l = append(l, appendNTimes(n+13, n/10))
	l = append(l, appendNTimes(n+14, n/10))
	l = append(l, appendNTimes(n+15, n/10))
	l = append(l, appendNTimes(n+16, n/10))
	l = append(l, appendNTimes(n+17, n/10))
	l = append(l, appendNTimes(n+18, n/10))
	l = append(l, appendNTimes(n+19, n/10))
	l = append(l, appendNTimes(n+20, n/10))
	return NewTestOperator(l)
}

func TestTopDesc(t *testing.T) {
	list := setUpTop()
	a := assert.New(t)

	op := NewTopOperator(list, 5, false)

	r := op.Next()
	k := r[0].(storage.IntVector).ValuesAsInt()
	v := r[1].(storage.IntVector).ValuesAsInt()

	a.Len(k, 5, "length is wrong ")
	a.Len(v, 5, "length is wrong ")
	c := 21
	for i := 0; i < len(k); i++ {
		a.Equal(k[i], c)
		a.Equal(v[i], c)
		c--
	}
}

func TestTopAsc(t *testing.T) {
	list := setUpTop()
	a := assert.New(t)

	op := NewTopOperator(list, 5, true)

	r := op.Next()
	k := r[0].(storage.IntVector).ValuesAsInt()
	v := r[1].(storage.IntVector).ValuesAsInt()

	a.Len(k, 5, "length is wrong ")
	a.Len(v, 5, "length is wrong ")

	c := 10
	for i := 0; i < len(k); i++ {
		a.Equal(k[i], c)
		a.Equal(v[i], c)
		c++
	}

}

//TODO: make a benchmark more expensive for timsort to check how it goes
func BenchmarkTop(b *testing.B) {
	list := setUpTopN(b.N)

	op := NewTopOperator(list, 5, false)

	b.ResetTimer()
	op.Next()
}
