package main

import (
	"fmt"
	"sync"

	"github.com/wa-lang/ssago/05-branch/wabuildin"
	"github.com/wa-lang/ssago/05-branch/waops"
	"github.com/wa-lang/ssago/05-branch/watypes"
	"golang.org/x/tools/go/ssa"
)

type Engine struct {
	main     *ssa.Package
	initOnce sync.Once

	// 全局变量
	globals map[string]*watypes.Value
}

func NewEngine(mainpkg *ssa.Package) *Engine {
	p := &Engine{
		main:    mainpkg,
		globals: make(map[string]*watypes.Value),
	}
	return p
}

// 读全局变量
func (p *Engine) getGlobal(key *ssa.Global) (v *watypes.Value, ok bool) {
	v, ok = p.globals[key.RelString(nil)]
	return
}

// 设置全局变量
func (p *Engine) setGlobal(key *ssa.Global, v *watypes.Value) {
	p.globals[key.RelString(nil)] = v
}

// 全局变量零值初始化
func (p *Engine) initGlobals() *Engine {
	p.initOnce.Do(func() {
		for _, pkg := range p.main.Prog.AllPackages() {
			for _, m := range pkg.Members {
				switch v := m.(type) {
				case *ssa.Global:
					cell := waops.Zero(waops.Deref(v.Type()))
					p.setGlobal(v, &cell)
				}
			}
		}
	})
	return p
}

type Frame struct {
	//局部变量、虚拟寄存器等：
	env map[ssa.Value]watypes.Value

	//返回值
	result watypes.Value

	//当前块
	block *ssa.BasicBlock

	//上一个块
	prevBlock *ssa.BasicBlock
}

func NewFrame() *Frame {
	f := &Frame{
		env: make(map[ssa.Value]watypes.Value),
	}
	return f
}

// 读取值(nil/全局变量/虚拟寄存器等)
func (p *Engine) getValue(fr *Frame, key ssa.Value) watypes.Value {
	switch key := key.(type) {
	case *ssa.Global:
		if r, ok := p.getGlobal(key); ok {
			return r
		}
	case *ssa.Const:
		return waops.ConstValue(key)
	case nil:
		return nil
	}

	if r, ok := fr.env[key]; ok {
		return r
	}

	panic(fmt.Sprintf("get: no value for %T: %v", key, key.Name()))
}

func (p *Engine) runFunc(fn watypes.Value, args []watypes.Value) watypes.Value {
	if fn, ok := fn.(*ssa.Builtin); ok {
		return callBuiltin(fn, args)
	}

	if fn, ok := fn.(*ssa.Function); ok {
		fr := NewFrame()
		fr.block = fn.Blocks[0]
		// 函数的参数添加到上下文环境
		for i, p := range fn.Params {
			fr.env[p] = args[i]
		}

		for fr.block != nil {
			p.runFrame(fr) // 核心逻辑
		}

		return fr.result
	}

	panic(fmt.Sprintf("Unknown function: %v", fn))
}

func (p *Engine) runFrame(fr *Frame) {
	for i := 0; i < len(fr.block.Instrs); i++ {
		switch ins := fr.block.Instrs[i].(type) {
		case *ssa.Store:
			watypes.Store(waops.Deref(ins.Addr.Type()), p.getValue(fr, ins.Addr).(*watypes.Value), p.getValue(fr, ins.Val))

		case *ssa.UnOp:
			fr.env[ins] = waops.UnOp(ins, p.getValue(fr, ins.X))

		case *ssa.BinOp:
			fr.env[ins] = waops.BinOp(ins.Op, ins.X.Type(), p.getValue(fr, ins.X), p.getValue(fr, ins.Y))

		case *ssa.Call:
			args := p.prepareCall(fr, &ins.Call)
			fr.env[ins] = p.runFunc(ins.Call.Value, args)

		case *ssa.Return:
			switch len(ins.Results) {
			case 0:
			case 1:
				fr.result = p.getValue(fr, ins.Results[0])
			default:
				panic("multi-return is not supported")
			}
			fr.block = nil
			return

		case *ssa.If:
			if p.getValue(fr, ins.Cond).(bool) {
				//println("if:true, goto block:", fr.block.Succs[0].String())
				fr.prevBlock, fr.block = fr.block, fr.block.Succs[0] // true
			} else {
				//println("if:false, goto block:", fr.block.Succs[1].String())
				fr.prevBlock, fr.block = fr.block, fr.block.Succs[1] // false
			}
			return

		case *ssa.Jump:
			//println("jump to block:", fr.block.Succs[0].String())
			fr.prevBlock, fr.block = fr.block, fr.block.Succs[0]
			return

		case *ssa.Phi:
			for i, pred := range ins.Block().Preds {
				if fr.prevBlock == pred {
					fr.env[ins] = p.getValue(fr, ins.Edges[i])
					break
				}
			}

		default:
			panic(fmt.Sprintf("Unknown instruction: %v", ins))
		}
	}

	fr.block = nil
}

func (p *Engine) prepareCall(fr *Frame, call *ssa.CallCommon) (args []watypes.Value) {
	// 普通函数调用
	if call.Method != nil {
		panic("method is not supported")
	}

	for _, arg := range call.Args {
		args = append(args, p.getValue(fr, arg))
	}
	return
}

func callBuiltin(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	switch fn.Name() {
	case "print", "println": // print(any, ...)
		return wabuiltin.Print(fn, args)
	}

	panic("unknown built-in: " + fn.Name())
}
