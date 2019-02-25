package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/oauth2more/metabase"
)

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		panic(err)
	}

	_, authResponse, err := metabase.NewClientPasswordWithSessionId(
		os.Getenv("METABASE_BASE_URL"),
		os.Getenv("METABASE_USERNAME"),
		os.Getenv("METABASE_PASSWORD"),
		os.Getenv("METABASE_SESSION_ID"),
		true)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("AUTH_RESPONSE:")
	fmtutil.PrintJSON(authResponse)

	fmt.Println("DONE")
}
