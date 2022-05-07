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
