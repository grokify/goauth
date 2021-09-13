package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/grokify/oauth2more/credentials"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/net/http/httpsimple"
	"github.com/grokify/simplego/net/httputilmore"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
)

type Options struct {
	credentials.Options
	URL    string   `short:"U" long:"url" description:"URL"`
	Method string   `short:"X" long:"request" description:"Method"`
	Header []string `short:"H" long:"header" description:"HTTP Headers"`
	Body   string   `short:"d" long:"data" description:"HTTP Body"`
}

func (opts *Options) SimpleRequest() (httpsimple.SimpleRequest, error) {
	sr := httpsimple.SimpleRequest{
		URL:     opts.URL,
		Headers: map[string][]string{},
	}
	// Method
	if len(strings.TrimSpace(opts.Method)) > 0 {
		m, err := httputilmore.ParseHTTPMethod(opts.Method)
		if err != nil {
			return sr, err
		}
		sr.Method = string(m)
	} else {
		sr.Method = http.MethodGet
	}
	for _, h := range opts.Header {
		hparts := strings.SplitN(h, ":", 2)
		if len(hparts) != 2 {
			return sr, fmt.Errorf("could not parse header [%s]", h)
		}
		hname := strings.TrimSpace(hparts[0])
		if sr.Headers[hname] == nil {
			sr.Headers[hname] = []string{}
		}
		sr.Headers[hname] = append(sr.Headers[hname], strings.TrimSpace(hparts[1]))
	}
	if len(opts.Body) > 0 {
		sr.Body = opts.Body
	}
	return sr, nil
}

var rxParseHeader = regexp.MustCompile(`^([^:]+):(.+)$`)

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal().Err(err).Msg("required properties not present")
		panic("Z")
	}
	fmtutil.PrintJSON(opts)

	if 1 == 1 {
		sr, err := opts.SimpleRequest()
		if err != nil {
			log.Fatal().Err(err).Msg("simple request failure")
			panic("Z")
		}
		fmtutil.PrintJSON(sr)
	}

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
