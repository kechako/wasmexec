package runtime

import (
	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
	"github.com/kechako/wasmexec/mod/types"
)

type Local struct {
	Index types.Index
	Value Value
}

type FuncContext struct {
	// 関数
	f *mod.Function
	// 実行中の命令位置
	pos    int
	locals map[string]Value

	original *FuncContext
}

func newFuncContext(f *mod.Function, locals []Local, original *FuncContext) *FuncContext {
	localMap := make(map[string]Value)
	for _, l := range locals {
		key := makeIndexKey(l.Index)
		localMap[key] = l.Value
	}

	return &FuncContext{
		f:        f,
		pos:      0,
		locals:   localMap,
		original: original,
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

func (funcCtx *FuncContext) SetLocal(idx types.Index, value any) error {
	key := makeIndexKey(idx)
	if _, ok := funcCtx.locals[key]; !ok {
		return errLocalVariableInconsistent
	}

	funcCtx.locals[key] = NewValue(value)

	return nil
}

func (funcCtx *FuncContext) GetLocalInt32(idx types.Index) (int32, error) {
	key := makeIndexKey(idx)
	value, ok := funcCtx.locals[key]
	if !ok {
		return 0, errLocalVariableInconsistent
	}

	v, ok := value.Int32()
	if !ok {
		return 0, errLocalVariableInconsistent
	}

	return v, nil
}
