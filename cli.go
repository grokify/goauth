package goauth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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

func (opts *Options) NewClient(ctx context.Context, state string) (*http.Client, error) {
	if creds, err := opts.Credentials(); err != nil {
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

func (cli CLIRequest) Do(ctx context.Context, w io.Writer) error {
	if creds, err := cli.Options.Credentials(); err != nil {
		return errorsutil.NewErrorWithLocation(err.Error())
	} else if tok, err := creds.NewToken(ctx); err != nil {
		return errorsutil.NewErrorWithLocation(err.Error())
	} else if sr, err := cli.Request.Request(); err != nil {
		return errorsutil.NewErrorWithLocation(err.Error())
	} else {
		if at := strings.TrimSpace(tok.AccessToken); at != "" {
			sr.Headers.Add(httputilmore.HeaderAuthorization, authutil.TokenBearer+" "+at)
		}
		resp, err := sr.Do(ctx)
		if err != nil {
			return errorsutil.NewErrorWithLocation(err.Error())
		}
		if w != nil {
			if _, err := w.Write([]byte(fmt.Sprintf("Response Status Code: %d\n", resp.StatusCode))); err != nil {
				return errorsutil.NewErrorWithLocation(err.Error())
			}
		}
		b, err := httputilmore.ResponseBodyMore(resp, "", "  ")
		if err != nil {
			return errorsutil.NewErrorWithLocation(err.Error())
		} else {
			if w != nil {
				if _, err := w.Write([]byte(fmt.Sprintf("===== BEGIN BODY =====\n%s\n===== END BODY =====", string(b)))); err != nil {
					return errorsutil.NewErrorWithLocation(err.Error())
				}
			}
			return nil
		}
	}
}
