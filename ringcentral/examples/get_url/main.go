package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/grokify/oauth2more/credentials"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
)

type Options struct {
	CredsPath string `short:"c" long:"credspath" description:"Environment File Path" required:"true"`
	Account   string `short:"a" long:"account" description:"Environment Variable Name" required:"false"`
	URL       string `short:"u" long:"url" description:"URL"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal().Err(err).Msg("required properties not present")
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
	sclient, err := cset.NewSimpleClient(opts.Account)
	if err != nil {
		fmt.Println(string(err.Error()))
		log.Fatal().Err(err).
			Msg("cannot create simpleclient")
	}

	resp, err := sclient.Get(opts.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("get URL error")
	}
	fmt.Printf("STATUS [%d]", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("parse body error")
	}
	fmt.Println(string(body))

	fmt.Println("DONE")
}
