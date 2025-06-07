package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/grokify/goauth"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	cli := goauth.CLIRequest{}
	_, err := flags.Parse(&cli)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = cli.Do(context.Background(), "", os.Stdout)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}

	fmt.Println("\nDONE")
	os.Exit(0)
}
