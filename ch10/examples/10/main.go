package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

var x int

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
		return
	}

	ast.Print(nil, f.Decls[0].(*ast.FuncDecl).Body)
}

const src = `package pkgname
func main() {
	go hello("光谷码农")
}
`
