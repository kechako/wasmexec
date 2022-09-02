package types

import "regexp"

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

func NewIndex(idx int) Index {
	return Index{
		Index: idx,
	}
}

func NewIndexWithID(id ID) Index {
	return Index{
		ID: id,
	}
}

func (idx Index) IsIndex() bool {
	return !idx.IsID()
}

func (idx Index) IsID() bool {
	return !idx.ID.IsEmpty()
}
