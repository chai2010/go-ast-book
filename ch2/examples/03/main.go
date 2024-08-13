package main

import (
	"go/ast"
	"go/parser"
)

func main() {
	expr, _ := parser.ParseExpr(`x`)
	ast.Print(nil, expr)

	ast.Print(nil, ast.NewIdent(`x`))
}
