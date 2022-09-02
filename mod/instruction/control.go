package instruction

type Control struct {
	Instruction InstructionName
	Values      []any
}

func (c *Control) Name() InstructionName {
	return c.Instruction
}
