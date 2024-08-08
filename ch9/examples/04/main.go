package main

import (
	"go/ast"
	"go/parser"
)

func main() {
	expr, _ := parser.ParseExpr(`x[1:2:3]`)
	ast.Print(nil, expr)
}
