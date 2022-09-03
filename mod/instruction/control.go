package instruction

import "github.com/kechako/wasmexec/mod/types"

type ControlInstruction struct {
	Instruction InstructionName
}

func (i *ControlInstruction) Name() InstructionName {
	return i.Instruction
}

type CallInstruction struct {
	Instruction InstructionName
	Index       types.Index
}

func (i *CallInstruction) Name() InstructionName {
	return i.Instruction
}

type BlockInstruction struct {
	Instruction InstructionName
	Label       types.ID
}

func (i *BlockInstruction) Name() InstructionName {
	return i.Instruction
}
