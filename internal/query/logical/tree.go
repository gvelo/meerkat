package logical

import "meerkat/internal/query/parser"

type Operator byte

const (
	ADD               Operator = iota // +
	SUB                               // -
	MUL                               // *
	QUO                               // /
	REM                               // %
	EQL                               // ==
	EQL_CI                            // =~
	NEQ                               // !=
	NEQ_CI                            // !~
	LSS                               // <
	GTR                               // >
	LEQ                               // <=
	GEQ                               // >=
	AND                               // and
	OR                                // or
	IN                                // in
	NOT_IN                            // !in
	IN_CI                             // in~
	NOT_IN_CI                         // !in~
	HAS                               // has
	NOT_HAS                           // !has
	HAS_CS                            // has_cs
	NOT_HAS_CS                        // !has_cs
	HASPREFIX                         // hasprefix
	NOT_HASPREFIX                     // !hasprefix
	HASPREFIX_CS                      // hasprefix_cs
	NOT_HASPREFIX_CS                  // !hasprefix_cs
	HASSUFFIX                         // hassuffix
	NOT_HASSUFFIX                     // !hassuffix
	HASSUFFIX_CS                      // hassuffix_cs
	NOT_HASSUFFIX_CS                  // !hassuffix_cs
	CONTAINS                          // contains
	NOT_CONTAINS                      // !contains
	CONTAINS_CS                       // contains_cs
	NOT_CONTAINS_CS                   // !contains_cs
	STARTSWITH                        // startswith
	NOT_STARTSWITH                    // !startswith
	STARTSWITH_CS                     // startswith_cs
	NOT_STARTSWITH_CS                 // !startswith_cs
	ENDSWITH                          // endswith
	NOT_ENDSWITH                      // !endswith
	ENDSWITH_CS                       // endswith_cs
	NOT_ENDSWITH_CS                   // !endswith_cs
	MATCHES                           // matches
	HAS_ANY                           // has_any
	BETWEEN                           // between
	NOT_BETWEEN                       // !between
	RANGE                             // ..

)

var tokenToOp = map[parser.TokenType]Operator{
	parser.ADD:               ADD,
	parser.SUB:               SUB,
	parser.MUL:               MUL,
	parser.QUO:               QUO,
	parser.REM:               REM,
	parser.EQL:               EQL,
	parser.EQL_CI:            EQL_CI,
	parser.NEQ:               NEQ,
	parser.NEQ_CI:            NEQ_CI,
	parser.LSS:               LSS,
	parser.GTR:               GTR,
	parser.LEQ:               LEQ,
	parser.GEQ:               GEQ,
	parser.AND:               AND,
	parser.OR:                OR,
	parser.IN:                IN,
	parser.NOT_IN:            NOT_IN,
	parser.IN_CI:             IN_CI,
	parser.NOT_IN_CI:         NOT_IN_CI,
	parser.HAS:               HAS,
	parser.NOT_HAS:           NOT_HAS,
	parser.HAS_CS:            HAS_CS,
	parser.NOT_HAS_CS:        NOT_HAS_CS,
	parser.HASPREFIX:         HASPREFIX,
	parser.NOT_HASPREFIX:     NOT_HASPREFIX,
	parser.HASPREFIX_CS:      HASPREFIX_CS,
	parser.NOT_HASPREFIX_CS:  NOT_HASPREFIX_CS,
	parser.HASSUFFIX:         HASSUFFIX,
	parser.NOT_HASSUFFIX:     NOT_HASSUFFIX,
	parser.HASSUFFIX_CS:      HASSUFFIX_CS,
	parser.NOT_HASSUFFIX_CS:  NOT_HASSUFFIX_CS,
	parser.CONTAINS:          CONTAINS,
	parser.NOT_CONTAINS:      NOT_CONTAINS,
	parser.CONTAINS_CS:       CONTAINS_CS,
	parser.NOT_CONTAINS_CS:   NOT_CONTAINS_CS,
	parser.STARTSWITH:        STARTSWITH,
	parser.NOT_STARTSWITH:    NOT_STARTSWITH,
	parser.STARTSWITH_CS:     STARTSWITH_CS,
	parser.NOT_STARTSWITH_CS: NOT_STARTSWITH_CS,
	parser.ENDSWITH:          ENDSWITH,
	parser.NOT_ENDSWITH:      NOT_ENDSWITH,
	parser.ENDSWITH_CS:       ENDSWITH_CS,
	parser.NOT_ENDSWITH_CS:   NOT_ENDSWITH_CS,
	parser.MATCHES:           MATCHES,
	parser.HAS_ANY:           HAS_ANY,
	parser.BETWEEN:           BETWEEN,
	parser.NOT_BETWEEN:       NOT_BETWEEN,
	parser.RANGE:             RANGE,
}

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

type Op interface {
	// stage
	// paralelizable
	Node
}

// Expressions

type BinaryExpr struct {
	LeftExpr  Node
	Op        Operator
	RightExpr Node
}

func (n *BinaryExpr) Accept(v Visitor) {
	n.LeftExpr = Walk(n.LeftExpr, v)
	n.RightExpr = Walk(n.RightExpr, v)
}

type UnaryExpr struct {
	Op   Operator
	Expr Node
}

func (n *UnaryExpr) Accept(v Visitor) { n.Expr = Walk(n.Expr, v) }

type CallExpr struct {
	FuncName string
	ArgList  []Node
}

func (n *CallExpr) Accept(v Visitor) {
	for i, node := range n.ArgList {
		n.ArgList[i] = Walk(node, v)
	}
}

type ColRefExpr struct {
	// type ??
	Name string
}

func (n *ColRefExpr) Accept(Visitor) {}

type LiteralExpr struct {
	Value interface{}
}

func (n *LiteralExpr) Accept(Visitor) {}

type AggExpr struct {
	ColName string
	Expr    *CallExpr
}

func (n *AggExpr) Accept(v Visitor) { n.Expr = Walk(n.Expr, v).(*CallExpr) }

type ColumnExpr struct {
	ColName string
	Expr    Node
}

func (n *ColumnExpr) Accept(v Visitor) { n.Expr = Walk(n.Expr, v) }

// Operators

type SourceOp struct {
	TableName string
}

func (n *SourceOp) Accept(Visitor) {}

type FilterOp struct {
	Predicate Node
	Child     Node
}

func (n *FilterOp) Accept(v Visitor) {
	n.Predicate = Walk(n.Predicate, v)
	n.Child = Walk(n.Child, v)
}

type LimitOp struct {
	NumOfRows int
	Child     Node
}

func (n *LimitOp) Accept(v Visitor) { n.Child = Walk(n.Child, v) }

type SummarizeOp struct {
	Agg   []*AggExpr
	By    []*ColumnExpr
	Child Node
}

func (n *SummarizeOp) Accept(v Visitor) {
	for i, expr := range n.Agg {
		n.Agg[i] = Walk(expr, v).(*AggExpr)
	}
	for i, expr := range n.By {
		n.By[i] = Walk(expr, v).(*ColumnExpr)
	}
	n.Child = Walk(n.Child, v)
}

type NodeOutOp struct {
	Dst       string
	StreamMap map[string]int64
	Child     Node
}

func (n *NodeOutOp) Accept(v Visitor) { n.Child = Walk(n.Child, v) }

type MergeSortOp struct {
	StreamMap map[string]int64
}

func (n *MergeSortOp) Accept(Visitor) {}

type HashExchangeOutOp struct {
	Child Node
}

func (n *HashExchangeOutOp) Accept(v Visitor) { n.Child = Walk(n.Child, v) }

type HashExchangeInOp struct {
	Child Node
}

func (n *HashExchangeInOp) Accept(v Visitor) { n.Child = Walk(n.Child, v) }

type OutputOp struct {
	Child Node
}

func (n *OutputOp) Accept(v Visitor) { n.Child = Walk(n.Child, v) }

type DistSummaryOp struct {
	Child Node
}

func (n *DistSummaryOp) Accept(v Visitor) { n.Child = Walk(n.Child, v) }

type SummaryCollector struct {
	Child Node
}

func (n *SummaryCollector) Accept(v Visitor) { n.Child = Walk(n.Child, v) }

type ShuffleSummaryOp struct {
	Child Node
}

func (n *ShuffleSummaryOp) Accept(v Visitor) { n.Child = Walk(n.Child, v) }

type LocalSummaryOp struct {
	Child Node
}

func (n *LocalSummaryOp) Accept(v Visitor) { n.Child = Walk(n.Child, v) }

type MergeOp struct {
	Child Node
}

func (n *MergeOp) Accept(Visitor) {}

// TODO(gvelo)
var _ Node = &OutputOp{}
var _ Node = &SourceOp{}
var _ Node = &BinaryExpr{}
