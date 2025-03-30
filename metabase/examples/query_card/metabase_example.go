package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grokify/goauth/metabase"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
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
		log.Fatal(err)
	}

	cardID := 1

	metabase.TLSInsecureSkipVerify = true

	baseURL := os.Getenv("METABASE_BASE_URL")

	client, _, err := metabase.NewClientPassword(
		baseURL,
		os.Getenv("METABASE_USERNAME"),
		os.Getenv("METABASE_PASSWORD"),
		metabase.TLSInsecureSkipVerify)
	if err != nil {
		log.Fatal(err)
	}

	userURL := urlutil.JoinAbsolute(baseURL, "api/user/current")

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	err = httputilmore.PrintResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	cardURL := fmt.Sprintf("api/card/%v/query/%s", cardID, "json")
	cardURL = urlutil.JoinAbsolute(baseURL, cardURL)

	fmt.Println(cardURL)

	req, err = http.NewRequest(http.MethodPost, cardURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	err = httputilmore.PrintResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
