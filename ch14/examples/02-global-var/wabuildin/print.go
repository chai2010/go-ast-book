// 版权 @2019 凹语言 作者。保留所有权利。

package wabuiltin

import (
	"bytes"
	"os"

	"github.com/wa-lang/ssago/02-hello/watypes"
	"golang.org/x/tools/go/ssa"
)

func Print(fn *ssa.Builtin, args []watypes.Value) ssa.Value {
	ln := fn.Name() == "println"
	var buf bytes.Buffer

	for i, arg := range args {
		if i > 0 && ln {
			buf.WriteRune(' ')
		}
		buf.WriteString(watypes.ToString(arg))
	}
	if ln {
		buf.WriteRune('\n')
	}

	os.Stdout.Write(buf.Bytes())
	return nil
}
