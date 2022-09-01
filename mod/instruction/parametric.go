package instruction

type ParametricInstruction string

const (
	Drop ParametricInstruction = "drop"
)

func (i ParametricInstruction) IsValid() bool {
	switch i {
	case Drop:
		return true
	}

	return false
}

type Parametric struct {
	Instruction ParametricInstruction
}
