package main

import (
	"go/ast"
	"go/token"
)

func main() {
	var lit9527 = &ast.BasicLit{
		Kind:  token.INT,
		Value: "9527",
	}
	ast.Print(nil, lit9527)
}
