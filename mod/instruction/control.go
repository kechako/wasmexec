package instruction

type Control struct {
	Instruction InstructionName
}

func (c *Control) Name() InstructionName {
	return c.Instruction
}
