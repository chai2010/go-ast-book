# 附录A goyacc

yacc是用于构造编译器的工具，而goyacc是Go语言版本的yacc，是从早期的C语言版本yacc移植到Go语言的。早期的goyacc是Go语言标准命令之一，也是构建Go自身编译器的必备工具链之一，后来被逐步移出了内置工具。但是goyacc依然是一个开发语法分析器的利器。本章简单展示如何用goyacc构建一个命令行计算器小程序。

## A.1 计算器的特性

特性简介：

- 支持整数四则运算
- 支持小括弧提升优先级
- 支持临时变量保存结果

安装和使用(需要有GCC环境)：

```shell
$ go get github.com/chai2010/calculator
$ calculator
1+2*3
= 7
x=3-(2-1)
= 2
x*2
= 4
```

## A.2 词法符号

先创建`tok.h`文件，包含词法符号：

```c
enum {
	ILLEGAL = 10000,
	EOL = 10001,

	ID = 258,
	NUMBER = 259,

	ADD = 260, // +
	SUB = 261, // -
	MUL = 262, // *
	DIV = 263, // /
	ABS = 264, // |

	LPAREN = 265, // (
	RPAREN = 266, // )
	ASSIGN = 267, // =
};
```

其中`ILLEGAL`表示不能识别的无效的符号，`EOL`表示行的结尾，其它的符号与字面含义相同。

## A.3 词法解析

然后创建`calc.l`文件，定义每种词法的正则表达式：

```lex
%option noyywrap

%{
#include "tok.h"
%}

%%

[_a-zA-Z]+ { return ID; }
[0-9]+     { return NUMBER; }

"+"    { return ADD; }
"-"    { return SUB; }
"*"    { return MUL; }
"/"    { return DIV; }
"|"    { return ABS; }

"("    { return LPAREN; }
")"    { return RPAREN; }
"="    { return ASSIGN; }

\n     { return EOL; }
[ \t]  { /* ignore whitespace */ }
.      { return ILLEGAL; }

%%
```

最开始的`noyywrap`选项表示关闭`yywrap`特性，也就是去掉对flex库的依赖，生成可移植的词法分析器代码。然后在`%{`和`%}`中间是原生的C语言代码，通过包含`tok.h`引入了每种记号对应的枚举类型。在两组`%%`中间的部分是每种记号对应的正则表达式，先出现的优先匹配，如果匹配失败则继续尝试后面的规则。每个正则表达式后面跟着一组动作代码，也就是普通的C语言代码，这里都是返回记号的类型。

然后通过flex工具生成C语言词法解析器文件：

```shell
$ flex --prefix=yy --header-file=calc.lex.h -o calc.lex.c calc.l
```

其中`--prefix`表示生成的代码中标识符都是以`yy`前缀。在一个项目有多个flex生成代码时，可通过前缀区分。`--header-file`表示生成头问题，这样方便在其它代码中引用生成的词法分析函数。`-o`指定输出源代码文件的名字。

生成的词法分析器中，最重要的有以下几个：

```c
extern int yylineno;
extern char *yytext;

extern int yylex (void);
```

其中`yylineno`表示当前的行号，`yytext`表示当前记号对应的字符串。而`yylex`函数每次从标准输入读取一个记号，返回记号类型的值（在`tok.h`文件定义），如果遇到文件结尾则返回0。

如果需要从字符串解析，则需使用以下的导出函数：

```c
YY_BUFFER_STATE yy_scan_bytes (yyconst char *bytes,yy_size_t len  );
```

通过`yy_scan_bytes`函数，可以设置字符串作为要解析的目标，然后每次调用`yylex`函数就会从字符串读取数据。这些函数都在`calc.lex.h`文件中声明。

## A.4 将C语言词法分析器包装为Go函数

创建`lex.go`文件，内容如下：

```go
package main

//#include "tok.h"
//#include "calc.lex.h"
import "C"

type calcLex struct {}

func newCalcLexer(data []byte) *calcLex {
	p := new(calcLex)
	C.yy_scan_bytes((*C.char)(C.CBytes(data)), C.yy_size_t(len(data)))
	return p
}

func (p *calcLex) Lex(yylval *calcSymType) int {
	var tok = C.yylex()
	var yylineno = int(C.yylineno)
	var yytext = C.GoString(C.yytext)

	switch tok {
	case C.ID:
		// yylval.id = yytext
		return ID

	case C.NUMBER:
		//yylval.value, _ = strconv.Atoi(yytext)
		return NUMBER

	case C.ADD:
		return ADD
	// ...

	case C.EOL:
		return EOL
	}

	if tok == C.ILLEGAL {
		log.Printf("lex: ILLEGAL token, yytext = %q, yylineno = %d", yytext, yylineno)
	}

	return 0 // eof
}
```

新建的`calcLex`类型对应Go语言版本的词法分析器，底层工作通过CGO调用flex生成的C语言函数完成。首先`newCalcLexer`创建一个词法分析器，参数是要分析的数据，通过`C.yy_scan_bytes`函数调用表示从字符串解析记号。然后`calcLex`类型的`Lex`方法表示每次需要解析一个记号（暂时忽略方法的`calcSymType`参数），内部通过调用`C.yylex()`读取一个记号，同时记录行号和记号对应的字符串。最后将C语言的记号转为Go语言的记号值返回，比如`C.ID`对应Go语言的`ID`。

对应`ID`类型，`yytext`表示变量的名字。对于`NUMBER`类型，`yytext`保护数字对应的字符串，可以从字符串解析出数值。但是，Go语言的词法分析器如何返回变量的名字或者是数字的值呢？答案是通过`Lex`的`*calcSymType`类型的参数可以记录记号额外的属性值。而`calcSymType`类型是由`goyacc`工具生成的代码，在下面我们将介绍yacc的内容。

## A.5 `goyacc`生成语法解析器

`goyacc`是Go语言版本的yacc工具，是由Go语言官方团队维护的扩展包工具。

创建`calc.y`文件：

```yacc
%{
package main

var idValueMap = map[string]int{}
%}

%union {
	value int
	id    string
}

%type  <value> exp factor term
%token <value> NUMBER
%token <id>    ID

%token ADD SUB MUL DIV ABS
%token LPAREN RPAREN ASSIGN
%token EOL

%%
calclist
	: // nothing
	| calclist exp EOL {
		idValueMap["_"] = $2
		fmt.Printf("= %v\n", $2)
	}
	| calclist ID ASSIGN exp EOL {
		idValueMap["_"] = $4
		idValueMap[$2] = $4
		fmt.Printf("= %v\n", $4)
	}
	;

exp
	: factor         { $$ = $1 }
	| exp ADD factor { $$ = $1 + $3 }
	| exp SUB factor { $$ = $1 - $3 }
	;

factor
	: term            { $$ = $1 }
	| factor MUL term { $$ = $1 * $3 }
	| factor DIV term { $$ = $1 / $3 }
	;

term
	: NUMBER            { $$ = $1 }
	| ID                { $$ = idValueMap[$1] }
	;

%%
```

和flex工具类型，首先在`%{`和`%}`中间是原生的Go语言代码。然后`%union`定义了属性值，用于记录语法解析中每个规则额外的属性值。通过`%type`定义BNF规则中非终结的名字，`%token`定义终结记号名字（和flex定义的记号类型是一致的）。而`%type`和`%token`就可以通过`<value>`或`<id>`的可选语法，将后面的名字绑定到属性。就是后续代码中`$$`对应的属性，比如`%token <id> ID`表示`ID`对应的属性为`id`，因此在后面的`ID { $$ = idValueMap[$1] }`表示数值`id`属性的值，其中`idValueMap`用于管理变量的值。

然后通过goyacc工具生成代码：

```shell
$ goyacc -o calc.y.go -p "calc" calc.y
```

其中`-o`指定输出的文件名，`-p`指定标识符名字前缀（和flex的`--prefix`用法类似）。在生成的`calc.y.go`文件中将包含最重要的`calcParse`函数，该函数从指定的词法解析器中读取词法，然后进行语法分析。同时将包含`calcSymType`类型的定义，它是`Lex`词法函数的输出参数的类型。

在绑定了属性之后，还需要继续完善`Lex`词法函数的代码：

```go
func (p *calcLex) Lex(yylval *calcSymType) int {
	var tok = C.yylex()
	var yylineno = int(C.yylineno)
	var yytext = C.GoString(C.yytext)

	switch tok {
	case C.ID:
		yylval.id = yytext
		return ID

	case C.NUMBER:
		yylval.value, _ = strconv.Atoi(yytext)
		return NUMBER

	...
}
```

其中`yylval.id = yytext`表示词法将解析得到的变量名字填充到`id`属性中。而数字部分则是通过`yylval.value`属性保存。


## A.6 运行计算器

创建main函数：

```go
func main() {
	calcParse(newCalcLexer([]byte("1+2*3")))
}
```

`newCalcLexer`构造一个词法解析器，然后`calcParse`语法解析器将从词法解析器依次读取记号并解析语法，在解析语法的同时将进行表达式求值运算，同时更新`idValueMap`全局的变量。
