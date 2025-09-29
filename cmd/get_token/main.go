package main

import (
	"context"
	"fmt"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/fmt/fmtutil"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func main() {
	opts := goauth.Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal().Err(err)
	}
	fmtutil.MustPrintJSON(opts)

	creds, err := goauth.NewCredentialsFromSetFile(
		opts.CredsPath, opts.Account, true)
	if err != nil {
		log.Fatal().Err(err).
			Str("credsPath", opts.CredsPath).
			Str("accountKey", opts.Account).
			Msg("cannot read credentials")
	}

	var token *oauth2.Token

	if len(opts.CLI) > 0 {
		token, err = creds.NewTokenCLI(context.Background(), "mystate")
	} else {
		token, err = creds.NewToken(context.Background())
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
