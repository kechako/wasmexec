package runtime

import (
	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
)

type FuncContext struct {
	// 関数
	f *mod.Function
	// 実行中の命令位置
	pos int
}

func newFuncContext(f *mod.Function) *FuncContext {
	return &FuncContext{
		f:   f,
		pos: 0,
	}
}

func (funcCtx *FuncContext) GetInstruction() instruction.Instruction {
	if funcCtx.pos >= len(funcCtx.f.Instructions) {
		return nil
	}

	i := funcCtx.f.Instructions[funcCtx.pos]
	funcCtx.pos++
	return i
}
