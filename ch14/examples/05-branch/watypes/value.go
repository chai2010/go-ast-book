package watypes

import (
	"bytes"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/ssa"
)

// Wa支持的值类型的抽象接口
type Value interface{}

// 返回地址addr处存储的T类型的值
func Load(T types.Type, addr *Value) Value {
	return *addr
}

// 将类型为T的值v存入地址addr中
func Store(T types.Type, addr *Value, v Value) {
	*addr = v
}

// ToString 输出Value的可读字符串
func ToString(v Value) string {
	var b bytes.Buffer
	writeValue(&b, v)
	return b.String()
}

func writeValue(buf *bytes.Buffer, v Value) {
	switch v := v.(type) {
	case nil, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, complex64, complex128, string:
		fmt.Fprintf(buf, "%v", v)

	case *Value:
		if v == nil {
			buf.WriteString("<nil>")
		} else {
			fmt.Fprintf(buf, "%p", v)
		}

	case *ssa.Function, *ssa.Builtin:
		fmt.Fprintf(buf, "%p", v) // (an address)

	default:
		fmt.Fprintf(buf, "<%T>", v)
	}
}

func Equals(t types.Type, x, y Value) bool {
	switch x := x.(type) {
	case bool:
		return x == y.(bool)
	case int:
		return x == y.(int)
	case int8:
		return x == y.(int8)
	case int16:
		return x == y.(int16)
	case int32:
		return x == y.(int32)
	case int64:
		return x == y.(int64)
	case uint:
		return x == y.(uint)
	case uint8:
		return x == y.(uint8)
	case uint16:
		return x == y.(uint16)
	case uint32:
		return x == y.(uint32)
	case uint64:
		return x == y.(uint64)
	case uintptr:
		return x == y.(uintptr)
	case float32:
		return x == y.(float32)
	case float64:
		return x == y.(float64)
	case complex64:
		return x == y.(complex64)
	case complex128:
		return x == y.(complex128)
	case string:
		return x == y.(string)
	case *Value:
		return x == y.(*Value)
	}

	panic(fmt.Sprintf("comparing uncomparable type %s", t))
}
