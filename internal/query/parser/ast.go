// Copyright 2020 The Meerkat Authors
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

package parser

type Visitor interface {
	VisitPre(n Node)
	VisitPost(n Node) Node
}

type Node interface {
	Accept(v Visitor) Node
}

// Query Statements

// at the moment we only support Tabular Statement

type TabularStmt struct {
	TabularExpr *TabularExpr
}

func (s *TabularStmt) Accept(v Visitor) Node {
	v.VisitPre(s)
	s.TabularExpr = s.TabularExpr.Accept(v).(*TabularExpr)
	return v.VisitPost(s)
}

// Tabular Operators

type WhereOp struct {
	Predicate Node //  CallExpr || BinaryExpr
}

func (op *WhereOp) Accept(v Visitor) Node {
	v.VisitPre(op)
	op.Predicate = op.Predicate.Accept(v)
	return v.VisitPost(op)
}

type CountOp struct {
}

func (op *CountOp) Accept(Visitor) Node {
	panic("implement me")
}

type ExtendOp struct {
	Columns []*ColumnExpr
}

func (op *ExtendOp) Accept(Visitor) Node {
	panic("implement me")
}

type ProjectOp struct {
	Columns []*ColumnExpr
}

func (op *ProjectOp) Accept(Visitor) Node {
	panic("implement me")
}

type LimitOp struct {
	NumberOfRows *LitExpr
}

func (op *LimitOp) Accept(Visitor) Node {
	panic("implement me")
}

type SortOp struct {
	SortExpr []*SortExpr
}

func (op *SortOp) Accept(Visitor) Node {
	panic("implement me")
}

type SummarizeOp struct {
	Agg []*AggExpr
	By  []*ColumnExpr
}

func (op *SummarizeOp) Accept(Visitor) Node {
	panic("implement me")
}

type TopOp struct {
	NumberOfRows *LitExpr
	By           *SortExpr
}

// Expressions

// TabularExpr =  StringLiteral , { "|" TabularOperator }
type TabularExpr struct {
	Source    *LitExpr
	TabularOp []Node // tabular operators
}

func (e *TabularExpr) Accept(v Visitor) Node {
	v.VisitPre(e)
	for i, op := range e.TabularOp {
		e.TabularOp[i] = op.Accept(v)
	}
	return v.VisitPost(e)
}

type BinaryExpr struct {
	LeftExpr  interface{}
	Op        Token
	RightExpr interface{}
}

func (e *BinaryExpr) Accept(Visitor) Node {
	panic("implement me")
}

// -1
type UnaryExpr struct {
	Op   Token
	Expr interface{}
}

func (e *UnaryExpr) Accept(Visitor) Node {
	panic("implement me")
}

type LitExpr struct {
	Token Token
	Value interface{}
}

func (e *LitExpr) Accept(Visitor) Node {
	panic("implement me")
}

type CallExpr struct {
	FuncName *LitExpr
	ArgList  []Node // Expr
}

func (e *CallExpr) Accept(Visitor) Node {
	panic("implement me")
}

type ColumnExpr struct {
	ColName *LitExpr // IDENT
	Expr    Node
}

func (e *ColumnExpr) Accept(Visitor) Node {
	panic("implement me")
}

type AggExpr struct {
	ColName *LitExpr // IDENT
	Expr    *CallExpr
}

func (e *AggExpr) Accept(Visitor) Node {
	panic("implement me")
}

type SortExpr struct {
	Expr      Node // LitExpr( IDENT ) || BinaryExpr || callExpr
	Asc       bool
	NullFirst bool
}

func (e *SortExpr) Accept(Visitor) Node {
	panic("implement me")
}
