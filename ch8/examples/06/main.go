package main

import (
	"go/ast"
	"go/parser"
)

func main() {
	expr, _ := parser.ParseExpr(`map[int]int{1:2}`)
	ast.Print(nil, expr)
}
