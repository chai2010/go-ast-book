# Go语法树入门——开启自制编程语言和编译器之旅！

Go语法树是Go语言源文件的另一种语义等价的表现形式。而Go语言自带的`go fmt`和`go doc`等命令都是在Go语法树基础之上的分析工具。因此将Go语言程序作为输入数据，从语法树这个维度重新审视Go语言程序，我们将得到创建Go语言本身的技术。Go语法树由标准库的`go/ast`包定义，它是在`go/token`包定义的词法基础之上抽象的语法树结构。本书简单介绍语法树相关包的使用（出版社已经约稿出版本书，并在开源版本的基础之上增加了语义信息、SSA形式、LLVM和凹语言等内容，因为出版社版权问题不方便全部公开新增内容）。

![](cover.png)

- 作者：柴树杉，Github [@chai2010](https://github.com/chai2010)，Twitter [@chaishushan](https://twitter.com/chaishushan)
- 作者：史斌，Github [@benshi001](https://github.com/benshi001)
- 作者：丁尔男，Github [@3dgen](https://github.com/benshi001)
- 主页：https://github.com/chai2010/go-ast-book


# 版权

版权 [柴树杉](https://github.com/chai2010)、[史斌](https://github.com/benshi001)和[丁尔男](https://github.com/3dgen)，保留相关权力。针对Github注册用户提供以下的优惠权利：

1. Github平台免费在线阅读。
1. 关注本书项目(Star)，同时关注任意一个作者的 Github 或 推特账号，自动获得下载的权利。

**禁止非 Github 平台转载，作者保留相关法律权力。**

