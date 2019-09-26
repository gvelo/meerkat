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

package mql_parser

import (
	"github.com/stretchr/testify/assert"
	"meerkat/internal/query/logical"
	"testing"
)

func TestParseQuery(t *testing.T) {

	a := assert.New(t)

	// str := "index=Index campo1=100 and ( campo2=we or campo3=s123 )"
	str := "index=Index campo1=100 and ( campo2=we or campo3=s123 )"

	projection := Parse(str)

	a.Equal("Index", projection.IndexName)

	a.NotNil(projection.GetChildren())
	a.NotNil(projection.GetChildren()[0])

	n := projection.GetChildren()[0]
	f := n.(*logical.Filter)

	a.Nil(f.GetChildren())

	a.Equal(logical.AND, f.Op)
	a.Equal(false, f.Group)

	a.Equal(logical.DECIMAL, f.Left.(*logical.Filter).Right.(*logical.Exp).ExpType)
	a.Equal("100", f.Left.(*logical.Filter).Right.(*logical.Exp).Value)

	c1comp := f.Left.(*logical.Filter)
	a.Equal(logical.EQ, c1comp.Op)
	a.Equal(logical.IDENTIFIER, c1comp.Left.(*logical.Exp).ExpType)
	a.Equal("campo1", c1comp.Left.(*logical.Exp).Value)

	c1c3comp := f.Right.(*logical.Filter)
	a.Equal(logical.OR, c1c3comp.Op)
	a.Equal(true, c1c3comp.Group)

	a.Equal(logical.EQ, c1c3comp.Left.(*logical.Filter).Op)

	a.Equal(logical.EQ, c1c3comp.Right.(*logical.Filter).Op)

	a.NotNil(f.Right)
	a.NotNil(f.Op)
	a.NotNil(f.Left)

}

func TestParseQuery2(t *testing.T) {

	a := assert.New(t)

	str := "index=Index campo1=100 | top 10 | sort campo1 desc, campo3"

	p := Parse(str)

	a.Equal("Index", p.IndexName)
	a.Equal(10, p.Limit)

	a.Equal(2, len(p.Order))

	a.Equal("campo1", p.Order[0].Field)
	a.Equal("desc", p.Order[0].Direction)

	a.Equal("campo3", p.Order[1].Field)
	a.Equal("asc", p.Order[1].Direction)

}

func TestParseQuery3(t *testing.T) {

	a := assert.New(t)

	// TODO: implentar el earlier como un filtro mas asi entra en el optimizador
	str := "earlier=-1d request_id=\"a37cacc3-71d5-40f0-a329-a051a3949ced\" "

	p := Parse(str)

	a.Equal("", p.IndexName)

	a.Equal(1, len(p.GetChildren()))

	str = "request_id=\"a37cacc3-71d5-40f0-a329-a051a3949ced\" earlier=-1d  "

	p = Parse(str)

	a.Equal("", p.IndexName)

	a.Equal(1, len(p.GetChildren()))

}
