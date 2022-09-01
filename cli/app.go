package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
	"github.com/kechako/wasmexec/mod/text"
)

type App struct {
}

func (app *App) Run(ctx context.Context) error {
	args := os.Args[1:]
	if len(args) == 0 {
		return errors.New("invalid arguments")
	}

	m, err := app.decode(args[0])
	if err != nil {
		return err
	}

	dumpModule(m)

	return nil
}

func (app *App) decode(name string) (*mod.Module, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open WASM file: %w", err)
	}
	defer file.Close()

	return text.NewDecoder(file).Decode()
}

func dumpModule(m *mod.Module) {
	if m.ID.IsEmpty() {
		fmt.Println("Module:")
	} else {
		fmt.Printf("Module[%s]:\n", m.ID)
	}

	fmt.Println("  Functions:")
	for i, f := range m.Functions {
		dumpFunction(f, i, "    ")
	}
	fmt.Println("  Exports:")
	for _, e := range m.Exports {
		if e.Index.IsIndex() {
			fmt.Printf("    %q: %s %d\n", e.Name, e.Target, e.Index.Index)
		} else {
			fmt.Printf("    %q: %s %s\n", e.Name, e.Target, e.Index.ID)
		}
	}
}

func dumpFunction(f *mod.Function, index int, indent string) {
	if f.ID.IsEmpty() {
		fmt.Printf("%sFunc[%d]:\n", indent, index)
	} else {
		fmt.Printf("%sFunc[%s]:\n", indent, f.ID)
	}
	fmt.Println("      Parameters:")
	for i, p := range f.Parameters {
		dumpParameter(p, i, indent+"    ")
	}
	fmt.Println("      Results:")
	for i, r := range f.Results {
		dumpResult(r, i, indent+"    ")
	}
	fmt.Println("      Instructions:")
	for _, i := range f.Instructions {
		dumpInstruction(i, indent+"    ")
	}
}

func dumpParameter(p *mod.Parameter, index int, indent string) {
	if p.ID.IsEmpty() {
		fmt.Printf("%sParam[%d]: %s\n", indent, index, p.Type)
	} else {
		fmt.Printf("%sParam[%s]: %s\n", indent, p.ID, p.Type)
	}
}

func dumpResult(r *mod.Result, index int, indent string) {
	fmt.Printf("%sResult: %s\n", indent, r.Type)
}

func dumpInstruction(i instruction.Instruction, indent string) {
	switch v := i.(type) {
	case *instruction.I32:
		dumpI32Instruction(v, indent)
	case *instruction.Parametric:
		dumpParametricInstruction(v, indent)
	case *instruction.Control:
		dumpControlInstruction(v, indent)
	default:
		fmt.Printf("%s(unknown)\n", indent)
	}
}

func dumpI32Instruction(i *instruction.I32, indent string) {
	fmt.Printf("%s%s", indent, i.Instruction)
	for _, v := range i.Values {
		fmt.Printf(" %d", v)
	}
	fmt.Println()
}

func dumpParametricInstruction(i *instruction.Parametric, indent string) {
	fmt.Printf("%s%s\n", indent, i.Instruction)
}

func dumpControlInstruction(i *instruction.Control, indent string) {
	fmt.Printf("%s%s\n", indent, i.Instruction)
}
