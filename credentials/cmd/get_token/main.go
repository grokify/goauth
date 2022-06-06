package main

import (
	"fmt"

	"github.com/grokify/goauth/credentials"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func main() {
	opts := credentials.Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal().Err(err)
	}
	fmtutil.MustPrintJSON(opts)

	creds, err := credentials.ReadCredentialsFromFile(
		opts.CredsPath, opts.Account, true)
	if err != nil {
		log.Fatal().Err(err).
			Str("credsPath", opts.CredsPath).
			Str("accountKey", opts.Account).
			Msg("cannot read credentials")
	}

	var token *oauth2.Token

	if len(opts.CLI) > 0 {
		token, err = creds.NewTokenCLI("mystate")
	} else {
		token, err = creds.NewToken()
	}
	if err != nil {
		log.Fatal().
			Err(err).
			Str("filepath", opts.CredsPath).
			Str("account", opts.Account).
			Msg("failed to get new token")
		panic("failed")
	}

	token.Expiry = token.Expiry.UTC()

	fmtutil.MustPrintJSON(token)

	fmt.Println("DONE")
}
