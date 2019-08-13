parser grammar MqlParser;

options { tokenVocab=MqlLexer; }

//stats sum(b) BY index
// Rules
start : completeCommand EOF ;

identifierList: IDENTIFIER (',' IDENTIFIER )* ;

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

literal
 : STRING_LITERAL                                 #stringLiteral
 | DECIMAL_LITERAL                                #decimalLiteral
 | FLOAT_LITERAL                                  #floatLiteral
 | BOOL_LITERAL                                   #boolLiteral
 | IDENTIFIER                                     #identifier
 ;

expression
 : LPAREN expression RPAREN                             #parenExpression
 | left=IDENTIFIER op=comparator right=literal          #comparatorExpression
 | left=expression op=binary right=expression           #binaryExpression
 ;

comparator
 : GT | GE | LT | LE | ASSIGN
 ;

binary
 : AND | OR
 ;

commands
    : whereCommand
    | renameCommand
    | statCommand
    | fieldCommand
    | dedupCommand
    | sortCommand
    | headCommand
    | binCommand
    ;

// where command
whereCommand : WHERE expression;

// select expresion
selectCommand: SOURCE_TYPE ASSIGN IDENTIFIER expression*;

// rename expresion
renameCommand: RENAME (IDENTIFIER AS IDENTIFIER)+;

// stats expresion
statCommand : STATS agrupCall (AS IDENTIFIER)? (COMMA agrupCall (AS IDENTIFIER)? )*  BY IDENTIFIER;

// bin expresion
binCommand : BIN (IDENTIFIER | TIME_FIELD)? BIN_SPAN ASSIGN TIME_LITERAL;

// fields expresion
fieldCommand: FIELDS (ADD|SUB)? identifierList;

// dedup expresion
dedupCommand: DEDUP identifierList;

// sort expresion
sortCommand: SORT identifierList;

// head expresion
headCommand: HEAD DECIMAL_LITERAL;

completeCommand: selectCommand ( BITOR commands )* ;
