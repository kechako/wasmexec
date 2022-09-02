package instruction

import "github.com/kechako/wasmexec/mod/types"

type Variable struct {
	Instruction InstructionName
	Index       types.Index
}

func (v *Variable) Name() InstructionName {
	return v.Instruction
}
