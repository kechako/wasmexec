package mod

import (
	"regexp"

	"github.com/kechako/wasmexec/mod/instruction"
)

type Type string

const (
	Unkown Type = ""
	I32    Type = "i32"
	I64    Type = "i64"
	F32    Type = "f32"
	F64    Type = "f64"
)

type ID string

var regexpID = regexp.MustCompile("^\\$[0-9A-Za-z!#$%&'*+\\-,/:<=>?@\\\\^_`|~]+$")

func (id ID) IsValid() bool {
	return regexpID.MatchString(string(id))
}

func (id ID) IsEmpty() bool {
	return id == ""
}

type Index struct {
	Index int
	ID    ID
}

func (idx Index) IsIndex() bool {
	return !idx.IsID()
}

func (idx Index) IsID() bool {
	return !idx.ID.IsEmpty()
}

type Module struct {
	ID        ID
	Functions []*Function
	Exports   []*Export
}

type Function struct {
	ID           ID
	Parameters   []*Parameter
	Results      []*Result
	Instructions []instruction.Instruction
}

type Parameter struct {
	ID   ID
	Type Type
}

type Result struct {
	Type Type
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
	Index  Index
}
