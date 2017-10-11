package main

import (
	"fmt"
	"os"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/net/httputilmore"
	ou "github.com/grokify/oauth2util-go"
	"github.com/grokify/oauth2util-go/services/salesforce"
)

func main() {
	err := config.LoadDotEnv()
	if err != nil {
		panic(err)
	}

	client, err := salesforce.NewClientPassword(
		ou.ApplicationCredentials{
			ClientID:     os.Getenv("SALESFORCE_CLIENT_ID"),
			ClientSecret: os.Getenv("SALESFORCE_CLIENT_SECRET")},
		ou.UserCredentials{
			Username: os.Getenv("SALESFORCE_USERNAME"),
			Password: fmt.Sprintf("%v%v",
				os.Getenv("SALESFORCE_PASSWORD"),
				os.Getenv("SALESFORCE_SECURITY_KEY"))})

	if err != nil {
		panic(err)
	}

	urlBuilder := salesforce.NewURLBuilder(os.Getenv("SALESFORCE_INSTANCE_NAME"))

	apiURL := urlBuilder.Build("services/data")

	resp, err := client.Get(apiURL.String())
	if err != nil {
		panic(err)
	}

	httputilmore.PrintResponse(resp, true)

	fmt.Println("DONE")
}
