package goauth

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/jessevdk/go-flags"
)

// Options is a struct to be used with `ParseOptions()` or `github.com/jessevdk/go-flags`.
// It can be embedded in another struct and used directly with `github.com/jessevdk/go-flags`.
type Options struct {
	CredsPath string `long:"creds" description:"Environment File Path" required:"true"`
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

func (opts *Options) NewClient(ctx context.Context) (*http.Client, error) {
	if creds, err := opts.Credentials(); err != nil {
		return nil, errorsutil.Wrap(err, "error in `goauth.Options.NewClient() call to self as `goauth.Options.Credentials()`")
	} else {
		return creds.NewClient(ctx)
	}
}

func (opts *Options) UseCLI() bool {
	return len(opts.CLI) > 0
}

// CLIRequest will get a token using `goauth` and then execute the provided request
// paramters with the credential, e.g. OAuth 2.0 access token.
type CLIRequest struct {
	Options
	Request httpsimple.CLI
}

func (cli CLIRequest) Do(ctx context.Context, w io.Writer) error {
	if creds, err := cli.Options.Credentials(); err != nil {
		return err
	} else if tok, err := creds.NewToken(ctx); err != nil {
		return err
	} else if sr, err := cli.Request.Request(); err != nil {
		return err
	} else {
		if at := strings.TrimSpace(tok.AccessToken); at != "" {
			sr.Headers.Add(httputilmore.HeaderAuthorization, authutil.TokenBearer+" "+at)
		}
		if resp, err := sr.Do(ctx); err != nil {
			return err
		} else if b, err := httputilmore.ResponseBodyMore(resp, "", "  "); err != nil {
			return err
		} else {
			if w != nil {
				if _, err := w.Write(b); err != nil {
					return err
				}
			}
			return nil
		}
	}
}
