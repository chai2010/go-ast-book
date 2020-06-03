# 第9章 复合表达式

在基础面值和基础表达式章节，我们已经见过一些简单的表达式。本章我们将继续讨论复杂表达式，包含基于复杂面值和点号选择运算、索引运算、切片运算和函数调用等相互组合而成的表达式。

## 9.1 表达式语法

简单来说，表达式是指所有可以产生一个值的语句的集合。表达式的语法由PrimaryExpr定义：

```
PrimaryExpr = Operand
            | Conversion
            | MethodExpr
            | PrimaryExpr Selector
            | PrimaryExpr Index
            | PrimaryExpr Slice
            | PrimaryExpr TypeAssertion
            | PrimaryExpr Arguments
            .

Selector       = "." identifier .
Index          = "[" Expression "]" .
Slice          = "[" [ Expression ] ":" [ Expression ] "]" 
               | "[" [ Expression ] ":" Expression ":" Expression "]" .

TypeAssertion  = "." "(" Type ")" .
Arguments      = "(" [ ( ExpressionList | Type [ "," ExpressionList ] ) [ "..." ] [ "," ] ] ")" .
```

其中Operand是由一元或二元算术运算符组成的算术运算表达式。Conversion是强制类型转换，形式和函数调用有一定的相似性。MethodExpr是方法表达式。然后是点选择运算、索引运算、切片运算、类型断言和函数调用参数等高阶运算符。

## 9.2 转型和函数调用

二元算术运算符我们已经讲过，因此我们从转型操作和函数调用开始。下面是转型操作和函数参数的语法规范：

```
Conversion = Type "(" Expression [ "," ] ")" .
Arguments  = "(" [ ( ExpressionList | Type [ "," ExpressionList ] ) [ "..." ] [ "," ] ] ")" .
```

需要注意的是转型操作和只有一个参数的函数调用非常相似，但是转型操作是以一个类型开始（函数调用是以一个函数开始），然后小括号内是要转型的表达式。下面的例子是将x变量转型为int类型：

```go
int(x)
```

如果int被重新定义为一个函数，那么转型操作就会变成函数调用。我们先看看转型操作的语法树是如何表示的：

```go
func main() {
	expr, _ := parser.ParseExpr(`int(x)`)
	ast.Print(nil, expr)
}
```

输出的语法树如下：

```
 0  *ast.CallExpr {
 1  .  Fun: *ast.Ident {
 2  .  .  NamePos: 1
 3  .  .  Name: "int"
 4  .  .  Obj: *ast.Object {
 5  .  .  .  Kind: bad
 6  .  .  .  Name: ""
 7  .  .  }
 8  .  }
 9  .  Lparen: 4
10  .  Args: []ast.Expr (len = 1) {
11  .  .  0: *ast.Ident {
12  .  .  .  NamePos: 5
13  .  .  .  Name: "x"
14  .  .  .  Obj: *(obj @ 4)
15  .  .  }
16  .  }
17  .  Ellipsis: 0
18  .  Rparen: 6
19  }
```

转型操作居然是用`ast.CallExpr`表示，这说明在语法树中转型和函数调用的结构是完全一样的。这是因为在语法树解析阶段，解析器并不知道`int(x)`中的`int`是一个类型还是一个函数，因此也无法知晓这是一个转型操作还是一个函数调用。

`ast.CallExpr`结构体定义如下：

```go
type CallExpr struct {
	Fun      Expr      // function expression
	Lparen   token.Pos // position of "("
	Args     []Expr    // function arguments; or nil
	Ellipsis token.Pos // position of "..." (token.NoPos if there is no "...")
	Rparen   token.Pos // position of ")"
}
```

其中Fun如果是类型表达式，则表示这是一个转型操作。Fun之所以被定义为一个表达式，是因为Go语言中函数是第一类对象，可以像普通值一样被传递，通过表达式可以获取结构体、数组或map中保存的函数。而Args参数部分表示要转型的表达式或者是函数调用的参数列表。如果是函数调用，并且是可变参数函数调用，那么Ellipsis表示省略号位置（否则是一个无效的位置）。

## 9.3 点选择运算

点选择运算主要用于结构体选择其成员，或者是对象选择其方法。点选择运算语法如下：

```
PrimaryExpr = PrimaryExpr Selector .
Selector    = "." identifier .
```

如果有表达式`x`，则可以通过`x.y`访问其成员或方法函数。如果是`x`导入包，那么`x.y`将变成标识符含义。同样，在语法树解析阶段并无法区分一个选择表达式和导入包中的标识符。

下面是`x.y`解析的语法树结果：

```
 0  *ast.SelectorExpr {
 1  .  X: *ast.Ident {
 2  .  .  NamePos: 1
 3  .  .  Name: "x"
 4  .  .  Obj: *ast.Object {
 5  .  .  .  Kind: bad
 6  .  .  .  Name: ""
 7  .  .  }
 8  .  }
 9  .  Sel: *ast.Ident {
10  .  .  NamePos: 3
11  .  .  Name: "y"
12  .  }
13  }
```

其中X成员表示主体、Sel是被选择的成员（也可能是其它包的标识符）。`ast.SelectorExpr`结构体定义如下：

```go
type SelectorExpr struct {
	X   Expr   // expression
	Sel *Ident // field selector
}
```

其中X被定义为ast.Expr表达式类型，Sel是一个普通的标识符。

## 9.4 索引运算

索引运算主要用于数组、切片或map选择元素，其语法规范如下：

```
PrimaryExpr = PrimaryExpr Index .
Index       = "[" Expression "]" .
```

索引运算通过在主体表达式后面的中括弧中包含索引表达式。同样在语法树解析阶段无法区别索引运算主体的具体类型。下面是`x[y]`索引运算的语法树解析结果：

```
 0  *ast.IndexExpr {
 1  .  X: *ast.Ident {
 2  .  .  NamePos: 1
 3  .  .  Name: "x"
 4  .  .  Obj: *ast.Object {
 5  .  .  .  Kind: bad
 6  .  .  .  Name: ""
 7  .  .  }
 8  .  }
 9  .  Lbrack: 2
10  .  Index: *ast.Ident {
11  .  .  NamePos: 3
12  .  .  Name: "y"
13  .  .  Obj: *(obj @ 4)
14  .  }
15  .  Rbrack: 4
16  }
```

其中X是主体表达式，其中的标识符是x。而Index为索引表达式，在这个例子中是y。`ast.IndexExpr`结构体定义如下：

```go
type IndexExpr struct {
	X      Expr      // expression
	Lbrack token.Pos // position of "["
	Index  Expr      // index expression
	Rbrack token.Pos // position of "]"
}
```

其中X和Index成员都是表达式，具体的语义需要根据上下文判断X表达式的类型才能决定Index索引表达式的类型。

## 9.5 切片运算

切片运算是在数组或切片基础上生成新的切片，其语法规范如下：

```
PrimaryExpr =  PrimaryExpr Slice
Slice       = "[" [ Expression ] ":" [ Expression ] "]" 
            | "[" [ Expression ] ":" Expression ":" Expression "]"
            .
```

切片运算也是在一个主体表达式之后的中括弧中表示，不过切片运算至少有一个冒号分隔符，或者是两个冒号分隔符。切片运算主要包含开始索引、结束索引和最大范围三个部分。下面是`x[1:2:3]`切片运算的语法树：

```
 0  *ast.SliceExpr {
 1  .  X: *ast.Ident {
 2  .  .  NamePos: 1
 3  .  .  Name: "x"
 4  .  .  Obj: *ast.Object {
 5  .  .  .  Kind: bad
 6  .  .  .  Name: ""
 7  .  .  }
 8  .  }
 9  .  Lbrack: 2
10  .  Low: *ast.BasicLit {
11  .  .  ValuePos: 3
12  .  .  Kind: INT
13  .  .  Value: "1"
14  .  }
15  .  High: *ast.BasicLit {
16  .  .  ValuePos: 5
17  .  .  Kind: INT
18  .  .  Value: "2"
19  .  }
20  .  Max: *ast.BasicLit {
21  .  .  ValuePos: 7
22  .  .  Kind: INT
23  .  .  Value: "3"
24  .  }
25  .  Slice3: true
26  .  Rbrack: 8
27  }
```

切片运算通过`ast.SliceExpr`结构体表示，其中X、Low、High、Max分别表示切片运算的主体、开始索引、结束索引和最大范围。`ast.SliceExpr`结构体定义如下：

```go
type SliceExpr struct {
	X      Expr      // expression
	Lbrack token.Pos // position of "["
	Low    Expr      // begin of slice range; or nil
	High   Expr      // end of slice range; or nil
	Max    Expr      // maximum capacity of slice; or nil
	Slice3 bool      // true if 3-index slice (2 colons present)
	Rbrack token.Pos // position of "]"
}
```

其中X、Low、High、Max是我们已经熟悉的成员，都是表达式类型。另外Slice3标注是否为三索引的切片语法（不过这个字段对语义没有影响，因为可以从Max程序推导出最大的容量信息）。

## 9.6 类型断言

类型断言是判断一个接口对象是否满足另一个接口、或者接口持有的对象是否是一个确定的非接口类型。类型断言的语法规范如下：

```
PrimaryExpr    = PrimaryExpr TypeAssertion .
TypeAssertion  = "." "(" Type ")" .
```

在主体表达式之后通过点选择一个类型，类型放在小括弧中间。比如`x.(y)`就是将x接口断言为y接口或y类型，下面是它们的语法树：

```
 0  *ast.TypeAssertExpr {
 1  .  X: *ast.Ident {
 2  .  .  NamePos: 1
 3  .  .  Name: "x"
 4  .  .  Obj: *ast.Object {
 5  .  .  .  Kind: bad
 6  .  .  .  Name: ""
 7  .  .  }
 8  .  }
 9  .  Lparen: 3
10  .  Type: *ast.Ident {
11  .  .  NamePos: 4
12  .  .  Name: "y"
13  .  .  Obj: *(obj @ 4)
14  .  }
15  .  Rparen: 5
16  }
```

断言运算由`ast.TypeAssertExpr`表示，其中X是接口表达式，Type是要断言的类型表达式。`ast.TypeAssertExpr`结构体的定义如下：

```go
type TypeAssertExpr struct {
	X      Expr      // expression
	Lparen token.Pos // position of "("
	Type   Expr      // asserted type; nil means type switch X.(type)
	Rparen token.Pos // position of ")"
}
```

需要注意的是`x.(type)`也是一种特殊的类型断言，这时候`ast.TypeAssertExpr.Type`成员值为nil，对应的是类型switch语句结构。

## 9.7 小结

此处我们已经学习了基于各种基础类型、复合类型的各种表达式基础构件，通过组合这些运算就能产生各种复杂的表达式。最终将表达式的结果通过和赋值语句或控制流语句相结合，就可以改成程序的环境状态。而编程的本质就是通过语句改变成员的状态，然后在根据不同的状态选择执行不同的语句。

