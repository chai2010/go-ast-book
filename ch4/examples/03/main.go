package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast.Print(nil, f)
}

const src = `
// package doc
package foo

import (
	"fmt"
	"time"
)

// doc1

// bar doc
func bar() {
	fmt.Println(time.Now())
}
`

// mode: parser.ParseComments

//   0  *ast.File {
//   1  .  Doc: *ast.CommentGroup {
//   2  .  .  List: []*ast.Comment (len = 1) {
//   3  .  .  .  0: *ast.Comment {
//   4  .  .  .  .  Slash: 2
//   5  .  .  .  .  Text: "// package doc"
//   6  .  .  .  }
//   7  .  .  }
//   8  .  }
//   9  .  Package: 17
//  10  .  Name: *ast.Ident {
//  11  .  .  NamePos: 25
//  12  .  .  Name: "foo"
//  13  .  }
//  14  .  Decls: []ast.Decl (len = 2) {
//  15  .  .  0: *ast.GenDecl {
//  16  .  .  .  TokPos: 30
//  17  .  .  .  Tok: import
//  18  .  .  .  Lparen: 37
//  19  .  .  .  Specs: []ast.Spec (len = 2) {
//  20  .  .  .  .  0: *ast.ImportSpec {
//  21  .  .  .  .  .  Path: *ast.BasicLit {
//  22  .  .  .  .  .  .  ValuePos: 40
//  23  .  .  .  .  .  .  Kind: STRING
//  24  .  .  .  .  .  .  Value: "\"fmt\""
//  25  .  .  .  .  .  }
//  26  .  .  .  .  .  EndPos: 0
//  27  .  .  .  .  }
//  28  .  .  .  .  1: *ast.ImportSpec {
//  29  .  .  .  .  .  Path: *ast.BasicLit {
//  30  .  .  .  .  .  .  ValuePos: 47
//  31  .  .  .  .  .  .  Kind: STRING
//  32  .  .  .  .  .  .  Value: "\"time\""
//  33  .  .  .  .  .  }
//  34  .  .  .  .  .  EndPos: 0
//  35  .  .  .  .  }
//  36  .  .  .  }
//  37  .  .  .  Rparen: 54
//  38  .  .  }
//  39  .  .  1: *ast.FuncDecl {
//  40  .  .  .  Doc: *ast.CommentGroup {
//  41  .  .  .  .  List: []*ast.Comment (len = 1) {
//  42  .  .  .  .  .  0: *ast.Comment {
//  43  .  .  .  .  .  .  Slash: 66
//  44  .  .  .  .  .  .  Text: "// bar doc"
//  45  .  .  .  .  .  }
//  46  .  .  .  .  }
//  47  .  .  .  }
//  48  .  .  .  Name: *ast.Ident {
//  49  .  .  .  .  NamePos: 82
//  50  .  .  .  .  Name: "bar"
//  51  .  .  .  .  Obj: *ast.Object {
//  52  .  .  .  .  .  Kind: func
//  53  .  .  .  .  .  Name: "bar"
//  54  .  .  .  .  .  Decl: *(obj @ 39)
//  55  .  .  .  .  }
//  56  .  .  .  }
//  57  .  .  .  Type: *ast.FuncType {
//  58  .  .  .  .  Func: 77
//  59  .  .  .  .  Params: *ast.FieldList {
//  60  .  .  .  .  .  Opening: 85
//  61  .  .  .  .  .  Closing: 86
//  62  .  .  .  .  }
//  63  .  .  .  }
//  64  .  .  .  Body: *ast.BlockStmt {
//  65  .  .  .  .  Lbrace: 88
//  66  .  .  .  .  List: []ast.Stmt (len = 1) {
//  67  .  .  .  .  .  0: *ast.ExprStmt {
//  68  .  .  .  .  .  .  X: *ast.CallExpr {
//  69  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr {
//  70  .  .  .  .  .  .  .  .  X: *ast.Ident {
//  71  .  .  .  .  .  .  .  .  .  NamePos: 91
//  72  .  .  .  .  .  .  .  .  .  Name: "fmt"
//  73  .  .  .  .  .  .  .  .  }
//  74  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
//  75  .  .  .  .  .  .  .  .  .  NamePos: 95
//  76  .  .  .  .  .  .  .  .  .  Name: "Println"
//  77  .  .  .  .  .  .  .  .  }
//  78  .  .  .  .  .  .  .  }
//  79  .  .  .  .  .  .  .  Lparen: 102
//  80  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
//  81  .  .  .  .  .  .  .  .  0: *ast.CallExpr {
//  82  .  .  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr {
//  83  .  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
//  84  .  .  .  .  .  .  .  .  .  .  .  NamePos: 103
//  85  .  .  .  .  .  .  .  .  .  .  .  Name: "time"
//  86  .  .  .  .  .  .  .  .  .  .  }
//  87  .  .  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
//  88  .  .  .  .  .  .  .  .  .  .  .  NamePos: 108
//  89  .  .  .  .  .  .  .  .  .  .  .  Name: "Now"
//  90  .  .  .  .  .  .  .  .  .  .  }
//  91  .  .  .  .  .  .  .  .  .  }
//  92  .  .  .  .  .  .  .  .  .  Lparen: 111
//  93  .  .  .  .  .  .  .  .  .  Ellipsis: 0
//  94  .  .  .  .  .  .  .  .  .  Rparen: 112
//  95  .  .  .  .  .  .  .  .  }
//  96  .  .  .  .  .  .  .  }
//  97  .  .  .  .  .  .  .  Ellipsis: 0
//  98  .  .  .  .  .  .  .  Rparen: 113
//  99  .  .  .  .  .  .  }
// 100  .  .  .  .  .  }
// 101  .  .  .  .  }
// 102  .  .  .  .  Rbrace: 115
// 103  .  .  .  }
// 104  .  .  }
// 105  .  }
// 106  .  Scope: *ast.Scope {
// 107  .  .  Objects: map[string]*ast.Object (len = 1) {
// 108  .  .  .  "bar": *(obj @ 51)
// 109  .  .  }
// 110  .  }
// 111  .  Imports: []*ast.ImportSpec (len = 2) {
// 112  .  .  0: *(obj @ 20)
// 113  .  .  1: *(obj @ 28)
// 114  .  }
// 115  .  Unresolved: []*ast.Ident (len = 2) {
// 116  .  .  0: *(obj @ 70)
// 117  .  .  1: *(obj @ 83)
// 118  .  }
// 119  .  Comments: []*ast.CommentGroup (len = 3) {
// 120  .  .  0: *(obj @ 1)
// 121  .  .  1: *ast.CommentGroup {
// 122  .  .  .  List: []*ast.Comment (len = 1) {
// 123  .  .  .  .  0: *ast.Comment {
// 124  .  .  .  .  .  Slash: 57
// 125  .  .  .  .  .  Text: "// doc1"
// 126  .  .  .  .  }
// 127  .  .  .  }
// 128  .  .  }
// 129  .  .  2: *(obj @ 40)
// 130  .  }
// 131  }
