/* Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved. */
/* Use of this source code is governed by a Apache */
/* license that can be found in the LICENSE file. */

%option noyywrap

%{
#include "tok.h"
%}

%%

[_a-zA-Z]+ { return ID; }
[0-9]+     { return NUMBER; }

"+"    { return ADD; }
"-"    { return SUB; }
"*"    { return MUL; }
"/"    { return DIV; }
"|"    { return ABS; }

"("    { return LPAREN; }
")"    { return RPAREN; }
"="    { return ASSIGN; }

\n     { return EOL; }
[ \t]  { /* ignore whitespace */ }
.      { return ILLEGAL; }

%%

