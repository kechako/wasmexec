package instruction

import "github.com/kechako/wasmexec/mod/types"

type VariableInstruction struct {
	Instruction InstructionName
	Index       types.Index
}

func (i *VariableInstruction) Name() InstructionName {
	return i.Instruction
}
