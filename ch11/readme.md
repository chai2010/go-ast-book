# 第11章 类型检查

主流的编译器前端遵循词法解析、语法解析、语义解析等流程，然后才是基于中间表示的层层优化并最终产生目标代码。在得到抽象的语法树之后就表示完成了语法解析的工作。不过在进行中间优化或代码生成之前还需要对抽象语法树进行语义分析。语义分析需要更深层次理解代码的语义，比如两个变量相加是否合法，外层作用域有多个同名的变量时如何选择等。本章简单讨论`go/types`包的用法，展示如果通过该包实现语法树的类型检查功能。

## 11.1 语义错误

虽然Go语言是基于包和目录来组织代码，但是Go语言在语法树解析阶段并不关心包之间的依赖关系。这是因为在语法树解析阶段并不对代码本身做语义检测，因此很多语法正确但是语义错误的代码也可以生成语法树。

比如以下这个例子：

```go
func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	ast.Print(fset, f)
}

const src = `package pkg

func hello() {
	var _ = "a" + 1
}
`
```

在被解析的代码的`hello`函数，可以正常生成语法树。但是`hello`函数中唯一的语句`var _ = "a" + 1`的语义却是错误的，因为Go语言中不能将一个字符串和一个数字进行相加。如何识别这种语义层面的错误是`go/types`包需要完成的工作。

## 11.2 `go/types`包

`go/types`包是Go语言之父Robert Griesemer大神（发明了Go语言的接口等特性）开发的类型检查工具.该包从Go1.5时代开始被添加到标准库，是Go语言自举过程中的一个额外成果。据说这个包是Go语言标准库中代码量最大的一个包，也是功能最复杂的一个包（在使用之前需要对Go语法树有一定的基础知识）。这里我们将使用`go/types`包来检查之前例子中的语法错误。

重新调整代码如下：

```go
func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	pkg, err := new(types.Config).Check("hello.go", fset, []*ast.File{f}, nil)
	if err != nil {
		log.Fatal(err)
	}

	_ = pkg
}

const src = `package pkg

func hello() {
	var _ = "a" + 1
}
`
```

在通过`parser.ParseFile`解析单文件得到语法树之后，通过`new(types.Config).Check`函数来解析语法树中的语义错误。`new(types.Config)`首先是构造一个用于类型检查的配置对象，然后调用其唯一的`Check`方法检测语法树的语义。`Check`方法的签名如下：

```go
func (conf *Config) Check(path string, fset *token.FileSet, files []*ast.File, info *Info) (*Package, error)
```

第一个参数表示要检查包的路径，第二个参数表示全部的文件集合（用于将语法树中元素的位置信息解析为文件名和行列号），第三个参数是该包中所有文件对应的语法树，最后一个参数可用于存储检查过程中产生的分析结果。如果成功该方法返回一个`types.Package`对象，表示当前包的信息。

运行这个程序将产生以下的错误信息：

```
$ go run .
hello.go:4:10: cannot convert "a" (untyped string constant) to untyped int
```

错误提示在`hello.go`文件的第4行第10个字符位置的`"a"`字符串语法错误，无法将字符串常量转化为无类型的`int`类型。这样我们就可以轻易定位代码中出现错误的位置和错误产生的原因。

## 11.3 跨包的类型检查

真实的代码总是由多个包组成的，而`go/parser`包只处理当前包，如何处理导入包的类型是一个重要问题。比如有以下的代码：

```go
package main

import "math"

func main() {
	var _ = "a" + math.Pi
}
```

代码导入的是`math`包，然后引用了其中的`math.Pi`元素。要验证当前代码是否语义正确的前提，首先需要获取`math.Pi`元素的类型，因此首先要处理包的导入问题。

如果依然采用`new(types.Config).Check`方式验证将得到以下的错误：

```
hello.go:3:8: could not import math (Config.Importer not installed)
```

错误产生的原因是`types.Config`类型的检查对象并不知道如何加载`math`包的信息。`types.Config`对象的`Importer`成员复杂导入依赖包，其定义如下：

```go
type Config struct {
	Importer Importer
}

type Importer interface {
	Import(path string) (*Package, error)
}
```

对于任何一个导入包都会调用`Import(path string) (*Package, error)`加载导入信息，然后才能获取包中导出元素的信息。

对于标准库的`math`包，可以采用`go/importer`提供的默认包导入实现。代码如下：

```go
	// import "go/importer"
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("hello.go", fset, []*ast.File{f}, nil)
	if err != nil {
		log.Fatal(err)
	}
```

其中`types.Config`对象的`Importer`成员对应包导入对象，由`importer.Default()`初始化。然后就可以正常处理输入代码了。

不过`importer.Default()`处理的是Go语义当前环境的代码结构。Go语义代码结构是比较复杂的，其中包含标准库和用户的模块代码，每个包还可能启动了CGO特性。为了便于理解，我们可以手工构造一个简单的`math`包，因此包导入过程也可以简化。

为了简化，我们继续假设每个包只有一个源代码文件。定义`Program`结构体表示一个完整的程序对象，代码如下：

```go
type Program struct {
	fs   map[string]string
	ast  map[string]*ast.File
	pkgs map[string]*types.Package
	fset *token.FileSet
}

func NewProgram(fs map[string]string) *Program {
	return &Program{
		fs:   fs,
		ast:  make(map[string]*ast.File),
		pkgs: make(map[string]*types.Package),
		fset: token.NewFileSet(),
	}
}
```

其中`fs`表示每个包对应的源代码字符串，`ast`表示每个包对应的语法树，`pkgs`表示经过语义检查的包对象，`fset`则表示文件的位置信息。

首先为`Program`类型增加包加载`LoadPackage`方法：

```go
func (p *Program) LoadPackage(path string) (pkg *types.Package, f *ast.File, err error) {
	if pkg, ok := p.pkgs[path]; ok {
		return pkg, p.ast[path], nil
	}

	f, err = parser.ParseFile(p.fset, path, p.fs[path], parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}

	conf := types.Config{Importer: nil}
	pkg, err = conf.Check(path, p.fset, []*ast.File{f}, nil)
	if err != nil {
		return nil, nil, err
	}

	p.ast[path] = f
	p.pkgs[path] = pkg
	return pkg, f, nil
}
```

因为没有初始化`types.Config`的`Importer`成员，因此目前该方法只能加载没有导入其他包的叶子类型的包（对应`math`包就是这种类型）。比如叶子类型的`math`包被加载成功之后，则会被记录到`Program`对象的`ast`和`pkgs`成员中。然后当遇到已经被记录过的叶子包被导入时，就可以复用这些信息。

因此可以为`Program`类型实现`types.Importer`接口，只有一个`Import`方法：

```go
func (p *Program) Import(path string) (*types.Package, error) {
	if pkg, ok := p.pkgs[path]; ok {
		return pkg, nil
	}
	return nil, fmt.Errorf("not found: %s", path)
}
```

现在`Program`类型实现了`types.Importer`接口，就可以用于`types.Config`的包加载工作：

```go
func (p *Program) LoadPackage(path string) (pkg *types.Package, f *ast.File, err error) {
	// ...

	conf := types.Config{Importer: p} // 用 Program 作为包导入器
	pkg, err = conf.Check(path, p.fset, []*ast.File{f}, nil)
	if err != nil {
		return nil, nil, err
	}

	// ...
}
```

然后可以通过手工方式先加载叶子类型的`math`包，然后再加载主包：

```go
func main() {
	prog := NewProgram(map[string]string{
		"hello": `
			package main
			import "math"
			func main() { var _ = 2 * math.Pi }
		`,
		"math": `
			package math
			const Pi = 3.1415926
		`,
	})

	_, _, err := prog.LoadPackage("math")
	if err != nil {
		log.Fatal(err)
	}

	pkg, f, err := prog.LoadPackage("hello")
	if err != nil {
		log.Fatal(err)
	}
}
```

这种依赖包的导入包的加载是递归的，因此可以在导入环节的`Import`方法增加递归处理：

```go
func (p *Program) Import(path string) (*types.Package, error) {
	if pkg, ok := p.pkgs[path]; ok {
		return pkg, nil
	}
	pkg, _, err := p.LoadPackage(path)
	return pkg, err
}
```

当`pkgs`成员没有包信息时，通过`LoadPackage`方法加载。如果`LoadPackage`要导入的包是非叶子类型的包，会再次递归回到`Import`方法。因为Go语义禁止循环包导入，因此最终会在导入叶子包的时刻由`LoadPackage`函数返回结束递归。当然在真实的代码中，需要额外记录一个状态用于检查递归导入类型的错误。

这样我们就实现了一个支持递归包导入的功能，从而可以实现对于任何一个加载的语法树进行完整的类型检查。
