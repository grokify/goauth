package goauth

import (
	"context"
	"net/http"

	"github.com/grokify/mogo/errors/errorsutil"
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
