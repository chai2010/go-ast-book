package main

import (
	"go/ast"
	"go/parser"
)

func main() {
	expr, _ := parser.ParseExpr(`func(){}`)
	ast.Print(nil, expr)
}
