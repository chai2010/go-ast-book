package main

import (
	"strconv"

	parser "hello-antlr/calc"
)

type calcListener struct {
	*parser.BaseCalcListener
	Stack
}

func (l *calcListener) ExitMulDiv(c *parser.MulDivContext) {
	right, left := l.pop(), l.pop()
	switch c.GetOp().GetTokenType() {
	case parser.CalcParserMUL:
		l.push(left * right)
	case parser.CalcParserDIV:
		l.push(left / right)
	}
}

func (l *calcListener) ExitAddSub(c *parser.AddSubContext) {
	right, left := l.pop(), l.pop()
	switch c.GetOp().GetTokenType() {
	case parser.CalcParserADD:
		l.push(left + right)
	case parser.CalcParserSUB:
		l.push(left - right)
	}
}

func (l *calcListener) ExitNumber(c *parser.NumberContext) {
	i, _ := strconv.Atoi(c.GetText())
	l.push(i)
}
