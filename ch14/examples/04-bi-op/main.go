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

var i int

func add(i int, j int) int{
	return i + j
}

func main() {
	i = 24
	println("The answer is:", i + 3 * add(2, 4))
}
`

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.AllErrors)
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
	pkg, err := conf.Check("test.go", fset, []*ast.File{f}, info)
	if err != nil {
		log.Fatal(err)
	}

	var ssaProg = ssa.NewProgram(fset, ssa.SanityCheckFunctions)
	var ssaPkg = ssaProg.CreatePackage(pkg, []*ast.File{f}, info, true)

	ssaPkg.Build()
	ssaPkg.WriteTo(os.Stdout)
	ssaPkg.Func("main").WriteTo(os.Stdout)
	ssaPkg.Func("add").WriteTo(os.Stdout)

	p := NewEngine(ssaPkg)
	p.initGlobals()

	p.runFunc(ssaPkg.Func("main"), nil)
}
