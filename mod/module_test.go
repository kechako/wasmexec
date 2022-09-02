package mod

import (
	"testing"

	"github.com/kechako/wasmexec/mod/types"
)

var idTests = []struct {
	id    types.ID
	valid bool
}{
	{
		id:    "$01234",
		valid: true,
	},
	{
		id:    "$ABCD",
		valid: true,
	},
	{
		id:    "$test",
		valid: true,
	},
	{
		id:    "$$$$",
		valid: true,
	},
	{
		id:    "$!#$%&'*+-,/:<=>?@\\^_`|~",
		valid: true,
	},
	{
		id:    "abcd",
		valid: false,
	},
	{
		id:    "$[abcd]",
		valid: false,
	},
}

func Test_ID(t *testing.T) {
	for _, tt := range idTests {
		tt := tt
		valid := tt.id.IsValid()
		if valid != tt.valid {
			t.Errorf("ID(%q).Valid(): got %v, want %v", tt.id, valid, tt.valid)
		}
	}
}
