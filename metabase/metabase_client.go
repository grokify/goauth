package metabase

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/gotilla/config"
	hum "github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/net/urlutil"
	om "github.com/grokify/oauth2more"
)

const (
	MetabaseSessionHeader = "X-Metabase-Session"
	RelPathApiSession     = "api/session"
	RelPathApiUserCurrent = "api/user/current"
	RelPathApiDatabase    = "api/database"

	EnvBaseURL   = "METABASE_BASE_URL"
	EnvSessionId = "METABASE_SESSION_ID"
	EnvUsername  = "METABASE_USERNAME"
	EnvPassword  = "METABASE_PASSWORD"
)

var (
	TLSInsecureSkipVerify = false
)

type authRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthResponse struct {
	Id string `json:"id,omitempty"`
}

// NewClient returns a *http.Client that will add the Metabase Session
// header to each request.
func NewClientPassword(baseUrl, username, password string, tlsSkipVerify bool) (*http.Client, *AuthResponse, error) {
	resp, err := AuthRequest(
		urlutil.JoinAbsolute(baseUrl, RelPathApiSession),
		username,
		password,
		tlsSkipVerify)
	if err != nil {
		return nil, nil, err
	}

	res := &AuthResponse{}
	err = hum.UnmarshalResponseJSON(resp, res)
	if err != nil {
		return nil, res, err
	}

	return NewClientSessionId(res.Id, tlsSkipVerify), res, nil
}

// NewClientPasswordWithSessionId returns a *http.Client first attempting to use
// the supplied `sessionId` with a fallback to `username` and `password`.
func NewClientPasswordWithSessionId(baseUrl, username, password, sessionId string, tlsSkipVerify bool) (*http.Client, *AuthResponse, error) {
	sessionId = strings.TrimSpace(sessionId)
	if len(sessionId) > 0 {
		httpClient := NewClientSessionId(sessionId, tlsSkipVerify)
		userUrl := urlutil.JoinAbsolute(baseUrl, RelPathApiUserCurrent)
		resp, err := httpClient.Get(userUrl)
		if err == nil && resp.StatusCode == 200 {
			return httpClient, nil, nil
		}
	}
	return NewClientPassword(baseUrl, username, password, tlsSkipVerify)
}

func NewClientSessionId(sessionId string, tlsSkipVerify bool) *http.Client {
	client := &http.Client{}

	header := http.Header{}
	header.Add(MetabaseSessionHeader, sessionId)

	if tlsSkipVerify {
		client = om.ClientTLSInsecureSkipVerify(client)
	}

	client.Transport = hum.TransportWithHeaders{
		Transport: client.Transport,
		Header:    header}

	return client
}

type InitConfig struct {
	LoadEnv              bool
	EnvPath              string
	EnvMetabaseBaseUrl   string
	EnvMetabaseSessionId string
	EnvMetabaseUsername  string
	EnvMetabasePassword  string
	TlsSkipVerify        bool
}

func NewClientEnv(cfg InitConfig) (*http.Client, *AuthResponse, error) {
	if cfg.LoadEnv && len(cfg.EnvPath) > 0 {
		err := config.LoadDotEnvSkipEmpty(os.Getenv(cfg.EnvPath), "./.env")
		if err != nil {
			return nil, nil, err
		}
	}

	var httpClient *http.Client
	var authResponse *AuthResponse

	if len(os.Getenv(cfg.EnvMetabaseSessionId)) > 0 {
		httpClient = NewClientSessionId(os.Getenv(cfg.EnvMetabaseSessionId), true)
	} else {
		httpClient2, res, err := NewClientPassword(
			os.Getenv(cfg.EnvMetabaseBaseUrl),
			os.Getenv(cfg.EnvMetabaseUsername),
			os.Getenv(cfg.EnvMetabasePassword),
			cfg.TlsSkipVerify)
		if err != nil {
			return nil, authResponse, err
		}
		authResponse = res
		httpClient = httpClient2
	}
	return httpClient, authResponse, nil
}

// AuthRequest creates an authentiation request that returns a id that is used
// in Metabase API requests. It follows the following curl command:
// curl -v -H "Content-Type: application/json" -d '{"username":"myusername","password":"mypassword"}' -XPOST 'http://example.com/api/session'
func AuthRequest(authUrl, username, password string, tlsSkipVerify bool) (*http.Response, error) {
	bodyBytes, err := json.Marshal(authRequest{Username: username, Password: password})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, authUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add(hum.HeaderContentType, hum.ContentTypeAppJsonUtf8)

	client := &http.Client{}

	if tlsSkipVerify {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsSkipVerify},
		}
	}

	return client.Do(req)
}
