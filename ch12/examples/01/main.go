package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func main() {
	prog := NewProgram(map[string]string{
		"hello.go": `
			package main

			import "fmt"

			const Pi = 3.14

			func main() {
				for i := 2; i <= 8; i++ {
					fmt.Printf("%d*Pi = %.2f\n", i, Pi*float64(i))
				}
			}
		`,
		"fmt": `
			package fmt

			func Printf(format string, a ...interface{}) (n int, err error) {
				return
			}
		`,
	})

	pkg, _, err := prog.LoadPackage("hello.go")
	if err != nil {
		log.Fatal(err)
	}

	pkg.Scope().WriteTo(os.Stdout, 0, true)
	pkg.Scope().Parent().WriteTo(os.Stdout, 0, true)
}

type Program struct {
	fs   map[string]string
	ast  map[string]*ast.File
	pkgs map[string]*types.Package
	fset *token.FileSet
}

func NewProgram(fs map[string]string) *Program {
	return &Program{
		fs:   fs,
		ast:  make(map[string]*ast.File),
		pkgs: make(map[string]*types.Package),
		fset: token.NewFileSet(),
	}
}

func (p *Program) LoadPackage(path string) (pkg *types.Package, f *ast.File, err error) {
	if pkg, ok := p.pkgs[path]; ok {
		return pkg, p.ast[path], nil
	}

	f, err = parser.ParseFile(p.fset, path, p.fs[path], parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}

	conf := types.Config{Importer: p}
	pkg, err = conf.Check(path, p.fset, []*ast.File{f}, nil)
	if err != nil {
		return nil, nil, err
	}

	p.ast[path] = f
	p.pkgs[path] = pkg
	return pkg, f, nil
}

func (p *Program) Import(path string) (*types.Package, error) {
	if pkg, ok := p.pkgs[path]; ok {
		return pkg, nil
	}
	pkg, _, err := p.LoadPackage(path)
	return pkg, err
}
