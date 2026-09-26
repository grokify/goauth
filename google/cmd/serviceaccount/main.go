package main

import (
	"context"
	"fmt"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/google"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/jessevdk/go-flags"
)

func main() {
	opts := goauth.Options{}
	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)

	ctx := context.Background()

	// The credentials include the service account private key: never print them.
	creds, err := goauth.NewCredentialsFromSetFile(opts.CredsPath, opts.Account, true)
	logutil.FatalErr(err)

	clt, err := creds.NewClient(ctx)
	logutil.FatalErr(err)

	cu := google.ClientUtil{Client: clt}

	ui, err := cu.GetUserinfo(ctx)
	logutil.FatalErr(err)
	logutil.FatalErr(fmtutil.PrintJSON(ui))

	fmt.Println("DONE")
}
