## 1.2 复合面值

Go语言中除了基础面值之外还有复合面值。其结构体、数组、切片、map和自定义类型都属于复合面值。本节简单讨论复合面值相关的内容。

### 1.2.1 复合面值语法规范

根据Go语言规范文档，复合面值的语法定义如下：

```
CompositeLit  = LiteralType LiteralValue .
LiteralType   = StructType | ArrayType | "[" "..." "]" ElementType |
                SliceType | MapType | TypeName .
```

简单来说，复合面值由一个复合类型开始，紧接着是复合面值。其中复合类型可以是结构体、数组、切片、map或自定义类型。

### 1.2.2 复合面值的语法树表示

在`go/ast`包中`CompositeLit`表示复合面值，其结构定义如下：

```go
type CompositeLit struct {
    Type       Expr      // literal type; or nil
    Lbrace     token.Pos // position of "{"
    Elts       []Expr    // list of composite elements; or nil
    Rbrace     token.Pos // position of "}"
    Incomplete bool      // true if (source) expressions are missing in the Elts list
}
```

结构比基础面值突然复杂了很多。根据复合面值的语法定义，其结构是"LiteralType LiteralValue"形式，也就是一个类型和值的组合体。因此可以猜测`CompositeLit`结构体中的`Type`成员对应复合面值的类型，然后的`Elts`切片对应的是每个复合面值成员的描述。其它的成员，Lbrace和Rbrace表示复合面值中左右花括弧的位置。Incomplete成员的含义暂时忽略。

需要注意的是，其中的`Expr`是表示表达式的抽象接口，可以是面值也可以是其它可以组成表达式的结构。因此`CompositeLit`是一种递归结构，是一种简单的语法树形式。

## 1.2.3 数组面值

最简单的数组是零长度的数组：

```go
func main() {
	expr, _ := parser.ParseExpr(`[1]int{1}`)
	ast.Print(nil, expr)
}
```

```
 0  *ast.CompositeLit {
 1  .  Type: *ast.ArrayType {
 2  .  .  Lbrack: 1
 3  .  .  Len: *ast.BasicLit {
 4  .  .  .  ValuePos: 2
 5  .  .  .  Kind: INT
 6  .  .  .  Value: "0"
 7  .  .  }
 8  .  .  Elt: *ast.Ident {
 9  .  .  .  NamePos: 4
10  .  .  .  Name: "int"
11  .  .  .  Obj: *ast.Object {
12  .  .  .  .  Kind: bad
13  .  .  .  .  Name: ""
14  .  .  .  }
15  .  .  }
16  .  }
17  .  Lbrace: 7
18  .  Rbrace: 8
19  .  Incomplete: false
20  }
```

```go
type ArrayType struct {
    Lbrack token.Pos // position of "["
    Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
    Elt    Expr      // element type
}
```
## 1.2.3 数组面值