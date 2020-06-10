# 第10章 语句块和语句

语句近似看作是函数体内可独立执行的代码，语句块是由大括弧定义的语句容器，语句块和语句只能在函数体内部定义。本章节我们学习语句块和语句的语法树构造。

## 10.1 语法规范

语句块和语句是在函数体部分定义，函数体就是一个语句块。语句块的语法规范如下：

```bnf
FunctionBody  = Block .

Block         = "{" StatementList "}" .
StatementList = { Statement ";" } .

Statement     = Declaration | LabeledStmt | SimpleStmt
              | GoStmt | ReturnStmt | BreakStmt | ContinueStmt | GotoStmt
              | FallthroughStmt | Block | IfStmt | SwitchStmt | SelectStmt | ForStmt
              | DeferStmt
              .
```

FunctionBody函数体对应一个Block语句块。每个Block语句块内部由多个语句列表StatementList组成，每个语句之间通过分号分隔。语句又可分为声明语句、标签语句、普通表达式语句和其它诸多控制流语句。需要注意的是，Block语句块也是一种合法的语句，因此函数体实际上是由Block组成的多叉树结构表示，每个Block结点又可以递归保存其他的可嵌套Block的控制流等语句。

## 10.2 空语句块

一个最简单的函数不仅仅没有任何的输入参数和返回值，函数体中也没有任何的语句。下面代码分析`func main() {}`函数体语法树结构：

```go
func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
		return
	}

	ast.Print(nil, f.Decls[0].(*ast.FuncDecl).Body)
}

const src = `package pkgname
func main() {}
`
```

函数的声明由`ast.FuncDecl`结构体定义，其中的Body成员是`ast.BlockStmt`类型。`ast.BlockStmt`类型的定义如下：

```go
type Stmt interface {
	Node
	// contains filtered or unexported methods
}
type BlockStmt struct {
	Lbrace token.Pos // position of "{"
	List   []Stmt
	Rbrace token.Pos // position of "}"
}
```

语句由ast.Stmt接口表示，各种具体的满足ast.Stmt接口的类型大多会以Stmt为后缀名。其中BlockStmt语句块也是一种语句，BlockStmt其实是一个语句容器，其中List成员是一个`[]ast.Stmt`语句列表。

`func main() {}`函数体部分输出的语法树结果如下：

```
0  *ast.BlockStmt {
1  .  Lbrace: 29
2  .  Rbrace: 30
3  }
```

表示函数体没有任何其它的语句。

因为由大括弧定义的语句块也是一种合法的语句，因此我们可以在函数体再定义任意个空的语句块：

```go
func main() {
	{}
	{}
}
```

再次分析函数体的语法树，可以得到以下的结果：

```
 0  *ast.BlockStmt {
 1  .  Lbrace: 29
 2  .  List: []ast.Stmt (len = 2) {
 3  .  .  0: *ast.BlockStmt {
 4  .  .  .  Lbrace: 32
 5  .  .  .  Rbrace: 33
 6  .  .  }
 7  .  .  1: *ast.BlockStmt {
 8  .  .  .  Lbrace: 36
 9  .  .  .  Rbrace: 37
10  .  .  }
11  .  }
12  .  Rbrace: 39
13  }
```

其中List部分有两个新定义的语句块，每个语句块依然是ast.BlockStmt类型。函数体中的语句块构成的语法树和类型中的语法树结构是很相似的，但是语句的语法树最大的特定是可以循环递归定义，而类型的语法树不能递归定义自身（语义层面禁止）。

## 10.3 表达式语句

实际上定义空的语句块并不能算真正的语句，它只是在编译阶段定义新的变量作用域，并没有产生新的语句或计算。最简单的语句是表达式语句，不管是简单表达式还是复杂的表达式都可以作为一个独立的语句。表达式语句语法规范如下：

```
ExpressionStmt = Expression .
```

其实一个表达式语句就是对应一个表达式，而关于表达式的语法我们已经学习过。我们这里以一个最简单的常量作为标识符，来研究表达式语句的语法结构。下面是只有一个常量表达式语句的main函数：

```go
func main() {
	42
}
```

输出的语句的语法树如下：

```
chai-mba:02 chai$ go run main.go 
 0  *ast.BlockStmt {
 1  .  Lbrace: 29
 2  .  List: []ast.Stmt (len = 1) {
 3  .  .  0: *ast.ExprStmt {
 4  .  .  .  X: *ast.BasicLit {
 5  .  .  .  .  ValuePos: 32
 6  .  .  .  .  Kind: INT
 7  .  .  .  .  Value: "42"
 8  .  .  .  }
 9  .  .  }
10  .  }
11  .  Rbrace: 35
12  }
```

表达式语句由`ast.ExprStmt`结构体定义：

```go
type ExprStmt struct {
	X Expr // expression
}
```

它只是`ast.Expr`表达式的再次包装，以满足`ast.Stmt`接口。因为`ast.Expr`表达式本身也是一个接口类型，因此可以包含任意复杂的表达式。表达式语句最终会产生一个值，但是表达式的值没有被赋值到变量，因此表达式的返回值会被丢弃。不过表达式中可能还有函数调用，而函数调用可能有其它的副作用，因此表达式语句一般常用于触发函数调用。

## 10.4 返回语句

表达式不仅仅可以作为独立的表达式语句，同时表达式也是其它更复杂控制流语句的组成单元。对于函数比较重要的控制流语句是返回语句，返回语句的语法规范如下：

```
ReturnStmt     = "return" [ ExpressionList ] .
ExpressionList = Expression { "," Expression } .
```

返回语句以return关键字开始，后面跟着多个以逗号分隔的表达式，当然也可以没有返回值。下面例子在main函数增加一个返回两个值的返回语句：

```go
func main() {
	return 42, err
}
```

当然，按照Go语言规范main函数是没有返回值的，因此return语句也不能有返回值。不过我们目前还处在语法树解析阶段，并不会检查返回语句和函数的返回值类型是否匹配，这种类型匹配检查要在语法树构建之后才会进行。

main函数体的语法树结果如下：

```
 0  *ast.BlockStmt {
 1  .  Lbrace: 29
 2  .  List: []ast.Stmt (len = 1) {
 3  .  .  0: *ast.ReturnStmt {
 4  .  .  .  Return: 32
 5  .  .  .  Results: []ast.Expr (len = 2) {
 6  .  .  .  .  0: *ast.BasicLit {
 7  .  .  .  .  .  ValuePos: 39
 8  .  .  .  .  .  Kind: INT
 9  .  .  .  .  .  Value: "42"
10  .  .  .  .  }
11  .  .  .  .  1: *ast.Ident {
12  .  .  .  .  .  NamePos: 43
13  .  .  .  .  .  Name: "err"
14  .  .  .  .  }
15  .  .  .  }
16  .  .  }
17  .  }
18  .  Rbrace: 47
19  }
```

返回语句由`ast.ReturnStmt`类型表示，其中Results成员对应返回值列表，这里分别是基础的数值常量42和标识符err。`ast.ReturnStmt`类型定义如下：

```go
type ReturnStmt struct {
	Return  token.Pos // position of "return" keyword
	Results []Expr    // result expressions; or nil
}
```

其中Return成员表示return关键字的位置，Results成员对应一个表达式列表，如果为nil表示没有返回值。

## 10.5 声明语句

函数中除了输入参数和返回值参数之外，还可以定义临时的局部变量保存函数的状态。如果临时变量被闭包函数捕获，那么临时变量维持的函数状态将伴随闭包函数的整个生命周期。因此声明变量和声明函数一样重要。声明变量的声明语法和顶级包变量的语法是类似的：

```
Declaration  = ConstDecl | TypeDecl | VarDecl .
TopLevelDecl = Declaration | FunctionDecl | MethodDecl .
```

其中Declaration就是函数体内部的声明语法，可以在函数内部声明常量、变量和类型，但是不能声明函数和方法。关于TopLevelDecl定义顶级常量、变量和类型声明我们已经讨论过，其中已经包含了函数内部的声明语法。我们这里以一个简单的例子展示如果在语句块中保存声明语句：

```go
func main() {
	var a int
}
```

在main函数内部定义一个int类型变量，这个语法格式和全局变量的定义是一样的。语法树解析输出如下：

```
 0  *ast.BlockStmt {
 1  .  Lbrace: 29
 2  .  List: []ast.Stmt (len = 1) {
 3  .  .  0: *ast.DeclStmt {
 4  .  .  .  Decl: *ast.GenDecl {
 5  .  .  .  .  TokPos: 32
 6  .  .  .  .  Tok: var
 7  .  .  .  .  Lparen: 0
 8  .  .  .  .  Specs: []ast.Spec (len = 1) {
 9  .  .  .  .  .  0: *ast.ValueSpec {
10  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
11  .  .  .  .  .  .  .  0: *ast.Ident {
12  .  .  .  .  .  .  .  .  NamePos: 36
13  .  .  .  .  .  .  .  .  Name: "a"
14  .  .  .  .  .  .  .  .  Obj: *ast.Object {...}
20  .  .  .  .  .  .  .  }
21  .  .  .  .  .  .  }
22  .  .  .  .  .  .  Type: *ast.Ident {
23  .  .  .  .  .  .  .  NamePos: 38
24  .  .  .  .  .  .  .  Name: "int"
25  .  .  .  .  .  .  }
26  .  .  .  .  .  }
27  .  .  .  .  }
28  .  .  .  .  Rparen: 0
29  .  .  .  }
30  .  .  }
31  .  }
32  .  Rbrace: 42
33  }
```

声明的变量在`ast.DeclStmt`结构体中表示，结构体定义如下：

```go
type DeclStmt struct {
	Decl Decl // *GenDecl with CONST, TYPE, or VAR token
}
```

虽然Decl成员是`ast.Decl`类型的接口，但是注释已经明确表示只有常量、类型和变量几种声明，并不包含函数和方法的声明。因此，Decl成员只能是`ast.GenDecl`类型。

## 10.6 短声明和多赋值语句

函数内变量还可以采用短声明方式。短声明语法和多赋值语句类似，它是在声明变量的同时进行多赋值初始化，变量的类型从赋值表达式自动推导。短声明和多赋值语句语法规范如下：

```
Assignment   = ExpressionList assign_op ExpressionList .
ShortVarDecl = IdentifierList ":=" ExpressionList .
```

其中多赋值语句的左边是一个表达式列表，而短声明语句的左边是一组标识符列表。短声明和多赋值语句的右边都是一组表达式列表。我们以一个短声明多个变量来展示短声明和多赋值语句的语法树：

```go
func main() {
	a, b := 1, 2
}
```

输出的语法树结果如下：

```
     0  *ast.BlockStmt {
     1  .  Lbrace: 29
     2  .  List: []ast.Stmt (len = 1) {
     3  .  .  0: *ast.AssignStmt {
     4  .  .  .  Lhs: []ast.Expr (len = 2) {
     5  .  .  .  .  0: *ast.Ident {
     6  .  .  .  .  .  NamePos: 32
     7  .  .  .  .  .  Name: "a"
     8  .  .  .  .  .  Obj: *ast.Object {
     9  .  .  .  .  .  .  Kind: var
    10  .  .  .  .  .  .  Name: "a"
    11  .  .  .  .  .  .  Decl: *(obj @ 3)
    12  .  .  .  .  .  }
    13  .  .  .  .  }
    14  .  .  .  .  1: *ast.Ident {
    15  .  .  .  .  .  NamePos: 35
    16  .  .  .  .  .  Name: "b"
    17  .  .  .  .  .  Obj: *ast.Object {
    18  .  .  .  .  .  .  Kind: var
    19  .  .  .  .  .  .  Name: "b"
    20  .  .  .  .  .  .  Decl: *(obj @ 3)
    21  .  .  .  .  .  }
    22  .  .  .  .  }
    23  .  .  .  }
    24  .  .  .  TokPos: 37
    25  .  .  .  Tok: :=
    26  .  .  .  Rhs: []ast.Expr (len = 2) {
    27  .  .  .  .  0: *ast.BasicLit {
    28  .  .  .  .  .  ValuePos: 40
    29  .  .  .  .  .  Kind: INT
    30  .  .  .  .  .  Value: "1"
    31  .  .  .  .  }
    32  .  .  .  .  1: *ast.BasicLit {
    33  .  .  .  .  .  ValuePos: 43
    34  .  .  .  .  .  Kind: INT
    35  .  .  .  .  .  Value: "2"
    36  .  .  .  .  }
    37  .  .  .  }
    38  .  .  }
    39  .  }
    40  .  Rbrace: 45
    41  }
```

短声明和多赋值语句都通过`ast.AssignStmt`结构体表达，其定义如下：

```go
type AssignStmt struct {
	Lhs    []Expr
	TokPos token.Pos   // position of Tok
	Tok    token.Token // assignment token, DEFINE
	Rhs    []Expr
}
```

其中Lhs表示左边的表达式或标识符列表，而Rhs表示右边的表达式列表。短声明和多赋值语句是通过Tok来进行区分。

## 10.7 if/else分支语句

顺序、分支和循环是编程语言中三种基本的控制流语句。Go语言的if语句语法规范如下：

```
IfStmt = "if" [ SimpleStmt ";" ] Expression Block [ "else" ( IfStmt | Block ) ] .
```

分支由if关键字开始，首先是可选的SimpleStmt简单初始化语句（可以是局部变量短声明、赋值或表达式等语句），然后是if的条件表达式，最后是分支的主体部分。分支的主体Block为一个语句块，其中可以包含多个语句或嵌套其它的语句块。同时if可以携带一个else分支，对应分支条件为假的情况。

我们以一个不带短声明的if/else为例：

```go
func main() {
	if true {} else {}
}
```

输出的语法树如下：

```
 0  *ast.BlockStmt {
 1  .  Lbrace: 29
 2  .  List: []ast.Stmt (len = 1) {
 3  .  .  0: *ast.IfStmt {
 4  .  .  .  If: 32
 5  .  .  .  Cond: *ast.Ident {
 6  .  .  .  .  NamePos: 35
 7  .  .  .  .  Name: "true"
 8  .  .  .  }
 9  .  .  .  Body: *ast.BlockStmt {
10  .  .  .  .  Lbrace: 40
11  .  .  .  .  Rbrace: 41
12  .  .  .  }
13  .  .  .  Else: *ast.BlockStmt {
14  .  .  .  .  Lbrace: 48
15  .  .  .  .  Rbrace: 49
16  .  .  .  }
17  .  .  }
18  .  }
19  .  Rbrace: 51
20  }
```

if由`ast.IfStmt`结构体表示，其中的Cond为分支的条件表达式，Body为分支的主体语句块，Else为补充的语句块。`ast.IfStmt`结构体完整定义如下：

```go
type IfStmt struct {
	If   token.Pos // position of "if" keyword
	Init Stmt      // initialization statement; or nil
	Cond Expr      // condition
	Body *BlockStmt
	Else Stmt // else branch; or nil
}
```

除了分支调整、主体块、补充块，还有Init用于初始化部分。需要注意的是Else都被定义为`ast.Stmt`接口类型，而Body被明确定义为`ast.BlockStmt`类型，是否是想以接口类型来暗示else可能为空的情况。

## 10.8 for循环

Go语言中只有一种for循环语句，但是for语句的语法却最为复杂。for语句的语法规范如下：

```
ForStmt     = "for" [ Condition | ForClause | RangeClause ] Block .

Condition   = Expression .

ForClause   = [ InitStmt ] ";" [ Condition ] ";" [ PostStmt ] .
InitStmt    = SimpleStmt .
PostStmt    = SimpleStmt .

RangeClause = [ ExpressionList "=" | IdentifierList ":=" ] "range" Expression .
```

分析语法规范，可以对应以下四种类型：

```go
for {}
for true {}
for i := 0; true; i++ {}
for i, v := range m {}
```

其中第一个没有循环条件，默认条件是true，因此和第二个循环语句一样都是死循环。而第一和第二个循环语句其实是第三个经典循环结构的特例，在第三个循环语句中增加的初始化语句和循环迭代语句。最后第四个循环语句是一种新的循环结构，终于用于数组、切片和map的迭代。以上四个循环语句可以再次归纳为以下两种：

```go
for x; y; z {}
for x, y := range z {}
```

除了map只能通过`for range`迭代之外（如果借助标准包，可以通过`reflect.MapKeys`或`reflect.MapRange`等方式迭代循环map），其它的`for range`格式的循环都可以通过`for x; y; z {}`经典风格的循环替代。

因此我们先分析经典风格的`for x; y; z {}`循环：

```go
func main() {
	for x; y; z {}
}
```

其语法树结构如下：

```
 0  *ast.BlockStmt {
 1  .  Lbrace: 29
 2  .  List: []ast.Stmt (len = 1) {
 3  .  .  0: *ast.ForStmt {
 4  .  .  .  For: 32
 5  .  .  .  Init: *ast.ExprStmt {
 6  .  .  .  .  X: *ast.Ident {
 7  .  .  .  .  .  NamePos: 36
 8  .  .  .  .  .  Name: "x"
 9  .  .  .  .  }
10  .  .  .  }
11  .  .  .  Cond: *ast.Ident {
12  .  .  .  .  NamePos: 39
13  .  .  .  .  Name: "y"
14  .  .  .  }
15  .  .  .  Post: *ast.ExprStmt {
16  .  .  .  .  X: *ast.Ident {
17  .  .  .  .  .  NamePos: 42
18  .  .  .  .  .  Name: "z"
19  .  .  .  .  }
20  .  .  .  }
21  .  .  .  Body: *ast.BlockStmt {
22  .  .  .  .  Lbrace: 44
23  .  .  .  .  Rbrace: 45
24  .  .  .  }
25  .  .  }
26  .  }
27  .  Rbrace: 47
28  }
```

`ast.ForStmt`结构体表示经典的for循环，其中Init、Cond、Post和Body分别对应初始化语句、条件语句、迭代语句和循环体语句。`ast.ForStmt`结构体的定义如下：

```go
type ForStmt struct {
	For  token.Pos // position of "for" keyword
	Init Stmt      // initialization statement; or nil
	Cond Expr      // condition; or nil
	Post Stmt      // post iteration statement; or nil
	Body *BlockStmt
}
```

其中条件部分必须是表达式，初始化和迭代部分可以是普通的语句（普通语句是短声明和多赋值等，不能包含分支等复杂语句）。

在了解了经典风格的循环之后，我们再来看看最简单的`for range`循环：

```go
func main() {
	for range ch {}
}
```

我们省略来循环中的Key和Value部分。其语法树如下：

```
 0  *ast.BlockStmt {
 1  .  Lbrace: 29
 2  .  List: []ast.Stmt (len = 1) {
 3  .  .  0: *ast.RangeStmt {
 4  .  .  .  For: 32
 5  .  .  .  TokPos: 0
 6  .  .  .  Tok: ILLEGAL
 7  .  .  .  X: *ast.Ident {
 8  .  .  .  .  NamePos: 42
 9  .  .  .  .  Name: "ch"
10  .  .  .  }
11  .  .  .  Body: *ast.BlockStmt {
12  .  .  .  .  Lbrace: 45
13  .  .  .  .  Rbrace: 46
14  .  .  .  }
15  .  .  }
16  .  }
17  .  Rbrace: 48
18  }
```

`for range`循环的语法树由`ast.RangeStmt`结构表示，其完整定义如下：

```go
type RangeStmt struct {
	For        token.Pos   // position of "for" keyword
	Key, Value Expr        // Key, Value may be nil
	TokPos     token.Pos   // position of Tok; invalid if Key == nil
	Tok        token.Token // ILLEGAL if Key == nil, ASSIGN, DEFINE
	X          Expr        // value to range over
	Body       *BlockStmt
}
```

其中Key和Value对应循环时的迭代位置和值，X成员是生成要循环对象的表达式（可能是数组、切片、map和管道等），Body表示循环体语句块。另外，Tok成员可以区别Key和Value是多赋值语句还是短变量声明语句。

## 10.9 类型断言

和分支语句类似，类型识别也有两种：类型断言和类型switch。类型断言类似分支的if语句，通过多个if/else组合类型断言就可以模拟出类型switch。因此我们重点学习类型断言部分，下面是类型断言的语法规范：

```
PrimaryExpr     = PrimaryExpr TypeAssertion.
TypeAssertion   = "." "(" Type ")" .
```

类型断言是在一个表达式之后加点和小括弧定义，其中小括弧中的是期望查询的类型。从Go语言语义角度看，类型断言开始的表达式必须产生一个接口类型的值。不过在语法树阶段并不会做详细的语义检查。

下面的例子在main函数中定义一个最简单的类型断言：

```go
func main() {
	x.(int)
}
```

对x做类型断言，如果成功则返回x里面存储的int类型的值，如果失败则抛出异常。生成的语法树如下：

```
 0  *ast.BlockStmt {
 1  .  Lbrace: 29
 2  .  List: []ast.Stmt (len = 1) {
 3  .  .  0: *ast.ExprStmt {
 4  .  .  .  X: *ast.TypeAssertExpr {
 5  .  .  .  .  X: *ast.Ident {
 6  .  .  .  .  .  NamePos: 32
 7  .  .  .  .  .  Name: "x"
 8  .  .  .  .  }
 9  .  .  .  .  Lparen: 34
10  .  .  .  .  Type: *ast.Ident {
11  .  .  .  .  .  NamePos: 35
12  .  .  .  .  .  Name: "int"
13  .  .  .  .  }
14  .  .  .  .  Rparen: 38
15  .  .  .  }
16  .  .  }
17  .  }
18  .  Rbrace: 40
19  }
```

需要注意语法树的结构：首先是`ast.ExprStmt`结构体表示的表达式语句，其中的X成员才是对应类型断言表达式。类型断言由`ast.TypeAssertExpr`结构体表示，其定义如下：

```go
type TypeAssertExpr struct {
	X      Expr      // expression
	Lparen token.Pos // position of "("
	Type   Expr      // asserted type; nil means type switch X.(type)
	Rparen token.Pos // position of ")"
}
```

其中X成员是类型断言的主体表达式（产生一个接口值），Type成员是类型的表达式。如果Type为nil，则表示对应`x.(type)`形式的断言，这是类型switch中使用的形式。

## 10.10 go和defer语句

go和defer语句是Go语言中最有特色的语句，它们的语法结构也是非常相似的。下面是go和defer语句的语法规范：

```
GoStmt    = "go" Expression .
DeferStmt = "defer" Expression .
```

简而言之，就是在go和defer关键字后面跟一个表达式，不过这个表达式必须是函数或方法调用。go和defer语句在语法树中分别以`ast.GoStmt`和`ast.DeferStmt`结构定义：

```go
type GoStmt struct {
	Go   token.Pos // position of "go" keyword
	Call *CallExpr
}
type DeferStmt struct {
	Defer token.Pos // position of "defer" keyword
	Call  *CallExpr
}
```

其中都有一个Call成员表示函数或方法调用。下面以go语句为例：

```go
func main() {
	go hello("光谷码农")
}
```

其对应的语法树结果：

```
 0  *ast.BlockStmt {
 1  .  Lbrace: 29
 2  .  List: []ast.Stmt (len = 1) {
 3  .  .  0: *ast.GoStmt {
 4  .  .  .  Go: 32
 5  .  .  .  Call: *ast.CallExpr {
 6  .  .  .  .  Fun: *ast.Ident {
 7  .  .  .  .  .  NamePos: 35
 8  .  .  .  .  .  Name: "hello"
 9  .  .  .  .  }
10  .  .  .  .  Lparen: 40
11  .  .  .  .  Args: []ast.Expr (len = 1) {
12  .  .  .  .  .  0: *ast.BasicLit {
13  .  .  .  .  .  .  ValuePos: 41
14  .  .  .  .  .  .  Kind: STRING
15  .  .  .  .  .  .  Value: "\"光谷码农\""
16  .  .  .  .  .  }
17  .  .  .  .  }
18  .  .  .  .  Ellipsis: 0
19  .  .  .  .  Rparen: 55
20  .  .  .  }
21  .  .  }
22  .  }
23  .  Rbrace: 57
24  }
```

除了`ast.GoStmt`结构体，Call成员部分和表达式中函数调用的语法树结构完全一样。

## 10.11 总结

数据结构是程序状态的载体，语句是程序算法的灵魂。在了解了语句的语法树之后，我们就可以基于语法树对代码做很多事情，比如特殊模式的BUG检查、生成文档或特定平台的可执行代码等，甚至我们可以基于语法树解释执行Go语言程序。

