package instruction

type I32Instruction string

const (
	I32Const I32Instruction = "i32.const"
	I32Add   I32Instruction = "i32.add"
	I32Sub   I32Instruction = "i32.sub"
)

func (i I32Instruction) IsValid() bool {
	switch i {
	case I32Const, I32Add, I32Sub:
		return true
	}

	return false
}

type I32 struct {
	Instruction I32Instruction
	Values      []int32
}
