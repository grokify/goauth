package main

import (
	"fmt"

	ou "github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/aha"
	"github.com/grokify/oauth2more/facebook"
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

func ChooseClient(s string) ou.OAuth2Util {
	if s == "aha" {
		return &aha.ClientUtil{}
	}
	return &facebook.ClientUtil{}
}

func main() {
	var item Interface
	item = Choose("b")
	PrintInterfaceFunction(item)
	//fmt.Println(item.Function())

	var clientUtil ou.OAuth2Util
	clientUtil = ChooseClient("aha")
	fmt.Printf("%v\n", clientUtil)
	fmt.Println("DONE'")
}
