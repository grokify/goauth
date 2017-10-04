package main

import (
	"fmt"
	"os"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/net/urlutil"
	"github.com/grokify/oauth2util-go/services/ringcentral"
)

func main() {
	err := config.LoadDotEnv()
	if err != nil {
		panic(err)
	}

	client, err := ringcentral.NewClientPassword(
		ringcentral.ApplicationCredentials{
			ClientID:     os.Getenv("RC_CLIENT_ID"),
			ClientSecret: os.Getenv("RC_CLIENT_SECRET"),
			ServerURL:    os.Getenv("RC_SERVER_URL")},
		ringcentral.UserCredentials{
			Username:  os.Getenv("RC_USER_USERNAME"),
			Extension: os.Getenv("RC_USER_EXTENSION"),
			Password:  os.Getenv("RC_USER_PASSWORD")})

	if err != nil {
		panic(err)
	}

	extURL := urlutil.JoinAbsolute(os.Getenv("RC_SERVER_URL"), "restapi/v1.0/account/~")

	resp, err := client.Get(extURL)
	if err != nil {
		panic(err)
	}

	httputilmore.PrintResponse(resp, true)

	fmt.Println("DONE")
}
