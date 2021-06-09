package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grokify/oauth2more/metabase"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/type/stringsutil"
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

	fmt.Printf("EXAMPLE_COMMAND:\ncurl -XGET '%s' -H '%s: %s'\n",
		metabase.BuildURL(cfg.BaseURL, metabase.RelPathApiDatabase),
		metabase.MetabaseSessionHeader,
		authResponse.Id)

	fmt.Println("DONE")
}
