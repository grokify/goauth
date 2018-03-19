package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	hum "github.com/grokify/gotilla/net/httputilmore"
	om "github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/zendesk"
	"golang.org/x/oauth2"
)

func MeURL(subdomain string) string {
	return fmt.Sprintf("https://%v.zendesk.com/api/v2/users/me.json", subdomain)
}

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		panic(err)
	}

	client, err := om.NewClientPassword(
		oauth2.Config{},
		context.Background(),
		os.Getenv("ZENDESK_USERNAME"),
		os.Getenv("ZENDESK_PASSWORD"),
	)
	if err != nil {
		panic(err)
	}

	subdomain := os.Getenv("ZENDESK_SUBDOMAIN")

	if 1 == 0 {
		meURL := MeURL(subdomain)
		resp, err := client.Get(meURL)
		if err != nil {
			panic(err)
		}

		err = hum.PrintResponse(resp, true)
		if err != nil {
			panic(err)
		}
	}
	if 1 == 1 {
		me, resp, err := zendesk.GetMe(client, subdomain)
		if err != nil {
			panic(err)
		} else if resp.StatusCode >= 300 {
			panic(fmt.Errorf("Status Code %v", resp.StatusCode))
		}

		fmtutil.PrintJSON(me)

		cu := zendesk.NewClientUtil(client, subdomain)

		scim, err := cu.GetSCIMUser()
		if err != nil {
			panic(err)
		}
		fmtutil.PrintJSON(scim)
	}

	fmt.Println("DONE")
}
