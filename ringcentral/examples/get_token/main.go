package main

import (
	"fmt"
	"strings"

	"github.com/grokify/oauth2more/credentials"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type Options struct {
	CredsPath string `short:"c" long:"credspath" description:"Environment File Path"`
	Account   string `short:"a" long:"account" description:"Environment Variable Name"`
	Token     string `short:"t" long:"token" description:"Token"`
	CLI       []bool `long:"cli" description:"CLI"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal().Err(err)
	}
	fmtutil.PrintJSON(opts)

	cset, err := credentials.ReadFileCredentialsSet(opts.CredsPath, true)
	if err != nil {
		log.Fatal().Err(err).
			Str("credentials_filepath", opts.CredsPath).
			Msg("cannot read credentials file")
	}
	if len(strings.TrimSpace(opts.Account)) == 0 {
		log.Fatal().Err(err).
			Strs("available accounts", cset.Accounts()).
			Msg("no account specified")
	}
	credentials, err := cset.Get(opts.Account)
	if err != nil {
		log.Fatal().Err(err).
			Str("credentials_account", opts.Account).
			Msg("cannot find credentials account")
		panic("fail1")
	}

	var token *oauth2.Token

	if len(opts.CLI) > 0 {
		token, err = credentials.NewTokenCli("mystate")
	} else {
		token, err = credentials.NewToken()
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

	fmtutil.PrintJSON(token)

	fmt.Println("DONE")
}
