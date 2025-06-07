package goauth

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/type/maputil"
	"github.com/jessevdk/go-flags"
)

// Options is a struct to be used with `ParseOptions()` or `github.com/jessevdk/go-flags`.
// It can be embedded in another struct and used directly with `github.com/jessevdk/go-flags`.
type Options struct {
	CredsPath string `long:"creds" description:"Environment File Path"`
	Account   string `long:"account" description:"Environment Variable Name"`
	Token     string `long:"token" description:"Token"`
	CLI       []bool `long:"cli" description:"CLI"`
}

func ParseOptions() (*Options, error) {
	opts := &Options{}
	_, err := flags.Parse(opts)
	return opts, err
}

func (opts *Options) Credentials() (Credentials, error) {
	return ReadCredentialsFromSetFile(opts.CredsPath, opts.Account, false)
}

func (opts *Options) CredentialsSet(inflateEndpoints bool) (*CredentialsSet, error) {
	return ReadFileCredentialsSet(opts.CredsPath, inflateEndpoints)
}

func (opts *Options) NewClient(ctx context.Context, state string) (*http.Client, error) {
	if opts.CredsPath == "" && opts.Account == "" {
		return &http.Client{}, nil
	} else if creds, err := opts.Credentials(); err != nil {
		return nil, errorsutil.Wrap(err, "error in `goauth.Options.NewClient() call to self as `goauth.Options.Credentials()`")
	} else if opts.UseCLI() {
		if state == "" {
			state = time.Now().Format(time.RFC3339)
		}
		return creds.NewClientCLI(ctx, state)
	} else {
		return creds.NewClient(ctx)
	}
}

func (opts *Options) UseCLI() bool {
	return len(opts.CLI) > 0
}

// CLIRequest will get a token using `goauth` and then execute the provided request
// parameters with the credential, e.g. OAuth 2.0 access token.
type CLIRequest struct {
	Options
	Request httpsimple.CLI
}

func (cli CLIRequest) Do(ctx context.Context, state string, w io.Writer) error {
	if clt, err := cli.NewClient(ctx, state); err != nil {
		return errorsutil.WrapWithLocation(err)
	} else if sr, err := cli.Request.Request(); err != nil {
		return errorsutil.WrapWithLocation(err)
	} else {
		sc := httpsimple.NewClient(clt, "")
		resp, err := sc.Do(ctx, sr)
		if err != nil {
			return errorsutil.WrapWithLocation(err)
		}
		if _, err := fmtutil.FprintfIf(w, "Response Status Code: %d\n", resp.StatusCode); err != nil {
			return errorsutil.WrapWithLocation(err)
		}
		if _, err := fmtutil.FprintfIf(w, "===== BEGIN RESPONSE META =====\nStatus Code: %d\n===== END RESPONSE META =====\n", resp.StatusCode); err != nil {
			return errorsutil.WrapWithLocation(err)
		}

		if _, err := fmtutil.FprintIf(w, "===== BEGIN RESPONSE HEADERS =====\n"); err != nil {
			return errorsutil.WrapWithLocation(err)
		}
		hkeys := maputil.Keys(resp.Header)
		for _, header := range hkeys {
			if v, ok := resp.Header[header]; ok {
				for _, vi := range v {
					if _, err := fmtutil.FprintfIf(w, "%s: %s\n", header, vi); err != nil {
						return errorsutil.WrapWithLocation(err)
					}
				}
			}
		}
		if _, err := fmtutil.FprintIf(w, "===== END RESPONSE HEADERS =====\n"); err != nil {
			return errorsutil.WrapWithLocation(err)
		}

		b, err := httputilmore.ResponseBodyMore(resp, "", "  ")
		if err != nil {
			return errorsutil.WrapWithLocation(err)
		} else {
			if w != nil {
				if _, err := fmtutil.FprintfIf(w, "===== BEGIN RESPONSE BODY =====\n%s\n===== END RESPONSE BODY =====", string(b)); err != nil {
					return errorsutil.WrapWithLocation(err)
				}
			}
			return nil
		}
	}
}
