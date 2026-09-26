package main

import (
	"context"
	"fmt"
	"time"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/authutil"
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

	creds, err := goauth.NewCredentialsFromSetFile(opts.CredsPath, opts.Account, true)
	logutil.FatalErr(err)

	tok, err := creds.NewOrExistingValidToken(ctx)
	logutil.FatalErr(err)
	// Print token metadata only; never print the token or credentials.
	fmt.Printf("Token type: %s, expires in: %s\n",
		tok.Type(), time.Until(tok.Expiry).Round(time.Second))

	clt := authutil.NewClientTokenOAuth2(tok)

	cu := google.ClientUtil{Client: clt}

	ui, err := cu.GetUserinfo(ctx)
	logutil.FatalErr(err)
	logutil.FatalErr(fmtutil.PrintJSON(ui))

	fmt.Println("DONE")
}
