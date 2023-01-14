package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/grokify/goauth/credentials"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
	flags "github.com/jessevdk/go-flags"
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
			sr.BodyType = httpsimple.BodyTypeJSON
		}
	}
	return sr, nil
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)
	fmtutil.MustPrintJSON(opts)

	creds, err := credentials.ReadCredentialsFromFile(
		opts.Options.CredsPath, opts.Options.Account, true)
	logutil.FatalErr(err)

	var httpClient *http.Client
	if opts.Options.UseCLI() {
		httpClient, err = creds.NewClientCLI("mystate")
	} else {
		httpClient, err = creds.NewClient(context.Background())
	}
	logutil.FatalErr(err)

	sr, err := opts.SimpleRequest()
	logutil.FatalErr(err)

	fmtutil.MustPrintJSON(sr)
	sclient, err := creds.NewSimpleClientHTTP(httpClient)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	resp, err := sclient.Do(sr)
	logutil.FatalErr(err)

	fmt.Printf("STATUS [%d]", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	logutil.FatalErr(err)

	fmt.Println(string(body))

	fmt.Println("DONE")
}
