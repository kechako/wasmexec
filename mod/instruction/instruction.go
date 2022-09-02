package instruction

type InstructionName string

const (
	// Numeric instruction
	I32Const InstructionName = "i32.const"
	I32Add   InstructionName = "i32.add"
	I32Sub   InstructionName = "i32.sub"
	I32Mul   InstructionName = "i32.mul"
	I32DivS  InstructionName = "i32.div_s"

	// Parametric instruction
	Drop InstructionName = "drop"

	// Variable Instructions
	LocalGet InstructionName = "local.get"
	LocalSet InstructionName = "local.set"
	LocalTee InstructionName = "local.tee"

	// ControlInstruction
	Return InstructionName = "return"
	Call   InstructionName = "call"
)

func (name InstructionName) IsValid() bool {
	return name.IsI32() || name.IsParametric() || name.IsControl()
}

func (name InstructionName) IsI32() bool {
	switch name {
	case I32Const, I32Add, I32Sub, I32Mul, I32DivS:
		return true
	}

	return false
}

func (name InstructionName) IsParametric() bool {
	switch name {
	case Drop:
		return true
	}

	return false
}

func (name InstructionName) IsVariable() bool {
	switch name {
	case LocalGet, LocalSet, LocalTee:
		return true
	}

	return false
}

func (name InstructionName) IsControl() bool {
	switch name {
	case Return, Call:
		return true
	}

	return false
}

type Instruction interface {
	Name() InstructionName
}
