package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

func init() {
	log.SetFlags(0)
}

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	pkg, err := new(types.Config).Check("hello.go", fset, []*ast.File{f}, nil)
	if err != nil {
		log.Fatal(err)
	}

	_ = pkg
}

const src = `package main

func main() {
	var _ = "a" + 1
}
`
