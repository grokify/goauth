package main

import (
	"fmt"
	"os"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/oauth2more/ringcentral"
)

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		panic(err)
	}

	token, err := ringcentral.NewTokenPassword(
		ringcentral.ApplicationCredentials{
			ClientID:     os.Getenv("RINGCENTRAL_CLIENT_ID"),
			ClientSecret: os.Getenv("RINGCENTRAL_CLIENT_SECRET"),
			ServerURL:    os.Getenv("RINGCENTRAL_SERVER_URL")},
		ringcentral.PasswordCredentials{
			Username:  os.Getenv("RINGCENTRAL_USERNAME"),
			Extension: os.Getenv("RINGCENTRAL_EXTENSION"),
			Password:  os.Getenv("RINGCENTRAL_PASSWORD")})
	if err != nil {
		panic(err)
	}

	token.Expiry = token.Expiry.UTC()

	fmtutil.PrintJSON(token)

	fmt.Println("DONE")
}
