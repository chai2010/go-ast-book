# 第7章 复合类型

这里讨论的复合类型是指无法用一个标识符表示的类型，它们包含其它包中的基础类型（需要通过点号选择操作符）、指针类型、
数组类型、切片类型、结构体类型、map类型、管道类型、函数类型和接口类型，以及它们之间再次组合产生的更复杂的类型。

## 7.1 类型的语法

在基础类型声明章节我们已经简要学习过类型的声明语法规范，不过当时只讨论了基于标识符的简单声明。本章我们将继续探讨复合类型声明的语法和语法树的表示。以下是更为完整的类型声明的语法规范：

```bnf
TypeDecl  = "type" ( TypeSpec | "(" { TypeSpec ";" } ")" ) .
TypeSpec  = AliasDecl | TypeDef .

AliasDecl = identifier "=" Type .
TypeDef   = identifier Type .

Type      = TypeName | TypeLit | "(" Type ")" .
TypeName  = identifier | PackageName "." identifier .
TypeLit   = PointerType | ArrayType | SliceType
          | StructType | MapType | ChannelType
          | FunctionType | InterfaceType
          .
```

增加的部分主要在TypeName和TypeLit。TypeName不仅仅可以从当前空间的标识符定义新类型，还支持从其它包导入的标识符定义类型。而TypeLit表示类型面值，比如基于已有类型的指针，或者是匿名的结构体都属于类型的面值。

如前文所描述，类型定义由`*ast.TypeSpec`结构体表示，复合类型也是如此。下面再来回顾下该结构体的定义：

```go
type TypeSpec struct {
	Doc     *CommentGroup // associated documentation; or nil
	Name    *Ident        // type name
	Assign  token.Pos     // position of '=', if any; added in Go 1.9
	Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of th *XxxTypes
	Comment *CommentGroup // line comments; or nil
}
```

其中Name成员表示给类型命名，Type通过特殊的类型表达式表示类型的定义，此外如果Assign被设置则表示声明的是类型的别名。

## 7.2 基础类型

基础类型是最简单的类型，就是基于已有的命名类型再次定义新类型，或者是为已有类型定义新的别名。该类型的语法规则比较简单，主要限制在Type部分：

```
TypeDecl  = "type" ( TypeSpec | "(" { TypeSpec ";" } ")" ) .
TypeSpec  = AliasDecl | TypeDef .

AliasDecl = identifier "=" Type .
TypeDef   = identifier Type .

Type      = identifier | PackageName "." identifier .
```

Type表示已有的命名类型，可以是当前包的类型，也是可以其它包的类型。下面是这些类型的例子：

```go
type Int1 int
type Int2 pkg.Int
```

其中第一个Int1类型是基于当前名字空间可以直接访问的int类型，而第二个Int2类型是基于导入的pkg包中的Int类型。我们可以用以下代码解析上面的类型声明：

```go
func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		ast.Print(nil, decl.(*ast.GenDecl).Specs[0])
	}
}

const src = `package foo
type Int1 int
type Int2 pkg.int
`
```

第一个类型的输出结果如下：

```
 0  *ast.TypeSpec {
 1  .  Name: *ast.Ident {
 2  .  .  NamePos: 18
 3  .  .  Name: "Int1"
 4  .  .  Obj: *ast.Object {
 5  .  .  .  Kind: type
 6  .  .  .  Name: "Int1"
 7  .  .  .  Decl: *(obj @ 0)
 8  .  .  }
 9  .  }
10  .  Assign: 0
11  .  Type: *ast.Ident {
12  .  .  NamePos: 23
13  .  .  Name: "int"
14  .  }
15  }
```

第二个类型的输出结果如下：

```
 0  *ast.TypeSpec {
 1  .  Name: *ast.Ident {
 2  .  .  NamePos: 32
 3  .  .  Name: "Int2"
 4  .  .  Obj: *ast.Object {
 5  .  .  .  Kind: type
 6  .  .  .  Name: "Int2"
 7  .  .  .  Decl: *(obj @ 0)
 8  .  .  }
 9  .  }
10  .  Assign: 0
11  .  Type: *ast.SelectorExpr {
12  .  .  X: *ast.Ident {
13  .  .  .  NamePos: 37
14  .  .  .  Name: "pkg"
15  .  .  }
16  .  .  Sel: *ast.Ident {
17  .  .  .  NamePos: 41
18  .  .  .  Name: "int"
19  .  .  }
20  .  }
21  }
```

对比两个结果可以发现，Int1的Type定义对应的是`*ast.Ident`表示一个标识符，而Int2的Type定义对应的时候`*ast.SelectorExpr`表示是其它包的命名类型。`*ast.SelectorExp`结构体定义如下：

```go
type SelectorExpr struct {
	X   Expr   // expression
	Sel *Ident // field selector
}
```

结构体X成员被定义为Expr接口类型，不过根据当前的语法必须是一个标识符类型（之所以被定义为表达式接口，是因为其它的表达式会复用这个结构体）。Sel成员被定义为标识符类型，表示被选择的标识符名字。

## 7.3 指针类型

指针是操作底层类型时最强有力的武器，只要有指针就可以操作内存上的所有数据。最简单的是一级指针，然后再扩展出二级和更多级指针。以下是Go语言指针类型的语法规范：

```
PointerType = "*" BaseType .
BaseType    = Type .

Type        = TypeName | TypeLit | "(" Type ")" .
...
```

指针类型以星号`*`开头，后面是BaseType定义的类型表达式。从语法规范角度看，Go语言没有单独定义多级指针，只有一种指向BaseType类型的一级指针。但是PointerType又可以作为TypeLit类型面值被重新用作BaseType，这就产生了多级指针的语法。

下面是一级指针语法树解析的例子：

```go
func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		ast.Print(nil, decl.(*ast.GenDecl).Specs[0])
	}
}

const src = `package foo
type IntPtr *int
`
```

解析的结果如下：

```
 0  *ast.TypeSpec {
 1  .  Name: *ast.Ident {
 2  .  .  NamePos: 18
 3  .  .  Name: "IntPtr"
 4  .  .  Obj: *ast.Object {
 5  .  .  .  Kind: type
 6  .  .  .  Name: "IntPtr"
 7  .  .  .  Decl: *(obj @ 0)
 8  .  .  }
 9  .  }
10  .  Assign: 0
11  .  Type: *ast.StarExpr {
12  .  .  Star: 25
13  .  .  X: *ast.Ident {
14  .  .  .  NamePos: 26
15  .  .  .  Name: "int"
16  .  .  }
17  .  }
18  }
```

新类型的名字依然是普通的`*ast.Ident`标识符类型，其值是新类型的名字“IntPtr”。而`ast.TypeSpec.Type`成员则是新的`*ast.StarExpr`类型，其结构体定义如下：

```go
type StarExpr struct {
	Star token.Pos // position of "*"
	X    Expr      // operand
}
```

指针指向的X类型是一个递归定义的类型表达式。在这个例子中X就是一个`*ast.Ident`标识符类型表示的int，因此IntPtr类型是一个指向int类型的指针类型。

指针是一种天然递归定义的类型。我们可以再定义一个指向IntPtr类型的指针，它又是一个指向int类型的二级指针。但是在语法树表示时，指向IntPtr类型的一级指针和指向int类型的二级指针结构是不一样的，因为语法树解析器会将IntPtr和int都作为普通类型同等对待（语法树解析器只知道这是指向IntPtr类型的一级指针，而不知道它也是指向int类型的二级指针）。

下面的例子依然是在int类型基础之上定义二级指针：

```go
type IntPtrPtr **int
```

解析后语法树发生的最大的差异在类型定义部分：

```
11  .  Type: *ast.StarExpr {
12  .  .  Star: 28
13  .  .  X: *ast.StarExpr {
14  .  .  .  Star: 29
15  .  .  .  X: *ast.Ident {
16  .  .  .  .  NamePos: 30
17  .  .  .  .  Name: "int"
18  .  .  .  }
19  .  .  }
20  .  }
```

现在`ast.StarExpr.X`不再是一个`*ast.Ident`标识符类型，而是变成了`*ast.StarExpr`类型的指针类型。对于多级指针的`*ast.StarExpr`类型很像一个单向的链表，其中X成员指向的是减一级指针的`*ast.StarExpr`结点，链表的尾结点是一个`*ast.Ident`标识符类型。

## 7.4 数组类型

在传统的C/C++语言中，数组是和指针近似等同的类型，特别在传递参数时只传递数组的首地址。Go语言的数组类型是一种值类型，每次传递数组参数或者赋值都是生成数组的拷贝。但是从数组的语法定义角度看，它和指针类型也是非常相似的。以下是数组类型的语法规范：

```
ArrayType   = "[" ArrayLength "]" ElementType .
ArrayLength = Expression .
ElementType = Type .
```

Go语言中数组的长度也是其类型的组成部分，数组长度是由一个表达式定义（在语义层面这个表达式必须是常量）。然后是数组元素的类型。如果抛开数组的长度部分的差异，数组类型和指针类型是非常相似的语法结构。数组元素部分的ElementType类型也可以是数组，这又构成了多级数组的语法规范。

下面是简单的一维整型数组的例子：

```go
type IntArray [1]int
```

解析结果如下：

```
 0  *ast.TypeSpec {
 1  .  Name: *ast.Ident {
 2  .  .  NamePos: 18
 3  .  .  Name: "IntArray"
 4  .  .  Obj: *ast.Object {
 5  .  .  .  Kind: type
 6  .  .  .  Name: "IntArray"
 7  .  .  .  Decl: *(obj @ 0)
 8  .  .  }
 9  .  }
10  .  Assign: 0
11  .  Type: *ast.ArrayType {
12  .  .  Lbrack: 27
13  .  .  Len: *ast.BasicLit {
14  .  .  .  ValuePos: 28
15  .  .  .  Kind: INT
16  .  .  .  Value: "1"
17  .  .  }
18  .  .  Elt: *ast.Ident {
19  .  .  .  NamePos: 30
20  .  .  .  Name: "int"
21  .  .  }
22  .  }
23  }
```

数组的类型主要由`*ast.ArrayType`类型定义。数组的长度是一个`*ast.BasicLit`类型的表达式，也就是长度为1的数组。数组元素的长度是`*ast.Ident`类型的标识符表示，数组的元素对应int类型。

完整的`*ast.ArrayType`结构体如下：

```go
type ArrayType struct {
	Lbrack token.Pos // position of "["
	Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
	Elt    Expr      // element type
}
```

其中`ast.ArrayType.Len`成员是一个表示数组长度的表达式，该表达式必须可以产生常量的整数结果（也可以是三个点省略号表示从元素个数提取）。数组的元素由`ast.ArrayType.Elt`定义，其值对应一个类型表达式。和指针类型一样，数组类型也是可以递归定义的，数组的元素类型可以数数组、指针等其它任何类型。

同样，我们可以定义一个二维数组：

```go
type IntArrayArray [1][2]int
```

解析结果如下：

```
11  .  Type: *ast.ArrayType {
12  .  .  Lbrack: 32
13  .  .  Len: *ast.BasicLit {
14  .  .  .  ValuePos: 33
15  .  .  .  Kind: INT
16  .  .  .  Value: "1"
17  .  .  }
18  .  .  Elt: *ast.ArrayType {
19  .  .  .  Lbrack: 35
20  .  .  .  Len: *ast.BasicLit {
21  .  .  .  .  ValuePos: 36
22  .  .  .  .  Kind: INT
23  .  .  .  .  Value: "2"
24  .  .  .  }
25  .  .  .  Elt: *ast.Ident {
26  .  .  .  .  NamePos: 38
27  .  .  .  .  Name: "int"
28  .  .  .  }
29  .  .  }
30  .  }
```

同样，数组元素的类型也变成了嵌套的数组类型。N维的数组类型的语法树也类似一个单向链表结构，后`N-1`维的数组的元素也是`*ast.ArrayType`类型，最后的尾结点对应一个`*ast.Ident`标识符（也可以是其它面值类型）。

## 7.5 切片类型

Go语言中切片是简化的数组，切片中引入了诸多数组不支持的语法。不过对于切片类型的定义来说，切片和数组的差异就是省略了数组的长度而已。切片类型声明的语法规则如下：

```
SliceType   = "[" "]" ElementType .
ElementType = Type .
```

下面例子是定义一个int切片：

```go
type IntSlice []int
```

对其解析语法树的输出如下：

```
 0  *ast.TypeSpec {
 1  .  Name: *ast.Ident {
 2  .  .  NamePos: 18
 3  .  .  Name: "IntSlice"
 4  .  .  Obj: *ast.Object {
 5  .  .  .  Kind: type
 6  .  .  .  Name: "IntSlice"
 7  .  .  .  Decl: *(obj @ 0)
 8  .  .  }
 9  .  }
10  .  Assign: 0
11  .  Type: *ast.ArrayType {
12  .  .  Lbrack: 27
13  .  .  Elt: *ast.Ident {
14  .  .  .  NamePos: 29
15  .  .  .  Name: "int"
16  .  .  }
17  .  }
18  } 
```

切片和数组一样，也是通过`*ast.ArrayType`结构表示切片，不过Len长度成员为nil类型（切片必须是nil，如果是0则表示是数组类型）。

## 7.6 结构体类型

结构体类型是数组类型的再次演进：数组是类型相同的元素的组合，并通过下标索引定位元素；而结构体类型是不同类型元素的组合，可以通过名字来定位元素。结构体类型这种可以组合异构元素类型的抽象能力极大地改进了数据结构编程的体验。结构体类型的语法规范定义如下：

```
StructType     = "struct" "{" { FieldDecl ";" } "}" .
FieldDecl      = (IdentifierList Type | EmbeddedField) [ Tag ] .
EmbeddedField  = [ "*" ] TypeName .
Tag            = string_lit .

IdentifierList = identifier { "," identifier } .
TypeName       = identifier | PackageName "." identifier .
```

结构体通过struct关键字开始定义，然后在大括弧中包含成员的定义。每一个FieldDecl表示一组有着相同类型和Tag字符串的标识符名字，或者是嵌入的匿名类型或类型指针。以下是结构体的例子：

```go
type MyStruct struct {
	a, b int "int value"
	string
}
```

其中a和b成员不仅仅有着相同的int类型，同时还有着相同的Tag字符串，最后的成员是嵌入一个匿名的字符串。

对其解析语法树的输出如下（为了简化省略了一些无关的信息）：

```
11  .  Type: *ast.StructType {
12  .  .  Struct: 27
13  .  .  Fields: *ast.FieldList {
14  .  .  .  Opening: 34
15  .  .  .  List: []*ast.Field (len = 2) {
16  .  .  .  .  0: *ast.Field {
17  .  .  .  .  .  Names: []*ast.Ident (len = 2) {
18  .  .  .  .  .  .  0: *ast.Ident {
19  .  .  .  .  .  .  .  NamePos: 37
20  .  .  .  .  .  .  .  Name: "a"
21  .  .  .  .  .  .  .  Obj: *ast.Object {...}
26  .  .  .  .  .  .  }
27  .  .  .  .  .  .  1: *ast.Ident {
28  .  .  .  .  .  .  .  NamePos: 40
29  .  .  .  .  .  .  .  Name: "b"
30  .  .  .  .  .  .  .  Obj: *ast.Object {...}
35  .  .  .  .  .  .  }
36  .  .  .  .  .  }
37  .  .  .  .  .  Type: *ast.Ident {
38  .  .  .  .  .  .  NamePos: 42
39  .  .  .  .  .  .  Name: "int"
40  .  .  .  .  .  }
41  .  .  .  .  .  Tag: *ast.BasicLit {
42  .  .  .  .  .  .  ValuePos: 46
43  .  .  .  .  .  .  Kind: STRING
44  .  .  .  .  .  .  Value: "\"int value\""
45  .  .  .  .  .  }
46  .  .  .  .  }
47  .  .  .  .  1: *ast.Field {
48  .  .  .  .  .  Type: *ast.Ident {
49  .  .  .  .  .  .  NamePos: 59
50  .  .  .  .  .  .  Name: "string"
51  .  .  .  .  .  }
52  .  .  .  .  }
53  .  .  .  }
54  .  .  .  Closing: 66
55  .  .  }
56  .  .  Incomplete: false
57  .  }
```

所有的结构体成员由`*ast.FieldList`表示，其中有三个`*ast.Field`元素。第一个`*ast.Field`对应`a, b int "int value"`的成员声明，包含了成员名字列表、类型和Tag信息。最后的`*ast.Field`是嵌入的string成员，只有普通的名字而没有类型信息（匿名嵌入成员也可以单独定义Tag字符串）。

其中`ast.StructType`等和结构体相关的语法树结构定义如下：

```go
type StructType struct {
	Struct     token.Pos  // position of "struct" keyword
	Fields     *FieldList // list of field declarations
	Incomplete bool       // true if (source) fields are missing in the Fields list
}
type FieldList struct {
	Opening token.Pos // position of opening parenthesis/brace, if any
	List    []*Field  // field list; or nil
	Closing token.Pos // position of closing parenthesis/brace, if any
}
type Field struct {
	Doc     *CommentGroup // associated documentation; or nil
	Names   []*Ident      // field/method/parameter names; or nil
	Type    Expr          // field/method/parameter type
	Tag     *BasicLit     // field tag; or nil
	Comment *CommentGroup // line comments; or nil
}
```

StructType中最重要的信息是FieldList类型的Fields成员声明列表信息。而每一组成员声明又由`ast.Field`表示，其中包含一组成员的名字，共享的类型和Tag字符串。需要注意的是，`ast.Field`不仅仅用于表示结构体成员的语法树结点，同时也用于表示接口的方法列表、函数或方法的各种参数列表（接收者参数、输入参数和返回值），因此这是一个异常重要的类型。

## 7.7 Map类型

Map其实是从数组和结构体的混合类型发展而来。Map支持根据元素的名字（也就是key）动态添加删除元素，但是其中的所有元素必须有着相同的类型。很多其它语言甚至用Map代替结构体和数组，比如Lua中以Table关联数组同时实现了数组和结构体的功能，而JavaScript中也是通过类似Map的对象来实现结构体。Go作为一个静态语言将Map直接作为语言内置的语法构造引入是一个比较大胆激进的行为，但同时也简化了相关数据结构的编程（因为内置的语法增加了部分泛型的功能，大大提升了编程体验）。

Map类型的语法规范定义比较简单：

```
MapType = "map" "[" KeyType "]" ElementType .
KeyType = Type .
```

首先以map关键字开始，然后通过中括弧包含Key的类型，最后是元素的类型。需要注意的是，Map中的Key必须是可进行相等比较的类型（典型的切片就不能作为Key类型），但是在语法树解析阶段并不会做这类检查。

下面是基于map定义的新类型：

```go
type IntStringMap map[int]string
```

解析的语法树输出如下：

```
11  .  Type: *ast.MapType {
12  .  .  Map: 31
13  .  .  Key: *ast.Ident {
14  .  .  .  NamePos: 35
15  .  .  .  Name: "int"
16  .  .  }
17  .  .  Value: *ast.Ident {
18  .  .  .  NamePos: 39
19  .  .  .  Name: "string"
20  .  .  }
21  .  }
```

虽然Map功能强大，但是表示其类型的语法树比较简单。其中Key和Value部分都是类型表达式，这个例子中分别是int和string标识符。

下面是`ast.MapType`语法树结点的定义：

```go
type MapType struct {
	Map   token.Pos // position of "map" keyword
	Key   Expr
	Value Expr
}
```

其中Key和Value部分都是类型表达式，可以是其它更复杂的组合类型。


## 7.8 管道类型

管道是Go语言比较有特色的类型，管道有双向管道、只写管道和只读管道之分，同时管道有元素类型。管道类型的语法规范如下：

```
ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .
```

在语法树中管道类型由`ast.ChanType`结构体定义：

```go
type ChanType struct {
	Begin token.Pos // position of "chan" keyword or "<-" (whichever comes first)
	Arrow token.Pos // position of "<-" (token.NoPos if there is no "<-"); added in Go 1.1
	Dir   ChanDir   // channel direction
	Value Expr      // value type
}

type ChanDir int

const (
	SEND ChanDir = 1 << iota
	RECV
)
```

其中`ast.ChanType.Dir`是管道的方向，SEND表示发送、RECV表示接收、`SEND|RECV`比特位组合表示双向管道。下面的例子是一个双向的int管道：

```go
type IntChan chan int
```

解析的语法树结果如下：

```
11  .  Type: *ast.ChanType {
12  .  .  Begin: 26
13  .  .  Arrow: 0
14  .  .  Dir: 3
15  .  .  Value: *ast.Ident {
16  .  .  .  NamePos: 31
17  .  .  .  Name: "int"
18  .  .  }
19  .  }
```

其中`ast.ChanType.Dir`值是3，也就是`SEND|RECV`比特位组合，表示这是一个双向管道。而`ast.ChanType.Value`部分表示管道值的类型，这里是一个`ast.Ident`表示的int类型。

## 7.9 函数类型

函数类型基本上是函数签名部分，包含函数的输入参数和返回值类型。在函数声明一节我们已经见过函数声明的语法规范，但是函数类型不包含函数的名字。函数类型的语法规范如下：

```
FunctionType   = "func" Signature .
Signature      = Parameters [ Result ] .
Result         = Parameters | Type .
Parameters     = "(" [ ParameterList [ "," ] ] ")" .
ParameterList  = ParameterDecl { "," ParameterDecl } .
ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
```

func关键字后面直接是输入参数和返回值列表组成的函数签名，不包含函数的名字。下面函数类型的一个例子：

```go
type FuncType func(a, b int) bool
```

函数类型中类型部分也是由`ast.FuncType`结构体定义。关于函数类型的细节请参考函数声明章节。

## 7.10 接口类型

从语法结构角度看，接口和结构体类型很像，不过接口的每个成员都是函数类型。接口类型的语法规则如下：

```
InterfaceType      = "interface" "{" { MethodSpec ";" } "}" .
MethodSpec         = MethodName Signature | InterfaceTypeName .
MethodName         = identifier .
InterfaceTypeName  = TypeName .

Signature          = Parameters [ Result ] .
Result             = Parameters | Type .
```

接口中每个成员都是函数类型，但是函数类型部分不包含func关键字。下面是只要一个方法成员的接口：

```go
type IntReader interface {
	Read() int
}
```

对齐分析语法树结果如下：

```
11  .  Type: *ast.InterfaceType {
12  .  .  Interface: 28
13  .  .  Methods: *ast.FieldList {
14  .  .  .  Opening: 38
15  .  .  .  List: []*ast.Field (len = 1) {
16  .  .  .  .  0: *ast.Field {
17  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
18  .  .  .  .  .  .  0: *ast.Ident {
19  .  .  .  .  .  .  .  NamePos: 41
20  .  .  .  .  .  .  .  Name: "Read"
21  .  .  .  .  .  .  .  Obj: *ast.Object {
22  .  .  .  .  .  .  .  .  Kind: func
23  .  .  .  .  .  .  .  .  Name: "Read"
24  .  .  .  .  .  .  .  .  Decl: *(obj @ 16)
25  .  .  .  .  .  .  .  }
26  .  .  .  .  .  .  }
27  .  .  .  .  .  }
28  .  .  .  .  .  Type: *ast.FuncType {
29  .  .  .  .  .  .  Func: 0
30  .  .  .  .  .  .  Params: *ast.FieldList {
31  .  .  .  .  .  .  .  Opening: 45
32  .  .  .  .  .  .  .  Closing: 46
33  .  .  .  .  .  .  }
34  .  .  .  .  .  .  Results: *ast.FieldList {
35  .  .  .  .  .  .  .  Opening: 0
36  .  .  .  .  .  .  .  List: []*ast.Field (len = 1) {
37  .  .  .  .  .  .  .  .  0: *ast.Field {
38  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
39  .  .  .  .  .  .  .  .  .  .  NamePos: 48
40  .  .  .  .  .  .  .  .  .  .  Name: "int"
41  .  .  .  .  .  .  .  .  .  }
42  .  .  .  .  .  .  .  .  }
43  .  .  .  .  .  .  .  }
44  .  .  .  .  .  .  .  Closing: 0
45  .  .  .  .  .  .  }
46  .  .  .  .  .  }
47  .  .  .  .  }
48  .  .  .  }
49  .  .  .  Closing: 52
50  .  .  }
51  .  .  Incomplete: false
52  .  }
```

接口的语法树是`ast.InterfaceType`类型，其`Methods`成员列表和结构体成员的`*ast.FieldList`类型一样。下面是`ast.InterfaceType`和`ast.StructType`语法树结构的定义：

```go
type InterfaceType struct {
	Interface  token.Pos  // position of "interface" keyword
	Methods    *FieldList // list of methods
	Incomplete bool       // true if (source) methods are missing in the Methods list
}
type StructType struct {
	Struct     token.Pos  // position of "struct" keyword
	Fields     *FieldList // list of field declarations
	Incomplete bool       // true if (source) fields are missing in the Fields list
}
```

对比可以发现，接口和结构体语法树结点中除了方法列表和成员列表的名字不同之外，方法和成员都是由`ast.FieldList`定义的。因此上述的接口例子和下面的结构体其实非常相似：

```go
type IntReader struct {
	Read func() int
}
```

如果是结构体，那么Read成员就是一个函数类型，函数是`func() int`类型。总之在语法树层面接口和结构体可以采用相似的代码处理。

## 7.11 组合类型

复合类型最强大的地方在于通过不同组合生成更复杂的类型。但是第一步需要搞清楚基于基础类型构造的复合类型，然后才是复合类型之间的组合。在掌握了基础类型和复合类型的语法树结构之后，我们就可以解析任意复杂的类型，同时也就很容易理解Go语言中反射的类型结构。不管是数据结构还是函数都需要和类型关联，因此理解类型之后就把握了整个程序的脉络，剩下的就是向函数体中填充语句而已。
