package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/credentials"
	"github.com/grokify/goauth/ringcentral"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/net/httputilmore"
	"github.com/grokify/simplego/net/urlutil"
)

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		panic(err)
	}

	// client := &http.Client{}
	var client *http.Client
	if len(os.Getenv("RINGCENTRAL_ACCESS_TOKEN")) > 0 {
		client = goauth.NewClientAuthzTokenSimple(
			goauth.TokenBearer,
			os.Getenv("RINGCENTRAL_ACCESS_TOKEN"))
	} else {
		client, err = ringcentral.NewClientPassword(
			credentials.OAuth2Credentials{
				ClientID:     os.Getenv("RINGCENTRAL_CLIENT_ID"),
				ClientSecret: os.Getenv("RINGCENTRAL_CLIENT_SECRET"),
				ServerURL:    os.Getenv("RINGCENTRAL_SERVER_URL"),
				Username:     os.Getenv("RINGCENTRAL_USERNAME"),
				Password:     os.Getenv("RINGCENTRAL_PASSWORD")})
	}
	if err != nil {
		panic(err)
	}

	urlPath := "restapi/v1.0/account/~"

	apiURL := urlutil.JoinAbsolute(os.Getenv("RINGCENTRAL_SERVER_URL"), urlPath)

	resp, err := client.Get(apiURL)
	if err != nil {
		panic(err)
	}

	httputilmore.PrintResponse(resp, true)

	fmt.Println("DONE")
}
