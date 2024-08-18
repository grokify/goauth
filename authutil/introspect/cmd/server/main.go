package main

import (
	"log"

	"github.com/grokify/goauth/authutil/introspect"
)

func main() {
	svr := introspect.NewMockServer(introspect.IntrospectResponse{
		Username: "foo"},
		[]string{"bar", "baz"},
	)
	err := svr.ListenAndServe(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
