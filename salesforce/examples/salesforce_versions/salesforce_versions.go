package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/grokify/oauth2more/credentials"
	"github.com/grokify/oauth2more/salesforce"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/net/httputilmore"

	su "github.com/grokify/go-salesforce/clientutil"
)

func main() {
	//err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	err := config.LoadDotEnvSkipEmpty("./.env")
	if err != nil {
		panic(err)
	}

	fmt.Printf(os.Getenv("SALESFORCE_CLIENT_SECRET"))

	client, err := salesforce.NewClientPassword(
		credentials.OAuth2Credentials{
			ClientID:     os.Getenv("SALESFORCE_CLIENT_ID"),
			ClientSecret: os.Getenv("SALESFORCE_CLIENT_SECRET"),
			Username:     os.Getenv("SALESFORCE_USERNAME"),
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

	if 1 == 1 {
		cu := su.ClientUtil{
			HTTPClient: client,
			Instance:   os.Getenv("SALESFORCE_INSTANCE_NAME"),
			Version:    "v43.0"}
		resp, err := cu.Describe("ACCOUNT")
		if err != nil {
			log.Fatal(err)
		}
		//httputilmore.PrintResponse(resp, true)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
		desc := su.Describe{}
		err = json.Unmarshal(body, &desc)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(desc)

		types := map[string]int{}
		for _, f := range desc.Fields {
			if v, ok := types[f.Type]; ok {
				types[f.Type] = v + 1
			} else {
				types[f.Type] = 1
			}
		}
		fmtutil.PrintJSON(types)

	}

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
