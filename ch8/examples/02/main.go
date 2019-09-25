package main

import (
	"go/ast"
	"go/parser"
)

func main() {
	expr, _ := parser.ParseExpr(`interface{Close() error}{}`)
	ast.Print(nil, expr)
}
