package runtime

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kechako/wasmexec/mod/text"
)

type valueTypes interface {
	int32 | int64 | float32 | float64
}

func newResults(values ...any) []any {
	return values
}

func newTypedResults[T valueTypes](values ...T) []any {
	results := make([]any, len(values))

	for i, v := range values {
		results[i] = v
	}

	return results
}

var execFuncTests = map[string]struct {
	results []any
}{
	"test01.wat": {
		results: newTypedResults[int32](25, 3, 28, 2),
	},
	"test02.wat": {
		results: newTypedResults[int32](1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 1),
	},
	"test03.wat": {
		results: newTypedResults[int32](4),
	},
	"test04.wat": {
		results: newTypedResults[int32](40, 30, 20, 10, 50, 60, 70, 80, 80, 70, 60, 50),
	},
	"test05.wat": {
		results: newTypedResults[int32](300),
	},
	"test06.wat": {
		results: newTypedResults[int32](180),
	},
	"test07.wat": {
		results: newTypedResults[int32](180),
	},
}

func Test_VM_ExecFunc(t *testing.T) {
	ctx := context.Background()
	for name, tt := range execFuncTests {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			vm, err := createVM(name)
			if err != nil {
				t.Fatal(err)
			}

			results, err := vm.ExecFunc(ctx, "main")
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(results, tt.results); diff != "" {
				t.Errorf("VM.ExecFunc(ctx, \"main\"), differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func createVM(name string) (*VM, error) {
	file, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m, err := text.NewDecoder(file).Decode()
	if err != nil {
		return nil, err
	}

	return New(m), nil
}
