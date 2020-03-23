package executor

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage"
	encoding2 "meerkat/internal/storage/encoding"
	"meerkat/internal/storage/vector"
	"testing"
	"time"
)

// FAKES

type fakeMultiVectorOperator struct {
	vec []vector.Vector
	idx int
}

func (f *fakeMultiVectorOperator) Init() {
	f.idx = 0
}

func (f *fakeMultiVectorOperator) Destroy() {
	// nothing to do
}

func (f *fakeMultiVectorOperator) Next() []vector.Vector {
	if f.idx == 1 {
		return nil
	}
	f.idx++
	return f.vec
}

type fakeSegment struct {
	sMap map[string]storage.Column
}

func (f *fakeSegment) IndexName() string {
	panic("implement me")
}

func (f *fakeSegment) IndexID() []byte {
	panic("implement me")
}

func (f *fakeSegment) From() time.Time {
	panic("implement me")
}

func (f *fakeSegment) To() time.Time {
	panic("implement me")
}

func (f *fakeSegment) Rows() int {
	panic("implement me")
}

func (f *fakeSegment) Col(id []byte) storage.Column {
	return f.sMap[string(id)]
}

type fakeFloatColumn struct {
}

func (f *fakeFloatColumn) Encoding() encoding2.EncodingType {
	panic("implement me")
}

func (f *fakeFloatColumn) Validity() *roaring.Bitmap {
	panic("implement me")
}

func (f *fakeFloatColumn) HasNulls() bool {
	panic("implement me")
}

func (f *fakeFloatColumn) Stats() *storage.Stats {
	panic("implement me")
}

func (f *fakeFloatColumn) Dict() storage.FloatDict {
	panic("implement me")
}

func (f *fakeFloatColumn) Index() storage.FloatIndex {
	panic("implement me")
}

func (f *fakeFloatColumn) Read(pos []uint32) vector.FloatVector {
	panic("implement me")
}

func (f *fakeFloatColumn) Iterator() storage.FloatIterator {
	panic("implement me")
}

type fakeIntColumn struct {
}

func (f *fakeIntColumn) Encoding() encoding2.EncodingType {
	panic("implement me")
}

func (f *fakeIntColumn) Validity() *roaring.Bitmap {
	panic("implement me")
}

func (f *fakeIntColumn) HasNulls() bool {
	panic("implement me")
}

func (f *fakeIntColumn) Stats() *storage.Stats {
	panic("implement me")
}

func (f *fakeIntColumn) Dict() storage.IntDict {
	panic("implement me")
}

func (f *fakeIntColumn) Index() storage.IntIndex {
	panic("implement me")
}

func (f *fakeIntColumn) Read(pos []uint32) vector.IntVector {
	panic("implement me")
}

func (f *fakeIntColumn) Iterator() storage.IntIterator {
	panic("implement me")
}

type fakeStringColumn struct {
}

func (f *fakeStringColumn) Encoding() encoding2.EncodingType {
	panic("implement me")
}

func (f *fakeStringColumn) Validity() *roaring.Bitmap {
	panic("implement me")
}

func (f *fakeStringColumn) HasNulls() bool {
	panic("implement me")
}

func (f *fakeStringColumn) Stats() *storage.Stats {
	panic("implement me")
}

func (f *fakeStringColumn) Dict() storage.ByteSliceDict {
	panic("implement me")
}

func (f *fakeStringColumn) Index() storage.ByteSliceIndex {
	panic("implement me")
}

func (f *fakeStringColumn) ReadDictEnc(pos []uint32) vector.IntVector {
	panic("implement me")
}

func (f *fakeStringColumn) Read(pos []uint32) vector.ByteSliceVector {
	panic("implement me")
}

func (f *fakeStringColumn) DictEncodedIterator() storage.IntIterator {
	panic("implement me")
}

func (f *fakeStringColumn) Iterator() storage.ByteSliceIterator {
	panic("implement me")
}

func TestHAggScenario1(t *testing.T) {

	a := assert.New(t)

	// Set up child
	vec := make([]vector.Vector, 0)
	v := vector.NewFloatVector([]float64{1.2, 1.2, 1.4, 1.5, 1.6}, []uint64{})
	vec = append(vec, &v)
	v1 := vector.NewIntVector([]int{2, 2, 4, 5, 6}, []uint64{})
	vec = append(vec, &v1)
	v2 := vector.NewByteSliceVector([]byte("123123123123123"), []uint64{}, []int{3, 6, 9, 12, 15})
	vec = append(vec, &v2)
	f := &fakeMultiVectorOperator{
		vec: vec,
	}

	// Set up segment
	sMap := make(map[string]storage.Column)
	sMap["c1"] = &fakeFloatColumn{}
	sMap["c2"] = &fakeIntColumn{}
	sMap["c3"] = &fakeStringColumn{}

	fs := &fakeSegment{sMap: sMap}

	ag := []Aggregation{{
		AggType: Sum,
		AggCol:  0,
	}, {
		AggType: Max,
		AggCol:  0,
	}}

	g := []int{2}

	// Create ctx
	ctx := NewContext(fs)
	m := make(map[int][]byte)
	m[0] = []byte("c1")
	m[1] = []byte("c2")
	m[2] = []byte("c3")
	ctx.Value(ColumnIndexToColumnName, m)

	op := NewHashAggregateOperator(ctx, f, ag, g)

	start := time.Now()

	op.Init()

	r := op.Next()

	elapsed := time.Since(start)
	t.Logf(" took %s", elapsed)

	a.NotNil(r, "This should not be nil")

	a.Equal("123", string(r[0].(*vector.ByteSliceVector).Get(0)))
	a.Equal(6.9, r[1].(*vector.FloatVector).Values()[0])
	a.Equal(1.6, r[2].(*vector.FloatVector).Values()[0])

}

func multiplyIntVector(v []int, n []uint64, times int) (rv []int, rn []uint64) {
	rv = make([]int, 0)

	for i := 0; i < times; i++ {
		rv = append(rv, v...)
	}

	if len(rn) > 0 {
		rn = make([]uint64, len(n)*times)

		for i := 0; i <= times; i++ {
			rn = append(rn, n...)
		}
	}
	return
}

func multiplyFloatVector(v []float64, n []uint64, times int) (rv []float64, rn []uint64) {
	rv = make([]float64, 0)

	for i := 0; i < times; i++ {
		rv = append(rv, v...)
	}

	if len(rn) > 0 {
		rn = make([]uint64, len(n)*times)

		for i := 0; i < times; i++ {
			rn = append(rn, n...)
		}
	}
	return
}

func multiplyBsVector(v []byte, n []uint64, o []int, times int) (rv []byte, rn []uint64, ro []int) {
	rv = make([]byte, 0)
	for i := 0; i < times; i++ {
		rv = append(rv, v...)
	}
	ro = make([]int, 0)
	ant := 0

	for i := 0; i < times; i++ {
		for _, it := range o {
			ro = append(ro, ant+it)
		}
		ant = ro[len(ro)-1]
	}

	if len(rn) > 0 {
		rn = make([]uint64, len(n)*times)

		for i := 0; i < times; i++ {
			rn = append(rn, n...)
		}
	}
	return
}

func TestHAggScenario2(t *testing.T) {

	a := assert.New(t)
	times := 1000000
	// Set up child
	vec := make([]vector.Vector, 0)
	rv, rn := multiplyFloatVector([]float64{1.2, 1.2, 1.4, 1.5, 1.6}, []uint64{}, times)
	v := vector.NewFloatVector(rv, rn)
	vec = append(vec, &v)
	rv1, rn1 := multiplyIntVector([]int{2, 2, 4, 5, 6}, []uint64{}, times)
	v1 := vector.NewIntVector(rv1, rn1)
	vec = append(vec, &v1)
	rv2, rn2, ro2 := multiplyBsVector([]byte("1123123123123"), []uint64{}, []int{1, 4, 7, 10, 13}, times)
	v2 := vector.NewByteSliceVector(rv2, rn2, ro2)
	vec = append(vec, &v2)
	f := &fakeMultiVectorOperator{
		vec: vec,
	}

	// Set up segment
	sMap := make(map[string]storage.Column)
	sMap["c1"] = &fakeFloatColumn{}
	sMap["c2"] = &fakeIntColumn{}
	sMap["c3"] = &fakeStringColumn{}

	fs := &fakeSegment{sMap: sMap}

	ag := []Aggregation{{
		AggType: Sum,
		AggCol:  0,
	}, {
		AggType: Max,
		AggCol:  0,
	}}

	g := []int{2}

	// Create ctx
	ctx := NewContext(fs)
	m := make(map[int][]byte)
	m[0] = []byte("c1")
	m[1] = []byte("c2")
	m[2] = []byte("c3")
	ctx.Value(ColumnIndexToColumnName, m)

	op := NewHashAggregateOperator(ctx, f, ag, g)
	start := time.Now()

	op.Init()

	r := op.Next()

	elapsed := time.Since(start)
	t.Logf(" took %s", elapsed)

	a.NotNil(r, "This should not be nil")

	a.Equal("1", string(r[0].(*vector.ByteSliceVector).Get(0)))
	a.InDelta(1200000, r[1].(*vector.FloatVector).Values()[0], 0.1)
	a.InDelta(1.2, r[2].(*vector.FloatVector).Values()[0], 0.1)

	a.Equal("123", string(r[0].(*vector.ByteSliceVector).Get(1)))
	a.InDelta(5700000, r[1].(*vector.FloatVector).Values()[1], 0.1)
	a.InDelta(1.6, r[2].(*vector.FloatVector).Values()[1], 0.1)

}

func TestSortScenario(t *testing.T) {

	a := assert.New(t)
	// Set up child
	vec := make([]vector.Vector, 0)
	v := vector.NewFloatVector([]float64{1.2, 1.2, 1.4, 1.5, 1.6}, []uint64{})
	vec = append(vec, &v)

	v1 := vector.NewIntVector([]int{2, 2, 4, 5, 6}, []uint64{})
	vec = append(vec, &v1)

	v2 := vector.NewByteSliceVector([]byte("1123123123123"), []uint64{}, []int{1, 4, 7, 10, 13})
	vec = append(vec, &v2)
	f := &fakeMultiVectorOperator{
		vec: vec,
	}

	// Set up segment
	sMap := make(map[string]storage.Column)
	sMap["c1"] = &fakeFloatColumn{}
	sMap["c2"] = &fakeIntColumn{}
	sMap["c3"] = &fakeStringColumn{}

	fs := &fakeSegment{sMap: sMap}

	ag := []Aggregation{{
		AggType: Sum,
		AggCol:  0,
	}, {
		AggType: Max,
		AggCol:  0,
	}}

	g := []int{2}

	// Create ctx
	ctx := NewContext(fs)
	m := make(map[int][]byte)
	m[0] = []byte("c1")
	m[1] = []byte("c2")
	m[2] = []byte("c3")
	ctx.Value(ColumnIndexToColumnName, m)

	op := NewSortedAggregateOperator(ctx, f, ag, g)
	start := time.Now()

	op.Init()

	r := op.Next()

	elapsed := time.Since(start)
	t.Logf(" took %s", elapsed)

	a.NotNil(r, "This should not be nil")

	a.Equal("1", string(r[0].(*vector.ByteSliceVector).Get(0)))
	a.InDelta(1.2, r[1].(*vector.FloatVector).Values()[0], 0.1)
	a.InDelta(1.2, r[2].(*vector.FloatVector).Values()[0], 0.1)

	a.Equal("123", string(r[0].(*vector.ByteSliceVector).Get(1)))
	a.InDelta(5.7, r[1].(*vector.FloatVector).Values()[1], 0.1)
	a.InDelta(1.6, r[2].(*vector.FloatVector).Values()[1], 0.1)

}

func TestSortScenario2(t *testing.T) {

	a := assert.New(t)
	// Set up child
	vec := make([]vector.Vector, 0)
	v := vector.NewFloatVector([]float64{1.2, 1.2, 1.4, 1.5, 1.6}, []uint64{})
	vec = append(vec, &v)

	v1 := vector.NewIntVector([]int{2, 2, 4, 5, 6}, []uint64{})
	vec = append(vec, &v1)

	v2 := vector.NewByteSliceVector([]byte("1123123123123"), []uint64{}, []int{1, 4, 7, 10, 13})
	vec = append(vec, &v2)
	f := &fakeMultiVectorOperator{
		vec: vec,
	}

	// Set up segment
	sMap := make(map[string]storage.Column)
	sMap["c1"] = &fakeFloatColumn{}
	sMap["c2"] = &fakeIntColumn{}
	sMap["c3"] = &fakeStringColumn{}

	fs := &fakeSegment{sMap: sMap}

	ag := []Aggregation{{
		AggType: Sum,
		AggCol:  0,
	}, {
		AggType: Max,
		AggCol:  0,
	}}

	g := []int{1, 2}

	// Create ctx
	ctx := NewContext(fs)
	m := make(map[int][]byte)
	m[0] = []byte("c1")
	m[1] = []byte("c2")
	m[2] = []byte("c3")
	ctx.Value(ColumnIndexToColumnName, m)

	op := NewSortedAggregateOperator(ctx, f, ag, g)
	start := time.Now()

	op.Init()

	r := op.Next()

	elapsed := time.Since(start)
	t.Logf(" took %s", elapsed)

	a.NotNil(r, "This should not be nil")

	// k1
	a.Equal(2, r[0].(*vector.IntVector).Values()[0])
	a.Equal(2, r[0].(*vector.IntVector).Values()[1])
	a.Equal(4, r[0].(*vector.IntVector).Values()[2])
	a.Equal(5, r[0].(*vector.IntVector).Values()[3])
	a.Equal(6, r[0].(*vector.IntVector).Values()[4])

	// k2
	a.Equal("1", string(r[1].(*vector.ByteSliceVector).Get(0)))
	a.Equal("123", string(r[1].(*vector.ByteSliceVector).Get(1)))
	a.Equal("123", string(r[1].(*vector.ByteSliceVector).Get(2)))
	a.Equal("123", string(r[1].(*vector.ByteSliceVector).Get(3)))
	a.Equal("123", string(r[1].(*vector.ByteSliceVector).Get(4)))

	a.InDelta(1.2, r[2].(*vector.FloatVector).Values()[0], 0.1)
	a.InDelta(1.2, r[2].(*vector.FloatVector).Values()[1], 0.1)
	a.InDelta(1.4, r[2].(*vector.FloatVector).Values()[2], 0.1)
	a.InDelta(1.5, r[2].(*vector.FloatVector).Values()[3], 0.1)
	a.InDelta(1.6, r[2].(*vector.FloatVector).Values()[4], 0.1)

	a.InDelta(1.2, r[3].(*vector.FloatVector).Values()[0], 0.1)
	a.InDelta(1.2, r[3].(*vector.FloatVector).Values()[1], 0.1)
	a.InDelta(1.4, r[3].(*vector.FloatVector).Values()[2], 0.1)
	a.InDelta(1.5, r[3].(*vector.FloatVector).Values()[3], 0.1)
	a.InDelta(1.6, r[3].(*vector.FloatVector).Values()[4], 0.1)

}