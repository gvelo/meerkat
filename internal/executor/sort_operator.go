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
)

func NewSortOperator(ctx Context, child MultiVectorOperator, keyNames []string) *SortOperator {
	return &SortOperator{
		ctx:      ctx,
		child:    child,
		keyNames: keyNames,
		log:      log.With().Str("src", "SortOperator").Logger(),
	}
}

// SortOperator
type SortOperator struct {
	ctx       Context
	child     MultiVectorOperator // (Positions to review)
	keyNames  []string
	sz        int
	sorted    [][]byte
	processed int
	log       zerolog.Logger
}

func (op *SortOperator) Init() {
	op.child.Init()

	colIds := make([]int, 0, len(op.keyNames))
	for _, it := range op.keyNames {
		_, id, err := op.ctx.GetFieldProcessed().FindField(it)
		if err != nil {
			log.Error().Err(err)
		}

		colIds = append(colIds, id)
	}

	n := op.child.Next()
	i := 0
	var k []byte

	// que hago agrando los vectores?

	keys := make([][]byte, 0)
	for ; n != nil; n = op.child.Next() {
		l1 := getLen(n[0])
		// iterate over all "rows"
		for ; i < l1; i++ {
			k = createKey(n, colIds, i)
			keys = append(keys, k)
		}
	}

	// Do the sort.

	op.sorted = keys
}

func (op *SortOperator) Destroy() {
	op.child.Destroy()
}

func (op *SortOperator) Next() []interface{} {

	// tngo que cerear un objeto key , idx ?
	if op.processed == len(op.sorted) {
		return nil
	}

	return nil
}
