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

type Node interface {
	Walk()
}

// Query Statements

// at the moment we only support Tabular Statement

type TabularStmt struct {
	TabularExpr *TabularExpr
}

// Tabular Operators

type WhereOp struct {
	Predicate interface{} //  CallExpr || BinaryExpr
}

type CountOp struct {
}

type ExtendOp struct {
	Columns []*ColumnExpr
}

type LimitOp struct {
	NumberOfRows *LitExpr
}

type SortOp struct {
	SortExpr []*SortExpr
}

type SummarizeOp struct {
	Agg *AggExpr
	By  *ColumnExpr
}

type TopOp struct {
	NumberOfRows *LitExpr
	By           *SortExpr
}

// Expressions

// TabularExpr =  StringLiteral , { "|" TabularOperator }
type TabularExpr struct {
	Source    *LitExpr
	TabularOp []interface{} // tabular operators
}

type BinaryExpr struct {
	LeftExpr  interface{}
	Op        Token
	RightExpr interface{}
}

// -1
type UnaryExpr struct {
	Op   Token
	Expr interface{}
}

type LitExpr struct {
	Token Token
	Value interface{}
}

type CallExpr struct {
	FuncName *LitExpr
	ArgList  []interface{} // Expr
}

type ColumnExpr struct {
	ColName *LitExpr // IDENT
	Expr    interface{}
}

type AggExpr struct {
	ColName *LitExpr // IDENT
	Expr    CallExpr
}

type SortExpr struct {
	Expr      interface{} // LitExpr( IDENT ) || BinaryExpr || callExpr
	Asc       bool
	NullFirst bool
}
