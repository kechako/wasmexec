package runtime

import (
	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
	"github.com/kechako/wasmexec/mod/types"
)

type BlockContext struct {
	block    *mod.Block
	pos      int
	original VMContext
}

var _ VMContext = (*BlockContext)(nil)

func newBlockContext(label types.ID, original VMContext) (VMContext, error) {
	block, ok := original.GetBlock(label)
	if !ok {
		return nil, errBlockNotFound
	}

	return &BlockContext{
		block:    block,
		pos:      0,
		original: original,
	}, nil
}

func (blockCtx *BlockContext) NewFuncContext(f *mod.Function, locals []Local) VMContext {
	return blockCtx.original.NewFuncContext(f, locals)
}

func (blockCtx *BlockContext) NewBlockContext(label types.ID) (VMContext, error) {
	return newBlockContext(label, blockCtx)
}

func (blockCtx *BlockContext) Results() []*mod.Result {
	return blockCtx.block.Results
}

func (blockCtx *BlockContext) Parameters() []*mod.Local {
	return blockCtx.block.Parameters
}

func (blockCtx *BlockContext) Original() VMContext {
	return blockCtx.original
}

func (blockCtx *BlockContext) GetBlock(label types.ID) (*mod.Block, bool) {
	return blockCtx.original.GetBlock(label)
}

func (blockCtx *BlockContext) GetInstruction() instruction.Instruction {
	if blockCtx.pos >= len(blockCtx.block.Instructions) {
		return nil
	}

	i := blockCtx.block.Instructions[blockCtx.pos]
	blockCtx.pos++
	return i
}

func (blockCtx *BlockContext) SetLocal(idx types.Index, value Value) error {
	return blockCtx.original.SetLocal(idx, value)
}

func (blockCtx *BlockContext) GetLocal(idx types.Index) (Value, error) {
	return blockCtx.original.GetLocal(idx)
}
