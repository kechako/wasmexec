package instruction

type I32Instruction string

type I32 struct {
	Instruction InstructionName
	Values      []int32
}

func (i32 *I32) Name() InstructionName {
	return i32.Instruction
}
