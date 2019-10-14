// Code generated from /Users/sebad/desa/workspace_go/eventdb/internal/query/MqlParser.g4 by ANTLR 4.7.2. DO NOT EDIT.

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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 83, 216,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 5, 5, 66, 10, 5, 3, 6, 3,
	6, 3, 6, 7, 6, 71, 10, 6, 12, 6, 14, 6, 74, 11, 6, 3, 7, 3, 7, 5, 7, 78,
	10, 7, 3, 8, 3, 8, 3, 8, 7, 8, 83, 10, 8, 12, 8, 14, 8, 86, 11, 8, 3, 9,
	3, 9, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10,
	3, 10, 3, 10, 3, 10, 3, 10, 5, 10, 104, 10, 10, 3, 10, 5, 10, 107, 10,
	10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 7, 10, 117,
	10, 10, 12, 10, 14, 10, 120, 11, 10, 3, 11, 3, 11, 3, 12, 3, 12, 3, 13,
	5, 13, 127, 10, 13, 3, 13, 7, 13, 130, 10, 13, 12, 13, 14, 13, 133, 11,
	13, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 5, 14,
	144, 10, 14, 3, 15, 3, 15, 3, 15, 3, 16, 7, 16, 150, 10, 16, 12, 16, 14,
	16, 153, 11, 16, 3, 16, 5, 16, 156, 10, 16, 3, 16, 7, 16, 159, 10, 16,
	12, 16, 14, 16, 162, 11, 16, 3, 17, 3, 17, 3, 17, 3, 17, 6, 17, 168, 10,
	17, 13, 17, 14, 17, 169, 3, 18, 3, 18, 3, 18, 3, 18, 3, 18, 3, 18, 7, 18,
	178, 10, 18, 12, 18, 14, 18, 181, 11, 18, 3, 19, 3, 19, 3, 19, 3, 19, 3,
	19, 3, 20, 3, 20, 3, 20, 3, 21, 3, 21, 3, 21, 3, 22, 3, 22, 3, 22, 3, 22,
	5, 22, 198, 10, 22, 3, 22, 3, 22, 3, 23, 3, 23, 3, 23, 3, 24, 3, 24, 3,
	24, 3, 25, 3, 25, 3, 25, 7, 25, 211, 10, 25, 12, 25, 14, 25, 214, 11, 25,
	3, 25, 2, 3, 18, 26, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28,
	30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 2, 7, 3, 2, 23, 38, 3, 2, 21, 22,
	3, 2, 59, 60, 4, 2, 48, 50, 56, 57, 3, 2, 17, 18, 2, 220, 2, 50, 3, 2,
	2, 2, 4, 53, 3, 2, 2, 2, 6, 55, 3, 2, 2, 2, 8, 65, 3, 2, 2, 2, 10, 67,
	3, 2, 2, 2, 12, 75, 3, 2, 2, 2, 14, 79, 3, 2, 2, 2, 16, 87, 3, 2, 2, 2,
	18, 106, 3, 2, 2, 2, 20, 121, 3, 2, 2, 2, 22, 123, 3, 2, 2, 2, 24, 131,
	3, 2, 2, 2, 26, 143, 3, 2, 2, 2, 28, 145, 3, 2, 2, 2, 30, 151, 3, 2, 2,
	2, 32, 163, 3, 2, 2, 2, 34, 171, 3, 2, 2, 2, 36, 182, 3, 2, 2, 2, 38, 187,
	3, 2, 2, 2, 40, 190, 3, 2, 2, 2, 42, 193, 3, 2, 2, 2, 44, 201, 3, 2, 2,
	2, 46, 204, 3, 2, 2, 2, 48, 207, 3, 2, 2, 2, 50, 51, 5, 48, 25, 2, 51,
	52, 7, 2, 2, 3, 52, 3, 3, 2, 2, 2, 53, 54, 9, 2, 2, 2, 54, 5, 3, 2, 2,
	2, 55, 56, 5, 4, 3, 2, 56, 57, 7, 39, 2, 2, 57, 58, 7, 80, 2, 2, 58, 59,
	7, 40, 2, 2, 59, 7, 3, 2, 2, 2, 60, 66, 7, 76, 2, 2, 61, 66, 7, 67, 2,
	2, 62, 66, 7, 72, 2, 2, 63, 66, 7, 74, 2, 2, 64, 66, 7, 80, 2, 2, 65, 60,
	3, 2, 2, 2, 65, 61, 3, 2, 2, 2, 65, 62, 3, 2, 2, 2, 65, 63, 3, 2, 2, 2,
	65, 64, 3, 2, 2, 2, 66, 9, 3, 2, 2, 2, 67, 72, 7, 80, 2, 2, 68, 69, 7,
	46, 2, 2, 69, 71, 7, 80, 2, 2, 70, 68, 3, 2, 2, 2, 71, 74, 3, 2, 2, 2,
	72, 70, 3, 2, 2, 2, 72, 73, 3, 2, 2, 2, 73, 11, 3, 2, 2, 2, 74, 72, 3,
	2, 2, 2, 75, 77, 7, 80, 2, 2, 76, 78, 9, 3, 2, 2, 77, 76, 3, 2, 2, 2, 77,
	78, 3, 2, 2, 2, 78, 13, 3, 2, 2, 2, 79, 84, 5, 12, 7, 2, 80, 81, 7, 46,
	2, 2, 81, 83, 5, 12, 7, 2, 82, 80, 3, 2, 2, 2, 83, 86, 3, 2, 2, 2, 84,
	82, 3, 2, 2, 2, 84, 85, 3, 2, 2, 2, 85, 15, 3, 2, 2, 2, 86, 84, 3, 2, 2,
	2, 87, 88, 7, 3, 2, 2, 88, 89, 7, 48, 2, 2, 89, 90, 7, 80, 2, 2, 90, 17,
	3, 2, 2, 2, 91, 92, 8, 10, 1, 2, 92, 93, 7, 39, 2, 2, 93, 94, 5, 18, 10,
	2, 94, 95, 7, 40, 2, 2, 95, 107, 3, 2, 2, 2, 96, 97, 5, 8, 5, 2, 97, 98,
	5, 20, 11, 2, 98, 99, 5, 8, 5, 2, 99, 107, 3, 2, 2, 2, 100, 101, 7, 19,
	2, 2, 101, 103, 7, 48, 2, 2, 102, 104, 9, 4, 2, 2, 103, 102, 3, 2, 2, 2,
	103, 104, 3, 2, 2, 2, 104, 105, 3, 2, 2, 2, 105, 107, 7, 71, 2, 2, 106,
	91, 3, 2, 2, 2, 106, 96, 3, 2, 2, 2, 106, 100, 3, 2, 2, 2, 107, 118, 3,
	2, 2, 2, 108, 109, 12, 5, 2, 2, 109, 110, 5, 22, 12, 2, 110, 111, 5, 18,
	10, 6, 111, 117, 3, 2, 2, 2, 112, 113, 12, 3, 2, 2, 113, 114, 7, 48, 2,
	2, 114, 115, 7, 52, 2, 2, 115, 117, 5, 18, 10, 4, 116, 108, 3, 2, 2, 2,
	116, 112, 3, 2, 2, 2, 117, 120, 3, 2, 2, 2, 118, 116, 3, 2, 2, 2, 118,
	119, 3, 2, 2, 2, 119, 19, 3, 2, 2, 2, 120, 118, 3, 2, 2, 2, 121, 122, 9,
	5, 2, 2, 122, 21, 3, 2, 2, 2, 123, 124, 9, 6, 2, 2, 124, 23, 3, 2, 2, 2,
	125, 127, 9, 4, 2, 2, 126, 125, 3, 2, 2, 2, 126, 127, 3, 2, 2, 2, 127,
	128, 3, 2, 2, 2, 128, 130, 7, 80, 2, 2, 129, 126, 3, 2, 2, 2, 130, 133,
	3, 2, 2, 2, 131, 129, 3, 2, 2, 2, 131, 132, 3, 2, 2, 2, 132, 25, 3, 2,
	2, 2, 133, 131, 3, 2, 2, 2, 134, 144, 5, 28, 15, 2, 135, 144, 5, 32, 17,
	2, 136, 144, 5, 34, 18, 2, 137, 144, 5, 38, 20, 2, 138, 144, 5, 40, 21,
	2, 139, 144, 5, 44, 23, 2, 140, 144, 5, 46, 24, 2, 141, 144, 5, 36, 19,
	2, 142, 144, 5, 42, 22, 2, 143, 134, 3, 2, 2, 2, 143, 135, 3, 2, 2, 2,
	143, 136, 3, 2, 2, 2, 143, 137, 3, 2, 2, 2, 143, 138, 3, 2, 2, 2, 143,
	139, 3, 2, 2, 2, 143, 140, 3, 2, 2, 2, 143, 141, 3, 2, 2, 2, 143, 142,
	3, 2, 2, 2, 144, 27, 3, 2, 2, 2, 145, 146, 7, 8, 2, 2, 146, 147, 5, 18,
	10, 2, 147, 29, 3, 2, 2, 2, 148, 150, 5, 18, 10, 2, 149, 148, 3, 2, 2,
	2, 150, 153, 3, 2, 2, 2, 151, 149, 3, 2, 2, 2, 151, 152, 3, 2, 2, 2, 152,
	155, 3, 2, 2, 2, 153, 151, 3, 2, 2, 2, 154, 156, 5, 16, 9, 2, 155, 154,
	3, 2, 2, 2, 155, 156, 3, 2, 2, 2, 156, 160, 3, 2, 2, 2, 157, 159, 5, 18,
	10, 2, 158, 157, 3, 2, 2, 2, 159, 162, 3, 2, 2, 2, 160, 158, 3, 2, 2, 2,
	160, 161, 3, 2, 2, 2, 161, 31, 3, 2, 2, 2, 162, 160, 3, 2, 2, 2, 163, 167,
	7, 4, 2, 2, 164, 165, 7, 80, 2, 2, 165, 166, 7, 15, 2, 2, 166, 168, 7,
	80, 2, 2, 167, 164, 3, 2, 2, 2, 168, 169, 3, 2, 2, 2, 169, 167, 3, 2, 2,
	2, 169, 170, 3, 2, 2, 2, 170, 33, 3, 2, 2, 2, 171, 172, 7, 14, 2, 2, 172,
	173, 5, 4, 3, 2, 173, 174, 7, 16, 2, 2, 174, 179, 7, 80, 2, 2, 175, 176,
	7, 46, 2, 2, 176, 178, 7, 80, 2, 2, 177, 175, 3, 2, 2, 2, 178, 181, 3,
	2, 2, 2, 179, 177, 3, 2, 2, 2, 179, 180, 3, 2, 2, 2, 180, 35, 3, 2, 2,
	2, 181, 179, 3, 2, 2, 2, 182, 183, 7, 11, 2, 2, 183, 184, 7, 12, 2, 2,
	184, 185, 7, 48, 2, 2, 185, 186, 7, 71, 2, 2, 186, 37, 3, 2, 2, 2, 187,
	188, 7, 13, 2, 2, 188, 189, 5, 24, 13, 2, 189, 39, 3, 2, 2, 2, 190, 191,
	7, 7, 2, 2, 191, 192, 5, 10, 6, 2, 192, 41, 3, 2, 2, 2, 193, 197, 7, 6,
	2, 2, 194, 195, 7, 20, 2, 2, 195, 196, 7, 48, 2, 2, 196, 198, 7, 80, 2,
	2, 197, 194, 3, 2, 2, 2, 197, 198, 3, 2, 2, 2, 198, 199, 3, 2, 2, 2, 199,
	200, 7, 78, 2, 2, 200, 43, 3, 2, 2, 2, 201, 202, 7, 9, 2, 2, 202, 203,
	5, 14, 8, 2, 203, 45, 3, 2, 2, 2, 204, 205, 7, 10, 2, 2, 205, 206, 7, 67,
	2, 2, 206, 47, 3, 2, 2, 2, 207, 212, 5, 30, 16, 2, 208, 209, 7, 64, 2,
	2, 209, 211, 5, 26, 14, 2, 210, 208, 3, 2, 2, 2, 211, 214, 3, 2, 2, 2,
	212, 210, 3, 2, 2, 2, 212, 213, 3, 2, 2, 2, 213, 49, 3, 2, 2, 2, 214, 212,
	3, 2, 2, 2, 20, 65, 72, 77, 84, 103, 106, 116, 118, 126, 131, 143, 151,
	155, 160, 169, 179, 197, 212,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'index'", "'rename'", "'seach'", "'rex'", "'dedup'", "'where'", "'sort'",
	"'top'", "'bucket'", "'span'", "'fields'", "'stats'", "'as'", "'by'", "'and'",
	"'or'", "'earlier'", "'field'", "'asc'", "'desc'", "'avg'", "'count'",
	"'distinct_count'", "'estdc'", "'estdc_error'", "'max'", "'median'", "'min'",
	"'mode'", "'range'", "'stdev'", "'stdevp'", "'sum'", "'sumsq'", "'var'",
	"'varp'", "'('", "')'", "'{'", "'}'", "'['", "']'", "';'", "','", "'.'",
	"'='", "'>'", "'<'", "'!'", "'~'", "'?'", "':'", "'=='", "'<='", "'>='",
	"'!='", "'+'", "'-'", "'*'", "'/'", "'&'", "'|'", "'^'", "'%'", "", "",
	"", "", "", "", "", "", "", "", "'null'", "", "'_time'",
}
var symbolicNames = []string{
	"", "INDEX", "RENAME", "SEARCH", "REX", "DEDUP", "WHERE", "SORT", "TOP",
	"BUCKET", "SPAN", "FIELDS", "STATS", "AS", "BY", "AND", "OR", "EARLIER",
	"FIELD", "ASC", "DESC", "AVG", "COUNT", "DISTINCT_COUNT", "ESTDC", "ESTDC_ERROR",
	"MAX", "MEDIAN", "MIN", "MODE", "RANGE", "STDEV", "STDEVP", "SUM", "SUMSQ",
	"VAR", "VARP", "LPAREN", "RPAREN", "LBRACE", "RBRACE", "LBRACK", "RBRACK",
	"SEMI", "COMMA", "DOT", "ASSIGN", "GT", "LT", "BANG", "TILDE", "QUESTION",
	"COLON", "EQUAL", "LE", "GE", "NOTEQUAL", "ADD", "SUB", "MUL", "DIV", "BITAND",
	"BITOR", "CARET", "MOD", "DECIMAL_LITERAL", "HEX_LITERAL", "OCT_LITERAL",
	"BINARY_LITERAL", "TIME_LITERAL", "FLOAT_LITERAL", "HEX_FLOAT_LITERAL",
	"BOOL_LITERAL", "CHAR_LITERAL", "STRING_LITERAL", "NULL_LITERAL", "REGEX",
	"TIME_FIELD", "IDENTIFIER", "WS", "COMMENT", "LINE_COMMENT",
}

var ruleNames = []string{
	"start", "agrupTypes", "agrupCall", "literal", "identifierList", "sort",
	"sortList", "indexExpression", "expression", "comparator", "binary", "fieldList",
	"commands", "whereCommand", "selectCommand", "renameCommand", "statCommand",
	"bucketCommand", "fieldCommand", "dedupCommand", "rexCommand", "sortCommand",
	"topCommand", "completeCommand",
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
	MqlParserINDEX             = 1
	MqlParserRENAME            = 2
	MqlParserSEARCH            = 3
	MqlParserREX               = 4
	MqlParserDEDUP             = 5
	MqlParserWHERE             = 6
	MqlParserSORT              = 7
	MqlParserTOP               = 8
	MqlParserBUCKET            = 9
	MqlParserSPAN              = 10
	MqlParserFIELDS            = 11
	MqlParserSTATS             = 12
	MqlParserAS                = 13
	MqlParserBY                = 14
	MqlParserAND               = 15
	MqlParserOR                = 16
	MqlParserEARLIER           = 17
	MqlParserFIELD             = 18
	MqlParserASC               = 19
	MqlParserDESC              = 20
	MqlParserAVG               = 21
	MqlParserCOUNT             = 22
	MqlParserDISTINCT_COUNT    = 23
	MqlParserESTDC             = 24
	MqlParserESTDC_ERROR       = 25
	MqlParserMAX               = 26
	MqlParserMEDIAN            = 27
	MqlParserMIN               = 28
	MqlParserMODE              = 29
	MqlParserRANGE             = 30
	MqlParserSTDEV             = 31
	MqlParserSTDEVP            = 32
	MqlParserSUM               = 33
	MqlParserSUMSQ             = 34
	MqlParserVAR               = 35
	MqlParserVARP              = 36
	MqlParserLPAREN            = 37
	MqlParserRPAREN            = 38
	MqlParserLBRACE            = 39
	MqlParserRBRACE            = 40
	MqlParserLBRACK            = 41
	MqlParserRBRACK            = 42
	MqlParserSEMI              = 43
	MqlParserCOMMA             = 44
	MqlParserDOT               = 45
	MqlParserASSIGN            = 46
	MqlParserGT                = 47
	MqlParserLT                = 48
	MqlParserBANG              = 49
	MqlParserTILDE             = 50
	MqlParserQUESTION          = 51
	MqlParserCOLON             = 52
	MqlParserEQUAL             = 53
	MqlParserLE                = 54
	MqlParserGE                = 55
	MqlParserNOTEQUAL          = 56
	MqlParserADD               = 57
	MqlParserSUB               = 58
	MqlParserMUL               = 59
	MqlParserDIV               = 60
	MqlParserBITAND            = 61
	MqlParserBITOR             = 62
	MqlParserCARET             = 63
	MqlParserMOD               = 64
	MqlParserDECIMAL_LITERAL   = 65
	MqlParserHEX_LITERAL       = 66
	MqlParserOCT_LITERAL       = 67
	MqlParserBINARY_LITERAL    = 68
	MqlParserTIME_LITERAL      = 69
	MqlParserFLOAT_LITERAL     = 70
	MqlParserHEX_FLOAT_LITERAL = 71
	MqlParserBOOL_LITERAL      = 72
	MqlParserCHAR_LITERAL      = 73
	MqlParserSTRING_LITERAL    = 74
	MqlParserNULL_LITERAL      = 75
	MqlParserREGEX             = 76
	MqlParserTIME_FIELD        = 77
	MqlParserIDENTIFIER        = 78
	MqlParserWS                = 79
	MqlParserCOMMENT           = 80
	MqlParserLINE_COMMENT      = 81
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
	MqlParserRULE_indexExpression = 7
	MqlParserRULE_expression      = 8
	MqlParserRULE_comparator      = 9
	MqlParserRULE_binary          = 10
	MqlParserRULE_fieldList       = 11
	MqlParserRULE_commands        = 12
	MqlParserRULE_whereCommand    = 13
	MqlParserRULE_selectCommand   = 14
	MqlParserRULE_renameCommand   = 15
	MqlParserRULE_statCommand     = 16
	MqlParserRULE_bucketCommand   = 17
	MqlParserRULE_fieldCommand    = 18
	MqlParserRULE_dedupCommand    = 19
	MqlParserRULE_rexCommand      = 20
	MqlParserRULE_sortCommand     = 21
	MqlParserRULE_topCommand      = 22
	MqlParserRULE_completeCommand = 23
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
		p.SetState(48)
		p.CompleteCommand()
	}
	{
		p.SetState(49)
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
		p.SetState(51)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-21)&-(0x1f+1)) == 0 && ((1<<uint((_la-21)))&((1<<(MqlParserAVG-21))|(1<<(MqlParserCOUNT-21))|(1<<(MqlParserDISTINCT_COUNT-21))|(1<<(MqlParserESTDC-21))|(1<<(MqlParserESTDC_ERROR-21))|(1<<(MqlParserMAX-21))|(1<<(MqlParserMEDIAN-21))|(1<<(MqlParserMIN-21))|(1<<(MqlParserMODE-21))|(1<<(MqlParserRANGE-21))|(1<<(MqlParserSTDEV-21))|(1<<(MqlParserSTDEVP-21))|(1<<(MqlParserSUM-21))|(1<<(MqlParserSUMSQ-21))|(1<<(MqlParserVAR-21))|(1<<(MqlParserVARP-21)))) != 0) {
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
		p.SetState(53)
		p.AgrupTypes()
	}
	{
		p.SetState(54)
		p.Match(MqlParserLPAREN)
	}
	{
		p.SetState(55)
		p.Match(MqlParserIDENTIFIER)
	}
	{
		p.SetState(56)
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

	p.SetState(63)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case MqlParserSTRING_LITERAL:
		localctx = NewStringLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(58)
			p.Match(MqlParserSTRING_LITERAL)
		}

	case MqlParserDECIMAL_LITERAL:
		localctx = NewDecimalLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(59)
			p.Match(MqlParserDECIMAL_LITERAL)
		}

	case MqlParserFLOAT_LITERAL:
		localctx = NewFloatLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(60)
			p.Match(MqlParserFLOAT_LITERAL)
		}

	case MqlParserBOOL_LITERAL:
		localctx = NewBoolLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(61)
			p.Match(MqlParserBOOL_LITERAL)
		}

	case MqlParserIDENTIFIER:
		localctx = NewIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(62)
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
		p.SetState(65)
		p.Match(MqlParserIDENTIFIER)
	}
	p.SetState(70)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserCOMMA {
		{
			p.SetState(66)
			p.Match(MqlParserCOMMA)
		}
		{
			p.SetState(67)
			p.Match(MqlParserIDENTIFIER)
		}

		p.SetState(72)
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

func (s *SortContext) ASC() antlr.TerminalNode {
	return s.GetToken(MqlParserASC, 0)
}

func (s *SortContext) DESC() antlr.TerminalNode {
	return s.GetToken(MqlParserDESC, 0)
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
		p.SetState(73)

		var _m = p.Match(MqlParserIDENTIFIER)

		localctx.(*SortContext).field = _m
	}
	p.SetState(75)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == MqlParserASC || _la == MqlParserDESC {
		{
			p.SetState(74)

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
		p.SetState(77)
		p.Sort()
	}
	p.SetState(82)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserCOMMA {
		{
			p.SetState(78)
			p.Match(MqlParserCOMMA)
		}
		{
			p.SetState(79)
			p.Sort()
		}

		p.SetState(84)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IIndexExpressionContext is an interface to support dynamic dispatch.
type IIndexExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetName returns the name token.
	GetName() antlr.Token

	// SetName sets the name token.
	SetName(antlr.Token)

	// IsIndexExpressionContext differentiates from other interfaces.
	IsIndexExpressionContext()
}

type IndexExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	name   antlr.Token
}

func NewEmptyIndexExpressionContext() *IndexExpressionContext {
	var p = new(IndexExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_indexExpression
	return p
}

func (*IndexExpressionContext) IsIndexExpressionContext() {}

func NewIndexExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IndexExpressionContext {
	var p = new(IndexExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_indexExpression

	return p
}

func (s *IndexExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *IndexExpressionContext) GetName() antlr.Token { return s.name }

func (s *IndexExpressionContext) SetName(v antlr.Token) { s.name = v }

func (s *IndexExpressionContext) INDEX() antlr.TerminalNode {
	return s.GetToken(MqlParserINDEX, 0)
}

func (s *IndexExpressionContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(MqlParserASSIGN, 0)
}

func (s *IndexExpressionContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, 0)
}

func (s *IndexExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IndexExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IndexExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterIndexExpression(s)
	}
}

func (s *IndexExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitIndexExpression(s)
	}
}

func (p *MqlParser) IndexExpression() (localctx IIndexExpressionContext) {
	localctx = NewIndexExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, MqlParserRULE_indexExpression)

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
		p.Match(MqlParserINDEX)
	}
	{
		p.SetState(86)
		p.Match(MqlParserASSIGN)
	}
	{
		p.SetState(87)

		var _m = p.Match(MqlParserIDENTIFIER)

		localctx.(*IndexExpressionContext).name = _m
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

type RegexExpressionContext struct {
	*ExpressionContext
	left  IExpressionContext
	op    antlr.Token
	right IExpressionContext
}

func NewRegexExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RegexExpressionContext {
	var p = new(RegexExpressionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *RegexExpressionContext) GetOp() antlr.Token { return s.op }

func (s *RegexExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *RegexExpressionContext) GetLeft() IExpressionContext { return s.left }

func (s *RegexExpressionContext) GetRight() IExpressionContext { return s.right }

func (s *RegexExpressionContext) SetLeft(v IExpressionContext) { s.left = v }

func (s *RegexExpressionContext) SetRight(v IExpressionContext) { s.right = v }

func (s *RegexExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegexExpressionContext) TILDE() antlr.TerminalNode {
	return s.GetToken(MqlParserTILDE, 0)
}

func (s *RegexExpressionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *RegexExpressionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *RegexExpressionContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(MqlParserASSIGN, 0)
}

func (s *RegexExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterRegexExpression(s)
	}
}

func (s *RegexExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitRegexExpression(s)
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

type TimeExpressionContext struct {
	*ExpressionContext
	left  antlr.Token
	op    antlr.Token
	right antlr.Token
}

func NewTimeExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TimeExpressionContext {
	var p = new(TimeExpressionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *TimeExpressionContext) GetLeft() antlr.Token { return s.left }

func (s *TimeExpressionContext) GetOp() antlr.Token { return s.op }

func (s *TimeExpressionContext) GetRight() antlr.Token { return s.right }

func (s *TimeExpressionContext) SetLeft(v antlr.Token) { s.left = v }

func (s *TimeExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *TimeExpressionContext) SetRight(v antlr.Token) { s.right = v }

func (s *TimeExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeExpressionContext) TIME_LITERAL() antlr.TerminalNode {
	return s.GetToken(MqlParserTIME_LITERAL, 0)
}

func (s *TimeExpressionContext) EARLIER() antlr.TerminalNode {
	return s.GetToken(MqlParserEARLIER, 0)
}

func (s *TimeExpressionContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(MqlParserASSIGN, 0)
}

func (s *TimeExpressionContext) ADD() antlr.TerminalNode {
	return s.GetToken(MqlParserADD, 0)
}

func (s *TimeExpressionContext) SUB() antlr.TerminalNode {
	return s.GetToken(MqlParserSUB, 0)
}

func (s *TimeExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterTimeExpression(s)
	}
}

func (s *TimeExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitTimeExpression(s)
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
	_startState := 16
	p.EnterRecursionRule(localctx, 16, MqlParserRULE_expression, _p)
	var _la int

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
	p.SetState(104)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case MqlParserLPAREN:
		localctx = NewParenExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(90)
			p.Match(MqlParserLPAREN)
		}
		{
			p.SetState(91)
			p.expression(0)
		}
		{
			p.SetState(92)
			p.Match(MqlParserRPAREN)
		}

	case MqlParserDECIMAL_LITERAL, MqlParserFLOAT_LITERAL, MqlParserBOOL_LITERAL, MqlParserSTRING_LITERAL, MqlParserIDENTIFIER:
		localctx = NewComparatorExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(94)

			var _x = p.Literal()

			localctx.(*ComparatorExpressionContext).left = _x
		}
		{
			p.SetState(95)

			var _x = p.Comparator()

			localctx.(*ComparatorExpressionContext).op = _x
		}
		{
			p.SetState(96)

			var _x = p.Literal()

			localctx.(*ComparatorExpressionContext).right = _x
		}

	case MqlParserEARLIER:
		localctx = NewTimeExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(98)

			var _m = p.Match(MqlParserEARLIER)

			localctx.(*TimeExpressionContext).left = _m
		}
		{
			p.SetState(99)

			var _m = p.Match(MqlParserASSIGN)

			localctx.(*TimeExpressionContext).op = _m
		}
		p.SetState(101)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == MqlParserADD || _la == MqlParserSUB {
			{
				p.SetState(100)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*TimeExpressionContext).right = _lt

				_la = p.GetTokenStream().LA(1)

				if !(_la == MqlParserADD || _la == MqlParserSUB) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*TimeExpressionContext).right = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

		}
		{
			p.SetState(103)
			p.Match(MqlParserTIME_LITERAL)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(116)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(114)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext()) {
			case 1:
				localctx = NewBinaryExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				localctx.(*BinaryExpressionContext).left = _prevctx

				p.PushNewRecursionContext(localctx, _startState, MqlParserRULE_expression)
				p.SetState(106)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(107)

					var _x = p.Binary()

					localctx.(*BinaryExpressionContext).op = _x
				}
				{
					p.SetState(108)

					var _x = p.expression(4)

					localctx.(*BinaryExpressionContext).right = _x
				}

			case 2:
				localctx = NewRegexExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				localctx.(*RegexExpressionContext).left = _prevctx

				p.PushNewRecursionContext(localctx, _startState, MqlParserRULE_expression)
				p.SetState(110)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				{
					p.SetState(111)

					var _m = p.Match(MqlParserASSIGN)

					localctx.(*RegexExpressionContext).op = _m
				}
				{
					p.SetState(112)
					p.Match(MqlParserTILDE)
				}
				{
					p.SetState(113)

					var _x = p.expression(2)

					localctx.(*RegexExpressionContext).right = _x
				}

			}

		}
		p.SetState(118)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 18, MqlParserRULE_comparator)
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
		p.SetState(119)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-46)&-(0x1f+1)) == 0 && ((1<<uint((_la-46)))&((1<<(MqlParserASSIGN-46))|(1<<(MqlParserGT-46))|(1<<(MqlParserLT-46))|(1<<(MqlParserLE-46))|(1<<(MqlParserGE-46)))) != 0) {
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
	p.EnterRule(localctx, 20, MqlParserRULE_binary)
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
		p.SetState(121)
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

// IFieldListContext is an interface to support dynamic dispatch.
type IFieldListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldListContext differentiates from other interfaces.
	IsFieldListContext()
}

type FieldListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldListContext() *FieldListContext {
	var p = new(FieldListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_fieldList
	return p
}

func (*FieldListContext) IsFieldListContext() {}

func NewFieldListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldListContext {
	var p = new(FieldListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_fieldList

	return p
}

func (s *FieldListContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldListContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(MqlParserIDENTIFIER)
}

func (s *FieldListContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, i)
}

func (s *FieldListContext) AllADD() []antlr.TerminalNode {
	return s.GetTokens(MqlParserADD)
}

func (s *FieldListContext) ADD(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserADD, i)
}

func (s *FieldListContext) AllSUB() []antlr.TerminalNode {
	return s.GetTokens(MqlParserSUB)
}

func (s *FieldListContext) SUB(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserSUB, i)
}

func (s *FieldListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterFieldList(s)
	}
}

func (s *FieldListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitFieldList(s)
	}
}

func (p *MqlParser) FieldList() (localctx IFieldListContext) {
	localctx = NewFieldListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, MqlParserRULE_fieldList)
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
	p.SetState(129)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-57)&-(0x1f+1)) == 0 && ((1<<uint((_la-57)))&((1<<(MqlParserADD-57))|(1<<(MqlParserSUB-57))|(1<<(MqlParserIDENTIFIER-57)))) != 0 {
		p.SetState(124)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == MqlParserADD || _la == MqlParserSUB {
			{
				p.SetState(123)
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
			p.SetState(126)
			p.Match(MqlParserIDENTIFIER)
		}

		p.SetState(131)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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

func (s *CommandsContext) BucketCommand() IBucketCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBucketCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBucketCommandContext)
}

func (s *CommandsContext) RexCommand() IRexCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRexCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRexCommandContext)
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
	p.EnterRule(localctx, 24, MqlParserRULE_commands)

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

	p.SetState(141)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case MqlParserWHERE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(132)
			p.WhereCommand()
		}

	case MqlParserRENAME:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(133)
			p.RenameCommand()
		}

	case MqlParserSTATS:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(134)
			p.StatCommand()
		}

	case MqlParserFIELDS:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(135)
			p.FieldCommand()
		}

	case MqlParserDEDUP:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(136)
			p.DedupCommand()
		}

	case MqlParserSORT:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(137)
			p.SortCommand()
		}

	case MqlParserTOP:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(138)
			p.TopCommand()
		}

	case MqlParserBUCKET:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(139)
			p.BucketCommand()
		}

	case MqlParserREX:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(140)
			p.RexCommand()
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
	p.EnterRule(localctx, 26, MqlParserRULE_whereCommand)

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
		p.SetState(143)
		p.Match(MqlParserWHERE)
	}
	{
		p.SetState(144)
		p.expression(0)
	}

	return localctx
}

// ISelectCommandContext is an interface to support dynamic dispatch.
type ISelectCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetIndex returns the index rule contexts.
	GetIndex() IIndexExpressionContext

	// SetIndex sets the index rule contexts.
	SetIndex(IIndexExpressionContext)

	// IsSelectCommandContext differentiates from other interfaces.
	IsSelectCommandContext()
}

type SelectCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	index  IIndexExpressionContext
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

func (s *SelectCommandContext) GetIndex() IIndexExpressionContext { return s.index }

func (s *SelectCommandContext) SetIndex(v IIndexExpressionContext) { s.index = v }

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

func (s *SelectCommandContext) IndexExpression() IIndexExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIndexExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIndexExpressionContext)
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
	p.EnterRule(localctx, 28, MqlParserRULE_selectCommand)
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

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(149)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(146)
				p.expression(0)
			}

		}
		p.SetState(151)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())
	}
	p.SetState(153)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == MqlParserINDEX {
		{
			p.SetState(152)

			var _x = p.IndexExpression()

			localctx.(*SelectCommandContext).index = _x
		}

	}
	p.SetState(158)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserEARLIER || _la == MqlParserLPAREN || (((_la-65)&-(0x1f+1)) == 0 && ((1<<uint((_la-65)))&((1<<(MqlParserDECIMAL_LITERAL-65))|(1<<(MqlParserFLOAT_LITERAL-65))|(1<<(MqlParserBOOL_LITERAL-65))|(1<<(MqlParserSTRING_LITERAL-65))|(1<<(MqlParserIDENTIFIER-65)))) != 0) {
		{
			p.SetState(155)
			p.expression(0)
		}

		p.SetState(160)
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
	p.EnterRule(localctx, 30, MqlParserRULE_renameCommand)
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
		p.Match(MqlParserRENAME)
	}
	p.SetState(165)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == MqlParserIDENTIFIER {
		{
			p.SetState(162)
			p.Match(MqlParserIDENTIFIER)
		}
		{
			p.SetState(163)
			p.Match(MqlParserAS)
		}
		{
			p.SetState(164)
			p.Match(MqlParserIDENTIFIER)
		}

		p.SetState(167)
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

	// GetField returns the field token.
	GetField() antlr.Token

	// SetField sets the field token.
	SetField(antlr.Token)

	// GetF returns the f rule contexts.
	GetF() IAgrupTypesContext

	// SetF sets the f rule contexts.
	SetF(IAgrupTypesContext)

	// IsStatCommandContext differentiates from other interfaces.
	IsStatCommandContext()
}

type StatCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	f      IAgrupTypesContext
	field  antlr.Token
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

func (s *StatCommandContext) GetField() antlr.Token { return s.field }

func (s *StatCommandContext) SetField(v antlr.Token) { s.field = v }

func (s *StatCommandContext) GetF() IAgrupTypesContext { return s.f }

func (s *StatCommandContext) SetF(v IAgrupTypesContext) { s.f = v }

func (s *StatCommandContext) STATS() antlr.TerminalNode {
	return s.GetToken(MqlParserSTATS, 0)
}

func (s *StatCommandContext) BY() antlr.TerminalNode {
	return s.GetToken(MqlParserBY, 0)
}

func (s *StatCommandContext) AgrupTypes() IAgrupTypesContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAgrupTypesContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAgrupTypesContext)
}

func (s *StatCommandContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(MqlParserIDENTIFIER)
}

func (s *StatCommandContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, i)
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
	p.EnterRule(localctx, 32, MqlParserRULE_statCommand)
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
		p.Match(MqlParserSTATS)
	}
	{
		p.SetState(170)

		var _x = p.AgrupTypes()

		localctx.(*StatCommandContext).f = _x
	}
	{
		p.SetState(171)
		p.Match(MqlParserBY)
	}
	{
		p.SetState(172)

		var _m = p.Match(MqlParserIDENTIFIER)

		localctx.(*StatCommandContext).field = _m
	}
	p.SetState(177)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserCOMMA {
		{
			p.SetState(173)
			p.Match(MqlParserCOMMA)
		}
		{
			p.SetState(174)
			p.Match(MqlParserIDENTIFIER)
		}

		p.SetState(179)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IBucketCommandContext is an interface to support dynamic dispatch.
type IBucketCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetSpan returns the span token.
	GetSpan() antlr.Token

	// SetSpan sets the span token.
	SetSpan(antlr.Token)

	// IsBucketCommandContext differentiates from other interfaces.
	IsBucketCommandContext()
}

type BucketCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	span   antlr.Token
}

func NewEmptyBucketCommandContext() *BucketCommandContext {
	var p = new(BucketCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_bucketCommand
	return p
}

func (*BucketCommandContext) IsBucketCommandContext() {}

func NewBucketCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BucketCommandContext {
	var p = new(BucketCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_bucketCommand

	return p
}

func (s *BucketCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *BucketCommandContext) GetSpan() antlr.Token { return s.span }

func (s *BucketCommandContext) SetSpan(v antlr.Token) { s.span = v }

func (s *BucketCommandContext) BUCKET() antlr.TerminalNode {
	return s.GetToken(MqlParserBUCKET, 0)
}

func (s *BucketCommandContext) SPAN() antlr.TerminalNode {
	return s.GetToken(MqlParserSPAN, 0)
}

func (s *BucketCommandContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(MqlParserASSIGN, 0)
}

func (s *BucketCommandContext) TIME_LITERAL() antlr.TerminalNode {
	return s.GetToken(MqlParserTIME_LITERAL, 0)
}

func (s *BucketCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BucketCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BucketCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterBucketCommand(s)
	}
}

func (s *BucketCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitBucketCommand(s)
	}
}

func (p *MqlParser) BucketCommand() (localctx IBucketCommandContext) {
	localctx = NewBucketCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, MqlParserRULE_bucketCommand)

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
		p.SetState(180)
		p.Match(MqlParserBUCKET)
	}
	{
		p.SetState(181)
		p.Match(MqlParserSPAN)
	}
	{
		p.SetState(182)
		p.Match(MqlParserASSIGN)
	}
	{
		p.SetState(183)

		var _m = p.Match(MqlParserTIME_LITERAL)

		localctx.(*BucketCommandContext).span = _m
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

func (s *FieldCommandContext) FieldList() IFieldListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
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
	p.EnterRule(localctx, 36, MqlParserRULE_fieldCommand)

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
		p.SetState(185)
		p.Match(MqlParserFIELDS)
	}
	{
		p.SetState(186)
		p.FieldList()
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
	p.EnterRule(localctx, 38, MqlParserRULE_dedupCommand)

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
		p.SetState(188)
		p.Match(MqlParserDEDUP)
	}
	{
		p.SetState(189)
		p.IdentifierList()
	}

	return localctx
}

// IRexCommandContext is an interface to support dynamic dispatch.
type IRexCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetRexfield returns the rexfield token.
	GetRexfield() antlr.Token

	// GetRegex returns the regex token.
	GetRegex() antlr.Token

	// SetRexfield sets the rexfield token.
	SetRexfield(antlr.Token)

	// SetRegex sets the regex token.
	SetRegex(antlr.Token)

	// IsRexCommandContext differentiates from other interfaces.
	IsRexCommandContext()
}

type RexCommandContext struct {
	*antlr.BaseParserRuleContext
	parser   antlr.Parser
	rexfield antlr.Token
	regex    antlr.Token
}

func NewEmptyRexCommandContext() *RexCommandContext {
	var p = new(RexCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MqlParserRULE_rexCommand
	return p
}

func (*RexCommandContext) IsRexCommandContext() {}

func NewRexCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RexCommandContext {
	var p = new(RexCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MqlParserRULE_rexCommand

	return p
}

func (s *RexCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *RexCommandContext) GetRexfield() antlr.Token { return s.rexfield }

func (s *RexCommandContext) GetRegex() antlr.Token { return s.regex }

func (s *RexCommandContext) SetRexfield(v antlr.Token) { s.rexfield = v }

func (s *RexCommandContext) SetRegex(v antlr.Token) { s.regex = v }

func (s *RexCommandContext) REX() antlr.TerminalNode {
	return s.GetToken(MqlParserREX, 0)
}

func (s *RexCommandContext) REGEX() antlr.TerminalNode {
	return s.GetToken(MqlParserREGEX, 0)
}

func (s *RexCommandContext) FIELD() antlr.TerminalNode {
	return s.GetToken(MqlParserFIELD, 0)
}

func (s *RexCommandContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(MqlParserASSIGN, 0)
}

func (s *RexCommandContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(MqlParserIDENTIFIER, 0)
}

func (s *RexCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RexCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RexCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.EnterRexCommand(s)
	}
}

func (s *RexCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MqlParserListener); ok {
		listenerT.ExitRexCommand(s)
	}
}

func (p *MqlParser) RexCommand() (localctx IRexCommandContext) {
	localctx = NewRexCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, MqlParserRULE_rexCommand)
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
		p.SetState(191)
		p.Match(MqlParserREX)
	}
	p.SetState(195)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == MqlParserFIELD {
		{
			p.SetState(192)
			p.Match(MqlParserFIELD)
		}
		{
			p.SetState(193)
			p.Match(MqlParserASSIGN)
		}
		{
			p.SetState(194)

			var _m = p.Match(MqlParserIDENTIFIER)

			localctx.(*RexCommandContext).rexfield = _m
		}

	}
	{
		p.SetState(197)

		var _m = p.Match(MqlParserREGEX)

		localctx.(*RexCommandContext).regex = _m
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
	p.EnterRule(localctx, 42, MqlParserRULE_sortCommand)

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
		p.SetState(199)
		p.Match(MqlParserSORT)
	}
	{
		p.SetState(200)

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
	p.EnterRule(localctx, 44, MqlParserRULE_topCommand)

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
		p.SetState(202)
		p.Match(MqlParserTOP)
	}
	{
		p.SetState(203)

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
	p.EnterRule(localctx, 46, MqlParserRULE_completeCommand)
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
		p.SetState(205)
		p.SelectCommand()
	}
	p.SetState(210)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == MqlParserBITOR {
		{
			p.SetState(206)
			p.Match(MqlParserBITOR)
		}
		{
			p.SetState(207)
			p.Commands()
		}

		p.SetState(212)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

func (p *MqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 8:
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
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
