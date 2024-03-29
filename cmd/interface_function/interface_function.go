package main

import (
	"fmt"

	"github.com/grokify/goauth/aha"
	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/facebook"
)

type Interface interface {
	Function() bool
}

type A struct{}

func (a *A) Function() bool {
	return true
}

type B struct{}

func (b *B) Function() bool {
	return false
}

func Choose(s string) Interface {
	if s == "a" {
		return &A{}
	}
	return &B{}
}

func PrintInterfaceFunction(i Interface) {
	fmt.Println(i.Function())
}

func ChooseClient(s string) authutil.OAuth2Util {
	if s == "aha" {
		return &aha.ClientUtil{}
	}
	return &facebook.ClientUtil{}
}

func main() {
	item := Choose("b")
	PrintInterfaceFunction(item)
	//fmt.Println(item.Function())

	clientUtil := ChooseClient("aha")
	fmt.Printf("%v\n", clientUtil)
	fmt.Println("DONE'")
}
