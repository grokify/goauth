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
	err := config.LoadDotEnv("ENV_PATH")
	if err != nil {
		panic(err)
	}

	fmt.Printf(os.Getenv("SALESFORCE_CLIENT_SECRET"))

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

	sc := salesforce.NewSalesforceClient(client, os.Getenv("SALESFORCE_INSTANCE_NAME"))

	apiURL := sc.URLBuilder.Build("services/data")

	resp, err := client.Get(apiURL.String())
	if err != nil {
		panic(err)
	}

	httputilmore.PrintResponse(resp, true)

	if 1 == 0 {
		resp, err = sc.ExecSOQL("select id from contact")
		if err != nil {
			panic(err)
		}

		httputilmore.PrintResponse(resp, true)
	}

	if 1 == 0 {
		err = sc.DeleteContactsAll()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("DONE")
}
