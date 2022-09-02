package runtime

import (
	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
	"github.com/kechako/wasmexec/mod/types"
)

type FuncContext struct {
	// 関数
	f *mod.Function
	// 実行中の命令位置
	pos    int
	locals map[string]Value
}

func newFuncContext(f *mod.Function) *FuncContext {
	return &FuncContext{
		f:      f,
		pos:    0,
		locals: make(map[string]Value),
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

func (funcCtx *FuncContext) AddLocal(idx types.Index, value any) {
	key := makeIndexKey(idx)
	funcCtx.locals[key] = NewValue(value)
}

func (funcCtx *FuncContext) GetLocalInt32(idx types.Index) (int32, bool) {
	key := makeIndexKey(idx)
	value, ok := funcCtx.locals[key]
	if !ok {
		return 0, false
	}

	return value.Int32()
}
