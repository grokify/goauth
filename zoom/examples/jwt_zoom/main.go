package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/zoom"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
)

func main() {
	files, err := config.LoadDotEnv(
		".env", os.Getenv("ENV_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(files)

	apiKey := os.Getenv(zoom.EnvZoomApiKey)
	apiSecret := os.Getenv(zoom.EnvZoomApiSecret)

	tokenString := ""
	if 1 == 1 {
		_, tokenString, err = zoom.CreateJwtToken(apiKey, apiSecret, time.Hour)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tokenString = "tmpToken"
	}

	fmt.Printf("TOK [%v]\n", tokenString)

	if 1 == 0 {
		token, err := oauth2more.ParseJwtTokenString(
			tokenString, apiSecret,
			&jwt.StandardClaims{Issuer: apiKey})
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(token.Claims)
	}

	client := zoom.NewClientToken(tokenString)

	resp, err := client.Get("https://api.zoom.us/v2/users/")
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("RESP: %s\n", string(bytes))

	cu := zoom.NewClientUtil(client)
	scimUser, err := cu.GetSCIMUser()
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(cu.UserNative)
	fmtutil.PrintJSON(scimUser)

	fmt.Println("DONE")
}
