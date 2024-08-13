package main

import (
	"go/ast"
	"go/parser"
)

func main() {
	expr, _ := parser.ParseExpr(`x.(y)`)
	ast.Print(nil, expr)
}
