package instruction

type Parametric struct {
	Instruction InstructionName
}

func (p *Parametric) Name() InstructionName {
	return p.Instruction
}
