// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a Apache
// license that can be found in the LICENSE file.

package main

//#include "tok.h"
//#include "calc.lex.h"
import "C"

import (
	"errors"
	"log"
	"strconv"
)

type calcLex struct {
	yylineno int
	yytext   string
	lastErr  error
}

var _ calcLexer = (*calcLex)(nil)

func newCalcLexer(data []byte) *calcLex {
	p := new(calcLex)

	C.yy_scan_bytes(
		(*C.char)(C.CBytes(data)),
		C.yy_size_t(len(data)),
	)

	return p
}

// The parser calls this method to get each new token. This
// implementation returns operators and NUM.
func (p *calcLex) Lex(yylval *calcSymType) int {
	p.lastErr = nil

	var tok = C.yylex()

	p.yylineno = int(C.yylineno)
	p.yytext = C.GoString(C.yytext)

	switch tok {
	case C.ID:
		yylval.id = p.yytext
		return ID

	case C.NUMBER:
		yylval.value, _ = strconv.Atoi(p.yytext)
		return NUMBER

	case C.ADD:
		return ADD
	case C.SUB:
		return SUB
	case C.MUL:
		return MUL
	case C.DIV:
		return DIV
	case C.ABS:
		return ABS

	case C.LPAREN:
		return LPAREN
	case C.RPAREN:
		return RPAREN
	case C.ASSIGN:
		return ASSIGN

	case C.EOL:
		return EOL
	}

	if tok == C.ILLEGAL {
		log.Printf("lex: ILLEGAL token, yytext = %q, yylineno = %d", p.yytext, p.yylineno)
	}

	return 0 // eof
}

// The parser calls this method on a parse error.
func (p *calcLex) Error(s string) {
	p.lastErr = errors.New("yacc: " + s)
	if err := p.lastErr; err != nil {
		log.Println(err)
	}
}
