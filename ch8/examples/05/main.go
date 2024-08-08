package main

import (
	"go/ast"
	"go/parser"
)

var x1 = struct{ X int }{1}
var x2 = struct{ X int }{X: 1}

func main() {
	expr, _ := parser.ParseExpr(`struct{X int}{X:1}`)
	ast.Print(nil, expr)
}
