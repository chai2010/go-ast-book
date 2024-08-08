package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.Trace)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast.Print(nil, f)
}

const src = `package foo

import (
	"fmt"
	"time"
)

func bar() {
	fmt.Println(time.Now())
}
`

// mode: parser.Trace

/*
    1:  1: File (
    1:  1: . "package"
    1:  9: . IDENT foo
    1: 12: . ";"
    3:  1: . GenDecl(import) (
    3:  1: . . "import"
    3:  8: . . "("
    4:  2: . . ImportSpec (
    4:  2: . . . STRING "fmt"
    4:  7: . . . ";"
    5:  2: . . )
    5:  2: . . ImportSpec (
    5:  2: . . . STRING "time"
    5:  8: . . . ";"
    6:  1: . . )
    6:  1: . . ")"
    6:  2: . . ";"
    8:  1: . )
    8:  1: . Declaration (
    8:  1: . . FunctionDecl (
    8:  1: . . . "func"
    8:  6: . . . IDENT bar
    8:  9: . . . Signature (
    8:  9: . . . . Parameters (
    8:  9: . . . . . "("
    8: 10: . . . . . ")"
    8: 12: . . . . )
    8: 12: . . . . Result (
    8: 12: . . . . )
    8: 12: . . . )
    8: 12: . . . Body (
    8: 12: . . . . "{"
    9:  2: . . . . StatementList (
    9:  2: . . . . . Statement (
    9:  2: . . . . . . SimpleStmt (
    9:  2: . . . . . . . ExpressionList (
    9:  2: . . . . . . . . Expression (
    9:  2: . . . . . . . . . BinaryExpr (
    9:  2: . . . . . . . . . . UnaryExpr (
    9:  2: . . . . . . . . . . . PrimaryExpr (
    9:  2: . . . . . . . . . . . . Operand (
    9:  2: . . . . . . . . . . . . . IDENT fmt
    9:  5: . . . . . . . . . . . . )
    9:  5: . . . . . . . . . . . . "."
    9:  6: . . . . . . . . . . . . Selector (
    9:  6: . . . . . . . . . . . . . IDENT Println
    9: 13: . . . . . . . . . . . . )
    9: 13: . . . . . . . . . . . . CallOrConversion (
    9: 13: . . . . . . . . . . . . . "("
    9: 14: . . . . . . . . . . . . . Expression (
    9: 14: . . . . . . . . . . . . . . BinaryExpr (
    9: 14: . . . . . . . . . . . . . . . UnaryExpr (
    9: 14: . . . . . . . . . . . . . . . . PrimaryExpr (
    9: 14: . . . . . . . . . . . . . . . . . Operand (
    9: 14: . . . . . . . . . . . . . . . . . . IDENT time
    9: 18: . . . . . . . . . . . . . . . . . )
    9: 18: . . . . . . . . . . . . . . . . . "."
    9: 19: . . . . . . . . . . . . . . . . . Selector (
    9: 19: . . . . . . . . . . . . . . . . . . IDENT Now
    9: 22: . . . . . . . . . . . . . . . . . )
    9: 22: . . . . . . . . . . . . . . . . . CallOrConversion (
    9: 22: . . . . . . . . . . . . . . . . . . "("
    9: 23: . . . . . . . . . . . . . . . . . . ")"
    9: 24: . . . . . . . . . . . . . . . . . )
    9: 24: . . . . . . . . . . . . . . . . )
    9: 24: . . . . . . . . . . . . . . . )
    9: 24: . . . . . . . . . . . . . . )
    9: 24: . . . . . . . . . . . . . )
    9: 24: . . . . . . . . . . . . . ")"
    9: 25: . . . . . . . . . . . . )
    9: 25: . . . . . . . . . . . )
    9: 25: . . . . . . . . . . )
    9: 25: . . . . . . . . . )
    9: 25: . . . . . . . . )
    9: 25: . . . . . . . )
    9: 25: . . . . . . )
    9: 25: . . . . . . ";"
   10:  1: . . . . . )
   10:  1: . . . . )
   10:  1: . . . . "}"
   10:  2: . . . )
   10:  2: . . . ";"
   10:  3: . . )
   10:  3: . )
   10:  3: )
     0  *ast.File {
     1  .  Package: 1
     2  .  Name: *ast.Ident {
     3  .  .  NamePos: 9
     4  .  .  Name: "foo"
     5  .  }
     6  .  Decls: []ast.Decl (len = 2) {
     7  .  .  0: *ast.GenDecl {
     8  .  .  .  TokPos: 14
     9  .  .  .  Tok: import
    10  .  .  .  Lparen: 21
    11  .  .  .  Specs: []ast.Spec (len = 2) {
    12  .  .  .  .  0: *ast.ImportSpec {
    13  .  .  .  .  .  Path: *ast.BasicLit {
    14  .  .  .  .  .  .  ValuePos: 24
    15  .  .  .  .  .  .  Kind: STRING
    16  .  .  .  .  .  .  Value: "\"fmt\""
    17  .  .  .  .  .  }
    18  .  .  .  .  .  EndPos: 0
    19  .  .  .  .  }
    20  .  .  .  .  1: *ast.ImportSpec {
    21  .  .  .  .  .  Path: *ast.BasicLit {
    22  .  .  .  .  .  .  ValuePos: 31
    23  .  .  .  .  .  .  Kind: STRING
    24  .  .  .  .  .  .  Value: "\"time\""
    25  .  .  .  .  .  }
    26  .  .  .  .  .  EndPos: 0
    27  .  .  .  .  }
    28  .  .  .  }
    29  .  .  .  Rparen: 38
    30  .  .  }
    31  .  .  1: *ast.FuncDecl {
    32  .  .  .  Name: *ast.Ident {
    33  .  .  .  .  NamePos: 46
    34  .  .  .  .  Name: "bar"
    35  .  .  .  .  Obj: *ast.Object {
    36  .  .  .  .  .  Kind: func
    37  .  .  .  .  .  Name: "bar"
    38  .  .  .  .  .  Decl: *(obj @ 31)
    39  .  .  .  .  }
    40  .  .  .  }
    41  .  .  .  Type: *ast.FuncType {
    42  .  .  .  .  Func: 41
    43  .  .  .  .  Params: *ast.FieldList {
    44  .  .  .  .  .  Opening: 49
    45  .  .  .  .  .  Closing: 50
    46  .  .  .  .  }
    47  .  .  .  }
    48  .  .  .  Body: *ast.BlockStmt {
    49  .  .  .  .  Lbrace: 52
    50  .  .  .  .  List: []ast.Stmt (len = 1) {
    51  .  .  .  .  .  0: *ast.ExprStmt {
    52  .  .  .  .  .  .  X: *ast.CallExpr {
    53  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr {
    54  .  .  .  .  .  .  .  .  X: *ast.Ident {
    55  .  .  .  .  .  .  .  .  .  NamePos: 55
    56  .  .  .  .  .  .  .  .  .  Name: "fmt"
    57  .  .  .  .  .  .  .  .  }
    58  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
    59  .  .  .  .  .  .  .  .  .  NamePos: 59
    60  .  .  .  .  .  .  .  .  .  Name: "Println"
    61  .  .  .  .  .  .  .  .  }
    62  .  .  .  .  .  .  .  }
    63  .  .  .  .  .  .  .  Lparen: 66
    64  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
    65  .  .  .  .  .  .  .  .  0: *ast.CallExpr {
    66  .  .  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr {
    67  .  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
    68  .  .  .  .  .  .  .  .  .  .  .  NamePos: 67
    69  .  .  .  .  .  .  .  .  .  .  .  Name: "time"
    70  .  .  .  .  .  .  .  .  .  .  }
    71  .  .  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
    72  .  .  .  .  .  .  .  .  .  .  .  NamePos: 72
    73  .  .  .  .  .  .  .  .  .  .  .  Name: "Now"
    74  .  .  .  .  .  .  .  .  .  .  }
    75  .  .  .  .  .  .  .  .  .  }
    76  .  .  .  .  .  .  .  .  .  Lparen: 75
    77  .  .  .  .  .  .  .  .  .  Ellipsis: 0
    78  .  .  .  .  .  .  .  .  .  Rparen: 76
    79  .  .  .  .  .  .  .  .  }
    80  .  .  .  .  .  .  .  }
    81  .  .  .  .  .  .  .  Ellipsis: 0
    82  .  .  .  .  .  .  .  Rparen: 77
    83  .  .  .  .  .  .  }
    84  .  .  .  .  .  }
    85  .  .  .  .  }
    86  .  .  .  .  Rbrace: 79
    87  .  .  .  }
    88  .  .  }
    89  .  }
    90  .  Scope: *ast.Scope {
    91  .  .  Objects: map[string]*ast.Object (len = 1) {
    92  .  .  .  "bar": *(obj @ 35)
    93  .  .  }
    94  .  }
    95  .  Imports: []*ast.ImportSpec (len = 2) {
    96  .  .  0: *(obj @ 12)
    97  .  .  1: *(obj @ 20)
    98  .  }
    99  .  Unresolved: []*ast.Ident (len = 2) {
   100  .  .  0: *(obj @ 54)
   101  .  .  1: *(obj @ 67)
   102  .  }
   103  }
*/
