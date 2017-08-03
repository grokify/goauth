package main

import (
	"encoding/json"
	"fmt"
	"github.com/grokify/oauth2-util-go/services/ringcentral"
	"golang.org/x/oauth2"
)

func main() {
	cfg := oauth2.Config{
		ClientID:     "my_client_id",
		ClientSecret: "my_client_secret",
		Endpoint:     ringcentral.NewEndpoint(ringcentral.SandboxHostname),
		Scopes:       []string{},
	}
	data, _ := json.Marshal(cfg)

	fmt.Println(string(data))

	fmt.Println("DONE")
}
