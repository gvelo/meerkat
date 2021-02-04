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

type BatchOperatorDummy struct {
	b []Batch
	i int
}

func (b BatchOperatorDummy) Init() {
	b.i = 0
}

func (b BatchOperatorDummy) Close() {
	panic("implement me")
}

func (b BatchOperatorDummy) Accept(v Visitor) {
	panic("implement me")
}

func (b BatchOperatorDummy) Next() Batch {
	if b.i < len(b.b) {
		batch := b.b[b.i]
		b.i++
		return batch
	}
	return Batch{Len: 0}
}

func CreateBatchOp(vector [][]vector.Vector, types []storage.ColumnType, names []string) BatchOperator {

	b := make([]Batch, 0)
	for i := 0; i < len(vector); i++ {

		b = append(b, NewBatch())

		for j := 0; j < len(vector[i]); j++ {

			b[i].Columns[names[j]] = Col{
				Group:      0,
				Order:      1,
				Vec:        vector[i][j],
				ColumnType: types[j],
			}

		}
		b[i].Len = vector[i][0].Len()
	}
	bo := &BatchOperatorDummy{b, 0}
	return bo
}
