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

package rel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseQuery(t *testing.T) {

	a := assert.New(t)

	str := "indexname=Index campo1=100 and ( campo2=we or campo3=s123 )"

	p := NewMqlParser()

	is := p.Parse(str)

	a.Equal("Index", is.IndexScan.GetIndexName())

	a.NotNil(is.IndexScan.GetFilter())

	f := is.IndexScan.GetFilter()

	a.Equal(AND, f.Op)
	a.Equal(false, f.Group)

	a.Equal(DECIMAL, f.Left.(*Filter).Right.(*Exp).ExpType)
	a.Equal("100", f.Left.(*Filter).Right.(*Exp).Value)

	c1comp := f.Left.(*Filter)
	a.Equal(EQ, c1comp.Op)
	a.Equal(IDENTIFIER, c1comp.Left.(*Exp).ExpType)
	a.Equal("campo1", c1comp.Left.(*Exp).Value)

	c1c3comp := f.Right.(*Filter)
	a.Equal(OR, c1c3comp.Op)
	a.Equal(true, c1c3comp.Group)

	a.Equal(EQ, c1c3comp.Left.(*Filter).Op)

	a.Equal(EQ, c1c3comp.Right.(*Filter).Op)

	a.NotNil(f.Right)
	a.NotNil(f.Op)
	a.NotNil(f.Left)

}
