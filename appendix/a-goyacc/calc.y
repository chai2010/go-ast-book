/* Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved. */
/* Use of this source code is governed by a Apache */
/* license that can be found in the LICENSE file. */

/* simplest version of calculator */

%{
package main

import (
	"fmt"
)

var idValueMap = map[string]int{}
%}

%union {
	value int
	id    string
}

%type  <value> exp factor term
%token <value> NUMBER
%token <id>    ID

%token ADD SUB MUL DIV ABS
%token LPAREN RPAREN ASSIGN
%token EOL

%%
calclist
	: // nothing
	| calclist exp EOL {
		idValueMap["_"] = $2
		fmt.Printf("= %v\n", $2)
	}
	| calclist ID ASSIGN exp EOL {
		idValueMap["_"] = $4
		idValueMap[$2] = $4
		fmt.Printf("= %v\n", $4)
	}
	;

exp
	: factor         { $$ = $1 }
	| exp ADD factor { $$ = $1 + $3 }
	| exp SUB factor { $$ = $1 - $3 }
	;

factor
	: term            { $$ = $1 }
	| factor MUL term { $$ = $1 * $3 }
	| factor DIV term { $$ = $1 / $3 }
	;

term
	: NUMBER            { $$ = $1 }
	| ID                { $$ = idValueMap[$1] }
	| ABS term          { if $2 >= 0 { $$ = $2 } else { $$ = -$2 } }
	| LPAREN exp RPAREN { $$ = $2 }
	;

%%
