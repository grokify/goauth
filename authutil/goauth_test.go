package authutil

import (
	"testing"
)

var authorizationTypeTests = []struct {
	v    AuthorizationType
	want string
}{
	{Anonymous, "Anonymous"},
	{Basic, "Basic"},
	{Bearer, "Bearer"},
	{Digest, "Digest"},
	{NTLM, "NTLM"},
	{Negotiate, "Negotiate"},
	{OAuth, "OAuth"},
}

func TestAuthorizationTypeString(t *testing.T) {
	for _, tt := range authorizationTypeTests {
		try := tt.v.String()
		if try != tt.want {
			t.Errorf("goauth.AuthorizationType(%v): want (%v), got (%v)", tt.v, tt.want, try)
		}
	}
}
