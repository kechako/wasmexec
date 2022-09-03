package mod

import (
	"github.com/kechako/wasmexec/mod/instruction"
	"github.com/kechako/wasmexec/mod/types"
)

type Module struct {
	ID        types.ID
	Functions []*Function
	Exports   []*Export
}

type Function struct {
	ID           types.ID
	Parameters   []*Local
	Results      []*Result
	Locals       []*Local
	Blocks       []*Block
	Instructions []instruction.Instruction
}

type Local struct {
	ID   types.ID
	Type types.Type
}

type Result struct {
	Type types.Type
}

type BlockResult Result

type Block struct {
	Label        types.ID
	Parameters   []*Local
	Results      []*Result
	Instructions []instruction.Instruction
}

type ExportTarget string

const (
	ExportFunction ExportTarget = "func"
	ExportTable    ExportTarget = "table"
	ExportMemory   ExportTarget = "memory"
	ExportGlobal   ExportTarget = "global"
)

type Export struct {
	Name   string
	Target ExportTarget
	Index  types.Index
}
