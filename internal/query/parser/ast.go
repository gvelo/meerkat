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
	VisitPre(n Node) Node
	VisitPost(n Node) Node
}

func Walk(n Node, v Visitor) Node {
	n = v.VisitPre(n)
	n.Accept(v)
	return v.VisitPost(n)
}

type Node interface {
	Accept(v Visitor)
}

// Query Statements

// at the moment we only support Tabular Statement

type TabularStmt struct {
	TabularExpr *TabularExpr
}

func (s *TabularStmt) Accept(v Visitor) {
	s.TabularExpr = Walk(s.TabularExpr, v).(*TabularExpr)
}

// Tabular Operators

type WhereOp struct {
	Predicate Node //  CallExpr || BinaryExpr
}

func (op *WhereOp) Accept(v Visitor) {
	op.Predicate = Walk(op.Predicate, v)
}

type CountOp struct {
}

func (op *CountOp) Accept(Visitor) {
}

type ExtendOp struct {
	Columns []*ColumnExpr
}

func (op *ExtendOp) Accept(v Visitor) {
	for i, colExpr := range op.Columns {
		op.Columns[i] = Walk(colExpr, v).(*ColumnExpr)
	}
}

type ProjectOp struct {
	Columns []*ColumnExpr
}

func (op *ProjectOp) Accept(v Visitor) {
	for i, colExpr := range op.Columns {
		op.Columns[i] = Walk(colExpr, v).(*ColumnExpr)
	}
}

type LimitOp struct {
	NumberOfRows *LitExpr
}

func (op *LimitOp) Accept(Visitor) {
}

type SortOp struct {
	SortExpr []*SortExpr
}

func (op *SortOp) Accept(v Visitor) {
	for i, expr := range op.SortExpr {
		op.SortExpr[i] = Walk(expr, v).(*SortExpr)
	}
}

type SummarizeOp struct {
	Agg []*AggExpr
	By  []*ColumnExpr
}

func (op *SummarizeOp) Accept(v Visitor) {

	for i, expr := range op.Agg {
		op.Agg[i] = Walk(expr, v).(*AggExpr)

	}

	for i, expr := range op.By {
		op.By[i] = Walk(expr, v).(*ColumnExpr)
	}

}

type TopOp struct {
	NumberOfRows *LitExpr
	By           *SortExpr
}

func (op *TopOp) Accept(v Visitor) {
	op.By = Walk(op.By, v).(*SortExpr)
}

// Expressions

// TabularExpr =  StringLiteral , { "|" TabularOperator }
type TabularExpr struct {
	Source    *LitExpr
	TabularOp []Node // tabular operators
}

func (e *TabularExpr) Accept(v Visitor) {
	for i, op := range e.TabularOp {
		e.TabularOp[i] = Walk(op, v)
	}
}

type BinaryExpr struct {
	LeftExpr  Node
	Op        Token
	RightExpr Node
}

func (e *BinaryExpr) Accept(v Visitor) {
	e.LeftExpr = Walk(e.LeftExpr, v)
	e.RightExpr = Walk(e.RightExpr, v)
}

// -1
type UnaryExpr struct {
	Op   Token
	Expr Node
}

func (e *UnaryExpr) Accept(v Visitor) {
	e.Expr = Walk(e.Expr, v)
}

type LitExpr struct {
	Token Token
	Value interface{}
}

func (e *LitExpr) Accept(Visitor) {
}

type CallExpr struct {
	FuncName *LitExpr
	ArgList  []Node // Expr
}

func (e *CallExpr) Accept(v Visitor) {
	for i, arg := range e.ArgList {
		e.ArgList[i] = Walk(arg, v)
	}
}

type ColumnExpr struct {
	ColName *LitExpr // IDENT
	Expr    Node
}

func (e *ColumnExpr) Accept(v Visitor) {
	e.Expr = Walk(e.Expr, v)
}

type AggExpr struct {
	ColName *LitExpr // IDENT
	Expr    *CallExpr
}

func (e *AggExpr) Accept(v Visitor) {
	e.Expr = Walk(e.Expr, v).(*CallExpr)
}

type SortExpr struct {
	Expr      Node // LitExpr( IDENT ) || BinaryExpr || callExpr
	Asc       bool
	NullFirst bool
}

func (e *SortExpr) Accept(v Visitor) {
	e.Expr = Walk(e.Expr, v)
}
