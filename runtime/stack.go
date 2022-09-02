package runtime

import "container/list"

type ElementType byte

const (
	ValueElement ElementType = iota
	LabelElement
	ActivationElement
)

type Element struct {
	Type  ElementType
	Value Value
}

func newValueElement(v any) *Element {
	return &Element{
		Type:  ValueElement,
		Value: NewValue(v),
	}
}

func newActivationElement(funcCtx *FuncContext) *Element {
	return &Element{
		Type:  ActivationElement,
		Value: NewValue(funcCtx),
	}
}

func (elm *Element) Int32() (int32, bool) {
	return getElementValue[int32](elm, ValueElement)
}

func (elm *Element) FuncContext() (*FuncContext, bool) {
	return getElementValue[*FuncContext](elm, ActivationElement)
}

func getElementValue[T any](elm *Element, typ ElementType) (value T, ok bool) {
	if elm.Type != typ {
		return value, false
	}

	return GetValue[T](elm.Value)
}

type Stack struct {
	l   *list.List
	cap int
}

func NewStack(cap int) *Stack {
	return &Stack{
		l:   list.New(),
		cap: cap,
	}
}

func (s *Stack) Push(elm *Element) {
	if s.l.Len() == s.cap {
		panic("stack overflow")
	}

	s.l.PushBack(elm)
}

func (s *Stack) Pop() *Element {
	e := s.l.Back()
	if e == nil {
		panic("stack is empty")
	}
	s.l.Remove(e)

	return e.Value.(*Element)
}
