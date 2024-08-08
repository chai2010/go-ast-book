// 作者：史斌（https://github.com/benshi001）
// 转载请注明原作者

package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

var lineNo int = 0                         // 行编号
var varNo int = 0                          // 临时变量编号
var vars map[string]int = map[string]int{} // 已定义的变量
var srcError bool = false                  // 源文件是否包含错误

// 如果出现错误，则删除生成的.ll文件
func remove(path string) {
	if srcError {
		os.Remove(path)
	}
}

func main() {
	// 一次只能编译一个文件
	if len(os.Args) != 2 {
		fmt.Printf("正确用法：%s XXX.w\n", os.Args[0])
		return
	}

	// 中间结果文件和最终目标文件的基础路径
	// 如果源文件名是xxx.w，则生成xxx.ll/xxx.s/xxx.exe
	// 否则直接在源文件名后面追加.ll/.s/.exe
	var basePath string
	if path.Ext(os.Args[1]) == ".w" {
		basePath = strings.ReplaceAll(os.Args[1], ".w", "")
	} else {
		basePath = os.Args[1]
	}
	defer remove(basePath + ".ll")

	// 打开源文件
	fSrc, e0 := os.Open(os.Args[1])
	if e0 != nil {
		fmt.Printf("无法读取源文件%s\n", os.Args[1])
		return
	}
	defer fSrc.Close()

	// 创建.ll文件并写入基本信息
	fLl, e1 := os.Create(basePath + ".ll")
	if e1 != nil {
		fmt.Printf("无法创建文件%s\n", basePath+".ll")
		return
	}
	defer fLl.Close()

	// 生成.ll文件的开头
	fLl.WriteString("; source file: " + os.Args[1])
	fLl.WriteString("\n@str = constant [4 x i8] c\"%d\\0A\\00\"\n")
	fLl.WriteString("declare i32 @printf(i8*, ...)\n")
	fLl.WriteString("define i32 @main() {\n")
	fLl.WriteString("  %fmt = getelementptr [4 x i8], [4 x i8]* @str, i32 0, i32 0\n")

	// 逐行读取源文件并生成LLVM-IR
	brSrc := bufio.NewReader(fSrc)
	for {
		line, _, c := brSrc.ReadLine()
		if c == io.EOF {
			break
		}
		lineNo++

		// 源代码以注释形式插入
		lineSrc := strings.TrimSpace(string(line))
		if len(lineSrc) == 0 {
			// 忽略空行
			continue
		}
		fLl.WriteString("  ; " + lineSrc + "\n")

		// 这里我耍了一个花招，把所有的=替换为<，因此"a=b+c"
		// 将被替换为"a<b+c"。
		// 原因是a=b+c是语句，而a<b+c是表达式，我希望使用更简单的
		// parser.ParseExpr来分析源代码。
		lineSrc = strings.ReplaceAll(lineSrc, "=", "<")

		// 分析整行源代码
		expr, e2 := parser.ParseExpr(lineSrc)
		if e2 != nil {
			srcError = true
			fmt.Printf("源文件%s第%d行包含语法错误\n", os.Args[1], lineNo)
			return
		}

		// 判断是赋值语句还是print语句
		if callExpr, b := expr.(*ast.CallExpr); b { // print语句
			if b := processPrint(callExpr, fLl); !b {
				return
			}
		} else if binExpr, b := expr.(*ast.BinaryExpr); b { // 赋值语句
			if _, b := processExpr(binExpr, fLl); !b {
				return
			}
		} else {
			srcError = true
			fmt.Printf("源文件%s第%d行包含不支持的语法\n", os.Args[1], lineNo)
			return
		}
	}

	// 生成.ll文件的结尾
	fLl.WriteString("  ret i32 0\n}\n")
	fLl.Close()

	// 调用clang
	cmd := exec.Command("clang", basePath+".ll", "-O0", "-o", basePath+".exe")
	if e3 := cmd.Run(); e3 != nil {
		srcError = true
		fmt.Printf("调用clang失败，可能原因：\n")
		fmt.Printf("1. 未正确安装LLVM；\n")
		fmt.Printf("2. 源代码%s存在其它语法错误。\n", os.Args[1])
	}
}

// 这个函数递归调用自己，分析表达式，对子表达式的结果生成临时变量来保存。
// 第一个返回值是保存输入表达式结果的临时变量编号，可能被上一级表达式引用。
// 第二个返回值是输入表达式是否已被正确解析。
func processExpr(expr interface{}, fLl *os.File) (int, bool) {
	if binExpr, b0 := expr.(*ast.BinaryExpr); b0 { // 二元表达式
		switch binExpr.Op {
		// 赋值语句的最顶级，前面我们把=替换成了<
		case token.LSS:
			x, b1 := binExpr.X.(*ast.Ident)
			if !b1 { // 错误：赋值语句左侧只能是变量
				srcError = true
				fmt.Printf("源文件%s第%d行：赋值语句左侧只能是变量\n", os.Args[1], lineNo)
				return -1, false
			}
			// 检查变量是否定义过
			if _, ok := vars[x.Name]; ok {
				srcError = true
				fmt.Printf("源文件%s第%d行：变量重复定义\n", os.Args[1], lineNo)
				return -1, false
			}
			// 生成赋值语句
			idx2, b2 := processExpr(binExpr.Y, fLl)
			if b2 {
				stmt := fmt.Sprintf("  %%%s = add i64 %%tmp%d, 0\n", x.Name, idx2)
				fLl.WriteString(stmt)
			}
			// 记录已定义的变量
			vars[x.Name] = lineNo
			return -1, b2

		// 加减乘除
		case token.ADD, token.SUB, token.MUL, token.QUO:
			idxLeft, bLeft := processExpr(binExpr.X, fLl)
			if !bLeft { // 左分支包含语法错误
				return -1, false
			}
			idxRight, bRight := processExpr(binExpr.Y, fLl)
			if !bRight { // 右分支包含语法错误
				return -1, false
			}
			varNo++
			opMap := map[token.Token]string{
				token.ADD: "add",
				token.SUB: "sub",
				token.MUL: "mul",
				token.QUO: "sdiv",
			}
			// 生成：tmpX = left <op> right
			stmt := fmt.Sprintf("  %%tmp%d = %s i64 %%tmp%d, %%tmp%d\n",
				varNo, opMap[binExpr.Op], idxLeft, idxRight)
			fLl.WriteString(stmt)
			return varNo, true

		// 不支持其它其它运算
		default:
			srcError = true
			fmt.Printf("源文件%s第%d行：不支持的运算\n", os.Args[1], lineNo)
			return -1, false
		}
	} else if vExpr, b0 := expr.(*ast.Ident); b0 { // 树形表达式的最末端，单个变量
		// 检查变量是否定义过
		if _, ok := vars[vExpr.Name]; !ok {
			srcError = true
			fmt.Printf("源文件%s第%d行：引用未定义的变量\n", os.Args[1], lineNo)
			return -1, false
		}
		// 生成赋值语句
		varNo++
		stmt := fmt.Sprintf("  %%tmp%d = add i64 %%%s, 0\n", varNo, vExpr.Name)
		fLl.WriteString(stmt)
		return varNo, true
	} else if cExpr, b0 := expr.(*ast.BasicLit); b0 { // 树形表达式的最末端，单个常量
		// 生成赋值语句
		varNo++
		stmt := fmt.Sprintf("  %%tmp%d = add i64 %s, 0\n", varNo, cExpr.Value)
		fLl.WriteString(stmt)
		return varNo, true
	} else if pExpr, b0 := expr.(*ast.ParenExpr); b0 { // 括号表达式
		idx, b1 := processExpr(pExpr.X, fLl)
		return idx, b1
	} else { // 不支持其它的表达式
		srcError = true
		fmt.Printf("源文件%s第%d行：不支持的表达式\n", os.Args[1], lineNo)
		return -1, false
	}
}

// 返回true表示正常生成代码
// 返回false表示源代码有语法错误
func processPrint(call *ast.CallExpr, fLl *os.File) bool {
	// 只能打印一个值
	if len(call.Args) != 1 {
		srcError = true
		fmt.Printf("源文件%s第%d行：print只能打印一个数值\n", os.Args[1], lineNo)
		return false
	} else if litExpr, b := call.Args[0].(*ast.BasicLit); b {
		// 生成打印常量的printf
		fLl.WriteString("  call i32 (i8*, ...) @printf(i8* %fmt, i64 ")
		fLl.WriteString(litExpr.Value + ")\n")
		return true
	} else if idExpr, b := call.Args[0].(*ast.Ident); b {
		// 生成打印变量的printf
		fLl.WriteString("  call i32 (i8*, ...) @printf(i8* %fmt, i64 %")
		fLl.WriteString(idExpr.Name + ")\n")
		return true
	} else {
		// 错误：不能打印表达式
		srcError = true
		fmt.Printf("源文件%s第%d行：print只能打印变量或常量\n", os.Args[1], lineNo)
		return false
	}
}
