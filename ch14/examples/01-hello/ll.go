package main

import (
	"go/types"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	llvmTypes "github.com/llir/llvm/ir/types"
	"golang.org/x/tools/go/ssa"
)

func llModule(ssafnMain *ssa.Function) *ir.Module {
	m := ir.NewModule()

	// printf
	i8Ptr := llvmTypes.NewPointer(llvmTypes.I8)
	printf := m.NewFunc("printf", llvmTypes.I32, ir.NewParam("format", i8Ptr))
	printf.Sig.Variadic = true

	// global printf format string
	printf_format_int := m.NewGlobalDef("printf_format_int", constant.NewCharArray([]byte("%d\n\x00")))
	printf_format_int.Immutable = true

	// main
	fnMain := m.NewFunc("main", llvmTypes.I32)
	firstBlock := fnMain.NewBlock("entry")

	// main body
	for _, ins := range ssafnMain.Blocks[0].Instrs {
		switch ins := ins.(type) {
		case *ssa.Call:
			if ins.Call.Method == nil {
				if fnBuiltin, ok := ins.Call.Value.(*ssa.Builtin); ok {
					if fnBuiltin.Name() == "println" {
						llPrint(firstBlock, printf, printf_format_int, ins.Call.Args...)
					}
				}
			}
		default:
			break
		}
	}

	firstBlock.NewRet(constant.NewInt(llvmTypes.I32, 0))
	return m
}

func llPrint(block *ir.Block, printf *ir.Func, format *ir.Global, args ...ssa.Value) {
	for _, arg := range args {
		if arg, ok := arg.(*ssa.Const); ok {
			if t, ok := arg.Type().Underlying().(*types.Basic); ok {
				switch t.Kind() {
				case types.Int, types.UntypedInt:
					block.NewCall(
						printf,
						block.NewGetElementPtr(
							format.Type().(*llvmTypes.PointerType).ElemType, format,
							constant.NewInt(llvmTypes.I32, 0),
							constant.NewInt(llvmTypes.I32, 0),
						),
						constant.NewInt(llvmTypes.I32, arg.Int64()),
					)
				}
			}
		}
	}
}
