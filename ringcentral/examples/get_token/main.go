package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/oauth2more/ringcentral"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	EnvPath string `short:"e" long:"envPath" description:"Environment File Path"`
	EnvVar  string `short:"v" long:"envVar" description:"Environment Variable Name"`
	Token   string `short:"t" long:"token" description:"Token"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	err = config.LoadDotEnvSkipEmpty(opts.EnvPath, os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		log.Fatal(errors.Wrap(err, "E_LOAD_DOT_ENV"))
	}

	if len(opts.EnvVar) > 0 {
		if len(os.Getenv(opts.EnvVar)) == 0 {
			log.Fatal("E_NO_VAR")
		}
		ac := ringcentral.ApplicationConfig{}
		err := json.Unmarshal([]byte(os.Getenv(opts.EnvVar)), &ac)
		if err != nil {
			log.Fatal(
				errors.Wrap(
					err, fmt.Sprintf("E_JSON_UNMARSHAL [%v]", os.Getenv(opts.EnvVar))))
		}
		fmtutil.PrintJSON(ac)
		token, err := ringcentral.NewTokenPassword(
			ac.ApplicationCredentials(),
			ac.PasswordCredentials())
		if err != nil {
			log.Fatal(err)
		}
		token.Expiry = token.Expiry.UTC()

		fmtutil.PrintJSON(token)
	}

	if 1 == 0 {
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
			log.Fatal(err)
		}
		token.Expiry = token.Expiry.UTC()

		fmtutil.PrintJSON(token)
	}

	fmt.Println("DONE")
}
