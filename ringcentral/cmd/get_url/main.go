package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/grokify/goauth/credentials"
	"github.com/grokify/gohttp/httpsimple"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/net/httputilmore"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	credentials.Options
	URL    string   `short:"U" long:"url" description:"URL" required:"true"`
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
		if strings.Index(strings.TrimSpace(opts.Body), "{") == 0 {
			sr.IsJSON = true
		}
	}
	return sr, nil
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	logutil.FatalOnError(err)
	fmtutil.MustPrintJSON(opts)

	creds, err := credentials.ReadCredentialsFromFile(
		opts.CredsPath, opts.Account, true)
	logutil.FatalOnError(err)

	var httpClient *http.Client
	if opts.UseCLI() {
		httpClient, err = creds.NewClientCli("mystate")
	} else {
		httpClient, err = creds.NewClient(context.Background())
	}
	logutil.FatalOnError(err)

	sr, err := opts.SimpleRequest()
	logutil.FatalOnError(err)

	fmtutil.MustPrintJSON(sr)
	sclient, err := creds.NewSimpleClient(httpClient)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	resp, err := sclient.Do(sr)
	logutil.FatalOnError(err)

	fmt.Printf("STATUS [%d]", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	logutil.FatalOnError(err)

	fmt.Println(string(body))

	fmt.Println("DONE")
}
