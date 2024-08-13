package main

import (
	"go/ast"
	"go/parser"
)

func main() {
	expr, _ := parser.ParseExpr(`int(x)`)
	ast.Print(nil, expr)
}
