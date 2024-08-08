// 版权 @2019 凹语言 作者。保留所有权利。

package waops

import (
	"fmt"
	"go/token"
	"go/types"

	"github.com/wa-lang/ssago/04-bi-op/watypes"
)

// 实现二元运算
func BinOp(op token.Token, t types.Type, x, y watypes.Value) watypes.Value {
	return binop(op, t, x, y)
}

func binop(op token.Token, t types.Type, x, y watypes.Value) watypes.Value {
	switch op {
	case token.EQL:
		return watypes.Equals(t, x, y)
	case token.NEQ:
		return !watypes.Equals(t, x, y)

	case token.ADD:
		switch x.(type) {
		case int:
			return x.(int) + y.(int)
		case int8:
			return x.(int8) + y.(int8)
		case int16:
			return x.(int16) + y.(int16)
		case int32:
			return x.(int32) + y.(int32)
		case int64:
			return x.(int64) + y.(int64)
		case uint:
			return x.(uint) + y.(uint)
		case uint8:
			return x.(uint8) + y.(uint8)
		case uint16:
			return x.(uint16) + y.(uint16)
		case uint32:
			return x.(uint32) + y.(uint32)
		case uint64:
			return x.(uint64) + y.(uint64)
		case uintptr:
			return x.(uintptr) + y.(uintptr)
		case float32:
			return x.(float32) + y.(float32)
		case float64:
			return x.(float64) + y.(float64)
		case complex64:
			return x.(complex64) + y.(complex64)
		case complex128:
			return x.(complex128) + y.(complex128)
		case string:
			return x.(string) + y.(string)
		}

	case token.SUB:
		switch x.(type) {
		case int:
			return x.(int) - y.(int)
		case int8:
			return x.(int8) - y.(int8)
		case int16:
			return x.(int16) - y.(int16)
		case int32:
			return x.(int32) - y.(int32)
		case int64:
			return x.(int64) - y.(int64)
		case uint:
			return x.(uint) - y.(uint)
		case uint8:
			return x.(uint8) - y.(uint8)
		case uint16:
			return x.(uint16) - y.(uint16)
		case uint32:
			return x.(uint32) - y.(uint32)
		case uint64:
			return x.(uint64) - y.(uint64)
		case uintptr:
			return x.(uintptr) - y.(uintptr)
		case float32:
			return x.(float32) - y.(float32)
		case float64:
			return x.(float64) - y.(float64)
		case complex64:
			return x.(complex64) - y.(complex64)
		case complex128:
			return x.(complex128) - y.(complex128)
		}

	case token.MUL:
		switch x.(type) {
		case int:
			return x.(int) * y.(int)
		case int8:
			return x.(int8) * y.(int8)
		case int16:
			return x.(int16) * y.(int16)
		case int32:
			return x.(int32) * y.(int32)
		case int64:
			return x.(int64) * y.(int64)
		case uint:
			return x.(uint) * y.(uint)
		case uint8:
			return x.(uint8) * y.(uint8)
		case uint16:
			return x.(uint16) * y.(uint16)
		case uint32:
			return x.(uint32) * y.(uint32)
		case uint64:
			return x.(uint64) * y.(uint64)
		case uintptr:
			return x.(uintptr) * y.(uintptr)
		case float32:
			return x.(float32) * y.(float32)
		case float64:
			return x.(float64) * y.(float64)
		case complex64:
			return x.(complex64) * y.(complex64)
		case complex128:
			return x.(complex128) * y.(complex128)
		}

	case token.QUO:
		switch x.(type) {
		case int:
			return x.(int) / y.(int)
		case int8:
			return x.(int8) / y.(int8)
		case int16:
			return x.(int16) / y.(int16)
		case int32:
			return x.(int32) / y.(int32)
		case int64:
			return x.(int64) / y.(int64)
		case uint:
			return x.(uint) / y.(uint)
		case uint8:
			return x.(uint8) / y.(uint8)
		case uint16:
			return x.(uint16) / y.(uint16)
		case uint32:
			return x.(uint32) / y.(uint32)
		case uint64:
			return x.(uint64) / y.(uint64)
		case uintptr:
			return x.(uintptr) / y.(uintptr)
		case float32:
			return x.(float32) / y.(float32)
		case float64:
			return x.(float64) / y.(float64)
		case complex64:
			return x.(complex64) / y.(complex64)
		case complex128:
			return x.(complex128) / y.(complex128)
		}

	case token.REM:
		switch x.(type) {
		case int:
			return x.(int) % y.(int)
		case int8:
			return x.(int8) % y.(int8)
		case int16:
			return x.(int16) % y.(int16)
		case int32:
			return x.(int32) % y.(int32)
		case int64:
			return x.(int64) % y.(int64)
		case uint:
			return x.(uint) % y.(uint)
		case uint8:
			return x.(uint8) % y.(uint8)
		case uint16:
			return x.(uint16) % y.(uint16)
		case uint32:
			return x.(uint32) % y.(uint32)
		case uint64:
			return x.(uint64) % y.(uint64)
		case uintptr:
			return x.(uintptr) % y.(uintptr)
		}

	case token.AND:
		switch x.(type) {
		case int:
			return x.(int) & y.(int)
		case int8:
			return x.(int8) & y.(int8)
		case int16:
			return x.(int16) & y.(int16)
		case int32:
			return x.(int32) & y.(int32)
		case int64:
			return x.(int64) & y.(int64)
		case uint:
			return x.(uint) & y.(uint)
		case uint8:
			return x.(uint8) & y.(uint8)
		case uint16:
			return x.(uint16) & y.(uint16)
		case uint32:
			return x.(uint32) & y.(uint32)
		case uint64:
			return x.(uint64) & y.(uint64)
		case uintptr:
			return x.(uintptr) & y.(uintptr)
		}

	case token.OR:
		switch x.(type) {
		case int:
			return x.(int) | y.(int)
		case int8:
			return x.(int8) | y.(int8)
		case int16:
			return x.(int16) | y.(int16)
		case int32:
			return x.(int32) | y.(int32)
		case int64:
			return x.(int64) | y.(int64)
		case uint:
			return x.(uint) | y.(uint)
		case uint8:
			return x.(uint8) | y.(uint8)
		case uint16:
			return x.(uint16) | y.(uint16)
		case uint32:
			return x.(uint32) | y.(uint32)
		case uint64:
			return x.(uint64) | y.(uint64)
		case uintptr:
			return x.(uintptr) | y.(uintptr)
		}

	case token.XOR:
		switch x.(type) {
		case int:
			return x.(int) ^ y.(int)
		case int8:
			return x.(int8) ^ y.(int8)
		case int16:
			return x.(int16) ^ y.(int16)
		case int32:
			return x.(int32) ^ y.(int32)
		case int64:
			return x.(int64) ^ y.(int64)
		case uint:
			return x.(uint) ^ y.(uint)
		case uint8:
			return x.(uint8) ^ y.(uint8)
		case uint16:
			return x.(uint16) ^ y.(uint16)
		case uint32:
			return x.(uint32) ^ y.(uint32)
		case uint64:
			return x.(uint64) ^ y.(uint64)
		case uintptr:
			return x.(uintptr) ^ y.(uintptr)
		}

	case token.AND_NOT:
		switch x.(type) {
		case int:
			return x.(int) &^ y.(int)
		case int8:
			return x.(int8) &^ y.(int8)
		case int16:
			return x.(int16) &^ y.(int16)
		case int32:
			return x.(int32) &^ y.(int32)
		case int64:
			return x.(int64) &^ y.(int64)
		case uint:
			return x.(uint) &^ y.(uint)
		case uint8:
			return x.(uint8) &^ y.(uint8)
		case uint16:
			return x.(uint16) &^ y.(uint16)
		case uint32:
			return x.(uint32) &^ y.(uint32)
		case uint64:
			return x.(uint64) &^ y.(uint64)
		case uintptr:
			return x.(uintptr) &^ y.(uintptr)
		}

	case token.LSS:
		switch x.(type) {
		case int:
			return x.(int) < y.(int)
		case int8:
			return x.(int8) < y.(int8)
		case int16:
			return x.(int16) < y.(int16)
		case int32:
			return x.(int32) < y.(int32)
		case int64:
			return x.(int64) < y.(int64)
		case uint:
			return x.(uint) < y.(uint)
		case uint8:
			return x.(uint8) < y.(uint8)
		case uint16:
			return x.(uint16) < y.(uint16)
		case uint32:
			return x.(uint32) < y.(uint32)
		case uint64:
			return x.(uint64) < y.(uint64)
		case uintptr:
			return x.(uintptr) < y.(uintptr)
		case float32:
			return x.(float32) < y.(float32)
		case float64:
			return x.(float64) < y.(float64)
		case string:
			return x.(string) < y.(string)
		}

	case token.LEQ:
		switch x.(type) {
		case int:
			return x.(int) <= y.(int)
		case int8:
			return x.(int8) <= y.(int8)
		case int16:
			return x.(int16) <= y.(int16)
		case int32:
			return x.(int32) <= y.(int32)
		case int64:
			return x.(int64) <= y.(int64)
		case uint:
			return x.(uint) <= y.(uint)
		case uint8:
			return x.(uint8) <= y.(uint8)
		case uint16:
			return x.(uint16) <= y.(uint16)
		case uint32:
			return x.(uint32) <= y.(uint32)
		case uint64:
			return x.(uint64) <= y.(uint64)
		case uintptr:
			return x.(uintptr) <= y.(uintptr)
		case float32:
			return x.(float32) <= y.(float32)
		case float64:
			return x.(float64) <= y.(float64)
		case string:
			return x.(string) <= y.(string)
		}

	case token.GTR:
		switch x.(type) {
		case int:
			return x.(int) > y.(int)
		case int8:
			return x.(int8) > y.(int8)
		case int16:
			return x.(int16) > y.(int16)
		case int32:
			return x.(int32) > y.(int32)
		case int64:
			return x.(int64) > y.(int64)
		case uint:
			return x.(uint) > y.(uint)
		case uint8:
			return x.(uint8) > y.(uint8)
		case uint16:
			return x.(uint16) > y.(uint16)
		case uint32:
			return x.(uint32) > y.(uint32)
		case uint64:
			return x.(uint64) > y.(uint64)
		case uintptr:
			return x.(uintptr) > y.(uintptr)
		case float32:
			return x.(float32) > y.(float32)
		case float64:
			return x.(float64) > y.(float64)
		case string:
			return x.(string) > y.(string)
		}

	case token.GEQ:
		switch x.(type) {
		case int:
			return x.(int) >= y.(int)
		case int8:
			return x.(int8) >= y.(int8)
		case int16:
			return x.(int16) >= y.(int16)
		case int32:
			return x.(int32) >= y.(int32)
		case int64:
			return x.(int64) >= y.(int64)
		case uint:
			return x.(uint) >= y.(uint)
		case uint8:
			return x.(uint8) >= y.(uint8)
		case uint16:
			return x.(uint16) >= y.(uint16)
		case uint32:
			return x.(uint32) >= y.(uint32)
		case uint64:
			return x.(uint64) >= y.(uint64)
		case uintptr:
			return x.(uintptr) >= y.(uintptr)
		case float32:
			return x.(float32) >= y.(float32)
		case float64:
			return x.(float64) >= y.(float64)
		case string:
			return x.(string) >= y.(string)
		}
	}
	panic(fmt.Sprintf("invalid binary op: %T %s %T", x, op, y))
}
