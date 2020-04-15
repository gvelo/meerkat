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
	"meerkat/internal/storage/vector"
	"testing"
)

type TestOperator struct {
	values [][]int
	i      int
}

func (op *TestOperator) Destroy() {
	panic("implement me")
}

func NewTestOperator(values [][]int) MultiVectorOperator {
	return &TestOperator{
		values: values,
	}
}

func (op *TestOperator) Init() {
	op.i = 0
}

func (op *TestOperator) Next() []interface{} {
	res := make([]interface{}, 0)
	if op.i < len(op.values) {
		sc := vector.NewIntVector(op.values[op.i], []uint64{})
		sc.SetLen(len(op.values[op.i]))
		op.i++
		return append(res, sc)
	}
	return nil
}

func setUpTop() MultiVectorOperator {
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

func setUpTopN(n int) MultiVectorOperator {
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
	k := r[0].(*vector.IntVector).Values()
	v := r[1].(*vector.IntVector).Values()

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
	k := r[0].(*vector.IntVector).Values()
	v := r[1].(*vector.IntVector).Values()

	a.Len(k, 5, "length is wrong ")
	a.Len(v, 5, "length is wrong ")

	c := 10
	for i := 0; i < len(k); i++ {
		a.Equal(k[i], c)
		a.Equal(v[i], c)
		c++
	}

}

//TODO: make a benchmark more complex for timsort to check how it goes
func BenchmarkTop(b *testing.B) {
	list := setUpTopN(b.N)

	op := NewTopOperator(list, 5, false)

	b.ResetTimer()
	op.Next()
}
