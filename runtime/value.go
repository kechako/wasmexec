package runtime

type Value struct {
	Value any
}

func NewValue(v any) Value {
	return Value{
		Value: v,
	}
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
