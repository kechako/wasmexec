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
	pos       int
	locals    map[int]Value
	idToIndex map[types.ID]int
	blocks    map[types.ID]*mod.Block

	original VMContext
}

var _ VMContext = (*FuncContext)(nil)

func newFuncContext(f *mod.Function, locals []Local, original VMContext) VMContext {
	localMap := make(map[int]Value)
	idToIndex := make(map[types.ID]int)
	for i, l := range locals {
		localMap[i] = l.Value
		if l.Index.IsID() {
			idToIndex[l.Index.ID] = i
		}
	}

	blockMap := make(map[types.ID]*mod.Block)
	for _, block := range f.Blocks {
		blockMap[block.Label] = block
	}

	return &FuncContext{
		f:         f,
		pos:       0,
		locals:    localMap,
		idToIndex: idToIndex,
		blocks:    blockMap,
		original:  original,
	}
}

func (funcCtx *FuncContext) NewFuncContext(f *mod.Function, locals []Local) VMContext {
	return newFuncContext(f, locals, funcCtx)
}

func (funcCtx *FuncContext) NewBlockContext(label types.ID) (VMContext, error) {
	return newBlockContext(label, funcCtx)
}

func (funcCtx *FuncContext) Results() []*mod.Result {
	return funcCtx.f.Results
}

func (funcCtx *FuncContext) Parameters() []*mod.Local {
	return funcCtx.f.Parameters
}

func (funcCtx *FuncContext) Original() VMContext {
	return funcCtx.original
}

func (funcCtx *FuncContext) GetBlock(label types.ID) (*mod.Block, bool) {
	block, ok := funcCtx.blocks[label]
	return block, ok
}

func (funcCtx *FuncContext) GetInstruction() instruction.Instruction {
	if funcCtx.pos >= len(funcCtx.f.Instructions) {
		return nil
	}

	i := funcCtx.f.Instructions[funcCtx.pos]
	funcCtx.pos++
	return i
}

func (funcCtx *FuncContext) SetLocal(idx types.Index, value Value) error {
	index, err := funcCtx.getIndex(idx)
	if err != nil {
		return err
	}

	funcCtx.locals[index] = value

	return nil
}

func (funcCtx *FuncContext) GetLocal(idx types.Index) (Value, error) {
	var value Value
	index, err := funcCtx.getIndex(idx)
	if err != nil {
		return value, err
	}

	value, ok := funcCtx.locals[index]
	if !ok {
		return value, errLocalVariableInconsistent
	}

	return value, nil
}

func (funcCtx *FuncContext) getIndex(idx types.Index) (int, error) {
	if idx.IsID() {
		i, ok := funcCtx.idToIndex[idx.ID]
		if !ok {
			return 0, errLocalVariableInconsistent
		}
		return i, nil
	}

	return idx.Index, nil
}
