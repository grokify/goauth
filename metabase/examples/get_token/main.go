package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/gotilla/type/stringsutil"
	"github.com/grokify/oauth2more/metabase"
)

func main() {
	loaded, err := config.LoadDotEnvSkipEmptyInfo(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(loaded)

	if len(os.Getenv(metabase.EnvMetabaseUsername)) == 0 {
		log.Fatal("E_NO_METABASE_USERNAME")
	}

	cfg := metabase.Config{
		BaseURL:       os.Getenv(metabase.EnvMetabaseBaseURL),
		SessionID:     os.Getenv(metabase.EnvMetabaseSessionID),
		Username:      os.Getenv(metabase.EnvMetabaseUsername),
		Password:      os.Getenv(metabase.EnvMetabasePassword),
		TLSSkipVerify: stringsutil.ToBool(os.Getenv(metabase.EnvMetabaseTLSSkipVerify))}
	fmtutil.PrintJSON(cfg)

	_, authResponse, err := metabase.NewClient(cfg)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("AUTH_RESPONSE:")
	fmtutil.PrintJSON(authResponse)

	fmt.Println("DONE")
}
