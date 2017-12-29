package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/oauth2more/aha"
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

	client := aha.NewClient(
		os.Getenv(aha.AhaAccountEnv),
		os.Getenv(aha.AhaApiKeyEnv),
	)

	clientUtil := aha.NewClientUtil(client)

	user, err := clientUtil.GetSCIMUser()
	if err != nil {
		panic(err)
	}
	fmtutil.PrintJSON(user)

	fmt.Println("DONE")
}
