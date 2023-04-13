package authutil

import (
	"strings"
	"testing"
)

var rfc7617UserPassTests = []struct {
	v    string
	want string
}{
	{"hello:world", "aGVsbG86d29ybGQ="}}

func TestRFC7617UserPass(t *testing.T) {
	for _, tt := range rfc7617UserPassTests {
		parts := strings.Split(tt.v, ":")
		enc, err := RFC7617UserPass(parts[0], parts[1])
		if err != nil {
			t.Errorf("base64.RFC7617UserPass(%v, %v): want %v, error %v", parts[0], parts[1], tt.want, err)
		}

		if enc != tt.want {
			t.Errorf("base64.RFC7617UserPass(%v, %v): want %v, got %v", parts[0], parts[1], tt.want, enc)
		}
	}
}
