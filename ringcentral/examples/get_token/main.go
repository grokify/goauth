package main

import (
	"fmt"
	"os"

	"github.com/grokify/oauth2more/ringcentral"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
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
	fmtutil.PrintJSON(opts)

	files, err := config.LoadDotEnv(opts.EnvPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "E_LOAD_DOT_ENV"))
	}
	fmtutil.PrintJSON(files)

	if len(opts.EnvVar) > 0 {
		if len(os.Getenv(opts.EnvVar)) == 0 {
			log.Fatal("E_NO_VAR")
		}

		fmt.Println(os.Getenv(opts.EnvVar))

		credentials, err := ringcentral.NewCredentialsJSON([]byte(os.Getenv(opts.EnvVar)))
		if err != nil {
			log.Fatal(
				errors.Wrap(
					err, fmt.Sprintf("E_JSON_UNMARSHAL [%v]", os.Getenv(opts.EnvVar))))
		}
		fmtutil.PrintJSON(credentials)
		token, err := credentials.NewToken()
		if err != nil {
			log.Fatal(err)
		}

		token.Expiry = token.Expiry.UTC()

		fmtutil.PrintJSON(token)
	} else {
		fmt.Printf("No EnvVar [-v]\n")
	}

	fmt.Println("DONE")
}
