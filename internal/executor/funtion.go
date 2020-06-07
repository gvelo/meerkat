package executor

import (
	"meerkat/internal/storage/vector"
	"regexp"
)

func NewUnaryOperator(ctx Context, child VectorOperator, fn Fn) *UnaryOperator {
	return &UnaryOperator{
		ctx:   ctx,
		child: child,
		fn:    fn,
	}
}

type UnaryOperator struct {
	ctx   Context
	child VectorOperator
	fn    Fn
}

func (r *UnaryOperator) Init() {
	r.child.Init()
}

func (r *UnaryOperator) Destroy() {
	r.child.Destroy()
}

func (r *UnaryOperator) Next() interface{} {

	n := r.child.Next()

	for ; n != nil; n = r.child.Next() {
		r.fn.process(r.ctx, n)
	}

	return nil
}

// TODO(sebad): may be I should use only this operator and
// eliminate the Unary passing only one Slice in the
// multi vector
type BinaryOperator struct {
	ctx   Context
	child MultiVectorOperator
	fn    Fn
}

func (r *BinaryOperator) Init() {
	r.child.Init()
}

func (r *BinaryOperator) Destroy() {
	r.child.Destroy()
}

func (r *BinaryOperator) Next() interface{} {

	n := r.child.Next()

	for ; n != nil; n = r.child.Next() {
		r.fn.process(r.ctx, n...)
	}

	return nil
}

type Fn interface {
	process(ctx Context, params ...interface{}) interface{}
}

type fRegex struct {
	v string
}

func (f *fRegex) process(ctx Context, params ...interface{}) interface{} {

	dst := vector.DefaultVectorPool().GetByteSliceVector()
	src := params[0].(vector.ByteSliceVector)
	var validID = regexp.MustCompile(f.v)
	for i := 0; i < src.Len(); i++ {
		if validID.Match(src.Get(i)) {
			dst.AppendSlice(src.Get(i))
		}
	}

	return dst
}

type fStrLen struct {
}

func (f *fStrLen) process(ctx Context, params ...interface{}) interface{} {

	dst := vector.DefaultVectorPool().GetIntVector()
	src := params[0].(vector.ByteSliceVector)

	for i := 0; i < src.Len(); i++ {
		dst.AppendInt(len(src.Get(i)))
	}

	return dst
}
