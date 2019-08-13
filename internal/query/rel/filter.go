package rel

import "fmt"

type Operator int
type ExpType int

const (
	AND Operator = iota
	OR
	EQ
	GT
	LT
	LEQT
	GEQT
	DST
)

type Filter struct {
	L  interface{} // field / value / var or filter
	Op Operator    // operation AND / OR / = / != / > / < / >= etc etc
	R  interface{} // field / value / var or filter
}

func (f *Filter) toString() string {
	return fmt.Sprintf("L %s op %d R %s", f.L, f.Op, f.R)
}

func NewFilter(l interface{}, operator Operator, r interface{}) *Filter {
	f := new(Filter)
	f.L = l
	f.Op = operator
	f.R = r
	return f
}

const (
	LITERAL ExpType = iota
	FIELD
	VARIABLE
)

type Exp struct {
	expType ExpType
	value   string
}
