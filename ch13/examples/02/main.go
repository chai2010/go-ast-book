package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/ssa"
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
	})

	pkg, f, info, _ := prog.LoadPackage("main")

	var ssaProg = ssa.NewProgram(prog.fset, ssa.SanityCheckFunctions)
	var ssaPkg = ssaProg.CreatePackage(pkg, []*ast.File{f}, info, true)
	ssaPkg.Build()

	ssaPkg.WriteTo(os.Stdout)

	ssaPkg.Func("main").WriteTo(os.Stdout)
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

func (p *Program) LoadPackage(path string) (pkg *types.Package, f *ast.File, info *types.Info, err error) {
	if pkg, ok := p.pkgs[path]; ok {
		return pkg, p.ast[path], p.infos[path], nil
	}

	f, err = parser.ParseFile(p.fset, path, p.fs[path], parser.AllErrors)
	if err != nil {
		return nil, nil, nil, err
	}

	info = &types.Info{
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
		return nil, nil, nil, err
	}

	p.ast[path] = f
	p.pkgs[path] = pkg
	p.infos[path] = info
	return pkg, f, info, nil
}

func (p *Program) Import(path string) (*types.Package, error) {
	if pkg, ok := p.pkgs[path]; ok {
		return pkg, nil
	}
	pkg, _, _, err := p.LoadPackage(path)
	return pkg, err
}
