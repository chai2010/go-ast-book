package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	ast.Print(nil, f.Decls[0].(*ast.GenDecl).Specs[0])
}

const src = `package foo
type Node struct {
	Next *Node
}
`
