package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"go/ast"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	src := `package foo

import (
	"pkg-a"
	"pkg-b"
)

func bar() {}`

	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, s := range f.Imports {
		fmt.Printf("%#v\n", s.Path)
	}
	
	fmt.Println("----")
	
	for _, v := range f.Decls {
		// import group
		if s, ok := v.(*ast.GenDecl); ok && s.Tok == token.IMPORT {
			for _, v := range s.Specs {
				fmt.Printf("%#v\n", v.(*ast. ImportSpec).Path)
			}
		}
	}
}
