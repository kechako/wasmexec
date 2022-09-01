package mod

import (
	"errors"
)

var ErrInvalidFormat = errors.New("format is not valid")

type Decoder interface {
	Decode() (*Module, error)
}
