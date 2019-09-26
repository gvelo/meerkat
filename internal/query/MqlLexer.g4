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

// lexer  mqllexer.g4
lexer grammar MqlLexer;

// Keywords

INDEX: 'index';
RENAME: 'rename';
SEARCH: 'seach';
DEDUP: 'dedup';
WHERE: 'where';
SORT: 'sort';
TOP: 'top';
BIN: 'bin';
BIN_SPAN: 'span';
FIELDS: 'fields';
STATS: 'stats';
AS: 'as';
BY: 'by';
AND: 'and';
OR:   'or';
EARLIER:   'earlier';

ASC:   'asc';
DESC:   'desc';

// Aggregate functions

AVG: 'avg';
COUNT: 'count';
DISTINCT_COUNT: 'distinct_count';
ESTDC: 'estdc';
ESTDC_ERROR: 'estdc_error';
MAX: 'max';
MEDIAN: 'median';
MIN: 'min';
MODE: 'mode';
RANGE: 'range';
STDEV: 'stdev';
STDEVP: 'stdevp';
SUM: 'sum';
SUMSQ: 'sumsq';
VAR: 'var';
VARP: 'varp';

// Separators
LPAREN:             '(';
RPAREN:             ')';
LBRACE:             '{';
RBRACE:             '}';
LBRACK:             '[';
RBRACK:             ']';
SEMI:               ';';
COMMA:              ',';
DOT:                '.';

// Operators

ASSIGN:             '=';
GT:                 '>';
LT:                 '<';
BANG:               '!';
TILDE:              '~';
QUESTION:           '?';
COLON:              ':';

EQUAL:              '==';
LE:                 '<=';
GE:                 '>=';
NOTEQUAL:           '!=';
ADD:                '+';
SUB:                '-';
MUL:                '*';
DIV:                '/';
BITAND:             '&';
BITOR:              '|';
CARET:              '^';
MOD:                '%';


// Literals
DECIMAL_LITERAL:    ('0' | [1-9] (Digits? | '_'+ Digits)) [lL]?;
HEX_LITERAL:        '0' [xX] [0-9a-fA-F] ([0-9a-fA-F_]* [0-9a-fA-F])? [lL]?;
OCT_LITERAL:        '0' '_'* [0-7] ([0-7_]* [0-7])? [lL]?;
BINARY_LITERAL:     '0' [bB] [01] ([01_]* [01])? [lL]?;
TIME_LITERAL:       DECIMAL_LITERAL SpanLength;
FLOAT_LITERAL:      (Digits '.' Digits? | '.' Digits) ExponentPart? [fFdD]?
             |       Digits (ExponentPart [fFdD]? | [fFdD])
             ;

HEX_FLOAT_LITERAL:  '0' [xX] (HexDigits '.'? | HexDigits? '.' HexDigits) [pP] [+-]? Digits [fFdD]?;

BOOL_LITERAL:       'true'
            |       'false'
            ;

CHAR_LITERAL:       '\'' (~['\\\r\n] | EscapeSequence) '\'';

STRING_LITERAL:     '"' (~["\\\r\n] | EscapeSequence)* '"';
NULL_LITERAL:       'null';


TIME_FIELD: '_time';

// Fragment rules

fragment Miliseconds
    : 'ms'
    ;

fragment Month
    : 'mon'
    | 'month'
    | 'months'
    ;

fragment Days
    : 'd'
    | 'day'
    | 'days'
    ;

fragment Hours
    : 'h'
    | 'hr'
    | 'hrs'
    | 'hour'
    | 'hours'
    ;

fragment Minutes
    : 'm'
    | 'min'
    | 'mins'
    | 'minute'
    | 'minutes'
    ;

fragment Seconds
    : 's'
    | 'sec'
    | 'secs'
    | 'second'
    | 'seconds'
    ;

fragment SpanLength
    : Miliseconds
    | Seconds
    | Minutes
    | Hours
    | Days
    | Month
    ;

fragment ExponentPart
    : [eE] [+-]? Digits
    ;

fragment Digits
    : [0-9] ([0-9_]* [0-9])?
    ;

fragment HexDigits
    : HexDigit ((HexDigit | '_')* HexDigit)?
    ;
fragment HexDigit
    : [0-9a-fA-F]
    ;

fragment EscapeSequence
    : '\\' [btnfr"'\\]
    | '\\' ([0-3]? [0-7])? [0-7]
    | '\\' 'u'+ HexDigit HexDigit HexDigit HexDigit
    ;

fragment Letter
    : [a-zA-Z$_] // these are the "java letters" below 0x7F
    | ~[\u0000-\u007F\uD800-\uDBFF] // covers all characters above 0x7F which are not a surrogate
    | [\uD800-\uDBFF] [\uDC00-\uDFFF] // covers UTF-16 surrogate pairs encodings for U+10000 to U+10FFFF
    ;

fragment LetterOrDigit
    : Letter
    | [0-9]
    ;

// Identifiers
IDENTIFIER:         Letter LetterOrDigit*;

// Whitespace and comments
WS:                 [ \t\r\n\u000C]+ -> channel(HIDDEN);
COMMENT:            '/*' .*? '*/'    -> channel(HIDDEN);
LINE_COMMENT:       '//' ~[\r\n]*    -> channel(HIDDEN);