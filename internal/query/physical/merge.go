package physical

import (
	"container/heap"
	"fmt"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
)

type MergeOp struct {
	heap       MinHeap
	input      []BatchOperator
	lists      []Batch
	idxInList  []int
	initilized bool
}

// An Item is something we manage in a value queue.
type Item struct {
	inputIdx int
	value    int64 // the time value
	// The index is needed by update and is maintained by the heap.Interface methods.
	//index int // The index of the item in the heap.
}

func (i Item) String() string {
	return fmt.Sprintf(" value: %d , inputIdx: %d ", i.value, i.inputIdx)
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
		inputIdx: item.inputIdx,
		value:    item.value,
	})
}

func NewMergeOp(input []BatchOperator) *MergeOp {
	return &MergeOp{
		input:      input,
		idxInList:  make([]int, len(input)),
		lists:      make([]Batch, len(input)),
		heap:       nil,
		initilized: false,
	}

}

func (m *MergeOp) Init() {
	// TODO(gvelo) Init() it is ok to call init multiple times so we need to
	// track initialization state.
	if !m.initilized {
		for _, operator := range m.input {
			operator.Init()
		}
		m.initHeap()
	}
	m.initilized = true
}

func (m *MergeOp) initHeap() {

	k := len(m.input)
	m.heap = make(MinHeap, 0)
	m.idxInList = make([]int, k)
	for i := range m.idxInList {
		m.idxInList[i] = 0
	}

	for i, it := range m.input {
		b := it.Next()
		m.lists[i] = b
		// push all first items.
		if b.Len > 0 {
			//TODO: It's ok to use "ts" ?
			ts := b.Columns["ts"].Vec.(*vector.Int64Vector).Values()
			heap.Push(&m.heap, &Item{
				inputIdx: i,
				value:    ts[0],
			})

			m.idxInList[i] = 0
		} else {
			i--
		}
	}
}

func (m *MergeOp) Close() {
	for _, operator := range m.input {
		operator.Close()
	}
	m.initilized = false
}

func (m *MergeOp) Next() Batch {

	// Take the batch with more columns.
	maxIdx := 0
	max := 0
	for x, it := range m.lists {
		if len(it.Columns) > max {
			maxIdx = x
			max = len(it.Columns)
		}
	}

	// Build new dest batch
	res := m.lists[maxIdx].Clone()

	// while we have items in the heap
	for m.heap.Len() > 0 {

		for i := 0; i < res.Len; i++ {

			it := m.heap.Pop().(*Item)
			println("inputIdx ", it.inputIdx, "value ", it.value, "m.idxInList[it.inputIdx] ", m.idxInList[it.inputIdx])

			m.copyRow(m.lists[it.inputIdx].Columns, m.idxInList[it.inputIdx], res.Columns)
			m.idxInList[it.inputIdx]++

			if m.idxInList[it.inputIdx] < res.Len {
				// TODO: Check the ts index.
				nv := m.lists[it.inputIdx].Columns["ts"].Vec.(*vector.Int64Vector).Get(m.idxInList[it.inputIdx])

				nt := &Item{
					inputIdx: it.inputIdx,
					value:    nv,
				}

				m.heap.Push(nt)
			}

		}
		return res

	}

	return Batch{
		Len: 0,
	}

}

func (m *MergeOp) copyValues(cs Col, si int, cd Col) {

	if cs.Vec.HasNulls() && !cs.Vec.IsValid(si) {
		cd.Vec.AppendNull()
		return
	}

	switch cs.ColumnType {
	case storage.ColumnType_BOOL:
		v := cs.Vec.(*vector.BoolVector).Get(si)
		cd.Vec.(*vector.BoolVector).AppendBool(v)
		break
	case storage.ColumnType_DATETIME, storage.ColumnType_INT64, storage.ColumnType_TIMESTAMP:
		v := cs.Vec.(*vector.Int64Vector).Get(si)
		cd.Vec.(*vector.Int64Vector).AppendInt64(v)
		break
	case storage.ColumnType_INT32:
		v := cs.Vec.(*vector.Int32Vector).Get(si)
		cd.Vec.(*vector.Int32Vector).AppendInt32(v)
		break
	case storage.ColumnType_FLOAT64:
		v := cs.Vec.(*vector.Float64Vector).Get(si)
		cd.Vec.(*vector.Float64Vector).AppendFloat64(v)
		break
	case storage.ColumnType_STRING:
		v := cs.Vec.(*vector.ByteSliceVector).Get(si)
		cd.Vec.(*vector.ByteSliceVector).AppendSlice(v)
		break
	default:
		panic("No vector found.")
	}

}

func (m *MergeOp) copyRow(src map[string]Col, srcIdx int, dst map[string]Col) {
	for key, cd := range dst {
		cs, ok := src[key]
		if ok {
			m.copyValues(cs, srcIdx, cd)
		} else {
			// if the vector doesn't handle nulls, we need to create the valid vector.
			cd.Vec.AppendNull()
		}
	}
}

func (m *MergeOp) Accept(v Visitor) {
	for i, operator := range m.input {
		m.input[i] = Walk(operator, v).(BatchOperator)
	}
}
