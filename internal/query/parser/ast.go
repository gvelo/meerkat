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

func (op *CountOp) Accept(v Visitor) Node {
	v.VisitPre(op)
	return v.VisitPost(op)
}

type ExtendOp struct {
	Columns []*ColumnExpr
}

func (op *ExtendOp) Accept(v Visitor) Node {
	v.VisitPre(op)
	for i, colExpr := range op.Columns {
		op.Columns[i] = colExpr.Accept(v).(*ColumnExpr)
	}
	return v.VisitPost(op)
}

type ProjectOp struct {
	Columns []*ColumnExpr
}

func (op *ProjectOp) Accept(v Visitor) Node {
	v.VisitPre(op)
	for i, colExpr := range op.Columns {
		op.Columns[i] = colExpr.Accept(v).(*ColumnExpr)
	}
	return v.VisitPost(op)
}

type LimitOp struct {
	NumberOfRows *LitExpr
}

func (op *LimitOp) Accept(v Visitor) Node {
	v.VisitPre(op)
	return v.VisitPost(op)
}

type SortOp struct {
	SortExpr []*SortExpr
}

func (op *SortOp) Accept(v Visitor) Node {
	v.VisitPre(op)
	for i, expr := range op.SortExpr {
		op.SortExpr[i] = expr.Accept(v).(*SortExpr)
	}
	return v.VisitPost(op)
}

type SummarizeOp struct {
	Agg []*AggExpr
	By  []*ColumnExpr
}

func (op *SummarizeOp) Accept(v Visitor) Node {

	v.VisitPre(op)

	for i, expr := range op.Agg {
		op.Agg[i] = expr.Accept(v).(*AggExpr)

	}

	for i, expr := range op.By {
		op.By[i] = expr.Accept(v).(*ColumnExpr)
	}

	return v.VisitPost(op)

}

type TopOp struct {
	NumberOfRows *LitExpr
	By           *SortExpr
}

func (op *TopOp) Accept(v Visitor) Node {
	v.VisitPre(op)
	op.By = op.By.Accept(v).(*SortExpr)
	return v.VisitPost(op)
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
	LeftExpr  Node
	Op        Token
	RightExpr Node
}

func (e *BinaryExpr) Accept(v Visitor) Node {
	v.VisitPre(e)
	e.LeftExpr = e.LeftExpr.Accept(v)
	e.RightExpr = e.RightExpr.Accept(v)
	return v.VisitPost(e)
}

// -1
type UnaryExpr struct {
	Op   Token
	Expr Node
}

func (e *UnaryExpr) Accept(v Visitor) Node {
	v.VisitPre(e)
	e.Expr = e.Expr.Accept(v)
	return v.VisitPost(e)
}

type LitExpr struct {
	Token Token
	Value interface{}
}

func (e *LitExpr) Accept(v Visitor) Node {
	v.VisitPre(e)
	return v.VisitPost(e)
}

type CallExpr struct {
	FuncName *LitExpr
	ArgList  []Node // Expr
}

func (e *CallExpr) Accept(v Visitor) Node {
	v.VisitPre(e)
	for i, expr := range e.ArgList {
		e.ArgList[i] = expr.Accept(v)
	}
	return v.VisitPost(e)
}

type ColumnExpr struct {
	ColName *LitExpr // IDENT
	Expr    Node
}

func (e *ColumnExpr) Accept(v Visitor) Node {
	v.VisitPre(e)
	e.Expr = e.Expr.Accept(v)
	return v.VisitPost(e)
}

type AggExpr struct {
	ColName *LitExpr // IDENT
	Expr    *CallExpr
}

func (e *AggExpr) Accept(v Visitor) Node {
	v.VisitPre(e)
	e.Expr = e.Expr.Accept(v).(*CallExpr)
	return v.VisitPost(e)
}

type SortExpr struct {
	Expr      Node // LitExpr( IDENT ) || BinaryExpr || callExpr
	Asc       bool
	NullFirst bool
}

func (e *SortExpr) Accept(v Visitor) Node {
	v.VisitPre(e)
	e.Expr = e.Expr.Accept(v)
	return v.VisitPost(e)
}
