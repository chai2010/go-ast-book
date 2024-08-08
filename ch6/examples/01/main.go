package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func hello(string,  int) {
	//println(string)
}

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		//ast.Print(nil, decl)
		if fn, ok := decl.(*ast.FuncDecl); ok {
			fmt.Println("func name: ", fn.Name)
			//ast.Print(nil, fn.Recv)
			ast.Print(nil, fn.Type.Params.List)
		}
	}
}

const src = `package foo
func hello(string,  string)
func hello2(a string, b string)
//func hello(s0, s1 string,  string, f func(a, b int))
`
