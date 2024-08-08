package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"

	"golang.org/x/tools/go/ssa"
)

const src = `
package main

var s = "hello ssa"

func main() {
	for i := 0; i < 3; i++ {
		println(s)
	}
}
`

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	conf := types.Config{Importer: nil}
	pkg, err := conf.Check("hello.go", fset, []*ast.File{f}, info)
	if err != nil {
		log.Fatal(err)
	}

	var ssaProg = ssa.NewProgram(fset, ssa.SanityCheckFunctions)
	var ssaPkg = ssaProg.CreatePackage(pkg, []*ast.File{f}, info, true)

	ssaPkg.Build()

	ssaPkg.WriteTo(os.Stdout)
	ssaPkg.Func("init").WriteTo(os.Stdout)
	ssaPkg.Func("main").WriteTo(os.Stdout)
}

/*
package hello.go:
  func  init       func()
  var   init$guard bool
  func  main       func()

# Name: hello.go.main
# Package: hello.go
# Location: hello.go:4:6
func main():
0:                                                                entry P:0 S:1
        jump 3
1:                                                             for.body P:1 S:1
        t0 = println("hello ssa -- chai...":string)                          ()
        t1 = t2 + 1:int                                                     int
        jump 3
2:                                                             for.done P:1 S:0
        return
3:                                                             for.loop P:2 S:2
        t2 = phi [0: 0:int, 1: t1] #i                                       int
        t3 = t2 < 3:int                                                    bool
        if t3 goto 1 else 2
*/
