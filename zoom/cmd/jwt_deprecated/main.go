package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/grokify/goauth/authutil/jwtutil"
	"github.com/grokify/goauth/zoom"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/urlutil"
)

func main() {
	files, err := config.LoadDotEnv([]string{".env", os.Getenv("ENV_PATH")}, -1)
	logutil.FatalErr(err)

	fmtutil.MustPrintJSON(files)

	apiKey := os.Getenv(zoom.EnvZoomAPIKey)
	apiSecret := os.Getenv(zoom.EnvZoomAPISecret)

	tokenString := ""

	_, tokenString, err = zoom.CreateJWTToken(apiKey, apiSecret, time.Hour)
	logutil.FatalErr(err)

	fmt.Printf("TOK [%v]\n", tokenString)

	if 1 == 0 {
		token, err := jwtutil.ParseJWTString(
			tokenString, apiSecret,
			&jwt.RegisteredClaims{Issuer: apiKey})
		logutil.FatalErr(err)

		fmtutil.MustPrintJSON(token.Claims)
	}

	client := zoom.NewClientToken(tokenString)

	resp, err := client.Get("https://api.zoom.us/v2/users/")
	logutil.FatalErr(err)

	bytes, err := io.ReadAll(resp.Body)
	logutil.FatalErr(err)

	fmt.Printf("RESP: %s\n", string(bytes))

	cu := zoom.NewClientUtil(client)
	scimUser, err := cu.GetSCIMUser()
	logutil.FatalErr(err)

	fmtutil.MustPrintJSON(cu.UserNative)
	fmtutil.MustPrintJSON(scimUser)

	resp, err = createMeeting(client)
	logutil.FatalErr(err)

	bytes, err = io.ReadAll(resp.Body)
	logutil.FatalErr(err)

	fmt.Println(string(bytes))

	fmt.Println("DONE")
}

func createMeeting(client *http.Client) (*http.Response, error) {
	sc := httpsimple.Client{
		BaseURL:    zoom.ZoomAPIURLBase,
		HTTPClient: client}
	req := httpsimple.Request{
		Method:   http.MethodPost,
		URL:      urlutil.JoinAbsolute(zoom.ZoomAPIURLUsersMe, "meetings"),
		Body:     []byte(reqBody),
		BodyType: httpsimple.BodyTypeJSON}
	return sc.Do(req)
}

const reqBody = `{
	"topic":"MeetingOne",
	"type":2,
	"start_time":"2021-07-04T00:00:00Z",
	"duration":30,
	"agenda":"meet and greet"
}`
