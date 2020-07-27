package executor

import (
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage"
	encoding2 "meerkat/internal/storage/encoding"
	"meerkat/internal/storage/vector"
)

// FAKES //TODO(sebad): make code more compact and legible, remove duplicates, use switchs.
type fakeMultiVectorOperator struct {
	vec [][]interface{}
	idx int
}

func (f *fakeMultiVectorOperator) Init() {
	f.idx = 0
}

func (f *fakeMultiVectorOperator) Destroy() {
	// nothing to do
}

func (f *fakeMultiVectorOperator) Next() []interface{} {
	if f.idx == len(f.vec) {
		return nil
	}
	v := f.vec[f.idx]
	f.idx++
	return v
}

type fakeVectorOperator struct {
	vec []interface{}
	idx int
}

func (f *fakeVectorOperator) Init() {
	f.idx = 0
}

func (f *fakeVectorOperator) Destroy() {
	// nothing to do
}

func (f *fakeVectorOperator) Next() interface{} {
	if f.idx == len(f.vec) {
		return nil
	}
	v := f.vec[f.idx]
	f.idx++
	return v
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

func (f *fakeFloatColumn) Index() storage.FloatIndex {
	panic("implement me")
}

func (f *fakeFloatColumn) Reader() storage.Float64ColumnReader {
	panic("implement me")
}

func (f *fakeFloatColumn) Iterator() storage.Float64Iterator {
	panic("implement me")
}

type fakeIntIterator struct {
	v        [][]int64
	i        int
	validity [][]uint64
	length   []int
}

func (f *fakeIntIterator) HasNext() bool {
	return f.i < len(f.v)
}

func (f *fakeIntIterator) Next() vector.Int64Vector {
	var v1 vector.Int64Vector
	if len(f.validity) > 0 {
		v1 = vector.NewInt64Vector(f.v[f.i], f.validity[f.i])
		v1.SetLen(f.length[f.i])
	} else {
		v1 = vector.NewInt64Vector(f.v[f.i], []uint64{})
	}

	f.i++
	return v1
}

func NewFakeIntIterator(values [][]int64, validity [][]uint64, length []int) storage.Int64Iterator {
	return &fakeIntIterator{
		v:        values,
		validity: validity,
		length:   length,
	}
}

func NewFakeColumn(values interface{}) storage.Column {
	in := values.(input)
	switch v := in.values.(type) {
	case [][]int64:
		return &fakeIntColumn{
			v:        v,
			validity: in.validity,
			length:   in.length,
		}
	case [][]string:
		return &fakeStringColumn{
			v:        v,
			validity: in.validity,
			length:   in.length,
		}
	}
	panic("Not implemented")
	return nil
}

type fakeIntColumn struct {
	v        [][]int64
	validity [][]uint64
	length   []int
}

func (f *fakeIntColumn) Encoding() encoding2.EncodingType {
	panic("implement me")
}

func (f *fakeIntColumn) Validity() *roaring.Bitmap {
	panic("implement me")
}

func (f *fakeIntColumn) HasNulls() bool {
	return len(f.validity) > 0
}

func (f *fakeIntColumn) Stats() *storage.Stats {
	panic("implement me")
}

func (f *fakeIntColumn) Index() storage.Int64Index {
	panic("implement me")
}

func (f *fakeIntColumn) Reader() storage.Int64ColumnReader {
	panic("implement me")
}

func (f *fakeIntColumn) Iterator() storage.Int64Iterator {
	return NewFakeIntIterator(f.v, f.validity, f.length)
}

type fakeBinaryIterator struct {
	v        [][]string
	i        int
	validity [][]uint64
	length   []int
}

func (f *fakeBinaryIterator) HasNext() bool {
	return f.i < len(f.v)

}

func (f *fakeBinaryIterator) Next() vector.ByteSliceVector {
	data := make([]byte, 0, 1000)
	offset := make([]int, 0, 1000)
	idx := 0
	for _, y := range f.v[f.i] {
		data = append(data, []byte(y)...)
		idx = idx + len(y)
		offset = append(offset, idx)

	}
	var val []uint64
	if len(f.validity) > 0 {
		val = f.validity[f.i]
	} else {
		val = nil
	}
	v1 := vector.NewByteSliceVector(data, offset, val)
	// v1.SetLen(len(f.offset))
	f.i++
	return v1
}

func NewFakeBinaryIterator(values [][]string, validity [][]uint64, length []int) storage.ByteSliceIterator {
	return &fakeBinaryIterator{
		v:        values,
		i:        0,
		validity: validity,
		length:   length,
	}

}

type fakeStringColumn struct {
	v        [][]string
	validity [][]uint64
	length   []int
}

func (f *fakeStringColumn) Encoding() encoding2.EncodingType {
	panic("implement me")
}

func (f *fakeStringColumn) Validity() *roaring.Bitmap {
	panic("implement me")
}

func (f *fakeStringColumn) HasNulls() bool {
	return len(f.validity) > 0
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

func (f *fakeStringColumn) DictEncReader() storage.Int64ColumnReader {
	panic("implement me")
}

func (f *fakeStringColumn) Reader() storage.ByteSliceReader {
	panic("implement me")
}

func (f *fakeStringColumn) Iterator() storage.ByteSliceIterator {
	return NewFakeBinaryIterator(f.v, f.validity, f.length)
}

func (f *fakeStringColumn) DictEncIterator() storage.Int64Iterator {
	panic("implement me")
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

/*
func multiplyBsVector(v []byte, o []int, n []uint64, times int) (rv []byte, ro []int, rn []uint64) {

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

	if len(n) > 0 {
		rn = make([]uint64, len(n)*times)

		for i := 0; i < times; i++ {
			rn = append(rn, n...)
		}
	}
	return
}

func setUp(vec [][]interface{}, ag []Aggregation, g []int, isHash bool) interface{} {

	f := &fakeMultiVectorOperator{
		vec: vec,
	}

	// Set up cf
	sMap := make(map[string]storage.Column)
	sMap["c1"] = &fakeFloatColumn{}
	sMap["c2"] = &fakeIntColumn{}
	sMap["c3"] = &fakeStringColumn{}

	fs := &fakeColFinder{sMap: sMap}

	// Create ctx // TODO: meter Segmento.
	ctx := NewContext( nil, 100)

	if isHash {
		return NewHashAggregateOperator(ctx, f, ag, g)
	} else {
		return NewSortedAggregateOperator(ctx, f, ag, g)
	}
}

func TestHAggScenario1(t *testing.T) {

	a := assert.New(t)

	// Set up child
	vec := make([][]interface{}, 0)

	lv := make([]interface{}, 0)
	v := vector.NewFloatVector([]float64{1.2, 1.2, 1.4, 1.5, 1.6}, []uint64{})
	v.SetLen(5)
	lv = append(lv, v)

	v1 := vector.NewIntVector([]int{2, 2, 4, 5, 6}, []uint64{})
	v1.SetLen(5)
	lv = append(lv, v1)

	v2 := vector.NewByteSliceVector([]byte("123123123123123"), []int{3, 6, 9, 12, 15}, []uint64{})
	v2.SetLen(5)
	lv = append(lv, v2)
	vec = append(vec, lv)

	ag := []Aggregation{{
		AggType: Sum,
		AggCol:  0,
	}, {
		AggType: Max,
		AggCol:  0,
	}}

	g := []int{2}

	op := setUp(vec, ag, g, true).(*HashAggregateOperator)

	start := time.Now()

	op.Init()

	r := op.Next()

	elapsed := time.Since(start)
	t.Logf(" took %s", elapsed)

	a.NotNil(r, "This should not be nil")

	a.Equal("123", string(r[0].(*vector.ByteSliceVector).Get(0)))
	a.Equal(6.9, r[1].(*vector.FloatVector).Value()[0])
	a.Equal(1.6, r[2].(*vector.FloatVector).Value()[0])

}

func TestHAggScenario2(t *testing.T) {

	a := assert.New(t)
	times := 1000000
	// Set up child
	vec := make([][]interface{}, 0)
	lv := make([]interface{}, 0)

	rv, rn := multiplyFloatVector([]float64{1.2, 1.2, 1.4, 1.5, 1.6}, []uint64{}, times)
	v := vector.NewFloatVector(rv, rn)
	v.SetLen(5 * times)

	lv = append(lv, v)

	rv1, rn1 := multiplyIntVector([]int{2, 2, 4, 5, 6}, []uint64{}, times)
	v1 := vector.NewIntVector(rv1, rn1)
	v1.SetLen(5 * times)

	lv = append(lv, v1)

	rv2, rn2, ro := multiplyBsVector([]byte("123123123123123"), []int{3, 6, 9, 12, 15}, []uint64{}, times)
	v2 := vector.NewByteSliceVector(rv2, rn2, ro)
	v2.SetLen(5 * times)

	lv = append(lv, v2)
	vec = append(vec, lv)

	ag := []Aggregation{{
		AggType: Sum,
		AggCol:  0,
	}, {
		AggType: Max,
		AggCol:  0,
	}}

	g := []int{2}

	op := setUp(vec, ag, g, true).(*HashAggregateOperator)

	start := time.Now()

	op.Init()

	r := op.Next()

	elapsed := time.Since(start)
	t.Logf(" took %s", elapsed)

	a.NotNil(r, "This should not be nil")

	a.Equal("123", string(r[0].(*vector.ByteSliceVector).Get(0)))
	a.InDelta(6900000, r[1].(*vector.FloatVector).Value()[0], 0.1)
	a.InDelta(1.6, r[2].(*vector.FloatVector).Value()[0], 0.1)

}

func TestSortScenario(t *testing.T) {

	a := assert.New(t)

	// Set up child
	vec := make([][]interface{}, 0)

	lv := make([]interface{}, 0)
	v := vector.NewFloatVector([]float64{1.2, 1.2, 1.4, 1.5, 1.6}, []uint64{})
	v.SetLen(5)
	lv = append(lv, v)

	v1 := vector.NewIntVector([]int{2, 2, 4, 5, 6}, []uint64{})
	v1.SetLen(5)
	lv = append(lv, v1)

	v2 := vector.NewByteSliceVector([]byte("1123123123123"), []int{1, 4, 7, 10, 13}, []uint64{})
	v2.SetLen(5)
	lv = append(lv, v2)

	vec = append(vec, lv)

	// Set up cf

	ag := []Aggregation{{
		AggType: Sum,
		AggCol:  0,
	}, {
		AggType: Max,
		AggCol:  0,
	}}

	g := []int{2}

	op := setUp(vec, ag, g, false).(*SortedAggregateOperator)
	start := time.Now()

	op.Init()

	r := op.Next()

	elapsed := time.Since(start)
	t.Logf(" took %s", elapsed)

	a.NotNil(r, "This should not be nil")

	a.Equal("1", string(r[0].(*vector.ByteSliceVector).Get(0)))
	a.InDelta(1.2, r[1].(*vector.FloatVector).Value()[0], 0.1)
	a.InDelta(1.2, r[2].(*vector.FloatVector).Value()[0], 0.1)

	a.Equal("123", string(r[0].(*vector.ByteSliceVector).Get(1)))
	a.InDelta(5.7, r[1].(*vector.FloatVector).Value()[1], 0.1)
	a.InDelta(1.6, r[2].(*vector.FloatVector).Value()[1], 0.1)

}

func TestSortScenario2(t *testing.T) {

	a := assert.New(t)

	// Set up child
	vec := make([][]interface{}, 0)
	lv := make([]interface{}, 0)

	v := vector.NewFloatVector([]float64{1.2, 1.2, 1.4, 1.5, 1.6}, []uint64{})
	v.SetLen(5)
	lv = append(lv, v)

	v1 := vector.NewIntVector([]int{2, 2, 4, 5, 6}, []uint64{})
	v1.SetLen(5)
	lv = append(lv, v1)

	v2 := vector.NewByteSliceVector([]byte("1123123123123"), []int{1, 4, 7, 10, 13}, []uint64{})
	v2.SetLen(5)

	lv = append(lv, v2)
	vec = append(vec, lv)

	ag := []Aggregation{{
		AggType: Sum,
		AggCol:  0,
	}, {
		AggType: Max,
		AggCol:  0,
	}}

	g := []int{1, 2}

	op := setUp(vec, ag, g, false).(*SortedAggregateOperator)
	start := time.Now()

	op.Init()

	r := op.Next()

	elapsed := time.Since(start)
	t.Logf(" took %s", elapsed)

	a.NotNil(r, "This should not be nil")

	// k1
	a.Equal(2, r[0].(*vector.IntVector).Value()[0])
	a.Equal(2, r[0].(*vector.IntVector).Value()[1])
	a.Equal(4, r[0].(*vector.IntVector).Value()[2])
	a.Equal(5, r[0].(*vector.IntVector).Value()[3])
	a.Equal(6, r[0].(*vector.IntVector).Value()[4])

	// k2
	a.Equal("1", string(r[1].(*vector.ByteSliceVector).Get(0)))
	a.Equal("123", string(r[1].(*vector.ByteSliceVector).Get(1)))
	a.Equal("123", string(r[1].(*vector.ByteSliceVector).Get(2)))
	a.Equal("123", string(r[1].(*vector.ByteSliceVector).Get(3)))
	a.Equal("123", string(r[1].(*vector.ByteSliceVector).Get(4)))

	a.InDelta(1.2, r[2].(*vector.FloatVector).Value()[0], 0.1)
	a.InDelta(1.2, r[2].(*vector.FloatVector).Value()[1], 0.1)
	a.InDelta(1.4, r[2].(*vector.FloatVector).Value()[2], 0.1)
	a.InDelta(1.5, r[2].(*vector.FloatVector).Value()[3], 0.1)
	a.InDelta(1.6, r[2].(*vector.FloatVector).Value()[4], 0.1)

	a.InDelta(1.2, r[3].(*vector.FloatVector).Value()[0], 0.1)
	a.InDelta(1.2, r[3].(*vector.FloatVector).Value()[1], 0.1)
	a.InDelta(1.4, r[3].(*vector.FloatVector).Value()[2], 0.1)
	a.InDelta(1.5, r[3].(*vector.FloatVector).Value()[3], 0.1)
	a.InDelta(1.6, r[3].(*vector.FloatVector).Value()[4], 0.1)

}
*/
