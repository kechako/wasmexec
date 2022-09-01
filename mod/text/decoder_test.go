package text

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kechako/wasmexec/mod"
	"github.com/kechako/wasmexec/mod/instruction"
)

var tests = map[string]struct {
	input string
	mod   *mod.Module
	err   error
}{
	"success 01": {
		input: `(module $testmod
  (func $main
  	(param $a i32) (param $b i64)
    (result i32)
   i32.const 5
   i32.const 20
   i32.add
   i32.const 4
   i32.sub
   return
     drop 
    i32.const 0
  )
  (export "main" (func $main))
)`,
		mod: &mod.Module{
			ID: "$testmod",
			Functions: []*mod.Function{
				{
					ID: "$main",
					Parameters: []*mod.Parameter{
						{ID: "$a", Type: mod.I32},
						{ID: "$b", Type: mod.I64},
					},
					Results: []*mod.Result{
						{Type: mod.I32},
					},
					Instructions: []instruction.Instruction{
						&instruction.I32{Instruction: instruction.I32Const, Values: []int32{5}},
						&instruction.I32{Instruction: instruction.I32Const, Values: []int32{20}},
						&instruction.I32{Instruction: instruction.I32Add},
						&instruction.I32{Instruction: instruction.I32Const, Values: []int32{4}},
						&instruction.I32{Instruction: instruction.I32Sub},
						&instruction.Control{Instruction: instruction.Return},
						&instruction.Parametric{Instruction: instruction.Drop},
						&instruction.I32{Instruction: instruction.I32Const, Values: []int32{0}},
					},
				},
			},
			Exports: []*mod.Export{
				{Name: "main", Target: mod.ExportFunction, Index: mod.Index{ID: "$main"}},
			},
		},
		err: nil,
	},
}

func Test_Decode(t *testing.T) {
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			d := NewDecoder(strings.NewReader(tt.input))
			m, err := d.Decode()
			if err != tt.err {
				t.Errorf("Decoder.Decode(): err: want: %v, got:%v", tt.err, err)
			}

			if diff := cmp.Diff(m, tt.mod); diff != "" {
				t.Errorf("Decoder.Decode(), differs: (-got +want)\n%s", diff)
			}
		})
	}
}
