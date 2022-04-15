# 附录B ANTLR4

ANTLR（ANother Tool for Language Recognition）是由Terence Parr博士开发的的语法分析器生成工具，可用于读取、处理、执行和翻译结构化的文件。ANTLR第一版采用C语言开发在1989年发布（第一版名字叫PCCTS），ANTLR2开始改用Java语言实现并在1997年发布第一版本，ANTLR3和ANTLR4分别在2005年和2013年发布，支持Go语言的ANTLR4.6在2016年底发布，从诞生到现在ANTLR已经有30多年的发展历史。目前ANTLR已经成为很多语言、工具和框架软件的基石，本节我们将简单展示如何通过ANTLR4构造一个Go语言版本的计算器。

## B.1 构造语法文件

首先创建`Calc.g4`文件描述表达式语法：

```antlr4
// Calc.g4
grammar Calc;

// Tokens
MUL: '*';
DIV: '/';
ADD: '+';
SUB: '-';
NUMBER: [0-9]+;
WHITESPACE: [ \r\n\t]+ -> skip;

// Rules
start : expression EOF;

expression
   : expression op=('*'|'/') expression # MulDiv
   | expression op=('+'|'-') expression # AddSub
   | NUMBER                             # Number
   ;
```

其中`MUL`、`DIV`、`ADD`、`SUB`、`NUMBER`和`WHITESPACE`采用大写名字书写的规则表示词法记号，每个词法规则采用类似正则表达式的语法书写，最后`WHITESPACE`后的`-> skip`动作表示跳过空白字符。

而以小写字母表示的`start`和`expression`则是采用BNF语法书写的语法规则：第一个`start`表示语法的开始，`expression`表示表达式语法。如果语法对应多个不同的规则，第一出现的规则优先级最高，同时每个规则后面的`# MulDiv`表示响应改规则的方法名字。

安装ANTLR4的jar包之后，可以用以下的命令输出Go语言版本的语法解析器：

```
$ java -jar antlr-4.8-complete.jar -Dlanguage=Go -o calc Calc.g4
```

其中`-Dlanguage=Go`表示输出Go语言版本的解析器代码，`-o calc`表示将Go代码输出到`calc`目录下。

## B.2 基于生成代码构造解析器

ANTLR4生成的代码已经包含了表达式完整的词法和语法分析器。比如可以按照如下方式为`"1+2*3"`表达式构造表达式分析器：

```go
import (
	"github.com/antlr/antlr4/runtime/Go/antlr"

	calc "./calc"
)

func main() {
	lexer := calc.NewCalcLexer(antlr.NewInputStream("1+2*3"))
	parser := calc.NewCalcParser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	...
}
```

首先生成的`calc.NewCalcLexer`函数基于`antlr.NewInputStream("1+2*3")`输入的文本流构建词分析器。然后`calc.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)`将解析到的词法记号转化为记号流传递给生成的`calc.NewCalcParser`语法解析函数。

查看生成的`CalcParser`结构体的导出方法：

```go
$ go doc CalcParser
package parser // import "."

type CalcParser struct {
	*antlr.BaseParser
}

func NewCalcParser(input antlr.TokenStream) *CalcParser
func (p *CalcParser) Expression() (localctx IExpressionContext)
func (p *CalcParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool
func (p *CalcParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool
func (p *CalcParser) Start() (localctx IStartContext)
```

其中`Start()`返回当前语法解析器得到的一个名字为`start`的规则对应的语法树。

## B.3 `calc.CalcListener`接口

要遍历语法树必须先实现`calc.CalcListener`接口：

```go
// CalcListener is a complete listener for a parse tree produced by CalcParser.
type CalcListener interface {
	antlr.ParseTreeListener

	EnterStart(c *StartContext)
	EnterNumber(c *NumberContext)
	EnterMulDiv(c *MulDivContext)
	EnterAddSub(c *AddSubContext)
	ExitStart(c *StartContext)
	ExitNumber(c *NumberContext)
	ExitMulDiv(c *MulDivContext)
	ExitAddSub(c *AddSubContext)
}
```

每个规则都有对应的进入和退出方法，比如`EnterMulDiv`和`ExitMulDiv`对应`# MulDiv`标柱的语法规则。这里我们只关心乘除法和加减法的规则，实现`calcListener`如下：

```go
type calcListener struct {
	*cacl.BaseCalcListener
}

func (l *calcListener) ExitMulDiv(c *parser.MulDivContext) { /* TODO */ }
func (l *calcListener) ExitAddSub(c *parser.AddSubContext) { /* TODO */ }
func (l *calcListener) ExitNumber(c *parser.NumberContext) { /* TODO */ }
```

其中`calc.BaseCalcListener`是ANTLR4实现的基础遍历者，基于这个对象继承的`calcListener`结构体只需要重新实现需要的方法即可满足`calc.CalcListener`接口。

## B.3 实现遍历规则的方法

对于带小括号的四则运算表达式需要一个临时栈用于保存中间结果：

```go
type calcListener struct {
	*cacl.BaseCalcListener
	stk []int
}

func (p *calcListener) push(i int) {
	p.stk = append(p.stk, i)
}
func (p *calcListener) pop() int {
	result := p.stk[len(p.stk)-1]
	p.stk = p.stk[:len(p.stk)-1]
	return result
}
```

其中`push`和`pop`分别对于入栈和出栈操作。然后就可以先实现`ExitNumber`方法，它在退出一个数字前被调用：

```go
func (l *calcListener) ExitNumber(c *parser.NumberContext) {
	i, _ := strconv.Atoi(c.GetText())
	l.push(i)
}
```

通过`c.GetText()`从当前上下文获取当前的数字，然后通过`l.push(i)`压入临时栈保存。然后分别在乘除法和加减法规则时从栈消费临时栈保存的中间结果，最终将运算的中间结果再压入临时栈中。

乘除法和加减法对应的`ExitMulDiv`和`ExitAddSub`方法实现如下：

```go
func (l *calcListener) ExitMulDiv(c *parser.MulDivContext) {
	right, left := l.pop(), l.pop()
	switch c.GetOp().GetTokenType() {
	case parser.CalcParserMUL:
		l.push(left * right)
	case parser.CalcParserDIV:
		l.push(left / right)
	}
}

func (l *calcListener) ExitAddSub(c *parser.AddSubContext) {
	right, left := l.pop(), l.pop()
	switch c.GetOp().GetTokenType() {
	case parser.CalcParserADD:
		l.push(left + right)
	case parser.CalcParserSUB:
		l.push(left - right)
	}
}
```

需要注意的是`l.pop()`从临时栈先弹出的值是二元表达式右边的值。`c.GetOp().GetTokenType()`是通过当前上下文获得当前运算符，获取运算符记号对应的函数名`GetOp`是根据`expression op=('*'|'/') expression`语法中的`op=('*'|'/')`名字生成。

## B.4 遍历语法树

现在就可以通过以下方法遍历针对返回的语法树：

```go
func main() {
	lexer := calc.NewCalcLexer(antlr.NewInputStream("1+2*3"))
	parser := calc.NewCalcParser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	var l calc.CalcListener = new(calcListener)
	antlr.NewParseTreeWalker().Walk(l, parser.Start())
	fmt.Println(l.(*calcListener).pop())
}
```

其中`l`是刚实现的满足`calc.CalcListener`接口的遍历者对象，`parser.Start()`是表达式解析后得到的语法树。如果表达式没有语法错误，遍历完成之后通过`pop()`方法从临时栈中弹出最后的运算结果。

## B.5 补充说明

ANTLR4对Go语言支持虽然只有几年时间但是已经足够稳定，比如Google基于Protobuf设计的CEL验证语言的Go语言版本就是基于ANTLR4生成语法解析器。ANTLR4是功能强大的语法解析器生成工具，不仅仅支持基于Listener模式的遍历，还支持通过Visitor模式支持更多定制操作的语法树遍历。此外ANTLR4在社区中配套的辅助工具也非常完善，比如VSCode就有对应的插件以铁路图等不同等方式展示语法文件。更详细的用法请参考Terence Parr博士的《编程语言实现模式》和《ANTLR4权威指南》，它们才是码农真正需要的屠龙刀和倚天剑。
