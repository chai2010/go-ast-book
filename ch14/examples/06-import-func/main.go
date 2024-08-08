package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"

	"github.com/wa-lang/ssago/06-import-func/watypes"
	"golang.org/x/tools/go/ssa"
)

const src = `
package main

func my_print(s string)
func main() {
	my_print("Hello, wa!")
}
`

func my_print(args ...watypes.Value) watypes.Value {
	fmt.Print("my_print: ")
	for _, a := range args {
		switch a := a.(type) {
		case []watypes.Value:
			for _, a := range a {
				fmt.Print(a)
			}
		default:
			fmt.Print(a)
		}
	}
	return nil
}

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

	user_funcs := make(map[string]UserFunc)
	user_funcs["my_print"] = my_print

	p := NewEngine(ssaPkg, user_funcs)
	p.initGlobals()

	p.runFunc(ssaPkg.Func("main"), nil)
}
