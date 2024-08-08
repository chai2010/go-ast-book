package main

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	calc "hello-antlr/calc"
)

func main() {
	if false {
		lexer := calc.NewCalcLexer(antlr.NewInputStream("1+2*3"))
		parser := calc.NewCalcParser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		var v calcVisitor
		parser.Start().Accept(&v)
		fmt.Println(v.pop())
	}
	if true {
		lexer := calc.NewCalcLexer(antlr.NewInputStream("2+2*5"))
		parser := calc.NewCalcParser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		var l calc.CalcListener = new(calcListener)
		antlr.NewParseTreeWalker().Walk(l, parser.Start())
		fmt.Println(l.(*calcListener).pop())
	}
}
