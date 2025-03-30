package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/salesforce"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/net/http/httputilmore"

	su "github.com/grokify/go-salesforce/clientutil"
)

func main() {
	//err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	_, err := config.LoadDotEnv([]string{"./.env"}, -1)
	if err != nil {
		panic(err)
	}

	fmt.Println(os.Getenv("SALESFORCE_CLIENT_SECRET"))

	client, err := salesforce.NewClientPassword(
		goauth.CredentialsOAuth2{
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
		log.Fatal(err)
	}

	err = httputilmore.PrintResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("SALESFORCE_DESCRIBE_ACCOUNT") == "true" {
		cu := su.ClientUtil{
			HTTPClient: client,
			Instance:   os.Getenv("SALESFORCE_INSTANCE_NAME"),
			Version:    "v43.0"}
		resp, err := cu.Describe("ACCOUNT")
		if err != nil {
			log.Fatal(err)
		}
		//httputilmore.PrintResponse(resp, true)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
		desc := su.Describe{}
		err = json.Unmarshal(body, &desc)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.MustPrintJSON(desc)

		types := map[string]int{}
		for _, f := range desc.Fields {
			if v, ok := types[f.Type]; ok {
				types[f.Type] = v + 1
			} else {
				types[f.Type] = 1
			}
		}
		fmtutil.MustPrintJSON(types)
	}

	if 1 == 0 {
		resp, err = sc.ExecSOQL("select id from contact")
		if err != nil {
			log.Fatal(err)
		}

		err = httputilmore.PrintResponse(resp, true)
		if err != nil {
			log.Fatal(err)
		}
	}

	if 1 == 0 {
		err = sc.DeleteContactsAll()
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("DONE")
}
