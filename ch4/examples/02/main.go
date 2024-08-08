package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.ImportsOnly)
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

// mode: parser.ImportsOnly

//  0  *ast.File {
//  1  .  Package: 1
//  2  .  Name: *ast.Ident {
//  3  .  .  NamePos: 9
//  4  .  .  Name: "foo"
//  5  .  }
//  6  .  Decls: []ast.Decl (len = 1) {
//  7  .  .  0: *ast.GenDecl {
//  8  .  .  .  TokPos: 14
//  9  .  .  .  Tok: import
// 10  .  .  .  Lparen: 21
// 11  .  .  .  Specs: []ast.Spec (len = 2) {
// 12  .  .  .  .  0: *ast.ImportSpec {
// 13  .  .  .  .  .  Path: *ast.BasicLit {
// 14  .  .  .  .  .  .  ValuePos: 24
// 15  .  .  .  .  .  .  Kind: STRING
// 16  .  .  .  .  .  .  Value: "\"fmt\""
// 17  .  .  .  .  .  }
// 18  .  .  .  .  .  EndPos: 0
// 19  .  .  .  .  }
// 20  .  .  .  .  1: *ast.ImportSpec {
// 21  .  .  .  .  .  Path: *ast.BasicLit {
// 22  .  .  .  .  .  .  ValuePos: 31
// 23  .  .  .  .  .  .  Kind: STRING
// 24  .  .  .  .  .  .  Value: "\"time\""
// 25  .  .  .  .  .  }
// 26  .  .  .  .  .  EndPos: 0
// 27  .  .  .  .  }
// 28  .  .  .  }
// 29  .  .  .  Rparen: 38
// 30  .  .  }
// 31  .  }
// 32  .  Scope: *ast.Scope {
// 33  .  .  Objects: map[string]*ast.Object (len = 0) {}
// 34  .  }
// 35  .  Imports: []*ast.ImportSpec (len = 2) {
// 36  .  .  0: *(obj @ 12)
// 37  .  .  1: *(obj @ 20)
// 38  .  }
// 39  }
