package titles

import (
	"testing"
)

var titleTests = []struct {
	v    string
	want string
}{
	{`Senior Product Manager`, `Product Management`},
	{`Senior Manager, Product Management`, `Product Management`},
}

func TestTitle(t *testing.T) {
	parser := NewParser()

	for _, tt := range titleTests {
		dept, err := parser.ParseTitle(tt.v)
		if err != nil {
			t.Errorf("ParseTitle Error: with %v, want %v, err %v", tt.v, tt.want, err)
		}
		if dept.String() != tt.want {
			t.Errorf("ParseTitle Mismatch: with %v, want %v, err %v", tt.v, tt.want, err)
		}
	}
}
