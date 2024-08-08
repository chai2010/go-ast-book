package main

import (
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	parser "hello-antlr/calc"
)

type calcVisitor struct {
	parser.BaseCalcVisitor
	Stack
}

func (v *calcVisitor) visitRule(node antlr.RuleNode) interface{} {
	node.Accept(v)
	return nil
}

func (v *calcVisitor) VisitStart(ctx *parser.StartContext) interface{} {
	return v.visitRule(ctx.Expression())
}

func (v *calcVisitor) VisitNumber(ctx *parser.NumberContext) interface{} {
	i, _ := strconv.Atoi(ctx.NUMBER().GetText())
	v.push(i)
	return nil
}

func (v *calcVisitor) VisitMulDiv(ctx *parser.MulDivContext) interface{} {
	v.visitRule(ctx.Expression(0))
	v.visitRule(ctx.Expression(1))

	right, left := v.pop(), v.pop()
	switch ctx.GetOp().GetTokenType() {
	case parser.CalcParserMUL:
		v.push(left * right)
	case parser.CalcParserDIV:
		v.push(left / right)
	}
	return nil
}

func (v *calcVisitor) VisitAddSub(ctx *parser.AddSubContext) interface{} {
	v.visitRule(ctx.Expression(0))
	v.visitRule(ctx.Expression(1))

	right, left := v.pop(), v.pop()
	switch ctx.GetOp().GetTokenType() {
	case parser.CalcParserADD:
		v.push(left + right)
	case parser.CalcParserSUB:
		v.push(left - right)
	}
	return nil
}
