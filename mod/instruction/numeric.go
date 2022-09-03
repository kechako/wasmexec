package instruction

type I32Instruction struct {
	Instruction InstructionName
	Values      []int32
}

func (i32 *I32Instruction) Name() InstructionName {
	return i32.Instruction
}
