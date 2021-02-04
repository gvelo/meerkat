package physical

import (
	"github.com/stretchr/testify/assert"
	"math"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
	"testing"
)

func newBatchOp() []BatchOperator {

	operators := make([]BatchOperator, 0)

	vectors1 := make([][]vector.Vector, 0)
	cols1 := make([]vector.Vector, 0)
	types := make([]storage.ColumnType, 0)
	names := make([]string, 0)

	v := vector.NewInt64Vector([]int64{1, 3, 5, 7, 9, 11, 13, 14, 15, 150, 200, 250, 300, 350}, nil)
	cols1 = append(cols1, &v)
	names = append(names, "ts")
	types = append(types, storage.ColumnType_TIMESTAMP)
	vectors1 = append(vectors1, cols1)

	bo := CreateBatchOp(vectors1, types, names)
	operators = append(operators, bo)

	vectors1 = make([][]vector.Vector, 0)
	cols1 = make([]vector.Vector, 0)
	types = make([]storage.ColumnType, 0)
	names = make([]string, 0)

	v = vector.NewInt64Vector([]int64{1, 3, 5, 7, 9, 11, 13, 14, 15, 150, 200, 250, 300, 350}, nil)
	cols1 = append(cols1, &v)
	names = append(names, "ts")
	types = append(types, storage.ColumnType_TIMESTAMP)
	vectors1 = append(vectors1, cols1)

	bo2 := CreateBatchOp(vectors1, types, names)

	operators = append(operators, bo2)

	return operators
}

func TestMerge(t *testing.T) {
	a := assert.New(t)
	mergeOp := NewMergeOp(newBatchOp())
	mergeOp.Init()
	compare := make([]int64, 2)
	compare[0] = math.MinInt64

	for b := mergeOp.Next(); b.Len != 0; b = mergeOp.Next() {

		for i := 0; i < b.Len; i++ {
			compare[1] = b.Columns["ts"].Vec.(*vector.Int64Vector).Get(i)
			t.Log("comparing compare[0] =  ", compare[0], "compare[1] =", compare[1])
			if compare[0] > compare[1] {
				a.Fail("Vector is not ")
			}
			compare[0] = compare[1]
		}

	}

}
