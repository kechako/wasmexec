package runtime

import (
	"context"
	"errors"
	"strconv"

	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
	"github.com/kechako/wasmexec/mod/types"
)

var (
	errExportNotFound            = errors.New("export is not found")
	errExportTargetNotFunction   = errors.New("export target is not a function")
	errFunctionNotFound          = errors.New("function is not found")
	errBlockNotFound             = errors.New("block is not found")
	errStackInconsistent         = errors.New("stack is inconsistent")
	errLocalVariableInconsistent = errors.New("local variables are inconsistent")
	errIntegerDivideByZero       = errors.New("integer divide by zero")
	errUnsupportedType           = errors.New("unsupported type")
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
		idxKey, idKey := makeIndexKeys(i, f.ID)
		vm.funcs[idxKey] = f
		if idKey != "" {
			vm.funcs[idKey] = f
		}
	}
}

func (vm *VM) makeExportTable() {
	for _, e := range vm.mod.Exports {
		vm.exports[e.Name] = e
	}
}

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

	results, err := vm.popContextResults(f.Results)
	if err != nil {
		return nil, err
	}

	return results, err
}

func (vm *VM) callFunc(ctx context.Context, f *mod.Function) error {
	vmCtx, err := vm.initFunction(f, nil)
	if err != nil {
		return err
	}

loop:
	for {
		i := vmCtx.GetInstruction()
		if i == nil {
			var err error
			vmCtx, err = vm.finalizeContext(vmCtx)
			if err != nil {
				return err
			}

			if vmCtx == nil {
				break loop
			}

			continue
		}

		switch i.Name() {
		case instruction.I32Const:
			i := i.(*instruction.I32Instruction)
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
			i := i.(*instruction.VariableInstruction)
			v, err := vmCtx.GetLocal(i.Index)
			if err != nil {
				return err
			}
			n, ok := v.Int32()
			if !ok {
				return errLocalVariableInconsistent
			}
			vm.stack.Push(newValueElement(n))
		case instruction.LocalSet:
			i := i.(*instruction.VariableInstruction)
			err := execLocalSet(vm, vmCtx, i.Index)
			if err != nil {
				return err
			}
		case instruction.LocalTee:
			i := i.(*instruction.VariableInstruction)
			elm := vm.stack.Pop()
			if elm.Type != ValueElement {
				return errStackInconsistent
			}
			vm.stack.Push(newValueElement(elm.Value.Value))
			vm.stack.Push(newValueElement(elm.Value.Value))

			err := execLocalSet(vm, vmCtx, i.Index)
			if err != nil {
				return err
			}
		case instruction.Block:
			i := i.(*instruction.BlockInstruction)
			var err error
			vmCtx, err = vm.initBlock(i.Label, vmCtx)
			if err != nil {
				return err
			}
		case instruction.Return:
			var err error
			vmCtx, err = vm.finalizeContext(vmCtx)
			if err != nil {
				return err
			}
			if vmCtx == nil {
				break loop
			}
		case instruction.Call:
			i := i.(*instruction.CallInstruction)
			index := i.Index
			key := makeIndexKey(index)
			f, ok := vm.funcs[key]
			if !ok {
				return errFunctionNotFound
			}

			var err error
			vmCtx, err = vm.initFunction(f, vmCtx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) initFunction(f *mod.Function, original VMContext) (VMContext, error) {
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

	var vmCtx VMContext
	if original == nil {
		vmCtx = newFuncContext(f, locals, nil)
	} else {
		vmCtx = original.NewFuncContext(f, locals)
	}

	vm.stack.Push(newActivationElement(vmCtx))

	return vmCtx, nil
}

func (vm *VM) initBlock(label types.ID, original VMContext) (VMContext, error) {
	vmCtx, err := original.NewBlockContext(label)
	if err != nil {
		return nil, err
	}

	parameters := vmCtx.Parameters()

	// parameters
	paramLen := len(parameters)
	values := make([]any, paramLen)
	for i := 0; i < paramLen; i++ {
		paramIdx := paramLen - i - 1
		p := parameters[paramIdx]
		switch p.Type {
		case types.I32:
			v, ok := vm.stack.Pop().Int32()
			if !ok {
				return nil, errStackInconsistent
			}
			values[i] = v
		default:
			return nil, errUnsupportedType
		}
	}

	vm.stack.Push(newActivationElement(vmCtx))

	for i := 0; i < paramLen; i++ {
		valueIdx := paramLen - i - 1
		vm.stack.Push(newValueElement(values[valueIdx]))
	}

	return vmCtx, nil
}

func (vm *VM) finalizeContext(vmCtx VMContext) (VMContext, error) {
	results, err := vm.popContextResults(vmCtx.Results())
	if err != nil {
		return nil, err
	}

	// pop func context
	popedCtx, ok := vm.stack.Pop().VMContext()
	if !ok {
		return nil, errStackInconsistent
	}
	if popedCtx != vmCtx {
		return nil, errStackInconsistent
	}

	for _, result := range results {
		vm.stack.Push(newValueElement(result))
	}

	return vmCtx.Original(), nil
}

func (vm *VM) popContextResults(results []*mod.Result) ([]any, error) {
	var values []any
	for _, r := range results {
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
		default:
			return nil, errStackInconsistent
		}

		values = append(values, result)
	}

	return values, nil
}

func execLocalSet(vm *VM, vmCtx VMContext, index types.Index) error {
	elm := vm.stack.Pop()
	if elm.Type != ValueElement {
		return errStackInconsistent
	}

	return vmCtx.SetLocal(index, elm.Value)
}

func makeIndexKey(idx types.Index) string {
	if idx.IsIndex() {
		return strconv.Itoa(idx.Index)
	}

	return string(idx.ID)
}

func makeIndexKeys(index int, id types.ID) (string, string) {
	idxKey := strconv.Itoa(index)
	if id.IsEmpty() {
		return idxKey, ""
	}
	return idxKey, string(id)
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
