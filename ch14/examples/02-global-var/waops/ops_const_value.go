// 版权 @2019 凹语言 作者。保留所有权利。

package waops

import (
	"fmt"
	"go/constant"
	"go/types"

	"github.com/wa-lang/ssago/02-hello/watypes"
	"golang.org/x/tools/go/ssa"
)

// 从常量构造值
func ConstValue(c *ssa.Const) watypes.Value {
	return constValue(c)
}

func constValue(c *ssa.Const) watypes.Value {
	if c.IsNil() {
		return Zero(c.Type()) // typed nil
	}

	if t, ok := c.Type().Underlying().(*types.Basic); ok {
		switch t.Kind() {
		case types.String, types.UntypedString:
			if c.Value.Kind() == constant.String {
				return constant.StringVal(c.Value)
			}
			return string(rune(c.Int64()))
		case types.Bool, types.UntypedBool:
			return constant.BoolVal(c.Value)
		case types.Int, types.UntypedInt:
			// Assume sizeof(int) is same on host and target.
			return int(c.Int64())
		case types.Int8:
			return int8(c.Int64())
		case types.Int16:
			return int16(c.Int64())
		case types.Int32, types.UntypedRune:
			return int32(c.Int64())
		case types.Int64:
			return c.Int64()
		case types.Uint:
			// Assume sizeof(uint) is same on host and target.
			return uint(c.Uint64())
		case types.Uint8:
			return uint8(c.Uint64())
		case types.Uint16:
			return uint16(c.Uint64())
		case types.Uint32:
			return uint32(c.Uint64())
		case types.Uint64:
			return c.Uint64()
		case types.Uintptr:
			// Assume sizeof(uintptr) is same on host and target.
			return uintptr(c.Uint64())
		case types.Float32:
			return float32(c.Float64())
		case types.Float64, types.UntypedFloat:
			return c.Float64()
		case types.Complex64:
			return complex64(c.Complex128())
		case types.Complex128, types.UntypedComplex:
			return c.Complex128()
		}
	}

	panic(fmt.Sprintf("constValue: %s", c))
}
