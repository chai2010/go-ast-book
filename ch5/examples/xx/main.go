package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast.Print(nil, f)
}

const src = `package foo

import (_ "fmt")

type MyInt int

var x MyInt

const Pi = 3.14

var (
	a = 1
	b = 2
)

func foo() {}
`

/*
$ go doc go/ast | grep Decl
func FilterDecl(decl Decl, f Filter) bool
type BadDecl struct{ ... }
type Decl interface{ ... }
type DeclStmt struct{ ... }
type FuncDecl struct{ ... }
type GenDecl struct{ ... }
*/
