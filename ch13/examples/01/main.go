package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/interp"
)

func main() {
	prog := NewProgram(map[string]string{
		"main": `
			package main

			func main() {
				for i := 0; i < 3; i++ {
					println(i, "hello chai2010")
				}
			}
		`,
		"runtime": `
			package runtime

			type errorString string

			func (e errorString) RuntimeError() {}
			func (e errorString) Error() string { return "runtime error: " + string(e) }

			type Error interface {
				error
				RuntimeError()
			}
		`,
	})

	prog.LoadPackage("main")
	prog.LoadPackage("runtime")

	var ssaProg = ssa.NewProgram(prog.fset, ssa.SanityCheckFunctions)
	var ssaMainPkg *ssa.Package

	for name, pkg := range prog.pkgs {
		ssaPkg := ssaProg.CreatePackage(pkg, []*ast.File{prog.ast[name]}, prog.infos[name], true)
		if name == "main" {
			ssaMainPkg = ssaPkg
		}
	}
	ssaProg.Build()

	exitCode := interp.Interpret(
		ssaMainPkg, 0, &types.StdSizes{8, 8},
		"main", []string{},
	)
	if exitCode != 0 {
		fmt.Println("exitCode:", exitCode)
	}
}

type Program struct {
	fs    map[string]string
	ast   map[string]*ast.File
	pkgs  map[string]*types.Package
	infos map[string]*types.Info
	fset  *token.FileSet
}

func NewProgram(fs map[string]string) *Program {
	return &Program{
		fs:    fs,
		ast:   make(map[string]*ast.File),
		pkgs:  make(map[string]*types.Package),
		infos: make(map[string]*types.Info),
		fset:  token.NewFileSet(),
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

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	conf := types.Config{Importer: p}
	pkg, err = conf.Check(path, p.fset, []*ast.File{f}, info)
	if err != nil {
		return nil, nil, err
	}

	p.ast[path] = f
	p.pkgs[path] = pkg
	p.infos[path] = info
	return pkg, f, nil
}

func (p *Program) Import(path string) (*types.Package, error) {
	if pkg, ok := p.pkgs[path]; ok {
		return pkg, nil
	}
	pkg, _, err := p.LoadPackage(path)
	return pkg, err
}
