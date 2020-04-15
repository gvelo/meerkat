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
	"meerkat/internal/schema"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
	"testing"
	"time"
)

func setUpIntBucket() MultiVectorOperator {
	l := make([][]int, 0)
	l = append(l, appendNTimes(10, 10))
	l = append(l, appendNTimes(11, 10))
	l = append(l, appendNTimes(12, 10))
	l = append(l, appendNTimes(13, 10))
	l = append(l, appendNTimes(14, 10))
	l = append(l, appendNTimes(15, 10))
	l = append(l, appendNTimes(16, 10))
	l = append(l, appendNTimes(17, 10))
	l = append(l, appendNTimes(18, 10))
	l = append(l, appendNTimes(19, 10))
	l = append(l, appendNTimes(20, 10))
	l = append(l, appendNTimes(21, 10))
	return NewTestOperator(l)
}

func setUpTimeBucket(t0, n int, d time.Duration) MultiVectorOperator {
	l := make([][]int, 0)
	l = append(l, appendNTimes(t0+int(10*d), n))
	l = append(l, appendNTimes(t0+int(11*d), n))
	l = append(l, appendNTimes(t0+int(12*d), n))
	l = append(l, appendNTimes(t0+int(13*d), n))
	l = append(l, appendNTimes(t0+int(14*d), n))
	l = append(l, appendNTimes(t0+int(15*d), n))
	l = append(l, appendNTimes(t0+int(16*d), n))
	l = append(l, appendNTimes(t0+int(17*d), n))
	l = append(l, appendNTimes(t0+int(18*d), n))
	l = append(l, appendNTimes(t0+int(19*d), n))
	l = append(l, appendNTimes(t0+int(20*d), n))
	return NewTestOperator(l)
}

func setUpTimeBucketParsing(d [][]string) MultiVectorOperator {
	l := make([][]int, 0)

	for _, items := range d {
		list := make([]int, 0)
		for _, item := range items {
			t, _ := time.Parse("02-01-2006 15:04:05 -07:00", item)
			list = append(list, int(t.UnixNano()))
		}
		l = append(l, list)
	}

	return NewTestOperator(l)
}

func appendNTimes(n int, times int) []int {
	l := make([]int, 0, 10)
	for i := 0; i < times; i++ {
		l = append(l, n)
	}
	return l
}

func TestTimeBucketInSecs(t *testing.T) {
	a := assert.New(t)
	t0, err := time.Parse("02-01-2006 15:04:05 -07:00", "01-01-2020 20:34:00 -03:00")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// Set up cf
	sMap := make(map[string]storage.Column)
	sMap["c1"] = &fakeFloatColumn{}
	sMap["c2"] = &fakeIntColumn{}
	sMap["c3"] = &fakeStringColumn{}

	fs := &fakeColFinder{sMap: sMap}

	fields := []field{
		{
			name:     "_ts",
			t:        schema.FieldType_TIMESTAMP,
			nullable: false,
		},
	}

	ii := createIndexInfo("Logs", fields...)
	// Create ctx
	ctx := NewContext(fs, ii)

	list := setUpTimeBucket(int(t0.UnixNano()), 10, time.Second)
	op := NewTimeBucketOperator(ctx, list, "5s")

	m := make(map[int]int)

	for r := op.Next(); r != nil; r = op.Next() {
		s := r[0].(vector.IntVector)
		k := s.Values()
		for i := 0; i < len(k); i++ {
			m[k[i]]++
		}
	}

	keys := make([]int, 0)
	for k, _ := range m {
		t.Log(time.Unix(0, int64(k)))
		keys = append(keys, k)
	}

	a.Equal(30, m[keys[0]])
	a.Equal(50, m[keys[1]])
	a.Equal(30, m[keys[2]])

}

func TestTimeBucketInHour(t *testing.T) {
	// a := assert.New(t)
	dates := [][]string{{"01-01-2020 20:34:00 -03:00", "01-01-2020 20:34:00 -03:00", "01-01-2020 20:34:00 -03:00",
		"01-01-2020 20:34:00 -03:00", "01-01-2020 20:34:00 -03:00", "01-01-2020 20:34:00 -03:00"}}
	list := setUpTimeBucketParsing(dates)

	// Set up cf
	sMap := make(map[string]storage.Column)
	sMap["_ts"] = &fakeIntColumn{}

	fs := &fakeColFinder{sMap: sMap}

	fields := []field{
		{
			name:     "_ts",
			t:        schema.FieldType_TIMESTAMP,
			nullable: false,
		},
	}

	ii := createIndexInfo("Logs", fields...)
	// Create ctx
	ctx := NewContext(fs, ii)

	op := NewTimeBucketOperator(ctx, list, "1h")

	m := make(map[int]int)

	for r := op.Next(); r != nil; r = op.Next() {
		s := r[0].(vector.IntVector)
		k := s.Values()
		for i := 0; i < len(k); i++ {
			m[k[i]]++
		}
	}

	keys := make([]int, 0)
	for k, _ := range m {
		t.Logf("Unix Time %v", time.Unix(0, int64(k)))
		keys = append(keys, k)
	}

}

func TestIntBucket(t *testing.T) {
	list := setUpIntBucket()
	a := assert.New(t)

	// Set up cf
	sMap := make(map[string]storage.Column)
	sMap["_ts"] = &fakeIntColumn{}

	fs := &fakeColFinder{sMap: sMap}

	fields := []field{
		{
			name:     "_ts",
			t:        schema.FieldType_TIMESTAMP,
			nullable: false,
		},
	}

	ii := createIndexInfo("Logs", fields...)
	// Create ctx
	ctx := NewContext(fs, ii)

	op := NewBucketOperator(ctx, list, 5)

	m := make(map[int]int)

	for r := op.Next(); r != nil; r = op.Next() {
		s := r[0].(vector.IntVector)
		k := s.Values()
		for i := 0; i < len(k); i++ {
			m[k[i]]++
		}
	}

	a.Equal(30, m[10]) // 10,11,12
	a.Equal(50, m[15]) // 13,14,15,16,17
	a.Equal(40, m[20]) // 18,19,20,21

}

func TestIntBucketList(t *testing.T) {
	a := assert.New(t)

	list := []int{95, 97, 98, 100, 102, 103, 105, 109, 110}
	e := []int{95, 95, 100, 100, 100, 105, 105, 110, 110}
	r := make([]int, 0, len(list))
	for _, i := range list {
		r = append(r, getNextSpan(i, 5, 0))
	}

	for i, _ := range list {
		a.Equal(r[i], e[i], "Error %d, is not equal than %d ", r[i], e[i])
	}
}
