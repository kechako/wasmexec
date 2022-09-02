package runtime

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
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

var (
	errExportNotFound          = errors.New("export is not found")
	errExportTargetNotFunction = errors.New("export target is not a function")
	errFunctionNotFound        = errors.New("function is not found")
	errStackInconsistent       = errors.New("stack is inconsistent")
)

func (vm *VM) ExecFunc(ctx context.Context, name string) error {
	// エクスポートを検索
	e, ok := vm.exports[name]
	if !ok {
		return errExportNotFound
	}

	// エクスポートは関数？
	if e.Target != mod.ExportFunction {
		return errExportTargetNotFunction
	}

	// 関数を検索
	f, ok := vm.funcs[makeFuncKey(e.Index)]
	if !ok {
		return errFunctionNotFound
	}

	typ := mod.I32
	if len(f.Results) > 0 {
		typ = f.Results[0].Type
	}

	var result any
	var err error
	switch typ {
	case mod.I32:
		result, err = callFunc[int32](ctx, vm, f)
	case mod.I64:
		result, err = callFunc[int64](ctx, vm, f)
	case mod.F32:
		result, err = callFunc[float32](ctx, vm, f)
	case mod.F64:
		result, err = callFunc[float64](ctx, vm, f)
	}
	if err != nil {
		return err
	}

	fmt.Println(result)

	return err
}

type resultValue interface {
	int32 | int64 | float32 | float64
}

func callFunc[T resultValue](ctx context.Context, vm *VM, f *mod.Function) (result T, err error) {
	funcCtx := newFuncContext(f)

	vm.stack.Push(newActivationElement(funcCtx))

loop:
	for {
		i := funcCtx.GetInstruction()
		if i == nil {
			break
		}

		switch i.Name() {
		case instruction.I32Const:
			i := i.(*instruction.I32)
			vm.stack.Push(newValueElement(i.Values[0]))
		case instruction.I32Add:
			c2, ok := vm.stack.Pop().Int32Value()
			if !ok {
				return result, errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32Value()
			if !ok {
				return result, errStackInconsistent
			}
			vm.stack.Push(newValueElement(c1 + c2))
		case instruction.I32Sub:
			c2, ok := vm.stack.Pop().Int32Value()
			if !ok {
				return result, errStackInconsistent
			}
			c1, ok := vm.stack.Pop().Int32Value()
			if !ok {
				return result, errStackInconsistent
			}
			vm.stack.Push(newValueElement(c1 - c2))
		case instruction.Return:
			break loop
		}
	}

	finalize := func() (T, error) {
		// pop result
		resultElm := vm.stack.Pop()
		result, ok := getElementValue[T](resultElm, ValueElement)
		if !ok {
			return result, errStackInconsistent
		}

		// pop func context
		_, ok = vm.stack.Pop().FuncContext()
		if !ok {
			return result, errStackInconsistent
		}

		return result, nil
	}

	return finalize()
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

func makeFuncKey(idx mod.Index) string {
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
