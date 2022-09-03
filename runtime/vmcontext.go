package runtime

import (
	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
	"github.com/kechako/wasmexec/mod/types"
)

type VMContext interface {
	NewFuncContext(f *mod.Function, locals []Local) VMContext
	NewBlockContext(label types.ID) (VMContext, error)
	Parameters() []*mod.Local
	Results() []*mod.Result
	Original() VMContext
	GetBlock(label types.ID) (*mod.Block, bool)
	GetInstruction() instruction.Instruction
	SetLocal(idx types.Index, value Value) error
	GetLocal(idx types.Index) (Value, error)
}

type Local struct {
	Index types.Index
	Value Value
}
