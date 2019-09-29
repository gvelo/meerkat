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

func (f *RootFilter) String() string {
	return fmt.Sprintf("Root Filter")
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

type Exp struct {
	ExpType  ExpType
	Value    string
	p        Node
	children []Node
}

func (f *Exp) String() string {
	return fmt.Sprintf("Exp %s ", f.ExpType)
}

func (f *Exp) GetParent() Node {
	return f.p
}

func (f *Exp) SetParent(n Node) {
	f.p = n
}
func (f *Exp) AddChild(n Node) {
	f.children = append(f.children, n)
}

func (f *Exp) GetChildren() []Node {
	return f.children
}
