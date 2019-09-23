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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 79, 195,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 3, 2, 3, 2,
	3, 2, 3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 5,
	3, 5, 5, 5, 60, 10, 5, 3, 6, 3, 6, 3, 6, 7, 6, 65, 10, 6, 12, 6, 14, 6,
	68, 11, 6, 3, 7, 3, 7, 7, 7, 72, 10, 7, 12, 7, 14, 7, 75, 11, 7, 3, 8,
	3, 8, 3, 8, 7, 8, 80, 10, 8, 12, 8, 14, 8, 83, 11, 8, 7, 8, 85, 10, 8,
	12, 8, 14, 8, 88, 11, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9,
	3, 9, 5, 9, 99, 10, 9, 3, 9, 3, 9, 3, 9, 3, 9, 7, 9, 105, 10, 9, 12, 9,
	14, 9, 108, 11, 9, 3, 10, 3, 10, 3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 3,
	12, 3, 12, 3, 12, 3, 12, 3, 12, 5, 12, 122, 10, 12, 3, 13, 3, 13, 3, 13,
	3, 14, 3, 14, 3, 14, 3, 14, 7, 14, 131, 10, 14, 12, 14, 14, 14, 134, 11,
	14, 3, 15, 3, 15, 3, 15, 3, 15, 6, 15, 140, 10, 15, 13, 15, 14, 15, 141,
	3, 16, 3, 16, 3, 16, 3, 16, 5, 16, 148, 10, 16, 3, 16, 3, 16, 3, 16, 3,
	16, 5, 16, 154, 10, 16, 7, 16, 156, 10, 16, 12, 16, 14, 16, 159, 11, 16,
	3, 16, 3, 16, 3, 16, 3, 17, 3, 17, 5, 17, 166, 10, 17, 3, 17, 3, 17, 3,
	17, 3, 17, 3, 18, 3, 18, 5, 18, 174, 10, 18, 3, 18, 3, 18, 3, 19, 3, 19,
	3, 19, 3, 20, 3, 20, 3, 20, 3, 21, 3, 21, 3, 21, 3, 22, 3, 22, 3, 22, 7,
	22, 190, 10, 22, 12, 22, 14, 22, 193, 11, 22, 3, 22, 2, 3, 16, 23, 2, 4,
	6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42,
	2, 8, 3, 2, 20, 35, 3, 2, 18, 19, 4, 2, 45, 47, 53, 54, 3, 2, 16, 17, 3,
	2, 75, 76, 3, 2, 56, 57, 2, 198, 2, 44, 3, 2, 2, 2, 4, 47, 3, 2, 2, 2,
	6, 49, 3, 2, 2, 2, 8, 59, 3, 2, 2, 2, 10, 61, 3, 2, 2, 2, 12, 69, 3, 2,
	2, 2, 14, 76, 3, 2, 2, 2, 16, 98, 3, 2, 2, 2, 18, 109, 3, 2, 2, 2, 20,
	111, 3, 2, 2, 2, 22, 121, 3, 2, 2, 2, 24, 123, 3, 2, 2, 2, 26, 126, 3,
	2, 2, 2, 28, 135, 3, 2, 2, 2, 30, 143, 3, 2, 2, 2, 32, 163, 3, 2, 2, 2,
	34, 171, 3, 2, 2, 2, 36, 177, 3, 2, 2, 2, 38, 180, 3, 2, 2, 2, 40, 183,
	3, 2, 2, 2, 42, 186, 3, 2, 2, 2, 44, 45, 5, 42, 22, 2, 45, 46, 7, 2, 2,
	3, 46, 3, 3, 2, 2, 2, 47, 48, 9, 2, 2, 2, 48, 5, 3, 2, 2, 2, 49, 50, 5,
	4, 3, 2, 50, 51, 7, 36, 2, 2, 51, 52, 7, 76, 2, 2, 52, 53, 7, 37, 2, 2,
	53, 7, 3, 2, 2, 2, 54, 60, 7, 73, 2, 2, 55, 60, 7, 64, 2, 2, 56, 60, 7,
	69, 2, 2, 57, 60, 7, 71, 2, 2, 58, 60, 7, 76, 2, 2, 59, 54, 3, 2, 2, 2,
	59, 55, 3, 2, 2, 2, 59, 56, 3, 2, 2, 2, 59, 57, 3, 2, 2, 2, 59, 58, 3,
	2, 2, 2, 60, 9, 3, 2, 2, 2, 61, 66, 7, 76, 2, 2, 62, 63, 7, 43, 2, 2, 63,
	65, 7, 76, 2, 2, 64, 62, 3, 2, 2, 2, 65, 68, 3, 2, 2, 2, 66, 64, 3, 2,
	2, 2, 66, 67, 3, 2, 2, 2, 67, 11, 3, 2, 2, 2, 68, 66, 3, 2, 2, 2, 69, 73,
	7, 76, 2, 2, 70, 72, 9, 3, 2, 2, 71, 70, 3, 2, 2, 2, 72, 75, 3, 2, 2, 2,
	73, 71, 3, 2, 2, 2, 73, 74, 3, 2, 2, 2, 74, 13, 3, 2, 2, 2, 75, 73, 3,
	2, 2, 2, 76, 86, 5, 12, 7, 2, 77, 81, 7, 43, 2, 2, 78, 80, 5, 12, 7, 2,
	79, 78, 3, 2, 2, 2, 80, 83, 3, 2, 2, 2, 81, 79, 3, 2, 2, 2, 81, 82, 3,
	2, 2, 2, 82, 85, 3, 2, 2, 2, 83, 81, 3, 2, 2, 2, 84, 77, 3, 2, 2, 2, 85,
	88, 3, 2, 2, 2, 86, 84, 3, 2, 2, 2, 86, 87, 3, 2, 2, 2, 87, 15, 3, 2, 2,
	2, 88, 86, 3, 2, 2, 2, 89, 90, 8, 9, 1, 2, 90, 91, 7, 36, 2, 2, 91, 92,
	5, 16, 9, 2, 92, 93, 7, 37, 2, 2, 93, 99, 3, 2, 2, 2, 94, 95, 5, 8, 5,
	2, 95, 96, 5, 18, 10, 2, 96, 97, 5, 8, 5, 2, 97, 99, 3, 2, 2, 2, 98, 89,
	3, 2, 2, 2, 98, 94, 3, 2, 2, 2, 99, 106, 3, 2, 2, 2, 100, 101, 12, 3, 2,
	2, 101, 102, 5, 20, 11, 2, 102, 103, 5, 16, 9, 4, 103, 105, 3, 2, 2, 2,
	104, 100, 3, 2, 2, 2, 105, 108, 3, 2, 2, 2, 106, 104, 3, 2, 2, 2, 106,
	107, 3, 2, 2, 2, 107, 17, 3, 2, 2, 2, 108, 106, 3, 2, 2, 2, 109, 110, 9,
	4, 2, 2, 110, 19, 3, 2, 2, 2, 111, 112, 9, 5, 2, 2, 112, 21, 3, 2, 2, 2,
	113, 122, 5, 24, 13, 2, 114, 122, 5, 28, 15, 2, 115, 122, 5, 30, 16, 2,
	116, 122, 5, 34, 18, 2, 117, 122, 5, 36, 19, 2, 118, 122, 5, 38, 20, 2,
	119, 122, 5, 40, 21, 2, 120, 122, 5, 32, 17, 2, 121, 113, 3, 2, 2, 2, 121,
	114, 3, 2, 2, 2, 121, 115, 3, 2, 2, 2, 121, 116, 3, 2, 2, 2, 121, 117,
	3, 2, 2, 2, 121, 118, 3, 2, 2, 2, 121, 119, 3, 2, 2, 2, 121, 120, 3, 2,
	2, 2, 122, 23, 3, 2, 2, 2, 123, 124, 7, 7, 2, 2, 124, 125, 5, 16, 9, 2,
	125, 25, 3, 2, 2, 2, 126, 127, 7, 3, 2, 2, 127, 128, 7, 45, 2, 2, 128,
	132, 7, 76, 2, 2, 129, 131, 5, 16, 9, 2, 130, 129, 3, 2, 2, 2, 131, 134,
	3, 2, 2, 2, 132, 130, 3, 2, 2, 2, 132, 133, 3, 2, 2, 2, 133, 27, 3, 2,
	2, 2, 134, 132, 3, 2, 2, 2, 135, 139, 7, 4, 2, 2, 136, 137, 7, 76, 2, 2,
	137, 138, 7, 14, 2, 2, 138, 140, 7, 76, 2, 2, 139, 136, 3, 2, 2, 2, 140,
	141, 3, 2, 2, 2, 141, 139, 3, 2, 2, 2, 141, 142, 3, 2, 2, 2, 142, 29, 3,
	2, 2, 2, 143, 144, 7, 13, 2, 2, 144, 147, 5, 6, 4, 2, 145, 146, 7, 14,
	2, 2, 146, 148, 7, 76, 2, 2, 147, 145, 3, 2, 2, 2, 147, 148, 3, 2, 2, 2,
	148, 157, 3, 2, 2, 2, 149, 150, 7, 43, 2, 2, 150, 153, 5, 6, 4, 2, 151,
	152, 7, 14, 2, 2, 152, 154, 7, 76, 2, 2, 153, 151, 3, 2, 2, 2, 153, 154,
	3, 2, 2, 2, 154, 156, 3, 2, 2, 2, 155, 149, 3, 2, 2, 2, 156, 159, 3, 2,
	2, 2, 157, 155, 3, 2, 2, 2, 157, 158, 3, 2, 2, 2, 158, 160, 3, 2, 2, 2,
	159, 157, 3, 2, 2, 2, 160, 161, 7, 15, 2, 2, 161, 162, 7, 76, 2, 2, 162,
	31, 3, 2, 2, 2, 163, 165, 7, 10, 2, 2, 164, 166, 9, 6, 2, 2, 165, 164,
	3, 2, 2, 2, 165, 166, 3, 2, 2, 2, 166, 167, 3, 2, 2, 2, 167, 168, 7, 11,
	2, 2, 168, 169, 7, 45, 2, 2, 169, 170, 7, 68, 2, 2, 170, 33, 3, 2, 2, 2,
	171, 173, 7, 12, 2, 2, 172, 174, 9, 7, 2, 2, 173, 172, 3, 2, 2, 2, 173,
	174, 3, 2, 2, 2, 174, 175, 3, 2, 2, 2, 175, 176, 5, 10, 6, 2, 176, 35,
	3, 2, 2, 2, 177, 178, 7, 6, 2, 2, 178, 179, 5, 10, 6, 2, 179, 37, 3, 2,
	2, 2, 180, 181, 7, 8, 2, 2, 181, 182, 5, 14, 8, 2, 182, 39, 3, 2, 2, 2,
	183, 184, 7, 9, 2, 2, 184, 185, 7, 64, 2, 2, 185, 41, 3, 2, 2, 2, 186,
	191, 5, 26, 14, 2, 187, 188, 7, 61, 2, 2, 188, 190, 5, 22, 12, 2, 189,
	187, 3, 2, 2, 2, 190, 193, 3, 2, 2, 2, 191, 189, 3, 2, 2, 2, 191, 192,
	3, 2, 2, 2, 192, 43, 3, 2, 2, 2, 193, 191, 3, 2, 2, 2, 18, 59, 66, 73,
	81, 86, 98, 106, 121, 132, 141, 147, 153, 157, 165, 173, 191,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'indexname'", "'rename'", "'seach'", "'dedup'", "'where'", "'sort'",
	"'top'", "'bin'", "'span'", "'fields'", "'stats'", "'as'", "'by'", "'and'",
	"'or'", "'asc'", "'desc'", "'avg'", "'count'", "'distinct_count'", "'estdc'",
	"'estdc_error'", "'max'", "'median'", "'min'", "'mode'", "'range'", "'stdev'",
	"'stdevp'", "'sum'", "'sumsq'", "'var'", "'varp'", "'('", "')'", "'{'",
	"'}'", "'['", "']'", "';'", "','", "'.'", "'='", "'>'", "'<'", "'!'", "'~'",
	"'?'", "':'", "'=='", "'<='", "'>='", "'!='", "'+'", "'-'", "'*'", "'/'",
	"'&'", "'|'", "'^'", "'%'", "", "", "", "", "", "", "", "", "", "", "'null'",
	"'_time'",
}
var symbolicNames = []string{
	"", "INDEX_NAME", "RENAME", "SEARCH", "DEDUP", "WHERE", "SORT", "TOP",
	"BIN", "BIN_SPAN", "FIELDS", "STATS", "AS", "BY", "AND", "OR", "ASC", "DESC",
	"AVG", "COUNT", "DISTINCT_COUNT", "ESTDC", "ESTDC_ERROR", "MAX", "MEDIAN",
	"MIN", "MODE", "RANGE", "STDEV", "STDEVP", "SUM", "SUMSQ", "VAR", "VARP",
	"LPAREN", "RPAREN", "LBRACE", "RBRACE", "LBRACK", "RBRACK", "SEMI", "COMMA",
	"DOT", "ASSIGN", "GT", "LT", "BANG", "TILDE", "QUESTION", "COLON", "EQUAL",
	"LE", "GE", "NOTEQUAL", "ADD", "SUB", "MUL", "DIV", "BITAND", "BITOR",
	"CARET", "MOD", "DECIMAL_LITERAL", "HEX_LITERAL", "OCT_LITERAL", "BINARY_LITERAL",
	"TIME_LITERAL", "FLOAT_LITERAL", "HEX_FLOAT_LITERAL", "BOOL_LITERAL", "CHAR_LITERAL",
	"STRING_LITERAL", "NULL_LITERAL", "TIME_FIELD", "IDENTIFIER", "WS", "COMMENT",
	"LINE_COMMENT",
}

var ruleNames = []string{
	"start", "agrupTypes", "agrupCall", "literal", "identifierList", "sort",
	"sortList", "expression", "comparator", "binary", "commands", "whereCommand",
	"selectCommand", "renameCommand", "statCommand", "binCommand", "fieldCommand",
	"dedupCommand", "sortCommand", "topCommand", "completeCommand",
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
	MqlParserINDEX_NAME        = 1
	MqlParserRENAME            = 2
	MqlParserSEARCH            = 3
	MqlParserDEDUP             = 4
	MqlParserWHERE             = 5
	MqlParserSORT              = 6
	MqlParserTOP               = 7
	MqlParserBIN               = 8
	MqlParserBIN_SPAN          = 9
	MqlParserFIELDS            = 10
	MqlParserSTATS             = 11
	MqlParserAS                = 12
	MqlParserBY                = 13
	MqlParserAND               = 14
	MqlParserOR                = 15
	MqlParserASC               = 16
	MqlParserDESC              = 17
	MqlParserAVG               = 18
	MqlParserCOUNT             = 19
	MqlParserDISTINCT_COUNT    = 20
	MqlParserESTDC             = 21
	MqlParserESTDC_ERROR       = 22
	MqlParserMAX               = 23
	MqlParserMEDIAN            = 24
	MqlParserMIN               = 25
	MqlParserMODE              = 26
	MqlParserRANGE             = 27
	MqlParserSTDEV             = 28
	MqlParserSTDEVP            = 29
	MqlParserSUM               = 30
	MqlParserSUMSQ             = 31
	MqlParserVAR               = 32
	MqlParserVARP              = 33
	MqlParserLPAREN            = 34
	MqlParserRPAREN            = 35
	MqlParserLBRACE            = 36
	MqlParserRBRACE            = 37
	MqlParserLBRACK            = 38
	MqlParserRBRACK            = 39
	MqlParserSEMI              = 40
	MqlParserCOMMA             = 41
	MqlParserDOT               = 42
	MqlParserASSIGN            = 43
	MqlParserGT                = 44
	MqlParserLT                = 45
	MqlParserBANG              = 46
	MqlParserTILDE             = 47
	MqlParserQUESTION          = 48
	MqlParserCOLON             = 49
	MqlParserEQUAL             = 50
	MqlParserLE                = 51
	MqlParserGE                = 52
	MqlParserNOTEQUAL          = 53
	MqlParserADD               = 54
	MqlParserSUB               = 55
	MqlParserMUL               = 56
	MqlParserDIV               = 57
	MqlParserBITAND            = 58
	MqlParserBITOR             = 59
	MqlParserCARET             = 60
	MqlParserMOD               = 61
	MqlParserDECIMAL_LITERAL   = 62
	MqlParserHEX_LITERAL       = 63
	MqlParserOCT_LITERAL       = 64
	MqlParserBINARY_LITERAL    = 65
	MqlParserTIME_LITERAL      = 66
	MqlParserFLOAT_LITERAL     = 67
	MqlParserHEX_FLOAT_LITERAL = 68
	MqlParserBOOL_LITERAL      = 69
	MqlParserCHAR_LITERAL      = 70
	MqlParserSTRING_LITERAL    = 71
	MqlParserNULL_LITERAL      = 72
	MqlParserTIME_FIELD        = 73
	MqlParserIDENTIFIER        = 74
	MqlParserWS                = 75
	MqlParserCOMMENT           = 76
	MqlParserLINE_COMMENT      = 77
)

// MqlParser rules.
const (
	MqlParserRULE_start           = 0
	MqlParserRULE_agrupTypes      = 1
	MqlParserRULE_agrupCall       = 2
	MqlParserRULE_literal         = 3
	MqlParserRULE_identifierList  = 4
	MqlParserRULE_sort            = 5
	MqlParserRULE_sortList        = 6
	MqlParserRULE_expression      = 7
	MqlParserRULE_comparator      = 8
	MqlParserRULE_binary          = 9
	MqlParserRULE_commands        = 10
	MqlParserRULE_whereCommand    = 11
	MqlParserRULE_selectCommand   = 12
	MqlParserRULE_renameCommand   = 13
	MqlParserRULE_statCommand     = 14
	MqlParserRULE_binCommand      = 15
	MqlParserRULE_fieldCommand    = 16
	MqlParserRULE_dedupCommand    = 17
	MqlParserRULE_sortCommand     = 18
	MqlParserRULE_topCommand      = 19
	MqlParserRULE_completeCommand = 20
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
		p.SetState(42)
		p.CompleteCommand()
	}
	{
		p.SetState(43)
		p.Match(MqlParserEOF)
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

func (p *MqlParser) AgrupTypes() (localctx IAgrupTypesContext) {
	localctx = NewAgrupTypesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, MqlParserRULE_agrupTypes)
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
		p.SetState(45)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-18)&-(0x1f+1)) == 0 && ((1<<uint((_la-18)))&((1<<(MqlParserAVG-18))|(1<<(MqlParserCOUNT-18))|(1<<(MqlParserDISTINCT_COUNT-18))|(1<<(MqlParserESTDC-18))|(1<<(MqlParserESTDC_ERROR-18))|(1<<(MqlParserMAX-18))|(1<<(MqlParserMEDIAN-18))|(1<<(MqlParserMIN-18))|(1<<(MqlParserMODE-18))|(1<<(MqlParserRANGE-18))|(1<<(MqlParserSTDEV-18))|(1<<(MqlParserSTDEVP-18))|(1<<(MqlParserSUM-18))|(1<<(MqlParserSUMSQ-18))|(1<<(MqlParserVAR-18))|(1<<(MqlParserVARP-18)))) != 0) {
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

func (p *MqlParser) AgrupCall() (localctx IAgrupCallContext) {
	localctx = NewAgrupCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, MqlParserRULE_agrupCall)

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
		p.SetState(47)
		p.AgrupTypes()
	}
	{
		p.SetState(48)
		p.Match(MqlParserLPAREN)
	}
	{
		p.SetState(49)
		p.Match(MqlParserIDENTIFIER)
	}
	{
		p.SetState(50)
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

func (p *MqlParser) Literal() (localctx ILiteralContext) {
	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, MqlParserRULE_literal)

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

	p.SetState(57)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case MqlParserSTRING_LITERAL:
		localctx = NewStringLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(52)
			p.Match(MqlParserSTRING_LITERAL)
		}

	case MqlParserDECIMAL_LITERAL:
		localctx = NewDecimalLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(53)
			p.Match(MqlParserDECIMAL_LITERAL)
		}

	case MqlParserFLOAT_LITERAL:
		localctx = NewFloatLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(54)
			p.Match(MqlParserFLOAT_LITERAL)
		}

	case MqlParserBOOL_LITERAL:
		localctx = NewBoolLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(55)
			p.Match(MqlParserBOOL_LITERAL)
		}

	case MqlParserIDENTIFIER:
		localctx = NewIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(56)
			p.Match(MqlParserIDENTIFIER)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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

func (p *MqlParser) IdentifierList() (localctx IIdentifierListContext) {
	localctx = NewIdentifierListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, MqlParserRULE_identifierList)
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
		p.SetState(59)
		p.Match(MqlParserIDENTIFIER)
	}
	p.SetState(64)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserCOMMA {
		{
			p.SetState(60)
			p.Match(MqlParserCOMMA)
		}
		{
			p.SetState(61)
			p.Match(MqlParserIDENTIFIER)
		}

		p.SetState(66)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ISortContext is an interface to support dynamic dispatch.
type ISortContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetField returns the field token.
	GetField() antlr.Token

	// GetDirection returns the direction token.
	GetDirection() antlr.Token

	// SetField sets the field token.
	SetField(antlr.Token)

	// SetDirection sets the direction token.
	SetDirection(antlr.Token)

	// IsSortContext differentiates from other interfaces.
	IsSortContext()
}

type SortContext struct {
	*antlr.BaseParserRuleContext
	parser    antlr.Parser
	field     antlr.Token
	direction antlr.Token
}

func NewEmptySortContext() *SortContext {
	var p = new(SortContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_sort
	return p
}

func (*SortContext) IsSortContext() {}

func NewSortContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortContext {
	var p = new(SortContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_sort

	return p
}

func (s *SortContext) GetParser() antlr.Parser { return s.parser }

func (s *SortContext) GetField() antlr.Token { return s.field }

func (s *SortContext) GetDirection() antlr.Token { return s.direction }

func (s *SortContext) SetField(v antlr.Token) { s.field = v }

func (s *SortContext) SetDirection(v antlr.Token) { s.direction = v }

func (s *SortContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, 0)
}

func (s *SortContext) AllASC() []antlr.TerminalNode {
	return s.GetTokens(MqlParserASC)
}

func (s *SortContext) ASC(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserASC, i)
}

func (s *SortContext) AllDESC() []antlr.TerminalNode {
	return s.GetTokens(MqlParserDESC)
}

func (s *SortContext) DESC(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserDESC, i)
}

func (s *SortContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterSort(s)
	}
}

func (s *SortContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitSort(s)
	}
}

func (p *MqlParser) Sort() (localctx ISortContext) {
	localctx = NewSortContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, MqlParserRULE_sort)
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
		p.SetState(67)

		var _m = p.Match(MqlParserIDENTIFIER)

		localctx.(*SortContext).field = _m
	}
	p.SetState(71)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserASC || _la == MqlParserDESC {
		{
			p.SetState(68)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*SortContext).direction = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == MqlParserASC || _la == MqlParserDESC) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*SortContext).direction = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(73)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ISortListContext is an interface to support dynamic dispatch.
type ISortListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSortListContext differentiates from other interfaces.
	IsSortListContext()
}

type SortListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortListContext() *SortListContext {
	var p = new(SortListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_sortList
	return p
}

func (*SortListContext) IsSortListContext() {}

func NewSortListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortListContext {
	var p = new(SortListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_sortList

	return p
}

func (s *SortListContext) GetParser() antlr.Parser { return s.parser }

func (s *SortListContext) AllSort() []ISortContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISortContext)(nil)).Elem())
	var tst = make([]ISortContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISortContext)
		}
	}

	return tst
}

func (s *SortListContext) Sort(i int) ISortContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISortContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISortContext)
}

func (s *SortListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(MqlParserCOMMA)
}

func (s *SortListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserCOMMA, i)
}

func (s *SortListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterSortList(s)
	}
}

func (s *SortListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitSortList(s)
	}
}

func (p *MqlParser) SortList() (localctx ISortListContext) {
	localctx = NewSortListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, MqlParserRULE_sortList)
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
		p.SetState(74)
		p.Sort()
	}
	p.SetState(84)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserCOMMA {
		{
			p.SetState(75)
			p.Match(MqlParserCOMMA)
		}
		p.SetState(79)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == MqlParserIDENTIFIER {
			{
				p.SetState(76)
				p.Sort()
			}

			p.SetState(81)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

		p.SetState(86)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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

type ComparatorExpressionContext struct {
	*ExpressionContext
	left  ILiteralContext
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

func (s *ComparatorExpressionContext) GetLeft() ILiteralContext { return s.left }

func (s *ComparatorExpressionContext) GetOp() IComparatorContext { return s.op }

func (s *ComparatorExpressionContext) GetRight() ILiteralContext { return s.right }

func (s *ComparatorExpressionContext) SetLeft(v ILiteralContext) { s.left = v }

func (s *ComparatorExpressionContext) SetOp(v IComparatorContext) { s.op = v }

func (s *ComparatorExpressionContext) SetRight(v ILiteralContext) { s.right = v }

func (s *ComparatorExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparatorExpressionContext) AllLiteral() []ILiteralContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ILiteralContext)(nil)).Elem())
	var tst = make([]ILiteralContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ILiteralContext)
		}
	}

	return tst
}

func (s *ComparatorExpressionContext) Literal(i int) ILiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILiteralContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *ComparatorExpressionContext) Comparator() IComparatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IComparatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IComparatorContext)
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

func (p *MqlParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *MqlParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 14
	p.EnterRecursionRule(localctx, 14, MqlParserRULE_expression, _p)

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
	p.SetState(96)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case MqlParserLPAREN:
		localctx = NewParenExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(88)
			p.Match(MqlParserLPAREN)
		}
		{
			p.SetState(89)
			p.expression(0)
		}
		{
			p.SetState(90)
			p.Match(MqlParserRPAREN)
		}

	case MqlParserDECIMAL_LITERAL, MqlParserFLOAT_LITERAL, MqlParserBOOL_LITERAL, MqlParserSTRING_LITERAL, MqlParserIDENTIFIER:
		localctx = NewComparatorExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(92)

			var _x = p.Literal()

			localctx.(*ComparatorExpressionContext).left = _x
		}
		{
			p.SetState(93)

			var _x = p.Comparator()

			localctx.(*ComparatorExpressionContext).op = _x
		}
		{
			p.SetState(94)

			var _x = p.Literal()

			localctx.(*ComparatorExpressionContext).right = _x
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(104)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewBinaryExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
			localctx.(*BinaryExpressionContext).left = _prevctx

			p.PushNewRecursionContext(localctx, _startState, MqlParserRULE_expression)
			p.SetState(98)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(99)

				var _x = p.Binary()

				localctx.(*BinaryExpressionContext).op = _x
			}
			{
				p.SetState(100)

				var _x = p.expression(2)

				localctx.(*BinaryExpressionContext).right = _x
			}

		}
		p.SetState(106)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())
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

func (p *MqlParser) Comparator() (localctx IComparatorContext) {
	localctx = NewComparatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, MqlParserRULE_comparator)
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
		p.SetState(107)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-43)&-(0x1f+1)) == 0 && ((1<<uint((_la-43)))&((1<<(MqlParserASSIGN-43))|(1<<(MqlParserGT-43))|(1<<(MqlParserLT-43))|(1<<(MqlParserLE-43))|(1<<(MqlParserGE-43)))) != 0) {
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

func (p *MqlParser) Binary() (localctx IBinaryContext) {
	localctx = NewBinaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, MqlParserRULE_binary)
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

func (s *CommandsContext) TopCommand() ITopCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITopCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITopCommandContext)
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

func (p *MqlParser) Commands() (localctx ICommandsContext) {
	localctx = NewCommandsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, MqlParserRULE_commands)

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

	p.SetState(119)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case MqlParserWHERE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(111)
			p.WhereCommand()
		}

	case MqlParserRENAME:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(112)
			p.RenameCommand()
		}

	case MqlParserSTATS:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(113)
			p.StatCommand()
		}

	case MqlParserFIELDS:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(114)
			p.FieldCommand()
		}

	case MqlParserDEDUP:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(115)
			p.DedupCommand()
		}

	case MqlParserSORT:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(116)
			p.SortCommand()
		}

	case MqlParserTOP:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(117)
			p.TopCommand()
		}

	case MqlParserBIN:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(118)
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

func (p *MqlParser) WhereCommand() (localctx IWhereCommandContext) {
	localctx = NewWhereCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, MqlParserRULE_whereCommand)

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
		p.SetState(121)
		p.Match(MqlParserWHERE)
	}
	{
		p.SetState(122)
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

func (s *SelectCommandContext) INDEX_NAME() antlr.TerminalNode {
	return s.GetToken(MqlParserINDEX_NAME, 0)
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

func (p *MqlParser) SelectCommand() (localctx ISelectCommandContext) {
	localctx = NewSelectCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, MqlParserRULE_selectCommand)
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
		p.SetState(124)
		p.Match(MqlParserINDEX_NAME)
	}
	{
		p.SetState(125)
		p.Match(MqlParserASSIGN)
	}
	{
		p.SetState(126)
		p.Match(MqlParserIDENTIFIER)
	}
	p.SetState(130)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserLPAREN || _la == MqlParserDECIMAL_LITERAL || (((_la-67)&-(0x1f+1)) == 0 && ((1<<uint((_la-67)))&((1<<(MqlParserFLOAT_LITERAL-67))|(1<<(MqlParserBOOL_LITERAL-67))|(1<<(MqlParserSTRING_LITERAL-67))|(1<<(MqlParserIDENTIFIER-67)))) != 0) {
		{
			p.SetState(127)
			p.expression(0)
		}

		p.SetState(132)
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

func (p *MqlParser) RenameCommand() (localctx IRenameCommandContext) {
	localctx = NewRenameCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, MqlParserRULE_renameCommand)
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
		p.SetState(133)
		p.Match(MqlParserRENAME)
	}
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == MqlParserIDENTIFIER {
		{
			p.SetState(134)
			p.Match(MqlParserIDENTIFIER)
		}
		{
			p.SetState(135)
			p.Match(MqlParserAS)
		}
		{
			p.SetState(136)
			p.Match(MqlParserIDENTIFIER)
		}

		p.SetState(139)
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

func (p *MqlParser) StatCommand() (localctx IStatCommandContext) {
	localctx = NewStatCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, MqlParserRULE_statCommand)
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
		p.SetState(141)
		p.Match(MqlParserSTATS)
	}
	{
		p.SetState(142)
		p.AgrupCall()
	}
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == MqlParserAS {
		{
			p.SetState(143)
			p.Match(MqlParserAS)
		}
		{
			p.SetState(144)
			p.Match(MqlParserIDENTIFIER)
		}

	}
	p.SetState(155)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserCOMMA {
		{
			p.SetState(147)
			p.Match(MqlParserCOMMA)
		}
		{
			p.SetState(148)
			p.AgrupCall()
		}
		p.SetState(151)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == MqlParserAS {
			{
				p.SetState(149)
				p.Match(MqlParserAS)
			}
			{
				p.SetState(150)
				p.Match(MqlParserIDENTIFIER)
			}

		}

		p.SetState(157)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(158)
		p.Match(MqlParserBY)
	}
	{
		p.SetState(159)
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

func (p *MqlParser) BinCommand() (localctx IBinCommandContext) {
	localctx = NewBinCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, MqlParserRULE_binCommand)
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
		p.SetState(161)
		p.Match(MqlParserBIN)
	}
	p.SetState(163)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == MqlParserTIME_FIELD || _la == MqlParserIDENTIFIER {
		{
			p.SetState(162)
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
		p.SetState(165)
		p.Match(MqlParserBIN_SPAN)
	}
	{
		p.SetState(166)
		p.Match(MqlParserASSIGN)
	}
	{
		p.SetState(167)
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

func (p *MqlParser) FieldCommand() (localctx IFieldCommandContext) {
	localctx = NewFieldCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, MqlParserRULE_fieldCommand)
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
		p.SetState(169)
		p.Match(MqlParserFIELDS)
	}
	p.SetState(171)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == MqlParserADD || _la == MqlParserSUB {
		{
			p.SetState(170)
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
		p.SetState(173)
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

func (p *MqlParser) DedupCommand() (localctx IDedupCommandContext) {
	localctx = NewDedupCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, MqlParserRULE_dedupCommand)

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
		p.SetState(175)
		p.Match(MqlParserDEDUP)
	}
	{
		p.SetState(176)
		p.IdentifierList()
	}

	return localctx
}

// ISortCommandContext is an interface to support dynamic dispatch.
type ISortCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetSList returns the sList rule contexts.
	GetSList() ISortListContext

	// SetSList sets the sList rule contexts.
	SetSList(ISortListContext)

	// IsSortCommandContext differentiates from other interfaces.
	IsSortCommandContext()
}

type SortCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	sList  ISortListContext
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

func (s *SortCommandContext) GetSList() ISortListContext { return s.sList }

func (s *SortCommandContext) SetSList(v ISortListContext) { s.sList = v }

func (s *SortCommandContext) SORT() antlr.TerminalNode {
	return s.GetToken(MqlParserSORT, 0)
}

func (s *SortCommandContext) SortList() ISortListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISortListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISortListContext)
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

func (p *MqlParser) SortCommand() (localctx ISortCommandContext) {
	localctx = NewSortCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, MqlParserRULE_sortCommand)

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
		p.SetState(178)
		p.Match(MqlParserSORT)
	}
	{
		p.SetState(179)

		var _x = p.SortList()

		localctx.(*SortCommandContext).sList = _x
	}

	return localctx
}

// ITopCommandContext is an interface to support dynamic dispatch.
type ITopCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLimit returns the limit token.
	GetLimit() antlr.Token

	// SetLimit sets the limit token.
	SetLimit(antlr.Token)

	// IsTopCommandContext differentiates from other interfaces.
	IsTopCommandContext()
}

type TopCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	limit  antlr.Token
}

func NewEmptyTopCommandContext() *TopCommandContext {
	var p = new(TopCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_topCommand
	return p
}

func (*TopCommandContext) IsTopCommandContext() {}

func NewTopCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopCommandContext {
	var p = new(TopCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_topCommand

	return p
}

func (s *TopCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *TopCommandContext) GetLimit() antlr.Token { return s.limit }

func (s *TopCommandContext) SetLimit(v antlr.Token) { s.limit = v }

func (s *TopCommandContext) TOP() antlr.TerminalNode {
	return s.GetToken(MqlParserTOP, 0)
}

func (s *TopCommandContext) DECIMAL_LITERAL() antlr.TerminalNode {
	return s.GetToken(MqlParserDECIMAL_LITERAL, 0)
}

func (s *TopCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TopCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TopCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterTopCommand(s)
	}
}

func (s *TopCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitTopCommand(s)
	}
}

func (p *MqlParser) TopCommand() (localctx ITopCommandContext) {
	localctx = NewTopCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, MqlParserRULE_topCommand)

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
		p.SetState(181)
		p.Match(MqlParserTOP)
	}
	{
		p.SetState(182)

		var _m = p.Match(MqlParserDECIMAL_LITERAL)

		localctx.(*TopCommandContext).limit = _m
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

func (p *MqlParser) CompleteCommand() (localctx ICompleteCommandContext) {
	localctx = NewCompleteCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, MqlParserRULE_completeCommand)
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
		p.SetState(184)
		p.SelectCommand()
	}
	p.SetState(189)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserBITOR {
		{
			p.SetState(185)
			p.Match(MqlParserBITOR)
		}
		{
			p.SetState(186)
			p.Commands()
		}

		p.SetState(191)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

func (p *MqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 7:
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
