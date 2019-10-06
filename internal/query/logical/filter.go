// Copyright 2019 The Meerkat Authors
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

//go:generate stringer -type=Operator
//go:generate stringer -type=ExpType

package logical

import (
	"fmt"
	"meerkat/internal/schema"
)

type Operator int

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

type ExpType int

const (
	FLOAT ExpType = iota
	DECIMAL
	BOOL
	STRING
	IDENTIFIER
	FUNCTION
)

type RootFilter struct {
	RootFilter *Filter
}

func (f *RootFilter) ResultString() string {
	return fmt.Sprintf("RF")
}

type Filter struct {
	// field / value / var or filter
	Left interface{}

	// operation AND / OR / = / != / > / < / >= etc etc
	Op Operator

	// field / value / var or filter
	Right interface{}

	Group bool
}

func (f *Filter) String() string {
	return fmt.Sprintf("Filter op %s", f.Op)
}

func NewFilter(l Node, operator Operator, r Node) *Filter {

	f := &Filter{
		Left:  l,
		Op:    operator,
		Right: r,
		Group: false,
	}

	return f
}

type Expression interface {
	Type() ExpType
	Value() string
}

type Exp struct {
	expType ExpType
	value   string
}

func NewExp(t ExpType, v string) *Exp {
	return &Exp{
		expType: t,
		value:   v,
	}
}

func (f *Exp) Value() string {
	return f.value
}

func (f *Exp) Type() ExpType {
	return f.expType
}

func (f *Exp) ResultString() string {
	return fmt.Sprintf("Expression %s value %s", f.expType, f.value)
}

type IdentifierExp struct {
	Exp
	Field *schema.Field
}

func NewIdentifier(t ExpType, v string, f *schema.Field) *IdentifierExp {
	return &Exp{
		expType: t,
		value:   v,
		Field:   f,
	}
}
