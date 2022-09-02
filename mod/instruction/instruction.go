package instruction

type InstructionName string

const (
	// Numeric instruction
	I32Const InstructionName = "i32.const"
	I32Add   InstructionName = "i32.add"
	I32Sub   InstructionName = "i32.sub"

	// Parametric instruction
	Drop InstructionName = "drop"

	// ControlInstruction
	Return InstructionName = "return"
)

func (name InstructionName) IsValid() bool {
	switch name {
	case I32Const, I32Add, I32Sub:
		return true
	case Drop:
		return true
	case Return:
		return true
	}

	return false
}

type Instruction interface {
	Name() InstructionName
}
