package text

import (
	"errors"
	"fmt"
	"io"

	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
	"github.com/kechako/wasmexec/mod/text/sexp"
	"github.com/kechako/wasmexec/mod/types"
)

type Decoder struct {
	p *sexp.Parser
}

var _ mod.Decoder = (*Decoder)(nil)

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		p: sexp.New(r),
	}
}

func (d *Decoder) Decode() (*mod.Module, error) {
	node, err := d.p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to decode wat: %v", err)
	}

	m, err := parseModule(node)
	if err != nil {
		return nil, fmt.Errorf("failed to decode wat: %v", err)
	}

	return m, nil
}

var errInvalidModuleFormat = errors.New("invalid module format")

func parseModule(node *sexp.Node) (*mod.Module, error) {
	if v, ok := node.Car.SymbolValue(); !ok || v != "module" {
		return nil, errInvalidModuleFormat
	}

	m := &mod.Module{}

	first := true
	for curr := node.Cdr; curr != nil; curr = curr.Cdr {
		car := curr.Car

		switch car.Type {
		case sexp.NodeCell:
			err := parseModuleField(m, car)
			if err != nil {
				return nil, err
			}
		case sexp.NodeSymbol:
			if first {
				v, _ := car.SymbolValue()
				id := types.ID(v)
				if !id.IsValid() {
					return nil, errInvalidModuleFormat
				}
				m.ID = id
			} else {
				return nil, errInvalidModuleFormat
			}
		}
		first = false
	}

	return m, nil
}

var errUnsupportedField = errors.New("unsupported module field")

func parseModuleField(m *mod.Module, node *sexp.Node) error {
	car := node.Car

	sym, ok := car.SymbolValue()
	if !ok {
		return errInvalidModuleFormat
	}

	switch sym {
	case "func":
		p := &functionParser{}
		f, err := p.Parse(node.Cdr)
		if err != nil {
			return err
		}
		m.Functions = append(m.Functions, f)
	case "export":
		e, err := parseExport(node.Cdr)
		if err != nil {
			return err
		}
		m.Exports = append(m.Exports, e)
	default:
		return errUnsupportedField
	}

	return nil
}

type functionParser struct {
	f *mod.Function
}

func (p *functionParser) Parse(node *sexp.Node) (*mod.Function, error) {
	if node == nil {
		return nil, errInvalidModuleFormat
	}

	// id (optional)
	var id types.ID
	if v, ok := node.Car.SymbolValue(); ok {
		id = types.ID(v)
		if !id.IsValid() {
			return nil, errInvalidModuleFormat
		}

		node = node.Cdr
		if node == nil {
			return nil, errInvalidModuleFormat
		}
	}

	f := &mod.Function{
		ID: id,
	}
	p.f = f

	curr := node

	// parse params
	for curr != nil && isFunctionParam(curr.Car) {
		car := curr.Car

		p, err := p.parseLocal(car.Cdr)
		if err != nil {
			return nil, err
		}

		f.Parameters = append(f.Parameters, p)

		curr = curr.Cdr
	}

	// parse results
	for curr != nil && isFunctionResult(curr.Car) {
		car := curr.Car

		r, err := p.parseResult(car.Cdr)
		if err != nil {
			return nil, err
		}

		f.Results = append(f.Results, r)

		curr = curr.Cdr
	}

	// parse locals
	for curr != nil && isFunctionLocal(curr.Car) {
		car := curr.Car

		l, err := p.parseLocal(car.Cdr)
		if err != nil {
			return nil, err
		}

		f.Locals = append(f.Locals, l)

		curr = curr.Cdr
	}

	instructions, err := p.parseInstructions(curr)
	if err != nil {
		return nil, err
	}
	f.Instructions = instructions

	return f, nil
}

func isFunctionParam(node *sexp.Node) bool {
	if node == nil || node.Type != sexp.NodeCell {
		return false
	}

	sym, ok := node.Car.SymbolValue()
	if !ok {
		return false
	}

	return sym == "param"
}

func isFunctionLocal(node *sexp.Node) bool {
	if node == nil || node.Type != sexp.NodeCell {
		return false
	}

	sym, ok := node.Car.SymbolValue()
	if !ok {
		return false
	}

	return sym == "local"
}

func isFunctionResult(node *sexp.Node) bool {
	if node == nil || node.Type != sexp.NodeCell {
		return false
	}

	sym, ok := node.Car.SymbolValue()
	if !ok {
		return false
	}

	return sym == "result"
}

func (p *functionParser) parseLocal(node *sexp.Node) (*mod.Local, error) {
	if node == nil {
		return nil, errInvalidModuleFormat
	}

	// id (optional)
	var id types.ID
	if node.Cdr != nil {
		if v, ok := node.Car.SymbolValue(); ok {
			id = types.ID(v)
			if !id.IsValid() {
				return nil, errInvalidModuleFormat
			}

			node = node.Cdr
			if node == nil {
				return nil, errInvalidModuleFormat
			}
		}
	}

	// type
	v, ok := node.Car.SymbolValue()
	if !ok {
		return nil, errInvalidModuleFormat
	}

	typ := parseType(v)
	if typ == types.Unkown {
		return nil, errInvalidModuleFormat
	}

	return &mod.Local{
		ID:   id,
		Type: typ,
	}, nil
}

func parseType(s string) types.Type {
	switch s {
	case "i32":
		return types.I32
	case "i64":
		return types.I64
	case "f32":
		return types.F32
	case "f64":
		return types.F64
	}

	return types.Unkown
}

func (p *functionParser) parseResult(node *sexp.Node) (*mod.Result, error) {
	if node == nil {
		return nil, errInvalidModuleFormat
	}

	// type
	v, ok := node.Car.SymbolValue()
	if !ok {
		return nil, errInvalidModuleFormat
	}

	typ := parseType(v)
	if typ == types.Unkown {
		return nil, errInvalidModuleFormat
	}

	return &mod.Result{
		Type: typ,
	}, nil
}

func (p *functionParser) parseInstructions(node *sexp.Node) ([]instruction.Instruction, error) {
	var instructions []instruction.Instruction

	curr := node
	for curr != nil {
		i, next, err := p.parseInstruction(curr)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, i)

		curr = next
	}

	return instructions, nil
}

var errUnsupportedInstruction = errors.New("unsupported instruction")

func (p *functionParser) parseInstruction(node *sexp.Node) (instruction.Instruction, *sexp.Node, error) {
	if node == nil {
		return nil, nil, errInvalidModuleFormat
	}

	if node.Car.Type == sexp.NodeCell {
		i, err := p.parseNestedInstruction(node.Car)
		if err != nil {
			return nil, nil, err
		}
		return i, node.Cdr, nil
	} else {
		iname, err := getInstructionName(node)
		if err != nil {
			return nil, nil, err
		}

		i, next, err := p.parseI32Instruction(iname, node.Cdr)
		if err == nil {
			return i, next, nil
		}
		if err != nil && err != errUnsupportedInstruction {
			return nil, nil, err
		}

		i, next, err = p.parseParametricInstruction(iname, node.Cdr)
		if err == nil {
			return i, next, nil
		}
		if err != nil && err != errUnsupportedInstruction {
			return nil, nil, err
		}

		i, next, err = p.parseVariableInstruction(iname, node.Cdr)
		if err == nil {
			return i, next, nil
		}
		if err != nil && err != errUnsupportedInstruction {
			return nil, nil, err
		}

		i, next, err = p.parseControlInstruction(iname, node.Cdr)
		if err == nil {
			return i, next, nil
		}
		if err != nil && err != errUnsupportedInstruction {
			return nil, nil, err
		}
	}

	return nil, nil, errUnsupportedInstruction
}

func getInstructionName(node *sexp.Node) (instruction.InstructionName, error) {
	if node == nil {
		return "", errInvalidModuleFormat
	}

	v, ok := node.Car.SymbolValue()
	if !ok {
		return "", errInvalidModuleFormat
	}

	return instruction.InstructionName(v), nil
}

func (p *functionParser) parseI32Instruction(iname instruction.InstructionName, node *sexp.Node) (instruction.Instruction, *sexp.Node, error) {
	if !iname.IsI32() {
		return nil, nil, errUnsupportedInstruction
	}

	switch iname {
	case instruction.I32Const:
		n, ok := node.Car.IntValue()
		if !ok {
			return nil, nil, errInvalidModuleFormat
		}
		return &instruction.I32Instruction{
			Instruction: iname,
			Values:      []int32{int32(n)},
		}, node.Cdr, nil
	}

	return &instruction.I32Instruction{
		Instruction: iname,
	}, node, nil
}

func (p *functionParser) parseParametricInstruction(iname instruction.InstructionName, node *sexp.Node) (instruction.Instruction, *sexp.Node, error) {
	if !iname.IsParametric() {
		return nil, nil, errUnsupportedInstruction
	}

	return &instruction.ParametricInstruction{
		Instruction: iname,
	}, node, nil
}

func (p *functionParser) parseVariableInstruction(iname instruction.InstructionName, node *sexp.Node) (instruction.Instruction, *sexp.Node, error) {
	if !iname.IsVariable() {
		return nil, nil, errUnsupportedInstruction
	}

	index, err := parseIndex(node)
	if err != nil {
		return nil, nil, err
	}
	return &instruction.VariableInstruction{
		Instruction: iname,
		Index:       index,
	}, node.Cdr, nil
}

func (p *functionParser) parseControlInstruction(iname instruction.InstructionName, node *sexp.Node) (instruction.Instruction, *sexp.Node, error) {
	if !iname.IsControl() {
		return nil, nil, errUnsupportedInstruction
	}

	switch iname {
	case instruction.Block:
	case instruction.Call:
		index, err := parseIndex(node)
		if err != nil {
			return nil, nil, err
		}
		return &instruction.CallInstruction{
			Instruction: iname,
			Index:       index,
		}, node.Cdr, nil
	}

	return &instruction.ControlInstruction{
		Instruction: iname,
	}, node, nil
}

func (p *functionParser) parseNestedInstruction(node *sexp.Node) (instruction.Instruction, error) {
	iname, err := getInstructionName(node)
	if err != nil {
		return nil, err
	}

	switch iname {
	case instruction.Block:
		return p.parseBlockInstruction(node.Cdr)
	}

	return nil, errUnsupportedInstruction
}

func (p *functionParser) parseBlockInstruction(node *sexp.Node) (instruction.Instruction, error) {
	if node == nil {
		return nil, errInvalidModuleFormat
	}

	// label
	var label types.ID
	if v, ok := node.Car.SymbolValue(); ok {
		label = types.ID(v)
		if !label.IsValid() {
			return nil, errInvalidModuleFormat
		}

		node = node.Cdr
		if node == nil {
			return nil, errInvalidModuleFormat
		}
	}

	block := &mod.Block{
		Label: label,
	}

	curr := node

	// parse params
	for curr != nil && isFunctionParam(curr.Car) {
		car := curr.Car

		p, err := p.parseBlockParam(car.Cdr)
		if err != nil {
			return nil, err
		}

		block.Parameters = append(block.Parameters, p)

		curr = curr.Cdr
	}

	// parse results
	for curr != nil && isFunctionResult(curr.Car) {
		car := curr.Car

		r, err := p.parseBlockResult(car.Cdr)
		if err != nil {
			return nil, err
		}

		block.Results = append(block.Results, r)

		curr = curr.Cdr
	}

	instructions, err := p.parseInstructions(curr)
	if err != nil {
		return nil, err
	}
	block.Instructions = instructions

	p.f.Blocks = append(p.f.Blocks, block)

	return &instruction.BlockInstruction{
		Instruction: instruction.Block,
		Label:       label,
	}, nil
}

func (p *functionParser) parseBlockParam(node *sexp.Node) (*mod.Local, error) {
	if node == nil {
		return nil, errInvalidModuleFormat
	}

	// type
	v, ok := node.Car.SymbolValue()
	if !ok {
		return nil, errInvalidModuleFormat
	}

	typ := parseType(v)
	if typ == types.Unkown {
		return nil, errInvalidModuleFormat
	}

	return &mod.Local{
		Type: typ,
	}, nil
}

func (p *functionParser) parseBlockResult(node *sexp.Node) (*mod.Result, error) {
	if node == nil {
		return nil, errInvalidModuleFormat
	}

	// type
	v, ok := node.Car.SymbolValue()
	if !ok {
		return nil, errInvalidModuleFormat
	}

	typ := parseType(v)
	if typ == types.Unkown {
		return nil, errInvalidModuleFormat
	}

	return &mod.Result{
		Type: typ,
	}, nil
}

func parseExport(node *sexp.Node) (*mod.Export, error) {
	if node == nil {
		return nil, errInvalidModuleFormat
	}

	// name
	name, ok := node.Car.StringValue()
	if !ok {
		return nil, errInvalidModuleFormat
	}

	node = node.Cdr
	if node == nil {
		return nil, errInvalidModuleFormat
	}

	// export description
	if node.Car.Type != sexp.NodeCell {
		return nil, errInvalidModuleFormat
	}
	node = node.Car

	// export target
	tv, ok := node.Car.SymbolValue()
	if !ok {
		return nil, errInvalidModuleFormat
	}
	target, err := parseExportTarget(tv)
	if err != nil {
		return nil, err
	}

	node = node.Cdr
	if node == nil {
		return nil, errInvalidModuleFormat
	}

	// index
	index, err := parseIndex(node)
	if err != nil {
		return nil, err
	}

	return &mod.Export{
		Name:   name,
		Target: target,
		Index:  index,
	}, nil
}

func parseExportTarget(s string) (mod.ExportTarget, error) {
	switch s {
	case "func":
		return mod.ExportFunction, nil
	case "table":
		return mod.ExportTable, nil
	case "memory":
		return mod.ExportMemory, nil
	case "global":
		return mod.ExportGlobal, nil
	}

	return "", errInvalidModuleFormat
}

func parseIndex(node *sexp.Node) (types.Index, error) {
	var index types.Index
	if idx, ok := node.Car.IntValue(); ok {
		index.Index = int(idx)
	} else if v, ok := node.Car.SymbolValue(); ok {
		id := types.ID(v)
		if !id.IsValid() {
			return index, errInvalidModuleFormat
		}
		index.ID = id
	} else {
		return index, errInvalidModuleFormat
	}

	return index, nil
}
