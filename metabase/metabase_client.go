package metabase

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"

	hum "github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/net/urlutil"
)

const (
	MetabaseSessionHeader = "X-Metabase-Session"
	RelPathApiSession     = "api/session"
	RelPathApiUserCurrent = "api/user/current"
)

var (
	TLSInsecureSkipVerify = false
)

type authRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type authResponse struct {
	Id string `json:"id,omitempty"`
}

// NewClient returns a *http.Client that will add the Metabase Session
// header to each request.
func NewClient(baseUrl, username, password string) (*http.Client, error) {
	resp, err := AuthRequest(
		urlutil.JoinAbsolute(baseUrl, RelPathApiSession),
		username,
		password)
	if err != nil {
		return nil, err
	}

	res := &authResponse{}
	err = hum.UnmarshalResponseJSON(resp, res)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Add(MetabaseSessionHeader, res.Id)

	client := &http.Client{}

	if TLSInsecureSkipVerify {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: TLSInsecureSkipVerify},
		}
	}

	client.Transport = hum.TransportWithHeaders{
		Transport: client.Transport,
		Header:    header}

	return client, nil
}

func AuthRequest(authUrl, username, password string) (*http.Response, error) {
	bodyBytes, err := json.Marshal(authRequest{Username: username, Password: password})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, authUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add(hum.ContentTypeHeader, hum.ContentTypeValueJSONUTF8)

	client := &http.Client{}

	if TLSInsecureSkipVerify {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: TLSInsecureSkipVerify},
		}
	}

	return client.Do(req)
}
