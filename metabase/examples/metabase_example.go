package main

import (
	"fmt"
	"net/http"
	"os"

	hum "github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/net/urlutil"
	"github.com/grokify/oauth2util/metabase"
	"github.com/joho/godotenv"
)

func loadEnv() error {
	if len(os.Getenv("ENV_PATH")) > 0 {
		return godotenv.Load(os.Getenv("ENV_PATH"))
	}
	return godotenv.Load()
}

func main() {
	err := loadEnv()
	if err != nil {
		panic(err)
	}

	baseUrl := os.Getenv("METABASE_BASE_URL")

	client, err := metabase.NewClient(baseUrl,
		os.Getenv("METABASE_USERNAME"),
		os.Getenv("METABASE_PASSWORD"),
	)
	if err != nil {
		panic(err)
	}

	userUrl := urlutil.JoinAbsolute(baseUrl, "api/user/current")

	req, err := http.NewRequest("GET", userUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	hum.PrintResponse(resp, true)

	fmt.Println("DONE")
}
