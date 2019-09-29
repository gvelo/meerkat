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

parser grammar MqlParser;

options { tokenVocab=MqlLexer; }

//stats sum(b) BY index
// Rules
start : completeCommand EOF ;

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

identifierList: IDENTIFIER (',' IDENTIFIER )* ;

sort: field=IDENTIFIER direction=(ASC|DESC)?;

sortList:  sort (COMMA sort )* ;

indexExpression: INDEX ASSIGN name=IDENTIFIER;

expression
 : LPAREN expression RPAREN                              #parenExpression
 | left=literal op=comparator right=literal              #comparatorExpression
 | left=expression op=binary right=expression            #binaryExpression
 | left=EARLIER op=ASSIGN right=(ADD|SUB)? TIME_LITERAL  #timeExpression
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
    | topCommand
    | bucketCommand
    | rexCommand
    ;

// where command
whereCommand : WHERE expression;

// select expresion
selectCommand: expression* (index=indexExpression)? expression* ;

// rename expresion
renameCommand: RENAME (IDENTIFIER AS IDENTIFIER)+;

// stats expresion
statCommand : STATS f=agrupTypes BY field= IDENTIFIER (COMMA  IDENTIFIER)*   ;

// bucket expresion
bucketCommand : BUCKET SPAN ASSIGN span=TIME_LITERAL;

// fields expresion
fieldCommand: FIELDS (ADD|SUB)? identifierList;

// dedup expresion
dedupCommand: DEDUP identifierList;

// rex expresion
rexCommand: REX (FIELD ASSIGN rexfield=IDENTIFIER)? regex=REGEX;

// sort expresion
sortCommand: SORT sList=sortList;

// head expresion
topCommand: TOP limit=DECIMAL_LITERAL;

completeCommand: selectCommand ( BITOR commands )* ;
