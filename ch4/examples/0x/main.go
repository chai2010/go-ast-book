package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {

	/*
	func ParseDir(fset *token.FileSet, path string, filter func(os.FileInfo) bool, mode Mode) (pkgs map[string]*ast.Package, first error)
	*/
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast.Print(nil, f)
}

const src = `package foo

import (
	"fmt"
	"time"
)

func bar() {
	fmt.Println(time.Now())
}
`

// 0  *ast.File {
// 1  .  Package: 1
// 2  .  Name: *ast.Ident {
// 3  .  .  NamePos: 9
// 4  .  .  Name: "foo"
// 5  .  }
// 6  .  Scope: *ast.Scope {
// 7  .  .  Objects: map[string]*ast.Object (len = 0) {}
// 8  .  }
// 9  }
