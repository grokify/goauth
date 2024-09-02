package goauth

// Options is a struct to be used with `github.com/jessevdk/go-flags`.
// It can be embedded in another struct.
type Options struct {
	CredsPath string `long:"creds" description:"Environment File Path" required:"true"`
	Account   string `long:"account" description:"Environment Variable Name"`
	Token     string `long:"token" description:"Token"`
	CLI       []bool `long:"cli" description:"CLI"`
}

func (opts *Options) Credentials() (Credentials, error) {
	return ReadFileCredentialsFromCredentialsSet(opts.CredsPath, opts.Account)
}

func (opts *Options) UseCLI() bool {
	return len(opts.CLI) > 0
}
