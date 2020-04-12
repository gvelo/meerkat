// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package parser

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOF-0]
	_ = x[COMMENT-1]
	_ = x[literal_beg-2]
	_ = x[IDENT-3]
	_ = x[INT-4]
	_ = x[FLOAT-5]
	_ = x[STRING-6]
	_ = x[TIME-7]
	_ = x[DATETIME-8]
	_ = x[BOOL-9]
	_ = x[literal_end-10]
	_ = x[operator_beg-11]
	_ = x[ADD-12]
	_ = x[SUB-13]
	_ = x[MUL-14]
	_ = x[QUO-15]
	_ = x[REM-16]
	_ = x[ASSIGN-17]
	_ = x[EQL-18]
	_ = x[EQL_CI-19]
	_ = x[NEQ-20]
	_ = x[NEQ_CI-21]
	_ = x[LSS-22]
	_ = x[GTR-23]
	_ = x[LEQ-24]
	_ = x[GEQ-25]
	_ = x[AND-26]
	_ = x[OR-27]
	_ = x[IN-28]
	_ = x[NOT_IN-29]
	_ = x[IN_CI-30]
	_ = x[NOT_IN_CI-31]
	_ = x[HAS-32]
	_ = x[NOT_HAS-33]
	_ = x[HAS_CS-34]
	_ = x[NOT_HAS_CS-35]
	_ = x[HASPREFIX-36]
	_ = x[NOT_HASPREFIX-37]
	_ = x[HASPREFIX_CS-38]
	_ = x[NOT_HASPREFIX_CS-39]
	_ = x[HASSUFFIX-40]
	_ = x[NOT_HASSUFFIX-41]
	_ = x[HASSUFFIX_CS-42]
	_ = x[NOT_HASSUFFIX_CS-43]
	_ = x[CONTAINS-44]
	_ = x[NOT_CONTAINS-45]
	_ = x[CONTAINS_CS-46]
	_ = x[NOT_CONTAINS_CS-47]
	_ = x[STARTSWITH-48]
	_ = x[NOT_STARTSWITH-49]
	_ = x[STARTSWITH_CS-50]
	_ = x[NOT_STARTSWITH_CS-51]
	_ = x[ENDSWITH-52]
	_ = x[NOT_ENDSWITH-53]
	_ = x[ENDSWITH_CS-54]
	_ = x[NOT_ENDSWITH_CS-55]
	_ = x[MATCHES-56]
	_ = x[HAS_ANY-57]
	_ = x[BETWEEN-58]
	_ = x[NOT_BETWEEN-59]
	_ = x[RANGE-60]
	_ = x[operator_end-61]
	_ = x[LPAREN-62]
	_ = x[LBRACK-63]
	_ = x[LBRACE-64]
	_ = x[COMMA-65]
	_ = x[PERIOD-66]
	_ = x[RPAREN-67]
	_ = x[RBRACK-68]
	_ = x[RBRACE-69]
	_ = x[SEMICOLON-70]
	_ = x[COLON-71]
	_ = x[PIPE-72]
}

const _TokenType_name = "EOFCOMMENTliteral_begIDENTINTFLOATSTRINGTIMEDATETIMEBOOLliteral_endoperator_begADDSUBMULQUOREMASSIGNEQLEQL_CINEQNEQ_CILSSGTRLEQGEQANDORINNOT_ININ_CINOT_IN_CIHASNOT_HASHAS_CSNOT_HAS_CSHASPREFIXNOT_HASPREFIXHASPREFIX_CSNOT_HASPREFIX_CSHASSUFFIXNOT_HASSUFFIXHASSUFFIX_CSNOT_HASSUFFIX_CSCONTAINSNOT_CONTAINSCONTAINS_CSNOT_CONTAINS_CSSTARTSWITHNOT_STARTSWITHSTARTSWITH_CSNOT_STARTSWITH_CSENDSWITHNOT_ENDSWITHENDSWITH_CSNOT_ENDSWITH_CSMATCHESHAS_ANYBETWEENNOT_BETWEENRANGEoperator_endLPARENLBRACKLBRACECOMMAPERIODRPARENRBRACKRBRACESEMICOLONCOLONPIPE"

var _TokenType_index = [...]uint16{0, 3, 10, 21, 26, 29, 34, 40, 44, 52, 56, 67, 79, 82, 85, 88, 91, 94, 100, 103, 109, 112, 118, 121, 124, 127, 130, 133, 135, 137, 143, 148, 157, 160, 167, 173, 183, 192, 205, 217, 233, 242, 255, 267, 283, 291, 303, 314, 329, 339, 353, 366, 383, 391, 403, 414, 429, 436, 443, 450, 461, 466, 478, 484, 490, 496, 501, 507, 513, 519, 525, 534, 539, 543}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
