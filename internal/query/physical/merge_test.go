package physical

import (
	"github.com/stretchr/testify/assert"
	"math"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
	"testing"
)

// TODO REFACTOR TESTS

func createTSTestData() []BatchOperator {

	operators := make([]BatchOperator, 0)

	bo := NewBatchOperatorTest()
	v := vector.NewInt64Vector([]int64{1, 3, 5, 7, 9, 11, 13, 14, 15, 150, 200, 250, 300, 350}, nil)
	bo.AddColumn(&v, storage.ColumnType_TIMESTAMP, "ts", 0)
	operators = append(operators, bo)

	bo = NewBatchOperatorTest()
	v = vector.NewInt64Vector([]int64{1, 3, 5, 7, 9, 11, 13, 14, 15, 150, 200, 250, 300, 350}, nil)
	bo.AddColumn(&v, storage.ColumnType_TIMESTAMP, "ts", 0)

	operators = append(operators, bo)

	return operators
}

func TestMergeTS(t *testing.T) {
	a := assert.New(t)
	mergeOp := NewMergeOp(createTSTestData())
	mergeOp.Init()
	compare := make([]int64, 2)
	compare[0] = math.MinInt64

	for b := mergeOp.Next(); b.Len != 0; b = mergeOp.Next() {

		for i := 0; i < b.Len; i++ {
			compare[1] = b.Columns["ts"].Vec.(*vector.Int64Vector).Get(i)
			if compare[0] > compare[1] {
				a.Fail("Vector is not ordered")
			}
			compare[0] = compare[1]
		}

	}

}

func createColumnsTestData() []BatchOperator {

	operators := make([]BatchOperator, 0)

	bo := NewBatchOperatorTest()
	bo.AddBatch()
	v := vector.NewInt64Vector([]int64{1, 3, 5, 7, 9, 11, 13, 14, 15, 150, 200, 250, 300, 350}, nil)
	bo.AddColumn(&v, storage.ColumnType_TIMESTAMP, "ts", 0)

	v2 := vector.NewInt64Vector([]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil)
	bo.AddColumn(&v2, storage.ColumnType_INT64, "other", 0)

	operators = append(operators, bo)

	bo = NewBatchOperatorTest()
	bo.AddBatch()
	v3 := vector.NewInt64Vector([]int64{1, 3, 5, 7, 9, 11, 13, 14, 15, 150, 200, 250, 300, 350}, nil)
	bo.AddColumn(&v3, storage.ColumnType_TIMESTAMP, "ts", 0)

	v4 := vector.NewInt64Vector([]int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, nil)
	bo.AddColumn(&v4, storage.ColumnType_INT64, "other", 0)

	operators = append(operators, bo)

	return operators
}

func TestMergeColumns(t *testing.T) {

	a := assert.New(t)
	mergeOp := NewMergeOp(createColumnsTestData())
	mergeOp.Init()
	compareTS := make([]int64, 2)
	compareTS[0] = math.MinInt64
	expect := int64(0)
	for b := mergeOp.Next(); b.Len != 0; b = mergeOp.Next() {

		for i := 0; i < b.Len; i++ {
			expect = int64(i % 2)
			compareTS[1] = b.Columns["ts"].Vec.(*vector.Int64Vector).Get(i)
			o := b.Columns["other"].Vec.(*vector.Int64Vector).Get(i)

			if expect != o {
				a.Fail("Vector is not well merged")
			}

			if compareTS[0] > compareTS[1] {
				a.Fail("Vector is not ordered")
			}
			compareTS[0] = compareTS[1]
		}

	}

}
