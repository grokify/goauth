package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/grokify/goauth"
	"github.com/grokify/goauth/zoom"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/net/http/httpsimple"
	"github.com/grokify/simplego/net/urlutil"
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
		token, err := goauth.ParseJwtTokenString(
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
	bytes, err := io.ReadAll(resp.Body)
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

	resp, err = createMeeting(client)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))

	fmt.Println("DONE")
}

func createMeeting(client *http.Client) (*http.Response, error) {
	sc := httpsimple.SimpleClient{
		BaseURL:    zoom.ZoomAPIBaseURL,
		HTTPClient: client}
	req := httpsimple.SimpleRequest{
		Method: http.MethodPost,
		URL:    urlutil.JoinAbsolute(zoom.ZoomAPIMeURL, "meetings"),
		Body:   []byte(reqBody),
		IsJSON: true}
	return sc.Do(req)
}

const reqBody = `{
	"topic":"MeetingOne",
	"type":2,
	"start_time":"2021-07-04T00:00:00Z",
	"duration":30,
	"agenda":"meet and greet"
}`
