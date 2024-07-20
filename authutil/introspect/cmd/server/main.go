package main

import (
	"github.com/grokify/goauth/authutil/introspect"
)

func main() {
	svr := introspect.NewMockServer(introspect.IntrospectResponse{
		Username: "foo"},
		[]string{"bar", "baz"},
	)
	svr.ListenAndServe(":8000")
}
