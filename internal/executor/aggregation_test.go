package executor

import (
	"meerkat/internal/storage"
	"testing"
)

func TestHAgg(t *testing.T) {
	//list := setUpTop()
	//a := assert.New(t)

	f := &fakeMultiVectorOperator{
		vec:
	}

	op := NewHashAggregateOperator()
	op.Next()
}

type fakeByteSliceVector struct {
	vec []string
}

func (f fakeByteSliceVector) Len() int {
	return len(f.vec)
}

func (f fakeByteSliceVector) Rid() []uint32 {
	return nil
}

func (f fakeByteSliceVector) Data() []byte {
	return nil
}

func (f fakeByteSliceVector) Offsets() []int {
	return nil
}

func (f fakeByteSliceVector) Get(i int) []byte {
	return []byte(f.vec[i])
}

type fakeFloatVector struct {
	vec []float64
}

func (f fakeFloatVector) Len() int {
	return len(f.vec)
}

func (f fakeFloatVector) Rid() []uint32 {
	return nil
}

func (f fakeFloatVector) ValuesAsFloat() []float64 {
	return f.vec
}

type fakeIntVector struct {
	vec []int
}

func (f fakeIntVector) Len() int {
	return len(f.vec)
}

func (f fakeIntVector) Rid() []uint32 {
	return nil
}

func (f fakeIntVector) ValuesAsInt() []int {
	return f.vec
}

type fakeMultiVectorOperator struct {
	vec []storage.Vector
}

func (f *fakeMultiVectorOperator) Init() {
	// Nothing to do
}

func (f *fakeMultiVectorOperator) Destroy() {
	// Nothing to do
}

func (f *fakeMultiVectorOperator) Next() []storage.Vector {
	return f.vec
}
