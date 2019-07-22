parser grammar MqlParser;

options { tokenVocab=MqlLexer; }

//stats sum(b) BY index
// Rules
start : a EOF ;

identifier_list: IDENTIFIER (',' IDENTIFIER )* ;

// stats expresion
agrupTypes
    : AVG
    | COUNT
    | DISTINCT_COUNT
    | ESTDC
    | ESTDC_ERROR
    | MAX
    | MEDIAN
    | MIN
    | MODE
    | RANGE
    | STDEV
    | STDEVP
    | SUM
    | SUMSQ
    | VAR
    | VARP
    ;

agrupCall: agrupTypes LPAREN IDENTIFIER RPAREN;

stat_expresion : STATS agrupCall (','agrupCall)* (AS IDENTIFIER)? BY IDENTIFIER;

// where expresion

comparators
    : EQUAL
    | LE
    | GE
    | NOTEQUAL
    ;

expression: IDENTIFIER comparators IDENTIFIER;

expressions
    : expression ( ( AND | OR )  expression)*
    ;

where_expresion : WHERE expressions;

// select expresion

select_expression: SOURCE_TYPE ASSIGN IDENTIFIER (expression)*;


// fields expresion

fields: FIELDS (ADD|SUB)? identifier_list;

a: select_expression;
