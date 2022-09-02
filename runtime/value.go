package runtime

import (
	"github.com/kechako/wasmexec/mod/types"
)

type Value struct {
	Value any
}

func NewValue(v any) Value {
	return Value{
		Value: v,
	}
}

func DefaultValue(typ types.Type) Value {
	switch typ {
	case types.I32:
		return NewValue(int32(0))
	case types.I64:
		return NewValue(int64(0))
	case types.F32:
		return NewValue(float32(0))
	case types.F64:
		return NewValue(float64(0))
	}
	panic("unsupported type")
}

func (value Value) Int32() (int32, bool) {
	return GetValue[int32](value)
}

func GetValue[T any](value Value) (v T, ok bool) {
	v, ok = value.Value.(T)
	if !ok {
		return v, false
	}

	return v, true
}
