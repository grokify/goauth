package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grokify/oauth2more/zendesk"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
	hum "github.com/grokify/simplego/net/httputilmore"
)

func MeURL(subdomain string) string {
	return fmt.Sprintf("https://%v.zendesk.com/api/v2/users/me.json", subdomain)
}

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		panic(err)
	}

	client, err := zendesk.NewClientPassword(
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
