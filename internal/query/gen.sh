#!/usr/bin/env bash

java -Xmx500M -cp "/usr/local/lib/antlr-4.7.1-complete.jar:$CLASSPATH" org.antlr.v4.Tool -o ./mql_parser -package mql_parser -listener -visitor -Dlanguage=Go -lib ../query ./MqlParser.g4 ./MqlLexer.g4