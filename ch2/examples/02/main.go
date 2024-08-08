package main

import (
	"go/ast"
	"go/parser"
)

func main() {
	expr, _ := parser.ParseExpr(`9527`)
	ast.Print(nil, expr)
}
