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
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"meerkat/internal/query/mql_parser"
	"meerkat/internal/tools"
)

type Builder interface {
	Scan(name string) Builder
	Filter(filter interface{}) Builder
	Project(e ...interface{}) Builder
	Aggregate(groupKey string, aggCall ...interface{}) Builder
	Distinct() Builder
	Sort(exp ...interface{}) Builder
	Limit(offset int) Builder
	SemiJoin(expr interface{}) Builder
	AntiJoin(expr interface{}) Builder
	Union(expr interface{}) Builder
	Intersect(expr interface{}) Builder
	Minus(expr interface{}) Builder
	Match(regex string) Builder
	And(a interface{}) Builder
	Or(o interface{}) Builder
	CreateExpresion(e interface{}) *Exp
	Build() *ParsedTree
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
	ts := NewIndexScan(name)
	r.push(ts)
	return r
}

func (r *relationalAlgBuilder) And(a interface{}) Builder {
	return r
}

func (r *relationalAlgBuilder) Or(o interface{}) Builder {
	return r
}

func (r *relationalAlgBuilder) Filter(f interface{}) Builder {
	is := r.peek().(*IndexScan)
	is.SetFilter(f.(*Filter))
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

func (r *relationalAlgBuilder) Sort(exp ...interface{}) Builder {
	return r
}

func (r *relationalAlgBuilder) Limit(offset int) Builder {
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

func (r *relationalAlgBuilder) Build() *ParsedTree {
	return &ParsedTree{IndexScan: r.peek().(*IndexScan)}
}

func (r *relationalAlgBuilder) CreateExpresion(l interface{}) *Exp {

	var e *Exp
	switch  l.(type) {
	case *mql_parser.DecimalLiteralContext:
		e = &Exp{
			ExpType: DECIMAL,
			Value:   l.(*mql_parser.DecimalLiteralContext).GetText(),
		}
	case *mql_parser.FloatLiteralContext:
		e = &Exp{
			ExpType: FLOAT,
			Value:   l.(*mql_parser.FloatLiteralContext).GetText(),
		}
	case *mql_parser.BoolLiteralContext:
		e = &Exp{
			ExpType: BOOL,
			Value:   l.(*mql_parser.BoolLiteralContext).GetText(),
		}
	case *mql_parser.IdentifierContext:
		e = &Exp{
			ExpType: IDENTIFIER,
			Value:   l.(*mql_parser.IdentifierContext).GetText(),
		}
	case *antlr.CommonToken: // string
		e = &Exp{
			ExpType: STRING,
			Value:   l.(*antlr.CommonToken).GetText(),
		}
	default:
		tools.Logf("Could not create expresion %s ", e.Value )
	}

	tools.Logf(" %s ", e.Value )
	return e
}
