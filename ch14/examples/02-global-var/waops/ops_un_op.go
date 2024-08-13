// 版权 @2019 凹语言 作者。保留所有权利。

package waops

import (
	"fmt"
	"go/token"

	"github.com/wa-lang/ssago/02-hello/watypes"
	"golang.org/x/tools/go/ssa"
)

// 一元运算符
func UnOp(instr *ssa.UnOp, x watypes.Value) watypes.Value {
	return unop(instr, x)
}

func unop(instr *ssa.UnOp, x watypes.Value) watypes.Value {
	switch instr.Op {
	case token.MUL: // 指针类型
		return watypes.Load(Deref(instr.X.Type()), x.(*watypes.Value))

	case token.NOT: // 非
		return !x.(bool)

	case token.SUB:
		switch x := x.(type) {
		case int:
			return -x
		case int8:
			return -x
		case int16:
			return -x
		case int32:
			return -x
		case int64:
			return -x
		case uint:
			return -x
		case uint8:
			return -x
		case uint16:
			return -x
		case uint32:
			return -x
		case uint64:
			return -x
		case uintptr:
			return -x
		case float32:
			return -x
		case float64:
			return -x
		case complex64:
			return -x
		case complex128:
			return -x
		}
	case token.XOR:
		switch x := x.(type) {
		case int:
			return ^x
		case int8:
			return ^x
		case int16:
			return ^x
		case int32:
			return ^x
		case int64:
			return ^x
		case uint:
			return ^x
		case uint8:
			return ^x
		case uint16:
			return ^x
		case uint32:
			return ^x
		case uint64:
			return ^x
		case uintptr:
			return ^x
		}
	}
	panic(fmt.Sprintf("invalid unary op %s %T", instr.Op, x))
}
