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
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"meerkat/internal/query/logical"
	"meerkat/internal/tools"
)

type Builder interface {
	Scan(name string) Builder
	Filter(filter interface{}) Builder
	Project(e ...interface{}) Builder
	Aggregate(groupKey string, aggCall ...interface{}) Builder
	Distinct() Builder
	Sort(exp ...string) Builder
	Limit(offset int) Builder
	Build() *logical.Projection

	// Not used yet
	SemiJoin(expr interface{}) Builder
	AntiJoin(expr interface{}) Builder
	Union(expr interface{}) Builder
	Intersect(expr interface{}) Builder
	Minus(expr interface{}) Builder
	Match(regex string) Builder
	CreateExpresion(e interface{}) *logical.Exp
}

func NewRelBuilder() Builder {
	b := new(relationalAlgBuilder)
	b.queue = make([]interface{}, 0)
	return b
}

type relationalAlgBuilder struct {
	queue []interface{}
}

func (r *relationalAlgBuilder) push(n interface{}) {
	r.queue = append(r.queue, n)
}

func (r *relationalAlgBuilder) peek() interface{} {
	return r.queue[0]
}

func (r *relationalAlgBuilder) pop() interface{} {
	e := r.queue[0]
	r.queue = r.queue[1:]
	return e
}

func (r *relationalAlgBuilder) Scan(name string) Builder {
	ts := logical.NewProjection(name)
	r.push(ts)
	return r
}

func (r *relationalAlgBuilder) Filter(f interface{}) Builder {
	is := r.peek().(*logical.Projection)
	is.AddChild(f.(logical.Node))
	return r
}

func (r *relationalAlgBuilder) Project(e ...interface{}) Builder {
	return r
}

func (r *relationalAlgBuilder) Aggregate(groupKey string, aggCall ...interface{}) Builder {
	return r
}

func (r *relationalAlgBuilder) Distinct() Builder {
	return r
}

func (r *relationalAlgBuilder) Sort(exp ...string) Builder {
	p := r.queue[0].(*logical.Projection) //projection
	p.Order = make([]*logical.Order, 0)

	for i := 0; i < len(exp); i++ {

		isOrder := exp[i] == "asc" || exp[i] == "desc"
		nextIsOrder := i+1 < len(exp) && (exp[i+1] == "asc" || exp[i+1] == "desc")
		if isOrder {
			continue
		} else {
			o := new(logical.Order)
			o.Field = exp[i]
			if nextIsOrder {
				o.Direction = exp[i+1]
			}
			p.Order = append(p.Order, o)
		}

	}

	return r
}

func (r *relationalAlgBuilder) Limit(limit int) Builder {
	p := r.queue[0].(*logical.Projection) //projection
	p.Limit = limit
	return r
}

func (r *relationalAlgBuilder) SemiJoin(expr interface{}) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) AntiJoin(expr interface{}) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) Union(expr interface{}) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) Intersect(expr interface{}) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) Minus(expr interface{}) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) Match(regex string) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) Build() *logical.Projection {
	return r.peek().(*logical.Projection)
}

func (r *relationalAlgBuilder) CreateExpresion(l interface{}) *logical.Exp {

	var e *logical.Exp
	switch l.(type) {
	case *DecimalLiteralContext:
		e = &logical.Exp{
			ExpType: logical.DECIMAL,
			Value:   l.(*DecimalLiteralContext).GetText(),
		}
	case *FloatLiteralContext:
		e = &logical.Exp{
			ExpType: logical.FLOAT,
			Value:   l.(*FloatLiteralContext).GetText(),
		}
	case *BoolLiteralContext:
		e = &logical.Exp{
			ExpType: logical.BOOL,
			Value:   l.(*BoolLiteralContext).GetText(),
		}
	case *IdentifierContext:
		e = &logical.Exp{
			ExpType: logical.IDENTIFIER,
			Value:   l.(*IdentifierContext).GetText(),
		}
	case *antlr.CommonToken: // string
		e = &logical.Exp{
			ExpType: logical.STRING,
			Value:   l.(*antlr.CommonToken).GetText(),
		}
	default:
		tools.Logf("Could not create expresion %s ", e.Value)
	}

	tools.Logf(" %s ", e.Value)
	return e
}
