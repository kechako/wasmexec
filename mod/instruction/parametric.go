package instruction

type ParametricInstruction struct {
	Instruction InstructionName
}

func (i *ParametricInstruction) Name() InstructionName {
	return i.Instruction
}
