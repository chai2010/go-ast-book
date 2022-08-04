package main

import (
	"fmt"
	"sync"

	"github.com/wa-lang/ssago/02-hello/wabuildin"
	"github.com/wa-lang/ssago/02-hello/waops"
	"github.com/wa-lang/ssago/02-hello/watypes"
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

func (p *Engine) runFunc(fr *Frame, fn *ssa.Function) {
	fmt.Println("--- runFunc begin ---")
	defer fmt.Println("--- runFunc end   ---")

	if len(fn.Blocks) > 0 {
		for blk := fn.Blocks[0]; blk != nil; {
			blk = p.runFuncBlock(fr, fn.Blocks[0])
		}
	}
}

// 运行Block
func (p *Engine) runFuncBlock(fr *Frame, block *ssa.BasicBlock) (nextBlock *ssa.BasicBlock) {
	for _, ins := range block.Instrs {
		switch ins := ins.(type) {
		case *ssa.Store:
			println("ssa.Store")
			watypes.Store(waops.Deref(ins.Addr.Type()), p.getValue(fr, ins.Addr).(*watypes.Value), p.getValue(fr, ins.Val))

		case *ssa.UnOp:
			println("ssa.UnOp")
			fr.env[ins] = waops.UnOp(ins, p.getValue(fr, ins.X))

		case *ssa.Call:
			println("ssa.Call")
			args := p.prepareCall(fr, &ins.Call)
			fr.env[ins] = p.call(ins, args)
		}
	}
	return nil
}

func (p *Engine) prepareCall(fr *Frame, call *ssa.CallCommon) (args []watypes.Value) {
	// 转换参数, getValue 是核心方法
	for _, arg := range call.Args {
		args = append(args, p.getValue(fr, arg))
	}

	return
}

func (p *Engine) call(ins *ssa.Call, args []watypes.Value) watypes.Value {
	switch {
	case ins.Call.Method == nil: // 普通函数调用
		switch callFn := ins.Call.Value.(type) {
		case *ssa.Builtin:
			return callBuiltin(callFn, args)
		}
	}

	panic("Unknown call")
}

func callBuiltin(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	switch fn.Name() {
	case "print", "println": // print(any, ...)
		return wabuiltin.Print(fn, args)
	}

	panic("unknown built-in: " + fn.Name())
}
