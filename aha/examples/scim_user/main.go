package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/grokify/goauth/aha"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
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

	sclient, err := aha.NewSimpleClient(
		os.Getenv(aha.AhaAccountEnv),
		os.Getenv(aha.AhaAPIKeyEnv),
	)
	logutil.FatalErr(err)

	clientUtil := aha.NewClientUtil(nil)
	clientUtil.SetSimpleClient(sclient)

	user, err := clientUtil.GetSCIMUser()
	if err != nil {
		panic(err)
	}
	fmtutil.MustPrintJSON(user)

	fmt.Println("DONE")
}
