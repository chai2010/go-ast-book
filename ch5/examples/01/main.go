package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.ImportsOnly)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range f.Imports {
		fmt.Printf("import: name = %v, path = %#v\n", s.Name, s.Path)
	}
}

const src = `package foo
import "pkg-a"
import pkg_b_v2 "pkg-b"
import _ "pkg-c"
import . "pkg-d"
`
