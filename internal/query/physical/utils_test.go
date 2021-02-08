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

package physical

import (
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
)

type BatchOperatorTest struct {
	batches []Batch
	i       int
}

func NewBatchOperatorTest() *BatchOperatorTest {
	return &BatchOperatorTest{}
}

func (b *BatchOperatorTest) Init() {
	b.i = 0
}

func (b *BatchOperatorTest) Close() {
	panic("implement me")
}

func (b *BatchOperatorTest) Accept(v Visitor) {
	panic("implement me")
}

func (b *BatchOperatorTest) Next() Batch {
	if b.i < len(b.batches) {
		batch := b.batches[b.i]
		b.i++
		return batch
	}
	return Batch{Len: 0}
}

func (b *BatchOperatorTest) AddBatch() {
	b.batches = append(b.batches, NewBatch())
}

func (b *BatchOperatorTest) AddColumn(vector vector.Vector, t storage.ColumnType, name string, g int64) {
	if len(b.batches) == 0 {
		panic("batch is null, you should use AddBatch, before call AddColumn")
	}
	cb := int64(len(b.batches) - 1)
	b.batches[cb].Columns[name] = Col{
		Group:      g,
		Order:      cb,
		Vec:        vector,
		ColumnType: t,
	}
	b.batches[cb].Len = vector.Len()
}
