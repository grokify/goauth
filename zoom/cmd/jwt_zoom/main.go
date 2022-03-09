package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/grokify/goauth"
	"github.com/grokify/goauth/zoom"
	"github.com/grokify/gohttp/httpsimple"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/net/urlutil"
)

func main() {
	files, err := config.LoadDotEnv(
		".env", os.Getenv("ENV_PATH"))
	logutil.FatalOnError(err)

	fmtutil.MustPrintJSON(files)

	apiKey := os.Getenv(zoom.EnvZoomApiKey)
	apiSecret := os.Getenv(zoom.EnvZoomApiSecret)

	tokenString := ""

	_, tokenString, err = zoom.CreateJwtToken(apiKey, apiSecret, time.Hour)
	logutil.FatalOnError(err)

	fmt.Printf("TOK [%v]\n", tokenString)

	if 1 == 0 {
		token, err := goauth.ParseJwtTokenString(
			tokenString, apiSecret,
			&jwt.StandardClaims{Issuer: apiKey})
		logutil.FatalOnError(err)

		fmtutil.MustPrintJSON(token.Claims)
	}

	client := zoom.NewClientToken(tokenString)

	resp, err := client.Get("https://api.zoom.us/v2/users/")
	logutil.FatalOnError(err)

	bytes, err := io.ReadAll(resp.Body)
	logutil.FatalOnError(err)

	fmt.Printf("RESP: %s\n", string(bytes))

	cu := zoom.NewClientUtil(client)
	scimUser, err := cu.GetSCIMUser()
	logutil.FatalOnError(err)

	fmtutil.MustPrintJSON(cu.UserNative)
	fmtutil.MustPrintJSON(scimUser)

	resp, err = createMeeting(client)
	logutil.FatalOnError(err)

	bytes, err = io.ReadAll(resp.Body)
	logutil.FatalOnError(err)

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
