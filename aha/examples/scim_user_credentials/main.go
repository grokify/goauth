package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/aha"
	"github.com/grokify/goauth/endpoints"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"golang.org/x/oauth2"
)

func loadEnv() error {
	if len(os.Getenv("ENV_PATH")) > 0 {
		return godotenv.Load(os.Getenv("ENV_PATH"))
	}
	return godotenv.Load()
}

func main() {
	err := loadEnv()
	if err != nil {
		panic(err)
	}

	creds := goauth.Credentials{
		Type:    goauth.TypeOAuth2,
		Service: endpoints.ServiceAha,
		OAuth2: &goauth.CredentialsOAuth2{
			ServerURL: os.Getenv(aha.AhaServerURL),
			Token: &oauth2.Token{
				AccessToken: os.Getenv(aha.AhaAPIKeyEnv),
			},
		},
	}

	sclient, err := creds.NewSimpleClient(context.Background())
	logutil.FatalErr(err)

	clientUtil := aha.NewClientUtil(nil)
	clientUtil.SetSimpleClient(sclient)

	user, err := clientUtil.GetSCIMUser()
	if err != nil {
		panic(err)
	}
	fmtutil.MustPrintJSON(user)

	fmt.Println("DONE")
}
