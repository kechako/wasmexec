package instruction

type ControlInstruction struct {
	Instruction InstructionName
	Values      []any
}

func (i *ControlInstruction) Name() InstructionName {
	return i.Instruction
}
