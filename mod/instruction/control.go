package instruction

type ControlInstruction string

const (
	Return ControlInstruction = "return"
)

func (i ControlInstruction) IsValid() bool {
	switch i {
	case Return:
		return true
	}

	return false
}

type Control struct {
	Instruction ControlInstruction
}
