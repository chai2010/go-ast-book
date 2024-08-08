package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast.Print(nil, f)
}

const src = `package foo

import  "bar"
import  "bar2"

type MyInt int

const PI = 3.14

var size = 10

func hello() func() {return func(){size++}}
`

/*
     0  *ast.File {
     1  .  Package: 1
     2  .  Name: *ast.Ident {
     3  .  .  NamePos: 9
     4  .  .  Name: "foo"
     5  .  }
     6  .  Decls: []ast.Decl (len = 5) {
     7  .  .  0: *ast.GenDecl {
     8  .  .  .  TokPos: 14
     9  .  .  .  Tok: import
    10  .  .  .  Lparen: 0
    11  .  .  .  Specs: []ast.Spec (len = 1) {
    12  .  .  .  .  0: *ast.ImportSpec {
    13  .  .  .  .  .  Path: *ast.BasicLit {
    14  .  .  .  .  .  .  ValuePos: 22
    15  .  .  .  .  .  .  Kind: STRING
    16  .  .  .  .  .  .  Value: "\"bar\""
    17  .  .  .  .  .  }
    18  .  .  .  .  .  EndPos: 0
    19  .  .  .  .  }
    20  .  .  .  }
    21  .  .  .  Rparen: 0
    22  .  .  }
    23  .  .  1: *ast.GenDecl {
    24  .  .  .  TokPos: 29
    25  .  .  .  Tok: type
    26  .  .  .  Lparen: 0
    27  .  .  .  Specs: []ast.Spec (len = 1) {
    28  .  .  .  .  0: *ast.TypeSpec {
    29  .  .  .  .  .  Name: *ast.Ident {
    30  .  .  .  .  .  .  NamePos: 34
    31  .  .  .  .  .  .  Name: "MyInt"
    32  .  .  .  .  .  .  Obj: *ast.Object {
    33  .  .  .  .  .  .  .  Kind: type
    34  .  .  .  .  .  .  .  Name: "MyInt"
    35  .  .  .  .  .  .  .  Decl: *(obj @ 28)
    36  .  .  .  .  .  .  }
    37  .  .  .  .  .  }
    38  .  .  .  .  .  Assign: 0
    39  .  .  .  .  .  Type: *ast.Ident {
    40  .  .  .  .  .  .  NamePos: 40
    41  .  .  .  .  .  .  Name: "int"
    42  .  .  .  .  .  }
    43  .  .  .  .  }
    44  .  .  .  }
    45  .  .  .  Rparen: 0
    46  .  .  }
    47  .  .  2: *ast.GenDecl {
    48  .  .  .  TokPos: 45
    49  .  .  .  Tok: const
    50  .  .  .  Lparen: 0
    51  .  .  .  Specs: []ast.Spec (len = 1) {
    52  .  .  .  .  0: *ast.ValueSpec {
    53  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    54  .  .  .  .  .  .  0: *ast.Ident {
    55  .  .  .  .  .  .  .  NamePos: 51
    56  .  .  .  .  .  .  .  Name: "PI"
    57  .  .  .  .  .  .  .  Obj: *ast.Object {
    58  .  .  .  .  .  .  .  .  Kind: const
    59  .  .  .  .  .  .  .  .  Name: "PI"
    60  .  .  .  .  .  .  .  .  Decl: *(obj @ 52)
    61  .  .  .  .  .  .  .  .  Data: 0
    62  .  .  .  .  .  .  .  }
    63  .  .  .  .  .  .  }
    64  .  .  .  .  .  }
    65  .  .  .  .  .  Values: []ast.Expr (len = 1) {
    66  .  .  .  .  .  .  0: *ast.BasicLit {
    67  .  .  .  .  .  .  .  ValuePos: 56
    68  .  .  .  .  .  .  .  Kind: FLOAT
    69  .  .  .  .  .  .  .  Value: "3.14"
    70  .  .  .  .  .  .  }
    71  .  .  .  .  .  }
    72  .  .  .  .  }
    73  .  .  .  }
    74  .  .  .  Rparen: 0
    75  .  .  }
    76  .  .  3: *ast.GenDecl {
    77  .  .  .  TokPos: 62
    78  .  .  .  Tok: var
    79  .  .  .  Lparen: 0
    80  .  .  .  Specs: []ast.Spec (len = 1) {
    81  .  .  .  .  0: *ast.ValueSpec {
    82  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    83  .  .  .  .  .  .  0: *ast.Ident {
    84  .  .  .  .  .  .  .  NamePos: 66
    85  .  .  .  .  .  .  .  Name: "size"
    86  .  .  .  .  .  .  .  Obj: *ast.Object {
    87  .  .  .  .  .  .  .  .  Kind: var
    88  .  .  .  .  .  .  .  .  Name: "size"
    89  .  .  .  .  .  .  .  .  Decl: *(obj @ 81)
    90  .  .  .  .  .  .  .  .  Data: 0
    91  .  .  .  .  .  .  .  }
    92  .  .  .  .  .  .  }
    93  .  .  .  .  .  }
    94  .  .  .  .  .  Values: []ast.Expr (len = 1) {
    95  .  .  .  .  .  .  0: *ast.BasicLit {
    96  .  .  .  .  .  .  .  ValuePos: 73
    97  .  .  .  .  .  .  .  Kind: INT
    98  .  .  .  .  .  .  .  Value: "10"
    99  .  .  .  .  .  .  }
   100  .  .  .  .  .  }
   101  .  .  .  .  }
   102  .  .  .  }
   103  .  .  .  Rparen: 0
   104  .  .  }
   105  .  .  4: *ast.FuncDecl {
   106  .  .  .  Name: *ast.Ident {
   107  .  .  .  .  NamePos: 82
   108  .  .  .  .  Name: "hello"
   109  .  .  .  .  Obj: *ast.Object {
   110  .  .  .  .  .  Kind: func
   111  .  .  .  .  .  Name: "hello"
   112  .  .  .  .  .  Decl: *(obj @ 105)
   113  .  .  .  .  }
   114  .  .  .  }
   115  .  .  .  Type: *ast.FuncType {
   116  .  .  .  .  Func: 77
   117  .  .  .  .  Params: *ast.FieldList {
   118  .  .  .  .  .  Opening: 87
   119  .  .  .  .  .  Closing: 88
   120  .  .  .  .  }
   121  .  .  .  }
   122  .  .  .  Body: *ast.BlockStmt {
   123  .  .  .  .  Lbrace: 90
   124  .  .  .  .  Rbrace: 91
   125  .  .  .  }
   126  .  .  }
   127  .  }
   128  .  Scope: *ast.Scope {
   129  .  .  Objects: map[string]*ast.Object (len = 4) {
   130  .  .  .  "hello": *(obj @ 109)
   131  .  .  .  "MyInt": *(obj @ 32)
   132  .  .  .  "PI": *(obj @ 57)
   133  .  .  .  "size": *(obj @ 86)
   134  .  .  }
   135  .  }
   136  .  Imports: []*ast.ImportSpec (len = 1) {
   137  .  .  0: *(obj @ 12)
   138  .  }
   139  .  Unresolved: []*ast.Ident (len = 1) {
   140  .  .  0: *(obj @ 39)
   141  .  }
   142  }
*/
