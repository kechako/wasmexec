package runtime

import (
	"context"
	"errors"
	"strconv"

	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
	"github.com/kechako/wasmexec/mod/types"
)

type VM struct {
	mod   *mod.Module
	stack *Stack

	funcs   map[string]*mod.Function
	exports map[string]*mod.Export
}

func New(m *mod.Module, opts ...Option) *VM {
	vmOpts := vmOptions{
		stackCapacity: 1024,
	}
	for _, opt := range opts {
		opt.apply(&vmOpts)
	}

	vm := &VM{
		mod:     m,
		stack:   NewStack(vmOpts.stackCapacity),
		funcs:   make(map[string]*mod.Function),
		exports: make(map[string]*mod.Export),
	}
	vm.init()

	return vm
}

func (vm *VM) init() {
	vm.makeFuncTable()
	vm.makeExportTable()
}

func (vm *VM) makeFuncTable() {
	for i, f := range vm.mod.Functions {
		idx := strconv.Itoa(i)
		vm.funcs[idx] = f
		if !f.ID.IsEmpty() {
			vm.funcs[string(f.ID)] = f
		}
	}
}

func (vm *VM) makeExportTable() {
	for _, e := range vm.mod.Exports {
		vm.exports[e.Name] = e
	}
}

var (
	errExportNotFound            = errors.New("export is not found")
	errExportTargetNotFunction   = errors.New("export target is not a function")
	errFunctionNotFound          = errors.New("function is not found")
	errStackInconsistent         = errors.New("stack is inconsistent")
	errLocalVariableInconsistent = errors.New("local variables are inconsistent")
	errIntegerDivideByZero       = errors.New("integer divide by zero")
	errUnsupportedType           = errors.New("unsupported type")
)

func (vm *VM) ExecFunc(ctx context.Context, name string) ([]any, error) {
	// エクスポートを検索
	e, ok := vm.exports[name]
	if !ok {
		return nil, errExportNotFound
	}

	// エクスポートは関数？
	if e.Target != mod.ExportFunction {
		return nil, errExportTargetNotFunction
	}

	// 関数を検索
	f, ok := vm.funcs[makeIndexKey(e.Index)]
	if !ok {
		return nil, errFunctionNotFound
	}

	err := vm.callFunc(ctx, f)
	if err != nil {
		return nil, err
	}

	results, err := vm.popFuncResults(f)
	if err != nil {
		return nil, err
	}

	return results, err
}

func (vm *VM) callFunc(ctx context.Context, f *mod.Function) error {
	funcCtx, err := vm.initFunction(f, nil)
	if err != nil {
		return err
	}

loop:
	for {
		i := funcCtx.GetInstruction()
		if i == nil {
			var err error
			funcCtx, err = vm.finalizeFunction(funcCtx.f)
			if err != nil {
				return err
			}

			if funcCtx == nil {
				break loop
			}

			continue
		}

		switch i.Name() {
		case instruction.I32Const:
			i := i.(*instruction.I32)
			vm.stack.Push(newValueElement(i.Values[0]))
		case instruction.I32Add:
			c2, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			vm.stack.Push(newValueElement(c1 + c2))
		case instruction.I32Sub:
			c2, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			vm.stack.Push(newValueElement(c1 - c2))
		case instruction.I32Mul:
			c2, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			vm.stack.Push(newValueElement(c1 * c2))
		case instruction.I32DivS:
			c2, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			if c2 == 0 {
				return errIntegerDivideByZero
			}
			vm.stack.Push(newValueElement(c1 / c2))
		case instruction.I32Eqz:
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			var b int32
			if c1 == 0 {
				b = 1
			}
			vm.stack.Push(newValueElement(b))
		case instruction.I32Eq:
			c2, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			var b int32
			if c1 == c2 {
				b = 1
			}
			vm.stack.Push(newValueElement(b))
		case instruction.I32Ne:
			c2, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			var b int32
			if c1 != c2 {
				b = 1
			}
			vm.stack.Push(newValueElement(b))
		case instruction.I32LtS:
			c2, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			var b int32
			if c1 < c2 {
				b = 1
			}
			vm.stack.Push(newValueElement(b))
		case instruction.I32GtS:
			c2, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			var b int32
			if c1 > c2 {
				b = 1
			}
			vm.stack.Push(newValueElement(b))
		case instruction.I32LeS:
			c2, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			var b int32
			if c1 <= c2 {
				b = 1
			}
			vm.stack.Push(newValueElement(b))
		case instruction.I32GeS:
			c2, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32()
			if !ok {
				return errStackInconsistent
			}
			var b int32
			if c1 >= c2 {
				b = 1
			}
			vm.stack.Push(newValueElement(b))
		case instruction.Drop:
			elm := vm.stack.Pop()
			if elm.Type != ValueElement {
				return errStackInconsistent
			}
		case instruction.LocalGet:
			i := i.(*instruction.Variable)
			v, err := funcCtx.GetLocalInt32(i.Index)
			if err != nil {
				return err
			}
			vm.stack.Push(newValueElement(v))
		case instruction.LocalSet:
			i := i.(*instruction.Variable)
			err := execLocalSet(vm, funcCtx, i.Index)
			if err != nil {
				return err
			}
		case instruction.LocalTee:
			i := i.(*instruction.Variable)
			elm := vm.stack.Pop()
			if elm.Type != ValueElement {
				return errStackInconsistent
			}
			vm.stack.Push(newValueElement(elm.Value.Value))
			vm.stack.Push(newValueElement(elm.Value.Value))

			err := execLocalSet(vm, funcCtx, i.Index)
			if err != nil {
				return err
			}
		case instruction.Return:
			var err error
			funcCtx, err = vm.finalizeFunction(funcCtx.f)
			if err != nil {
				return err
			}
			if funcCtx == nil {
				break loop
			}
		case instruction.Call:
			i := i.(*instruction.Control)
			index := i.Values[0].(types.Index)
			key := makeIndexKey(index)
			f, ok := vm.funcs[key]
			if !ok {
				return errFunctionNotFound
			}

			var err error
			funcCtx, err = vm.initFunction(f, funcCtx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) initFunction(f *mod.Function, original *FuncContext) (*FuncContext, error) {
	var locals []Local

	// parameters
	if len(f.Parameters) > 0 {
		for i := len(f.Parameters) - 1; i >= 0; i-- {
			p := f.Parameters[i]
			switch p.Type {
			case types.I32:
				v, ok := vm.stack.Pop().Int32()
				if !ok {
					return nil, errStackInconsistent
				}
				idx := types.NewIndex(i)
				if !p.ID.IsEmpty() {
					idx = types.NewIndexWithID(p.ID)
				}
				locals = append(locals, Local{
					Index: idx,
					Value: NewValue(v),
				})
			default:
				return nil, errUnsupportedType
			}
		}
	}

	localStart := len(locals)
	for i, l := range f.Locals {
		idx := types.NewIndex(localStart + i)
		if !l.ID.IsEmpty() {
			idx = types.NewIndexWithID(l.ID)
		}
		locals = append(locals, Local{
			Index: idx,
			Value: DefaultValue(l.Type),
		})
	}

	funcCtx := newFuncContext(f, locals, original)

	vm.stack.Push(newActivationElement(funcCtx))

	return funcCtx, nil
}

func (vm *VM) finalizeFunction(f *mod.Function) (*FuncContext, error) {

	results, err := vm.popFuncResults(f)
	if err != nil {
		return nil, err
	}

	// pop func context
	funcCtx, ok := vm.stack.Pop().FuncContext()
	if !ok {
		return nil, errStackInconsistent
	}

	for _, result := range results {
		vm.stack.Push(newValueElement(result))
	}

	return funcCtx.original, nil
}

func (vm *VM) popFuncResults(f *mod.Function) ([]any, error) {
	var results []any
	for _, r := range f.Results {
		// pop result
		elm := vm.stack.Pop()
		if elm.Type != ValueElement {
			return nil, errStackInconsistent
		}

		result := elm.Value.Value

		switch result.(type) {
		case int32:
			if r.Type != types.I32 {
				return nil, errStackInconsistent
			}
		case int64:
			if r.Type != types.I64 {
				return nil, errStackInconsistent
			}
		case float32:
			if r.Type != types.F32 {
				return nil, errStackInconsistent
			}
		case float64:
			if r.Type != types.F64 {
				return nil, errStackInconsistent
			}
		}

		results = append(results, result)
	}

	return results, nil
}

func execLocalSet(vm *VM, funcCtx *FuncContext, index types.Index) error {
	elm := vm.stack.Pop()
	if elm.Type != ValueElement {
		return errStackInconsistent
	}

	return funcCtx.SetLocal(index, elm.Value.Value)
}

func makeIndexKey(idx types.Index) string {
	if idx.IsIndex() {
		return strconv.Itoa(idx.Index)
	}

	return string(idx.ID)
}

type vmOptions struct {
	stackCapacity int
}

type Option interface {
	apply(opts *vmOptions)
}

type optionFunc func(opts *vmOptions)

func (f optionFunc) apply(opts *vmOptions) {
	f(opts)
}

func StackCapacity(cap int) Option {
	return optionFunc(func(opts *vmOptions) {
		opts.stackCapacity = cap
	})
}
