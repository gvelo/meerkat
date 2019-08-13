// Code generated from /Users/sdominguez/desa/workspace_go/eventdb/internal/query/MqlParser.g4 by ANTLR 4.7.2. DO NOT EDIT.

package mql_parser // MqlParser
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 78, 171,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 7,
	3, 47, 10, 3, 12, 3, 14, 3, 50, 11, 3, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3,
	5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 5, 6, 64, 10, 6, 3, 7, 3, 7, 3,
	7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 5, 7, 75, 10, 7, 3, 7, 3, 7, 3,
	7, 3, 7, 7, 7, 81, 10, 7, 12, 7, 14, 7, 84, 11, 7, 3, 8, 3, 8, 3, 9, 3,
	9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 5, 10, 98, 10,
	10, 3, 11, 3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 3, 12, 7, 12, 107, 10, 12,
	12, 12, 14, 12, 110, 11, 12, 3, 13, 3, 13, 3, 13, 3, 13, 6, 13, 116, 10,
	13, 13, 13, 14, 13, 117, 3, 14, 3, 14, 3, 14, 3, 14, 5, 14, 124, 10, 14,
	3, 14, 3, 14, 3, 14, 3, 14, 5, 14, 130, 10, 14, 7, 14, 132, 10, 14, 12,
	14, 14, 14, 135, 11, 14, 3, 14, 3, 14, 3, 14, 3, 15, 3, 15, 5, 15, 142,
	10, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 16, 3, 16, 5, 16, 150, 10, 16, 3,
	16, 3, 16, 3, 17, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19,
	3, 20, 3, 20, 3, 20, 7, 20, 166, 10, 20, 12, 20, 14, 20, 169, 11, 20, 3,
	20, 2, 3, 12, 21, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30,
	32, 34, 36, 38, 2, 7, 3, 2, 19, 34, 4, 2, 44, 46, 52, 53, 3, 2, 17, 18,
	3, 2, 74, 75, 3, 2, 55, 56, 2, 173, 2, 40, 3, 2, 2, 2, 4, 43, 3, 2, 2,
	2, 6, 51, 3, 2, 2, 2, 8, 53, 3, 2, 2, 2, 10, 63, 3, 2, 2, 2, 12, 74, 3,
	2, 2, 2, 14, 85, 3, 2, 2, 2, 16, 87, 3, 2, 2, 2, 18, 97, 3, 2, 2, 2, 20,
	99, 3, 2, 2, 2, 22, 102, 3, 2, 2, 2, 24, 111, 3, 2, 2, 2, 26, 119, 3, 2,
	2, 2, 28, 139, 3, 2, 2, 2, 30, 147, 3, 2, 2, 2, 32, 153, 3, 2, 2, 2, 34,
	156, 3, 2, 2, 2, 36, 159, 3, 2, 2, 2, 38, 162, 3, 2, 2, 2, 40, 41, 5, 38,
	20, 2, 41, 42, 7, 2, 2, 3, 42, 3, 3, 2, 2, 2, 43, 48, 7, 75, 2, 2, 44,
	45, 7, 42, 2, 2, 45, 47, 7, 75, 2, 2, 46, 44, 3, 2, 2, 2, 47, 50, 3, 2,
	2, 2, 48, 46, 3, 2, 2, 2, 48, 49, 3, 2, 2, 2, 49, 5, 3, 2, 2, 2, 50, 48,
	3, 2, 2, 2, 51, 52, 9, 2, 2, 2, 52, 7, 3, 2, 2, 2, 53, 54, 5, 6, 4, 2,
	54, 55, 7, 35, 2, 2, 55, 56, 7, 75, 2, 2, 56, 57, 7, 36, 2, 2, 57, 9, 3,
	2, 2, 2, 58, 64, 7, 72, 2, 2, 59, 64, 7, 63, 2, 2, 60, 64, 7, 68, 2, 2,
	61, 64, 7, 70, 2, 2, 62, 64, 7, 75, 2, 2, 63, 58, 3, 2, 2, 2, 63, 59, 3,
	2, 2, 2, 63, 60, 3, 2, 2, 2, 63, 61, 3, 2, 2, 2, 63, 62, 3, 2, 2, 2, 64,
	11, 3, 2, 2, 2, 65, 66, 8, 7, 1, 2, 66, 67, 7, 35, 2, 2, 67, 68, 5, 12,
	7, 2, 68, 69, 7, 36, 2, 2, 69, 75, 3, 2, 2, 2, 70, 71, 7, 75, 2, 2, 71,
	72, 5, 14, 8, 2, 72, 73, 5, 10, 6, 2, 73, 75, 3, 2, 2, 2, 74, 65, 3, 2,
	2, 2, 74, 70, 3, 2, 2, 2, 75, 82, 3, 2, 2, 2, 76, 77, 12, 3, 2, 2, 77,
	78, 5, 16, 9, 2, 78, 79, 5, 12, 7, 4, 79, 81, 3, 2, 2, 2, 80, 76, 3, 2,
	2, 2, 81, 84, 3, 2, 2, 2, 82, 80, 3, 2, 2, 2, 82, 83, 3, 2, 2, 2, 83, 13,
	3, 2, 2, 2, 84, 82, 3, 2, 2, 2, 85, 86, 9, 3, 2, 2, 86, 15, 3, 2, 2, 2,
	87, 88, 9, 4, 2, 2, 88, 17, 3, 2, 2, 2, 89, 98, 5, 20, 11, 2, 90, 98, 5,
	24, 13, 2, 91, 98, 5, 26, 14, 2, 92, 98, 5, 30, 16, 2, 93, 98, 5, 32, 17,
	2, 94, 98, 5, 34, 18, 2, 95, 98, 5, 36, 19, 2, 96, 98, 5, 28, 15, 2, 97,
	89, 3, 2, 2, 2, 97, 90, 3, 2, 2, 2, 97, 91, 3, 2, 2, 2, 97, 92, 3, 2, 2,
	2, 97, 93, 3, 2, 2, 2, 97, 94, 3, 2, 2, 2, 97, 95, 3, 2, 2, 2, 97, 96,
	3, 2, 2, 2, 98, 19, 3, 2, 2, 2, 99, 100, 7, 7, 2, 2, 100, 101, 5, 12, 7,
	2, 101, 21, 3, 2, 2, 2, 102, 103, 7, 3, 2, 2, 103, 104, 7, 44, 2, 2, 104,
	108, 7, 75, 2, 2, 105, 107, 5, 12, 7, 2, 106, 105, 3, 2, 2, 2, 107, 110,
	3, 2, 2, 2, 108, 106, 3, 2, 2, 2, 108, 109, 3, 2, 2, 2, 109, 23, 3, 2,
	2, 2, 110, 108, 3, 2, 2, 2, 111, 115, 7, 4, 2, 2, 112, 113, 7, 75, 2, 2,
	113, 114, 7, 15, 2, 2, 114, 116, 7, 75, 2, 2, 115, 112, 3, 2, 2, 2, 116,
	117, 3, 2, 2, 2, 117, 115, 3, 2, 2, 2, 117, 118, 3, 2, 2, 2, 118, 25, 3,
	2, 2, 2, 119, 120, 7, 14, 2, 2, 120, 123, 5, 8, 5, 2, 121, 122, 7, 15,
	2, 2, 122, 124, 7, 75, 2, 2, 123, 121, 3, 2, 2, 2, 123, 124, 3, 2, 2, 2,
	124, 133, 3, 2, 2, 2, 125, 126, 7, 42, 2, 2, 126, 129, 5, 8, 5, 2, 127,
	128, 7, 15, 2, 2, 128, 130, 7, 75, 2, 2, 129, 127, 3, 2, 2, 2, 129, 130,
	3, 2, 2, 2, 130, 132, 3, 2, 2, 2, 131, 125, 3, 2, 2, 2, 132, 135, 3, 2,
	2, 2, 133, 131, 3, 2, 2, 2, 133, 134, 3, 2, 2, 2, 134, 136, 3, 2, 2, 2,
	135, 133, 3, 2, 2, 2, 136, 137, 7, 16, 2, 2, 137, 138, 7, 75, 2, 2, 138,
	27, 3, 2, 2, 2, 139, 141, 7, 11, 2, 2, 140, 142, 9, 5, 2, 2, 141, 140,
	3, 2, 2, 2, 141, 142, 3, 2, 2, 2, 142, 143, 3, 2, 2, 2, 143, 144, 7, 12,
	2, 2, 144, 145, 7, 44, 2, 2, 145, 146, 7, 67, 2, 2, 146, 29, 3, 2, 2, 2,
	147, 149, 7, 13, 2, 2, 148, 150, 9, 6, 2, 2, 149, 148, 3, 2, 2, 2, 149,
	150, 3, 2, 2, 2, 150, 151, 3, 2, 2, 2, 151, 152, 5, 4, 3, 2, 152, 31, 3,
	2, 2, 2, 153, 154, 7, 6, 2, 2, 154, 155, 5, 4, 3, 2, 155, 33, 3, 2, 2,
	2, 156, 157, 7, 8, 2, 2, 157, 158, 5, 4, 3, 2, 158, 35, 3, 2, 2, 2, 159,
	160, 7, 9, 2, 2, 160, 161, 7, 63, 2, 2, 161, 37, 3, 2, 2, 2, 162, 167,
	5, 22, 12, 2, 163, 164, 7, 60, 2, 2, 164, 166, 5, 18, 10, 2, 165, 163,
	3, 2, 2, 2, 166, 169, 3, 2, 2, 2, 167, 165, 3, 2, 2, 2, 167, 168, 3, 2,
	2, 2, 168, 39, 3, 2, 2, 2, 169, 167, 3, 2, 2, 2, 15, 48, 63, 74, 82, 97,
	108, 117, 123, 129, 133, 141, 149, 167,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'sourcetype'", "'rename'", "'seach'", "'dedup'", "'where'", "'sort'",
	"'head'", "'top'", "'bin'", "'span'", "'fields'", "'stats'", "'as'", "'by'",
	"'and'", "'or'", "'avg'", "'count'", "'distinct_count'", "'estdc'", "'estdc_error'",
	"'max'", "'median'", "'min'", "'mode'", "'range'", "'stdev'", "'stdevp'",
	"'sum'", "'sumsq'", "'var'", "'varp'", "'('", "')'", "'{'", "'}'", "'['",
	"']'", "';'", "','", "'.'", "'='", "'>'", "'<'", "'!'", "'~'", "'?'", "':'",
	"'=='", "'<='", "'>='", "'!='", "'+'", "'-'", "'*'", "'/'", "'&'", "'|'",
	"'^'", "'%'", "", "", "", "", "", "", "", "", "", "", "'null'", "'_time'",
}
var symbolicNames = []string{
	"", "SOURCE_TYPE", "RENAME", "SEARCH", "DEDUP", "WHERE", "SORT", "HEAD",
	"TOP", "BIN", "BIN_SPAN", "FIELDS", "STATS", "AS", "BY", "AND", "OR", "AVG",
	"COUNT", "DISTINCT_COUNT", "ESTDC", "ESTDC_ERROR", "MAX", "MEDIAN", "MIN",
	"MODE", "RANGE", "STDEV", "STDEVP", "SUM", "SUMSQ", "VAR", "VARP", "LPAREN",
	"RPAREN", "LBRACE", "RBRACE", "LBRACK", "RBRACK", "SEMI", "COMMA", "DOT",
	"ASSIGN", "GT", "LT", "BANG", "TILDE", "QUESTION", "COLON", "EQUAL", "LE",
	"GE", "NOTEQUAL", "ADD", "SUB", "MUL", "DIV", "BITAND", "BITOR", "CARET",
	"MOD", "DECIMAL_LITERAL", "HEX_LITERAL", "OCT_LITERAL", "BINARY_LITERAL",
	"TIME_LITERAL", "FLOAT_LITERAL", "HEX_FLOAT_LITERAL", "BOOL_LITERAL", "CHAR_LITERAL",
	"STRING_LITERAL", "NULL_LITERAL", "TIME_FIELD", "IDENTIFIER", "WS", "COMMENT",
	"LINE_COMMENT",
}

var ruleNames = []string{
	"start", "identifierList", "agrupTypes", "agrupCall", "literal", "expression",
	"comparator", "binary", "commands", "whereCommand", "selectCommand", "renameCommand",
	"statCommand", "binCommand", "fieldCommand", "dedupCommand", "sortCommand",
	"headCommand", "completeCommand",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type MqlParser struct {
	*antlr.BaseParser
}

func NewMqlParser(input antlr.TokenStream) *MqlParser {
	this := new(MqlParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "MqlParser.g4"

	return this
}

// MqlParser tokens.
const (
	MqlParserEOF               = antlr.TokenEOF
	MqlParserSOURCE_TYPE       = 1
	MqlParserRENAME            = 2
	MqlParserSEARCH            = 3
	MqlParserDEDUP             = 4
	MqlParserWHERE             = 5
	MqlParserSORT              = 6
	MqlParserHEAD              = 7
	MqlParserTOP               = 8
	MqlParserBIN               = 9
	MqlParserBIN_SPAN          = 10
	MqlParserFIELDS            = 11
	MqlParserSTATS             = 12
	MqlParserAS                = 13
	MqlParserBY                = 14
	MqlParserAND               = 15
	MqlParserOR                = 16
	MqlParserAVG               = 17
	MqlParserCOUNT             = 18
	MqlParserDISTINCT_COUNT    = 19
	MqlParserESTDC             = 20
	MqlParserESTDC_ERROR       = 21
	MqlParserMAX               = 22
	MqlParserMEDIAN            = 23
	MqlParserMIN               = 24
	MqlParserMODE              = 25
	MqlParserRANGE             = 26
	MqlParserSTDEV             = 27
	MqlParserSTDEVP            = 28
	MqlParserSUM               = 29
	MqlParserSUMSQ             = 30
	MqlParserVAR               = 31
	MqlParserVARP              = 32
	MqlParserLPAREN            = 33
	MqlParserRPAREN            = 34
	MqlParserLBRACE            = 35
	MqlParserRBRACE            = 36
	MqlParserLBRACK            = 37
	MqlParserRBRACK            = 38
	MqlParserSEMI              = 39
	MqlParserCOMMA             = 40
	MqlParserDOT               = 41
	MqlParserASSIGN            = 42
	MqlParserGT                = 43
	MqlParserLT                = 44
	MqlParserBANG              = 45
	MqlParserTILDE             = 46
	MqlParserQUESTION          = 47
	MqlParserCOLON             = 48
	MqlParserEQUAL             = 49
	MqlParserLE                = 50
	MqlParserGE                = 51
	MqlParserNOTEQUAL          = 52
	MqlParserADD               = 53
	MqlParserSUB               = 54
	MqlParserMUL               = 55
	MqlParserDIV               = 56
	MqlParserBITAND            = 57
	MqlParserBITOR             = 58
	MqlParserCARET             = 59
	MqlParserMOD               = 60
	MqlParserDECIMAL_LITERAL   = 61
	MqlParserHEX_LITERAL       = 62
	MqlParserOCT_LITERAL       = 63
	MqlParserBINARY_LITERAL    = 64
	MqlParserTIME_LITERAL      = 65
	MqlParserFLOAT_LITERAL     = 66
	MqlParserHEX_FLOAT_LITERAL = 67
	MqlParserBOOL_LITERAL      = 68
	MqlParserCHAR_LITERAL      = 69
	MqlParserSTRING_LITERAL    = 70
	MqlParserNULL_LITERAL      = 71
	MqlParserTIME_FIELD        = 72
	MqlParserIDENTIFIER        = 73
	MqlParserWS                = 74
	MqlParserCOMMENT           = 75
	MqlParserLINE_COMMENT      = 76
)

// MqlParser rules.
const (
	MqlParserRULE_start           = 0
	MqlParserRULE_identifierList  = 1
	MqlParserRULE_agrupTypes      = 2
	MqlParserRULE_agrupCall       = 3
	MqlParserRULE_literal         = 4
	MqlParserRULE_expression      = 5
	MqlParserRULE_comparator      = 6
	MqlParserRULE_binary          = 7
	MqlParserRULE_commands        = 8
	MqlParserRULE_whereCommand    = 9
	MqlParserRULE_selectCommand   = 10
	MqlParserRULE_renameCommand   = 11
	MqlParserRULE_statCommand     = 12
	MqlParserRULE_binCommand      = 13
	MqlParserRULE_fieldCommand    = 14
	MqlParserRULE_dedupCommand    = 15
	MqlParserRULE_sortCommand     = 16
	MqlParserRULE_headCommand     = 17
	MqlParserRULE_completeCommand = 18
)

// IStartContext is an interface to support dynamic dispatch.
type IStartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_start
	return p
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_start

	return p
}

func (s *StartContext) GetParser() antlr.Parser { return s.parser }

func (s *StartContext) CompleteCommand() ICompleteCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICompleteCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICompleteCommandContext)
}

func (s *StartContext) EOF() antlr.TerminalNode {
	return s.GetToken(MqlParserEOF, 0)
}

func (s *StartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterStart(s)
	}
}

func (s *StartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitStart(s)
	}
}

func (s *StartContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitStart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) Start() (localctx IStartContext) {
	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, MqlParserRULE_start)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(38)
		p.CompleteCommand()
	}
	{
		p.SetState(39)
		p.Match(MqlParserEOF)
	}

	return localctx
}

// IIdentifierListContext is an interface to support dynamic dispatch.
type IIdentifierListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIdentifierListContext differentiates from other interfaces.
	IsIdentifierListContext()
}

type IdentifierListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentifierListContext() *IdentifierListContext {
	var p = new(IdentifierListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_identifierList
	return p
}

func (*IdentifierListContext) IsIdentifierListContext() {}

func NewIdentifierListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifierListContext {
	var p = new(IdentifierListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_identifierList

	return p
}

func (s *IdentifierListContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentifierListContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(MqlParserIDENTIFIER)
}

func (s *IdentifierListContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, i)
}

func (s *IdentifierListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(MqlParserCOMMA)
}

func (s *IdentifierListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserCOMMA, i)
}

func (s *IdentifierListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifierListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentifierListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterIdentifierList(s)
	}
}

func (s *IdentifierListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitIdentifierList(s)
	}
}

func (s *IdentifierListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitIdentifierList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) IdentifierList() (localctx IIdentifierListContext) {
	localctx = NewIdentifierListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, MqlParserRULE_identifierList)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(41)
		p.Match(MqlParserIDENTIFIER)
	}
	p.SetState(46)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserCOMMA {
		{
			p.SetState(42)
			p.Match(MqlParserCOMMA)
		}
		{
			p.SetState(43)
			p.Match(MqlParserIDENTIFIER)
		}

		p.SetState(48)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IAgrupTypesContext is an interface to support dynamic dispatch.
type IAgrupTypesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAgrupTypesContext differentiates from other interfaces.
	IsAgrupTypesContext()
}

type AgrupTypesContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAgrupTypesContext() *AgrupTypesContext {
	var p = new(AgrupTypesContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_agrupTypes
	return p
}

func (*AgrupTypesContext) IsAgrupTypesContext() {}

func NewAgrupTypesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AgrupTypesContext {
	var p = new(AgrupTypesContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_agrupTypes

	return p
}

func (s *AgrupTypesContext) GetParser() antlr.Parser { return s.parser }

func (s *AgrupTypesContext) AVG() antlr.TerminalNode {
	return s.GetToken(MqlParserAVG, 0)
}

func (s *AgrupTypesContext) COUNT() antlr.TerminalNode {
	return s.GetToken(MqlParserCOUNT, 0)
}

func (s *AgrupTypesContext) DISTINCT_COUNT() antlr.TerminalNode {
	return s.GetToken(MqlParserDISTINCT_COUNT, 0)
}

func (s *AgrupTypesContext) ESTDC() antlr.TerminalNode {
	return s.GetToken(MqlParserESTDC, 0)
}

func (s *AgrupTypesContext) ESTDC_ERROR() antlr.TerminalNode {
	return s.GetToken(MqlParserESTDC_ERROR, 0)
}

func (s *AgrupTypesContext) MAX() antlr.TerminalNode {
	return s.GetToken(MqlParserMAX, 0)
}

func (s *AgrupTypesContext) MEDIAN() antlr.TerminalNode {
	return s.GetToken(MqlParserMEDIAN, 0)
}

func (s *AgrupTypesContext) MIN() antlr.TerminalNode {
	return s.GetToken(MqlParserMIN, 0)
}

func (s *AgrupTypesContext) MODE() antlr.TerminalNode {
	return s.GetToken(MqlParserMODE, 0)
}

func (s *AgrupTypesContext) RANGE() antlr.TerminalNode {
	return s.GetToken(MqlParserRANGE, 0)
}

func (s *AgrupTypesContext) STDEV() antlr.TerminalNode {
	return s.GetToken(MqlParserSTDEV, 0)
}

func (s *AgrupTypesContext) STDEVP() antlr.TerminalNode {
	return s.GetToken(MqlParserSTDEVP, 0)
}

func (s *AgrupTypesContext) SUM() antlr.TerminalNode {
	return s.GetToken(MqlParserSUM, 0)
}

func (s *AgrupTypesContext) SUMSQ() antlr.TerminalNode {
	return s.GetToken(MqlParserSUMSQ, 0)
}

func (s *AgrupTypesContext) VAR() antlr.TerminalNode {
	return s.GetToken(MqlParserVAR, 0)
}

func (s *AgrupTypesContext) VARP() antlr.TerminalNode {
	return s.GetToken(MqlParserVARP, 0)
}

func (s *AgrupTypesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AgrupTypesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AgrupTypesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterAgrupTypes(s)
	}
}

func (s *AgrupTypesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitAgrupTypes(s)
	}
}

func (s *AgrupTypesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitAgrupTypes(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) AgrupTypes() (localctx IAgrupTypesContext) {
	localctx = NewAgrupTypesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, MqlParserRULE_agrupTypes)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(49)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-17)&-(0x1f+1)) == 0 && ((1<<uint((_la-17)))&((1<<(MqlParserAVG-17))|(1<<(MqlParserCOUNT-17))|(1<<(MqlParserDISTINCT_COUNT-17))|(1<<(MqlParserESTDC-17))|(1<<(MqlParserESTDC_ERROR-17))|(1<<(MqlParserMAX-17))|(1<<(MqlParserMEDIAN-17))|(1<<(MqlParserMIN-17))|(1<<(MqlParserMODE-17))|(1<<(MqlParserRANGE-17))|(1<<(MqlParserSTDEV-17))|(1<<(MqlParserSTDEVP-17))|(1<<(MqlParserSUM-17))|(1<<(MqlParserSUMSQ-17))|(1<<(MqlParserVAR-17))|(1<<(MqlParserVARP-17)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IAgrupCallContext is an interface to support dynamic dispatch.
type IAgrupCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAgrupCallContext differentiates from other interfaces.
	IsAgrupCallContext()
}

type AgrupCallContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAgrupCallContext() *AgrupCallContext {
	var p = new(AgrupCallContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_agrupCall
	return p
}

func (*AgrupCallContext) IsAgrupCallContext() {}

func NewAgrupCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AgrupCallContext {
	var p = new(AgrupCallContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_agrupCall

	return p
}

func (s *AgrupCallContext) GetParser() antlr.Parser { return s.parser }

func (s *AgrupCallContext) AgrupTypes() IAgrupTypesContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAgrupTypesContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAgrupTypesContext)
}

func (s *AgrupCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(MqlParserLPAREN, 0)
}

func (s *AgrupCallContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, 0)
}

func (s *AgrupCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(MqlParserRPAREN, 0)
}

func (s *AgrupCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AgrupCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AgrupCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterAgrupCall(s)
	}
}

func (s *AgrupCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitAgrupCall(s)
	}
}

func (s *AgrupCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitAgrupCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) AgrupCall() (localctx IAgrupCallContext) {
	localctx = NewAgrupCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, MqlParserRULE_agrupCall)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(51)
		p.AgrupTypes()
	}
	{
		p.SetState(52)
		p.Match(MqlParserLPAREN)
	}
	{
		p.SetState(53)
		p.Match(MqlParserIDENTIFIER)
	}
	{
		p.SetState(54)
		p.Match(MqlParserRPAREN)
	}

	return localctx
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_literal
	return p
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) CopyFrom(ctx *LiteralContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type DecimalLiteralContext struct {
	*LiteralContext
}

func NewDecimalLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DecimalLiteralContext {
	var p = new(DecimalLiteralContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *DecimalLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DecimalLiteralContext) DECIMAL_LITERAL() antlr.TerminalNode {
	return s.GetToken(MqlParserDECIMAL_LITERAL, 0)
}

func (s *DecimalLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterDecimalLiteral(s)
	}
}

func (s *DecimalLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitDecimalLiteral(s)
	}
}

func (s *DecimalLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitDecimalLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type IdentifierContext struct {
	*LiteralContext
}

func NewIdentifierContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentifierContext {
	var p = new(IdentifierContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *IdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifierContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, 0)
}

func (s *IdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterIdentifier(s)
	}
}

func (s *IdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitIdentifier(s)
	}
}

func (s *IdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringLiteralContext struct {
	*LiteralContext
}

func NewStringLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringLiteralContext {
	var p = new(StringLiteralContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(MqlParserSTRING_LITERAL, 0)
}

func (s *StringLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterStringLiteral(s)
	}
}

func (s *StringLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitStringLiteral(s)
	}
}

func (s *StringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitStringLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type FloatLiteralContext struct {
	*LiteralContext
}

func NewFloatLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FloatLiteralContext {
	var p = new(FloatLiteralContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *FloatLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FloatLiteralContext) FLOAT_LITERAL() antlr.TerminalNode {
	return s.GetToken(MqlParserFLOAT_LITERAL, 0)
}

func (s *FloatLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterFloatLiteral(s)
	}
}

func (s *FloatLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitFloatLiteral(s)
	}
}

func (s *FloatLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitFloatLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type BoolLiteralContext struct {
	*LiteralContext
}

func NewBoolLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BoolLiteralContext {
	var p = new(BoolLiteralContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *BoolLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolLiteralContext) BOOL_LITERAL() antlr.TerminalNode {
	return s.GetToken(MqlParserBOOL_LITERAL, 0)
}

func (s *BoolLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterBoolLiteral(s)
	}
}

func (s *BoolLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitBoolLiteral(s)
	}
}

func (s *BoolLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitBoolLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) Literal() (localctx ILiteralContext) {
	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, MqlParserRULE_literal)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(61)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case MqlParserSTRING_LITERAL:
		localctx = NewStringLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(56)
			p.Match(MqlParserSTRING_LITERAL)
		}

	case MqlParserDECIMAL_LITERAL:
		localctx = NewDecimalLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(57)
			p.Match(MqlParserDECIMAL_LITERAL)
		}

	case MqlParserFLOAT_LITERAL:
		localctx = NewFloatLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(58)
			p.Match(MqlParserFLOAT_LITERAL)
		}

	case MqlParserBOOL_LITERAL:
		localctx = NewBoolLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(59)
			p.Match(MqlParserBOOL_LITERAL)
		}

	case MqlParserIDENTIFIER:
		localctx = NewIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(60)
			p.Match(MqlParserIDENTIFIER)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) CopyFrom(ctx *ExpressionContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type BinaryExpressionContext struct {
	*ExpressionContext
	left  IExpressionContext
	op    IBinaryContext
	right IExpressionContext
}

func NewBinaryExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BinaryExpressionContext {
	var p = new(BinaryExpressionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *BinaryExpressionContext) GetLeft() IExpressionContext { return s.left }

func (s *BinaryExpressionContext) GetOp() IBinaryContext { return s.op }

func (s *BinaryExpressionContext) GetRight() IExpressionContext { return s.right }

func (s *BinaryExpressionContext) SetLeft(v IExpressionContext) { s.left = v }

func (s *BinaryExpressionContext) SetOp(v IBinaryContext) { s.op = v }

func (s *BinaryExpressionContext) SetRight(v IExpressionContext) { s.right = v }

func (s *BinaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryExpressionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *BinaryExpressionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *BinaryExpressionContext) Binary() IBinaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBinaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBinaryContext)
}

func (s *BinaryExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterBinaryExpression(s)
	}
}

func (s *BinaryExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitBinaryExpression(s)
	}
}

func (s *BinaryExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitBinaryExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type ParenExpressionContext struct {
	*ExpressionContext
}

func NewParenExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenExpressionContext {
	var p = new(ParenExpressionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ParenExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenExpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(MqlParserLPAREN, 0)
}

func (s *ParenExpressionContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParenExpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(MqlParserRPAREN, 0)
}

func (s *ParenExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterParenExpression(s)
	}
}

func (s *ParenExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitParenExpression(s)
	}
}

func (s *ParenExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitParenExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type ComparatorExpressionContext struct {
	*ExpressionContext
	left  antlr.Token
	op    IComparatorContext
	right ILiteralContext
}

func NewComparatorExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ComparatorExpressionContext {
	var p = new(ComparatorExpressionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ComparatorExpressionContext) GetLeft() antlr.Token { return s.left }

func (s *ComparatorExpressionContext) SetLeft(v antlr.Token) { s.left = v }

func (s *ComparatorExpressionContext) GetOp() IComparatorContext { return s.op }

func (s *ComparatorExpressionContext) GetRight() ILiteralContext { return s.right }

func (s *ComparatorExpressionContext) SetOp(v IComparatorContext) { s.op = v }

func (s *ComparatorExpressionContext) SetRight(v ILiteralContext) { s.right = v }

func (s *ComparatorExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparatorExpressionContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, 0)
}

func (s *ComparatorExpressionContext) Comparator() IComparatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IComparatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IComparatorContext)
}

func (s *ComparatorExpressionContext) Literal() ILiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *ComparatorExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterComparatorExpression(s)
	}
}

func (s *ComparatorExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitComparatorExpression(s)
	}
}

func (s *ComparatorExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitComparatorExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *MqlParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 10
	p.EnterRecursionRule(localctx, 10, MqlParserRULE_expression, _p)

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(72)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case MqlParserLPAREN:
		localctx = NewParenExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(64)
			p.Match(MqlParserLPAREN)
		}
		{
			p.SetState(65)
			p.expression(0)
		}
		{
			p.SetState(66)
			p.Match(MqlParserRPAREN)
		}

	case MqlParserIDENTIFIER:
		localctx = NewComparatorExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(68)

			var _m = p.Match(MqlParserIDENTIFIER)

			localctx.(*ComparatorExpressionContext).left = _m
		}
		{
			p.SetState(69)

			var _x = p.Comparator()

			localctx.(*ComparatorExpressionContext).op = _x
		}
		{
			p.SetState(70)

			var _x = p.Literal()

			localctx.(*ComparatorExpressionContext).right = _x
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(80)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewBinaryExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
			localctx.(*BinaryExpressionContext).left = _prevctx

			p.PushNewRecursionContext(localctx, _startState, MqlParserRULE_expression)
			p.SetState(74)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(75)

				var _x = p.Binary()

				localctx.(*BinaryExpressionContext).op = _x
			}
			{
				p.SetState(76)

				var _x = p.expression(2)

				localctx.(*BinaryExpressionContext).right = _x
			}

		}
		p.SetState(82)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}

	return localctx
}

// IComparatorContext is an interface to support dynamic dispatch.
type IComparatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsComparatorContext differentiates from other interfaces.
	IsComparatorContext()
}

type ComparatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparatorContext() *ComparatorContext {
	var p = new(ComparatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_comparator
	return p
}

func (*ComparatorContext) IsComparatorContext() {}

func NewComparatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparatorContext {
	var p = new(ComparatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_comparator

	return p
}

func (s *ComparatorContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparatorContext) GT() antlr.TerminalNode {
	return s.GetToken(MqlParserGT, 0)
}

func (s *ComparatorContext) GE() antlr.TerminalNode {
	return s.GetToken(MqlParserGE, 0)
}

func (s *ComparatorContext) LT() antlr.TerminalNode {
	return s.GetToken(MqlParserLT, 0)
}

func (s *ComparatorContext) LE() antlr.TerminalNode {
	return s.GetToken(MqlParserLE, 0)
}

func (s *ComparatorContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(MqlParserASSIGN, 0)
}

func (s *ComparatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterComparator(s)
	}
}

func (s *ComparatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitComparator(s)
	}
}

func (s *ComparatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitComparator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) Comparator() (localctx IComparatorContext) {
	localctx = NewComparatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, MqlParserRULE_comparator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(83)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-42)&-(0x1f+1)) == 0 && ((1<<uint((_la-42)))&((1<<(MqlParserASSIGN-42))|(1<<(MqlParserGT-42))|(1<<(MqlParserLT-42))|(1<<(MqlParserLE-42))|(1<<(MqlParserGE-42)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IBinaryContext is an interface to support dynamic dispatch.
type IBinaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBinaryContext differentiates from other interfaces.
	IsBinaryContext()
}

type BinaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBinaryContext() *BinaryContext {
	var p = new(BinaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_binary
	return p
}

func (*BinaryContext) IsBinaryContext() {}

func NewBinaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinaryContext {
	var p = new(BinaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_binary

	return p
}

func (s *BinaryContext) GetParser() antlr.Parser { return s.parser }

func (s *BinaryContext) AND() antlr.TerminalNode {
	return s.GetToken(MqlParserAND, 0)
}

func (s *BinaryContext) OR() antlr.TerminalNode {
	return s.GetToken(MqlParserOR, 0)
}

func (s *BinaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BinaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterBinary(s)
	}
}

func (s *BinaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitBinary(s)
	}
}

func (s *BinaryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitBinary(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) Binary() (localctx IBinaryContext) {
	localctx = NewBinaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, MqlParserRULE_binary)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(85)
		_la = p.GetTokenStream().LA(1)

		if !(_la == MqlParserAND || _la == MqlParserOR) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// ICommandsContext is an interface to support dynamic dispatch.
type ICommandsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCommandsContext differentiates from other interfaces.
	IsCommandsContext()
}

type CommandsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommandsContext() *CommandsContext {
	var p = new(CommandsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_commands
	return p
}

func (*CommandsContext) IsCommandsContext() {}

func NewCommandsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandsContext {
	var p = new(CommandsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_commands

	return p
}

func (s *CommandsContext) GetParser() antlr.Parser { return s.parser }

func (s *CommandsContext) WhereCommand() IWhereCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWhereCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWhereCommandContext)
}

func (s *CommandsContext) RenameCommand() IRenameCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRenameCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRenameCommandContext)
}

func (s *CommandsContext) StatCommand() IStatCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStatCommandContext)
}

func (s *CommandsContext) FieldCommand() IFieldCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldCommandContext)
}

func (s *CommandsContext) DedupCommand() IDedupCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDedupCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDedupCommandContext)
}

func (s *CommandsContext) SortCommand() ISortCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISortCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISortCommandContext)
}

func (s *CommandsContext) HeadCommand() IHeadCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHeadCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHeadCommandContext)
}

func (s *CommandsContext) BinCommand() IBinCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBinCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBinCommandContext)
}

func (s *CommandsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommandsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommandsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterCommands(s)
	}
}

func (s *CommandsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitCommands(s)
	}
}

func (s *CommandsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitCommands(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) Commands() (localctx ICommandsContext) {
	localctx = NewCommandsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, MqlParserRULE_commands)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(95)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case MqlParserWHERE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(87)
			p.WhereCommand()
		}

	case MqlParserRENAME:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(88)
			p.RenameCommand()
		}

	case MqlParserSTATS:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(89)
			p.StatCommand()
		}

	case MqlParserFIELDS:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(90)
			p.FieldCommand()
		}

	case MqlParserDEDUP:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(91)
			p.DedupCommand()
		}

	case MqlParserSORT:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(92)
			p.SortCommand()
		}

	case MqlParserHEAD:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(93)
			p.HeadCommand()
		}

	case MqlParserBIN:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(94)
			p.BinCommand()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IWhereCommandContext is an interface to support dynamic dispatch.
type IWhereCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWhereCommandContext differentiates from other interfaces.
	IsWhereCommandContext()
}

type WhereCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhereCommandContext() *WhereCommandContext {
	var p = new(WhereCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_whereCommand
	return p
}

func (*WhereCommandContext) IsWhereCommandContext() {}

func NewWhereCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereCommandContext {
	var p = new(WhereCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_whereCommand

	return p
}

func (s *WhereCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereCommandContext) WHERE() antlr.TerminalNode {
	return s.GetToken(MqlParserWHERE, 0)
}

func (s *WhereCommandContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *WhereCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WhereCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterWhereCommand(s)
	}
}

func (s *WhereCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitWhereCommand(s)
	}
}

func (s *WhereCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitWhereCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) WhereCommand() (localctx IWhereCommandContext) {
	localctx = NewWhereCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, MqlParserRULE_whereCommand)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(97)
		p.Match(MqlParserWHERE)
	}
	{
		p.SetState(98)
		p.expression(0)
	}

	return localctx
}

// ISelectCommandContext is an interface to support dynamic dispatch.
type ISelectCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSelectCommandContext differentiates from other interfaces.
	IsSelectCommandContext()
}

type SelectCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectCommandContext() *SelectCommandContext {
	var p = new(SelectCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_selectCommand
	return p
}

func (*SelectCommandContext) IsSelectCommandContext() {}

func NewSelectCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectCommandContext {
	var p = new(SelectCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_selectCommand

	return p
}

func (s *SelectCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectCommandContext) SOURCE_TYPE() antlr.TerminalNode {
	return s.GetToken(MqlParserSOURCE_TYPE, 0)
}

func (s *SelectCommandContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(MqlParserASSIGN, 0)
}

func (s *SelectCommandContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, 0)
}

func (s *SelectCommandContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *SelectCommandContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SelectCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterSelectCommand(s)
	}
}

func (s *SelectCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitSelectCommand(s)
	}
}

func (s *SelectCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitSelectCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) SelectCommand() (localctx ISelectCommandContext) {
	localctx = NewSelectCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, MqlParserRULE_selectCommand)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(100)
		p.Match(MqlParserSOURCE_TYPE)
	}
	{
		p.SetState(101)
		p.Match(MqlParserASSIGN)
	}
	{
		p.SetState(102)
		p.Match(MqlParserIDENTIFIER)
	}
	p.SetState(106)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserLPAREN || _la == MqlParserIDENTIFIER {
		{
			p.SetState(103)
			p.expression(0)
		}

		p.SetState(108)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IRenameCommandContext is an interface to support dynamic dispatch.
type IRenameCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRenameCommandContext differentiates from other interfaces.
	IsRenameCommandContext()
}

type RenameCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRenameCommandContext() *RenameCommandContext {
	var p = new(RenameCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_renameCommand
	return p
}

func (*RenameCommandContext) IsRenameCommandContext() {}

func NewRenameCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RenameCommandContext {
	var p = new(RenameCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_renameCommand

	return p
}

func (s *RenameCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *RenameCommandContext) RENAME() antlr.TerminalNode {
	return s.GetToken(MqlParserRENAME, 0)
}

func (s *RenameCommandContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(MqlParserIDENTIFIER)
}

func (s *RenameCommandContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, i)
}

func (s *RenameCommandContext) AllAS() []antlr.TerminalNode {
	return s.GetTokens(MqlParserAS)
}

func (s *RenameCommandContext) AS(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserAS, i)
}

func (s *RenameCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RenameCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RenameCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterRenameCommand(s)
	}
}

func (s *RenameCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitRenameCommand(s)
	}
}

func (s *RenameCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitRenameCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) RenameCommand() (localctx IRenameCommandContext) {
	localctx = NewRenameCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, MqlParserRULE_renameCommand)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(109)
		p.Match(MqlParserRENAME)
	}
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == MqlParserIDENTIFIER {
		{
			p.SetState(110)
			p.Match(MqlParserIDENTIFIER)
		}
		{
			p.SetState(111)
			p.Match(MqlParserAS)
		}
		{
			p.SetState(112)
			p.Match(MqlParserIDENTIFIER)
		}

		p.SetState(115)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IStatCommandContext is an interface to support dynamic dispatch.
type IStatCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatCommandContext differentiates from other interfaces.
	IsStatCommandContext()
}

type StatCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatCommandContext() *StatCommandContext {
	var p = new(StatCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_statCommand
	return p
}

func (*StatCommandContext) IsStatCommandContext() {}

func NewStatCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatCommandContext {
	var p = new(StatCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_statCommand

	return p
}

func (s *StatCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *StatCommandContext) STATS() antlr.TerminalNode {
	return s.GetToken(MqlParserSTATS, 0)
}

func (s *StatCommandContext) AllAgrupCall() []IAgrupCallContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAgrupCallContext)(nil)).Elem())
	var tst = make([]IAgrupCallContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAgrupCallContext)
		}
	}

	return tst
}

func (s *StatCommandContext) AgrupCall(i int) IAgrupCallContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAgrupCallContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAgrupCallContext)
}

func (s *StatCommandContext) BY() antlr.TerminalNode {
	return s.GetToken(MqlParserBY, 0)
}

func (s *StatCommandContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(MqlParserIDENTIFIER)
}

func (s *StatCommandContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, i)
}

func (s *StatCommandContext) AllAS() []antlr.TerminalNode {
	return s.GetTokens(MqlParserAS)
}

func (s *StatCommandContext) AS(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserAS, i)
}

func (s *StatCommandContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(MqlParserCOMMA)
}

func (s *StatCommandContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserCOMMA, i)
}

func (s *StatCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterStatCommand(s)
	}
}

func (s *StatCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitStatCommand(s)
	}
}

func (s *StatCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitStatCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) StatCommand() (localctx IStatCommandContext) {
	localctx = NewStatCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, MqlParserRULE_statCommand)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(117)
		p.Match(MqlParserSTATS)
	}
	{
		p.SetState(118)
		p.AgrupCall()
	}
	p.SetState(121)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == MqlParserAS {
		{
			p.SetState(119)
			p.Match(MqlParserAS)
		}
		{
			p.SetState(120)
			p.Match(MqlParserIDENTIFIER)
		}

	}
	p.SetState(131)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserCOMMA {
		{
			p.SetState(123)
			p.Match(MqlParserCOMMA)
		}
		{
			p.SetState(124)
			p.AgrupCall()
		}
		p.SetState(127)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == MqlParserAS {
			{
				p.SetState(125)
				p.Match(MqlParserAS)
			}
			{
				p.SetState(126)
				p.Match(MqlParserIDENTIFIER)
			}

		}

		p.SetState(133)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(134)
		p.Match(MqlParserBY)
	}
	{
		p.SetState(135)
		p.Match(MqlParserIDENTIFIER)
	}

	return localctx
}

// IBinCommandContext is an interface to support dynamic dispatch.
type IBinCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBinCommandContext differentiates from other interfaces.
	IsBinCommandContext()
}

type BinCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBinCommandContext() *BinCommandContext {
	var p = new(BinCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_binCommand
	return p
}

func (*BinCommandContext) IsBinCommandContext() {}

func NewBinCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinCommandContext {
	var p = new(BinCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_binCommand

	return p
}

func (s *BinCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *BinCommandContext) BIN() antlr.TerminalNode {
	return s.GetToken(MqlParserBIN, 0)
}

func (s *BinCommandContext) BIN_SPAN() antlr.TerminalNode {
	return s.GetToken(MqlParserBIN_SPAN, 0)
}

func (s *BinCommandContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(MqlParserASSIGN, 0)
}

func (s *BinCommandContext) TIME_LITERAL() antlr.TerminalNode {
	return s.GetToken(MqlParserTIME_LITERAL, 0)
}

func (s *BinCommandContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, 0)
}

func (s *BinCommandContext) TIME_FIELD() antlr.TerminalNode {
	return s.GetToken(MqlParserTIME_FIELD, 0)
}

func (s *BinCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BinCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterBinCommand(s)
	}
}

func (s *BinCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitBinCommand(s)
	}
}

func (s *BinCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitBinCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) BinCommand() (localctx IBinCommandContext) {
	localctx = NewBinCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, MqlParserRULE_binCommand)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(137)
		p.Match(MqlParserBIN)
	}
	p.SetState(139)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == MqlParserTIME_FIELD || _la == MqlParserIDENTIFIER {
		{
			p.SetState(138)
			_la = p.GetTokenStream().LA(1)

			if !(_la == MqlParserTIME_FIELD || _la == MqlParserIDENTIFIER) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(141)
		p.Match(MqlParserBIN_SPAN)
	}
	{
		p.SetState(142)
		p.Match(MqlParserASSIGN)
	}
	{
		p.SetState(143)
		p.Match(MqlParserTIME_LITERAL)
	}

	return localctx
}

// IFieldCommandContext is an interface to support dynamic dispatch.
type IFieldCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldCommandContext differentiates from other interfaces.
	IsFieldCommandContext()
}

type FieldCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldCommandContext() *FieldCommandContext {
	var p = new(FieldCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_fieldCommand
	return p
}

func (*FieldCommandContext) IsFieldCommandContext() {}

func NewFieldCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldCommandContext {
	var p = new(FieldCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_fieldCommand

	return p
}

func (s *FieldCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldCommandContext) FIELDS() antlr.TerminalNode {
	return s.GetToken(MqlParserFIELDS, 0)
}

func (s *FieldCommandContext) IdentifierList() IIdentifierListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentifierListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentifierListContext)
}

func (s *FieldCommandContext) ADD() antlr.TerminalNode {
	return s.GetToken(MqlParserADD, 0)
}

func (s *FieldCommandContext) SUB() antlr.TerminalNode {
	return s.GetToken(MqlParserSUB, 0)
}

func (s *FieldCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterFieldCommand(s)
	}
}

func (s *FieldCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitFieldCommand(s)
	}
}

func (s *FieldCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitFieldCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) FieldCommand() (localctx IFieldCommandContext) {
	localctx = NewFieldCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, MqlParserRULE_fieldCommand)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(145)
		p.Match(MqlParserFIELDS)
	}
	p.SetState(147)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == MqlParserADD || _la == MqlParserSUB {
		{
			p.SetState(146)
			_la = p.GetTokenStream().LA(1)

			if !(_la == MqlParserADD || _la == MqlParserSUB) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(149)
		p.IdentifierList()
	}

	return localctx
}

// IDedupCommandContext is an interface to support dynamic dispatch.
type IDedupCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDedupCommandContext differentiates from other interfaces.
	IsDedupCommandContext()
}

type DedupCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDedupCommandContext() *DedupCommandContext {
	var p = new(DedupCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_dedupCommand
	return p
}

func (*DedupCommandContext) IsDedupCommandContext() {}

func NewDedupCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DedupCommandContext {
	var p = new(DedupCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_dedupCommand

	return p
}

func (s *DedupCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *DedupCommandContext) DEDUP() antlr.TerminalNode {
	return s.GetToken(MqlParserDEDUP, 0)
}

func (s *DedupCommandContext) IdentifierList() IIdentifierListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentifierListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentifierListContext)
}

func (s *DedupCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DedupCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DedupCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterDedupCommand(s)
	}
}

func (s *DedupCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitDedupCommand(s)
	}
}

func (s *DedupCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitDedupCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) DedupCommand() (localctx IDedupCommandContext) {
	localctx = NewDedupCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, MqlParserRULE_dedupCommand)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(151)
		p.Match(MqlParserDEDUP)
	}
	{
		p.SetState(152)
		p.IdentifierList()
	}

	return localctx
}

// ISortCommandContext is an interface to support dynamic dispatch.
type ISortCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSortCommandContext differentiates from other interfaces.
	IsSortCommandContext()
}

type SortCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortCommandContext() *SortCommandContext {
	var p = new(SortCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_sortCommand
	return p
}

func (*SortCommandContext) IsSortCommandContext() {}

func NewSortCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortCommandContext {
	var p = new(SortCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_sortCommand

	return p
}

func (s *SortCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *SortCommandContext) SORT() antlr.TerminalNode {
	return s.GetToken(MqlParserSORT, 0)
}

func (s *SortCommandContext) IdentifierList() IIdentifierListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentifierListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentifierListContext)
}

func (s *SortCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterSortCommand(s)
	}
}

func (s *SortCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitSortCommand(s)
	}
}

func (s *SortCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitSortCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) SortCommand() (localctx ISortCommandContext) {
	localctx = NewSortCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, MqlParserRULE_sortCommand)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(154)
		p.Match(MqlParserSORT)
	}
	{
		p.SetState(155)
		p.IdentifierList()
	}

	return localctx
}

// IHeadCommandContext is an interface to support dynamic dispatch.
type IHeadCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHeadCommandContext differentiates from other interfaces.
	IsHeadCommandContext()
}

type HeadCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHeadCommandContext() *HeadCommandContext {
	var p = new(HeadCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_headCommand
	return p
}

func (*HeadCommandContext) IsHeadCommandContext() {}

func NewHeadCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HeadCommandContext {
	var p = new(HeadCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_headCommand

	return p
}

func (s *HeadCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *HeadCommandContext) HEAD() antlr.TerminalNode {
	return s.GetToken(MqlParserHEAD, 0)
}

func (s *HeadCommandContext) DECIMAL_LITERAL() antlr.TerminalNode {
	return s.GetToken(MqlParserDECIMAL_LITERAL, 0)
}

func (s *HeadCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HeadCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HeadCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterHeadCommand(s)
	}
}

func (s *HeadCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitHeadCommand(s)
	}
}

func (s *HeadCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitHeadCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) HeadCommand() (localctx IHeadCommandContext) {
	localctx = NewHeadCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, MqlParserRULE_headCommand)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(157)
		p.Match(MqlParserHEAD)
	}
	{
		p.SetState(158)
		p.Match(MqlParserDECIMAL_LITERAL)
	}

	return localctx
}

// ICompleteCommandContext is an interface to support dynamic dispatch.
type ICompleteCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCompleteCommandContext differentiates from other interfaces.
	IsCompleteCommandContext()
}

type CompleteCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompleteCommandContext() *CompleteCommandContext {
	var p = new(CompleteCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_completeCommand
	return p
}

func (*CompleteCommandContext) IsCompleteCommandContext() {}

func NewCompleteCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompleteCommandContext {
	var p = new(CompleteCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_completeCommand

	return p
}

func (s *CompleteCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *CompleteCommandContext) SelectCommand() ISelectCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISelectCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISelectCommandContext)
}

func (s *CompleteCommandContext) AllBITOR() []antlr.TerminalNode {
	return s.GetTokens(MqlParserBITOR)
}

func (s *CompleteCommandContext) BITOR(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserBITOR, i)
}

func (s *CompleteCommandContext) AllCommands() []ICommandsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ICommandsContext)(nil)).Elem())
	var tst = make([]ICommandsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ICommandsContext)
		}
	}

	return tst
}

func (s *CompleteCommandContext) Commands(i int) ICommandsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICommandsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ICommandsContext)
}

func (s *CompleteCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompleteCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompleteCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterCompleteCommand(s)
	}
}

func (s *CompleteCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitCompleteCommand(s)
	}
}

func (s *CompleteCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MqlParserVisitor:
		return t.VisitCompleteCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MqlParser) CompleteCommand() (localctx ICompleteCommandContext) {
	localctx = NewCompleteCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, MqlParserRULE_completeCommand)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(160)
		p.SelectCommand()
	}
	p.SetState(165)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserBITOR {
		{
			p.SetState(161)
			p.Match(MqlParserBITOR)
		}
		{
			p.SetState(162)
			p.Commands()
		}

		p.SetState(167)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

func (p *MqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 5:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *MqlParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
