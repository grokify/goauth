package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/grokify/goauth/aha"
	"github.com/grokify/mogo/fmt/fmtutil"
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
		os.Getenv(aha.AhaAPIKeyEnv),
	)

	clientUtil := aha.NewClientUtil(client)

	user, err := clientUtil.GetSCIMUser()
	if err != nil {
		panic(err)
	}
	fmtutil.PrintJSON(user)

	fmt.Println("DONE")
}
