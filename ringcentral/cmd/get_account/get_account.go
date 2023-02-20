package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/credentials"
	"github.com/grokify/goauth/ringcentral"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
)

func main() {
	_, err := config.LoadDotEnv([]string{os.Getenv("ENV_PATH"), "./.env"}, -1)
	logutil.FatalErr(err)

	// client := &http.Client{}
	var client *http.Client
	if len(os.Getenv("RINGCENTRAL_ACCESS_TOKEN")) > 0 {
		client = goauth.NewClientAuthzTokenSimple(
			goauth.TokenBearer,
			os.Getenv("RINGCENTRAL_ACCESS_TOKEN"))
	} else {
		client, err = ringcentral.NewClientPassword(
			credentials.CredentialsOAuth2{
				ClientID:     os.Getenv("RINGCENTRAL_CLIENT_ID"),
				ClientSecret: os.Getenv("RINGCENTRAL_CLIENT_SECRET"),
				ServerURL:    os.Getenv("RINGCENTRAL_SERVER_URL"),
				Username:     os.Getenv("RINGCENTRAL_USERNAME"),
				Password:     os.Getenv("RINGCENTRAL_PASSWORD")})
	}
	logutil.FatalErr(err)

	urlPath := "restapi/v1.0/account/~"

	apiURL := urlutil.JoinAbsolute(os.Getenv("RINGCENTRAL_SERVER_URL"), urlPath)

	resp, err := client.Get(apiURL)
	logutil.FatalErr(err)

	logutil.FatalErr(httputilmore.PrintResponse(resp, true))

	fmt.Println("DONE")
}
