# 第8章 复合面值

在基础面值章节，我们已经见过整数、浮点数、复数、符文和字符串等一些简单的面值。除了基础面值之外，还有结构体面值、map面值和函数面值等。本节讨论复合面值的语法树表示。

## 8.1 面值的语法

在Go语言规范文档中，完整的面值语法由Literal定义，具体如下：

```
Literal       = BasicLit | CompositeLit | FunctionLit .

BasicLit      = int_lit | float_lit | imaginary_lit | rune_lit | string_lit .

CompositeLit  = LiteralType LiteralValue .
LiteralType   = StructType | ArrayType | "[" "..." "]" ElementType |
                SliceType | MapType | TypeName .
LiteralValue  = "{" [ ElementList [ "," ] ] "}" .
ElementList   = KeyedElement { "," KeyedElement } .
KeyedElement  = [ Key ":" ] Element .
Key           = FieldName | Expression | LiteralValue .
FieldName     = identifier .
Element       = Expression | LiteralValue .
```

其中BasicLit是基础面值，CompositeLit是复合面值，FunctionLit是函数面值。其中复合类型和函数类型均已经讨论过，而其面值正是在复合类型和函数类型基础之上扩展而来。

## 8.2 函数面值

虽然函数面值不在复合类型面值之类，但是函数面值和函数声明非常相似，因此我们先看下函数面值。函数面值的语法如下：

```
FunctionLit   = "func" Signature FunctionBody .
```

函数面值由FunctionLit定义，同样是由func关键字开始，后面是函数签名（输入参数和返回值）和函数体。函数面值和函数声明的最大差别是没有函数名字。

我们从最简单的函数面值开始：

```go
func(){}
```

该函数面值没有输入参数和返回值，同时函数体也没有任何语句，而且没有涉及上下文的变量引用，可以说是最简单的函数面值。因为面值也是一种表达式，因此可以用表达式的方式解析其语法树：

```go
func main() {
	expr, _ := parser.ParseExpr(`func(){}`)
	ast.Print(nil, expr)
}
```

输出的语法树结构如下：

```
 0  *ast.FuncLit {
 1  .  Type: *ast.FuncType {
 2  .  .  Func: 1
 3  .  .  Params: *ast.FieldList {
 4  .  .  .  Opening: 5
 5  .  .  .  Closing: 6
 6  .  .  }
 7  .  }
 8  .  Body: *ast.BlockStmt {
 9  .  .  Lbrace: 7
10  .  .  Rbrace: 8
11  .  }
12  }
```

函数面值的语法树由`*ast.FuncLit`结构体表示，其中再由Type成员表示类型，Body成员表示函数体语句。函数的类型和函数体分别由`ast.FuncType`和`ast.BlockStmt`结构体表示，它们和函数声明中的表示形式是一致的。

我们可以对比下`*ast.FuncLit`和`ast.FuncDecl`结构体的差异：

```go
type FuncLit struct {
	Type *FuncType  // function type
	Body *BlockStmt // function body
}
type FuncDecl struct {
	Doc  *CommentGroup // associated documentation; or nil
	Recv *FieldList    // receiver (methods); or nil (functions)
	Name *Ident        // function/method name
	Type *FuncType     // function signature: parameters, results, and position of "func" keyword
	Body *BlockStmt    // function body; or nil for external (non-Go) function
}
```

对比可以发现表示函数类型的Type成员和表示函数体语句的Body成员类型都是一样的，但是FuncDecl函数声明比FuncLit函数面值多了函数名字和接收者参数列表等信息。因此如果理解了函数声明的完整结构，就可以用相似的方式处理函数类型和函数语句。

需要注意的是函数有面值，但是接口没有面值。因为接口是在运行时表示其它满足接口的对象，我们无法直接构造接口面值。在需要通过面值构造接口变量的地方，一般可以通过结构体等其它类型构造的面值赋值给接口的方式实现。

## 8.3 复合类型面值语法

复合类型面值语法由类型和值组成，其语法规范如下：

```
CompositeLit  = LiteralType LiteralValue .
LiteralType   = StructType | ArrayType | "[" "..." "]" ElementType |
                SliceType | MapType | TypeName .
LiteralValue  = "{" [ ElementList [ "," ] ] "}" .
ElementList   = KeyedElement { "," KeyedElement } .
KeyedElement  = [ Key ":" ] Element .
Key           = FieldName | Expression | LiteralValue .
FieldName     = identifier .
Element       = Expression | LiteralValue .
```

复合类型主要包含结构体、数组、切片和map类型，此外还有基于这些类型命名的类型。结构体、数组、切片和map类型的面值在LiteralValue定义，对应一个大括号构成的语法结构。在LiteralValue描述的复合类型面值部分的大括号中，由一个可选的Key和对应的值组成，其中值可以是基础面值、生成值的表达式或者是LiteralValue类型。

以下是结构体、数组、切片和map类型常见的面值语法：

```go
[1]int{1}
[...]int{100:1,200:2}
[]int{1,2,3}
[]int{100:1,200:2}
struct {X int}{1}
struct {X int}{X:1}
map[int]int{1:1, 2:2}
```

其中数组和切片各有两种面值语法：一种是顺序指定初始值的列表，另一种是通过下标指定某个特定位置的初始值（两种格式可以混合使用）。结构体面值可以全部省略成员的名字，也可以指定成员的名字。map的面值必须完整指定Key和对应的值。

复合类型面值内元素的初始值又可能是复合面值，因此这也是一种递归语法结构。下面是一个嵌套复合类型的例子：

```go
[]image.Point{
	image.Point{X: 1, Y: 2},
	{X: 3, Y: 4},
	5: {6, 7},
}
```

最外层是`image.Point`类型的切片，第一个元素通过完整的面值语法`image.Point{X: 1, Y: 2}`指定初始值，第二个元素通过简化的`{X: 3, Y: 4}`语法初始化，第三四五个元素空缺为零值，最后一个元素通过下表语法结合`{6, 7}`指定。需要注意的是，虽然面值初始化有多种形式，但是在语法树中都是相似的，因此我们需要透过面值的表象理解其语法树的本质。

复合型面值的语法树通过`ast.CompositeLit`表示：

```go
type CompositeLit struct {
	Type       Expr      // literal type; or nil
	Lbrace     token.Pos // position of "{"
	Elts       []Expr    // list of composite elements; or nil
	Rbrace     token.Pos // position of "}"
	Incomplete bool      // true if (source) expressions are missing in the Elts list
}
```

其中`ast.CompositeLit.Type`对应复合类型的表达式，然后`ast.CompositeLit.Elts`是复合类型初始元素列表。初始元素列表可以是普通的值，也可能是Key-Value下标和值对，而且初始值还可能是其它的复合面值。

## 8.4 数组和切片面值

数组或切片面值是在数组类型后面的大括弧中包含数组的元素列表：

```go
[...]int{1,2:3}
```

因为数组面值也是一种表达式，因此可以直接通过解析表达式的方式生成语法树：

```go
func main() {
	expr, _ := parser.ParseExpr(`[...]int{1,2:3}`)
	ast.Print(nil, expr)
}
```

输出的语法树如下：

```
 0  *ast.CompositeLit {
 1  .  Type: *ast.ArrayType {
 2  .  .  Lbrack: 1
 3  .  .  Len: *ast.Ellipsis {
 4  .  .  .  Ellipsis: 2
 5  .  .  }
 6  .  .  Elt: *ast.Ident {
 7  .  .  .  NamePos: 6
 8  .  .  .  Name: "int"
 9  .  .  .  Obj: *ast.Object {
10  .  .  .  .  Kind: bad
11  .  .  .  .  Name: ""
12  .  .  .  }
13  .  .  }
14  .  }
15  .  Lbrace: 9
16  .  Elts: []ast.Expr (len = 2) {
17  .  .  0: *ast.BasicLit {
18  .  .  .  ValuePos: 10
19  .  .  .  Kind: INT
20  .  .  .  Value: "1"
21  .  .  }
22  .  .  1: *ast.KeyValueExpr {
23  .  .  .  Key: *ast.BasicLit {
24  .  .  .  .  ValuePos: 12
25  .  .  .  .  Kind: INT
26  .  .  .  .  Value: "2"
27  .  .  .  }
28  .  .  .  Colon: 13
29  .  .  .  Value: *ast.BasicLit {
30  .  .  .  .  ValuePos: 14
31  .  .  .  .  Kind: INT
32  .  .  .  .  Value: "3"
33  .  .  .  }
34  .  .  }
35  .  }
36  .  Rbrace: 15
37  .  Incomplete: false
38  }
```

复合面值语法树由`ast.CompositeLit`结构体表示，其中`ast.CompositeLit.Type`成员为`ast.ArrayType`表示这是数组或切片类型（如果没有长度信息则为切片类型，否则就是数组），而`ast.CompositeLit`Elts成员则是元素的值。初始元素是一个`[]ast.Expr`类型的切片，每个元素依然是一个表达式。数组的第一个元素是`ast.BasicLit`类型，表示这是一个基础面值类型。数组的第二个元素是`ast.KeyValueExpr`方式指定的，其中Key对应的数组下标是2，Value对应的值为3。

数组和切片语法的最大差别是数组有长度信息。在这个例子中数组是通过`...`省略号表达式自动计算数组的长度，在语法树中对应的是`ast.Ellipsis`表达式类型。如果`ast.ArrayType`结构体中的Len成员是空指针，则表示这是一个切片类型，否则对应可以生成数组长度的表达式。

## 8.5 结构体面值

结构体面值和数组面值类似，是在结构体类型后面的大括弧中包含结构体成员的初始值。下面是结构体例子：

```go
struct{X int}{X:1}
```

可以通过以下代码解析其语法树：

```go
func main() {
	expr, _ := parser.ParseExpr(`struct{X int}{X:1}`)
	ast.Print(nil, expr)
}
```

输出的语法树结果如下：

```
 0 *ast.CompositeLit {
 1  .  Type: *ast.StructType {...}
32  .  Lbrace: 14
33  .  Elts: []ast.Expr (len = 1) {
34  .  .  0: *ast.KeyValueExpr {
35  .  .  .  Key: *ast.Ident {
36  .  .  .  .  NamePos: 15
37  .  .  .  .  Name: "X"
38  .  .  .  }
39  .  .  .  Colon: 16
40  .  .  .  Value: *ast.BasicLit {
41  .  .  .  .  ValuePos: 17
42  .  .  .  .  Kind: INT
43  .  .  .  .  Value: "1"
44  .  .  .  }
45  .  .  }
46  .  }
47  .  Rbrace: 18
48  .  Incomplete: false
49  }
```

结构体面值依然是通过`ast.CompositeLit`结构体描述。结构体中成员的初始化通过`ast.KeyValueExpr`结构体初始化，Key部分为X表示成员名字，Value部分为X成员的初始值。

当然，结构体的初始化也可以不声明成员的名字：

```go
func main() {
	expr, _ := parser.ParseExpr(`struct{X int}{1}`)
	ast.Print(nil, expr)
}
```

现在的初始化方式生成的语法树变得更简单：

```
33  .  Elts: []ast.Expr (len = 1) {
34  .  .  0: *ast.BasicLit {
35  .  .  .  ValuePos: 15
36  .  .  .  Kind: INT
37  .  .  .  Value: "1"
38  .  .  }
39  .  }
```

只有一个元素是通过`ast.BasicLit`对应的基础面值表示，对应结构体的第一个成员。

## 8.6 map面值

map面值的表示方式和按成员名字初始化结构体的面值语法树基本一样：

```go
func main() {
	expr, _ := parser.ParseExpr(`map[int]int{1:2}`)
	ast.Print(nil, expr)
}
```

输出语法树中的初始化值列表部分（`ast.CompositeLit.Elts`）：

```
18  .  Elts: []ast.Expr (len = 1) {
19  .  .  0: *ast.KeyValueExpr {
20  .  .  .  Key: *ast.BasicLit {
21  .  .  .  .  ValuePos: 13
22  .  .  .  .  Kind: INT
23  .  .  .  .  Value: "1"
24  .  .  .  }
25  .  .  .  Colon: 14
26  .  .  .  Value: *ast.BasicLit {
27  .  .  .  .  ValuePos: 15
28  .  .  .  .  Kind: INT
29  .  .  .  .  Value: "2"
30  .  .  .  }
31  .  .  }
32  .  }
```

map的初始值只能通过`ast.KeyValueExpr`对应的键值对表示，因为缺少了key无法定位值对应的下标位置。

## 8.7 小结

非基础面值包含函数面值和复合类型面值。函数面值和顶级函数声明有着相似的语法，只是没有函数名部分，表示语法树的结构体都是一致的。而数组、切片、结构体和map等复合类型的初始化语法也是高度一致的，其中只有map必须通过键值对初始化，其它的复合类型同时支持键值对和顺序值列表初始化，因此初始化值对应的语法树有`ast.KeyValueExpr`和普通的`ast.Expr`类型。至此，和数据相关的类型和值已经全部讨论，在此基础之上构建数据的反射实现，也可以基于数据结构构建算法。类型和值是最基础的部分，因为它们是构成变量的基础。

