package executor

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage"
	"testing"
	"time"
)

// FAKES

type fakeMultiVectorOperator struct {
	vect []storage.Vector
	idx  int
}

func (f *fakeMultiVectorOperator) Init() {
	f.idx = 0
}

func (f *fakeMultiVectorOperator) Destroy() {
	// nothing to do
}

func (f *fakeMultiVectorOperator) Next() []storage.Vector {
	if f.idx == 1 {
		return nil
	}
	f.idx++
	return f.vect
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

func (f *fakeFloatColumn) Encoding() storage.Encoding {
	panic("implement me")
}

func (f fakeFloatColumn) Validity() *roaring.Bitmap {
	panic("implement me")
}

func (f *fakeFloatColumn) HasNulls() bool {
	panic("implement me")
}

func (f *fakeFloatColumn) Read(pos []uint32) (storage.Vector, error) {
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

func (f fakeFloatColumn) Iterator() storage.FloatIterator {
	panic("implement me")
}

type fakeIntColumn struct {
}

func (f *fakeIntColumn) Encoding() storage.Encoding {
	panic("implement me")
}

func (f *fakeIntColumn) Validity() *roaring.Bitmap {
	panic("implement me")
}

func (f *fakeIntColumn) HasNulls() bool {
	panic("implement me")
}

func (f *fakeIntColumn) Read(pos []uint32) (storage.Vector, error) {
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

func (f *fakeIntColumn) Iterator() storage.IntIterator {
	panic("implement me")
}

type fakeStringColumn struct {
}

func (f *fakeStringColumn) Encoding() storage.Encoding {
	panic("implement me")
}

func (f *fakeStringColumn) Validity() *roaring.Bitmap {
	panic("implement me")
}

func (f *fakeStringColumn) HasNulls() bool {
	panic("implement me")
}

func (f *fakeStringColumn) Read(pos []uint32) (storage.Vector, error) {
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

func (f *fakeStringColumn) ReadDictEnc(pos []uint32) (storage.IntVector, error) {
	panic("implement me")
}

func (f *fakeStringColumn) DictEncodedIterator() storage.IntIterator {
	panic("implement me")
}

func (f *fakeStringColumn) Iterator() storage.ByteSliceIterator {
	panic("implement me")
}

func TestHAgg(t *testing.T) {

	a := assert.New(t)

	// Set up vectors
	vec := make([]storage.Vector, 0)
	vec = append(vec, storage.NewFloatVectorFromSlice([]float64{1.2, 1.2, 1.4, 1.5, 1.6}))
	vec = append(vec, storage.NewIntVectorFromSlice([]int{2, 2, 4, 5, 6}))
	vec = append(vec, storage.NewByteSliceVectorSlice([][]byte{[]byte("123"), []byte("123"), []byte("123"), []byte("123"), []byte("123")}))

	sMap := make(map[string]storage.Column)
	sMap["c1"] = &fakeFloatColumn{}
	sMap["c2"] = &fakeIntColumn{}
	sMap["c3"] = &fakeStringColumn{}

	fs := &fakeSegment{sMap: sMap}

	f := &fakeMultiVectorOperator{
		vect: vec,
	}

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
	ctx.Value(ColumnIndexKeysKey, m)

	op := NewHashAggregateOperator(ctx, f, ag, g)
	op.Init()

	r := op.Next()
	a.NotNil(r, "This should not be nil")

	print(r)
}

func TestIt(t *testing.T) {
	result := make([][]interface{}, 0, 100)

	for i := 0; i < 10; i++ {
		k := make([]interface{}, 5, 5)
		result = append(result, k)
		for j := 0; j < 5; j++ {
			result[i][j] = i + j
		}
	}

	print(len(result))
}
