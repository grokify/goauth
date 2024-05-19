package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

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

	creds := zoom.ServerToServerOAuth2Credentials(
		os.Getenv(zoom.EnvZoomClientID),
		os.Getenv(zoom.EnvZoomCLientSecret),
		os.Getenv(zoom.EnvZoomApplicationID))

	tok, err := creds.NewToken(context.Background())
	logutil.FatalErr(err)

	fmt.Printf("TOK [%s]\n", tok.AccessToken)

	client := zoom.NewClientToken(tok.AccessToken)

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
