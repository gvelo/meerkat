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
	"container/heap"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage/vector"
)

func NewMergeOperator(ctx Context, children []MultiVectorOperator, result []interface{}) *MergeOperator {
	return &MergeOperator{
		ctx:      ctx,
		children: children,
		result:   result,
		log:      log.With().Str("src", "MergeOperator").Logger(),
	}
}

// An Item is something we manage in a value queue.
type Item struct {
	listFrom int
	value    int64 // the time value
	// The index is needed by update and is maintained by the heap.Interface methods.
	//index int // The index of the item in the heap.
}

func (i Item) String() string {
	return fmt.Sprintf(" value: %d , listFrom: %d ", i.value, i.listFrom)
}

// A MinHeap implements heap.Interface and holds Items.
type MinHeap []*Item

func (h MinHeap) Len() int {
	return len(h)
}
func (h MinHeap) Less(i, j int) bool {
	return h[i].value < h[j].value
}
func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x interface{}) {
	//n := len(*h)
	//item := x.(*Item)
	//item.index = n
	*h = append(*h, x.(*Item))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[0]
	old[0] = nil // avoid memory leak
	//item.index = -1 // for safety
	*h = old[1:n]
	return item
}

func (h *MinHeap) Peek() interface{} {
	old := *h
	item := old[0]
	return item
}

func (h *MinHeap) update(item *Item, value int64) {
	item.value = value
	heap.Push(h, &Item{
		listFrom: item.listFrom,
		value:    item.value,
	})
}

// MergeOperator
type MergeOperator struct {
	ctx       Context
	children  []MultiVectorOperator // (Positions to review)
	batchProc int
	result    []interface{}
	heap      MinHeap
	idxList   []int
	lists     [][]interface{}
	log       zerolog.Logger
}

func (op *MergeOperator) Init() {

	for _, it := range op.children {
		it.Init()
	}

	op.initHeap()

}

func (op *MergeOperator) initHeap() {

	k := len(op.children)

	op.heap = make(MinHeap, 0)
	op.idxList = make([]int, k)
	op.lists = make([][]interface{}, k)

	for i, it := range op.children {
		op.lists[i] = it.Next()
		// push all first items.
		if op.lists[i] != nil {
			vv := op.lists[i][TsIndex].(vector.Int64Vector)
			vp := &vv
			ts := vp.Values()

			heap.Push(&op.heap, &Item{
				listFrom: i,
				value:    ts[0],
			})

			op.idxList[i] = 0
		} else {
			i--
		}
	}
}

func (op *MergeOperator) Destroy() {

}

func (op *MergeOperator) Next() []interface{} {

	for i := 0; i < len(op.children); i++ {
		op.result[i] = reset(op.result[i])
	}

	if len(op.lists[TsIndex]) == 0 {
		return nil
	}

	// Until we reach the batch
	for i := 0; i < op.ctx.Sz(); i++ {

		it := op.heap.Pop().(*Item)

		// TODO(sebad) Handle different #columns
		for x, vec := range op.lists[it.listFrom] {
			op.result[x] = setItemInResultVector(op.result[x], vec, op.idxList[it.listFrom])
		}

		vv := op.lists[it.listFrom][TsIndex].(vector.Int64Vector)
		vp := &vv

		op.idxList[it.listFrom]++

		if op.idxList[it.listFrom] < vp.Len() {

			ts := vp.Get(op.idxList[it.listFrom])
			// insert newItem
			op.heap.update(it, ts)
		} else {
			if len(op.heap) == 0 {
				op.initHeap()

				if len(op.heap) == 0 {
					return op.result
				}

			}
		}

	}

	// batch.
	return op.result
}

func reset(v interface{}) interface{} {
	switch v.(type) {

	case vector.Int64Vector:

		vv2 := v.(vector.Int64Vector)
		vp2 := &vv2
		vp2.SetLen(0)
		return *vp2

	case vector.Float64Vector:

		vv2 := v.(vector.Float64Vector)
		vp2 := &vv2
		vp2.SetLen(0)
		return *vp2

	}

	panic("Item not found")
}

func setItemInResultVector(v interface{}, val interface{}, idx int) interface{} {
	switch v.(type) {

	case vector.Int64Vector:

		vv := val.(vector.Int64Vector)
		vp := &vv
		val := vp.Get(idx)

		vv2 := v.(vector.Int64Vector)
		vp2 := &vv2
		vp2.AppendInt64(val)

		return *vp2

	case vector.Float64Vector:

		vv := val.(vector.Float64Vector)
		vp := &vv
		val := vp.Get(idx)

		vv2 := v.(vector.Float64Vector)
		vp2 := &vv2
		vp2.AppendFloat64(val)

		return *vp2
	}

	panic("Item not found")
}
