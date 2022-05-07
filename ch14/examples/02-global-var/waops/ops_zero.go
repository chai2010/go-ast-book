// 版权 @2019 凹语言 作者。保留所有权利。

package waops

import (
	"fmt"
	"go/types"
	"unsafe"

	"github.com/wa-lang/ssago/02-hello/watypes"
)

// 生成零值
func Zero(t types.Type) watypes.Value {
	return zero(t)
}

func zero(t types.Type) watypes.Value {
	switch t := t.(type) {
	case *types.Basic:
		if t.Kind() == types.UntypedNil {
			panic("untyped nil has no zero value")
		}
		if t.Info()&types.IsUntyped != 0 {
			t = types.Default(t).(*types.Basic)
		}
		switch t.Kind() {
		case types.Bool:
			return false
		case types.Int:
			return int(0)
		case types.Int8:
			return int8(0)
		case types.Int16:
			return int16(0)
		case types.Int32:
			return int32(0)
		case types.Int64:
			return int64(0)
		case types.Uint:
			return uint(0)
		case types.Uint8:
			return uint8(0)
		case types.Uint16:
			return uint16(0)
		case types.Uint32:
			return uint32(0)
		case types.Uint64:
			return uint64(0)
		case types.Uintptr:
			return uintptr(0)
		case types.Float32:
			return float32(0)
		case types.Float64:
			return float64(0)
		case types.Complex64:
			return complex64(0)
		case types.Complex128:
			return complex128(0)
		case types.String:
			return ""
		case types.UnsafePointer:
			return unsafe.Pointer(nil)
		default:
			panic(fmt.Sprint("zero for unexpected type:", t))
		}
	case *types.Pointer:
		return (*watypes.Value)(nil)
	}
	panic(fmt.Sprint("zero: unexpected ", t))
}
