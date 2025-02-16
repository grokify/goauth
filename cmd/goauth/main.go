package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/log/logutil"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	cli := goauth.CLIRequest{}
	_, err := flags.Parse(&cli)
	logutil.FatalErr(err)

	err = cli.Do(context.Background(), os.Stdout)
	logutil.FatalErr(err)

	fmt.Println("\nDONE")
}
