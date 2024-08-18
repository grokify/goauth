package metabase

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	HeaderMetabaseSession = "X-Metabase-Session"
	RelPathAPIDatabase    = "api/database"
	RelPathAPISession     = "api/session"
	RelPathAPIUserCurrent = "api/user/current"

	// Example environment variables
	EnvMetabaseBaseURL       = "METABASE_BASE_URL"
	EnvMetabaseUsername      = "METABASE_USERNAME"
	EnvMetabasePassword      = "METABASE_PASSWORD" // #nosec G101
	EnvMetabaseSessionID     = "METABASE_SESSION_ID"
	EnvMetabaseTLSSkipVerify = "METABASE_TLS_SKIP_VERIFY" // #nosec G101
)

var (
	TLSInsecureSkipVerify = false
)

// Config is a basic struct to hold API access information for
// Metabase.
type Config struct {
	BaseURL       string
	SessionID     string
	Username      string
	Password      string
	TLSSkipVerify bool
}

func (cfg *Config) Validate() error {
	cfg.BaseURL = strings.TrimSpace(cfg.BaseURL)
	cfg.SessionID = strings.TrimSpace(cfg.SessionID)
	missing := []string{}
	if len(cfg.BaseURL) == 0 {
		missing = append(missing, "BaseURL")
	}

	if len(cfg.SessionID) == 0 && len(cfg.Username) == 0 && len(cfg.Password) == 0 {
		missing = append(missing, "SessionID or Username/Password")
	} else if len(cfg.SessionID) == 0 &&
		(len(cfg.Username) == 0 && len(cfg.Password) != 0) {
		missing = append(missing, "Username")
	} else if len(cfg.SessionID) == 0 &&
		(len(cfg.Username) != 0 && len(cfg.Password) == 0) {
		missing = append(missing, "Password")
	}

	if len(missing) > 0 {
		return fmt.Errorf("Config Missing: [%s]", strings.Join(missing, ","))
	}
	return nil
}

func NewClient(cfg Config) (*http.Client, *AuthResponse, error) {
	cfg.SessionID = strings.TrimSpace(cfg.SessionID)
	if len(cfg.SessionID) > 0 {
		httpClient := NewClientSessionID(cfg.SessionID, cfg.TLSSkipVerify)
		clientUtil := ClientUtil{
			HTTPClient: httpClient,
			BaseURL:    cfg.BaseURL}
		_, _, err := clientUtil.GetCurrentUser()
		if err == nil {
			return httpClient, nil, nil
		}
	}

	return NewClientPassword(
		cfg.BaseURL,
		cfg.Username,
		cfg.Password,
		cfg.TLSSkipVerify)
}

type authRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthResponse struct {
	ID string `json:"id,omitempty"`
}

// NewClient returns a *http.Client that will add the Metabase Session
// header to each request.
func NewClientPassword(baseURL, username, password string, allowInsecure bool) (*http.Client, *AuthResponse, error) {
	resp, err := AuthRequest(
		urlutil.JoinAbsolute(baseURL, RelPathAPISession),
		username,
		password,
		allowInsecure)
	if err != nil {
		return nil, nil, err
	}

	res := &AuthResponse{}
	_, err = jsonutil.UnmarshalReader(resp.Body, res)
	if err != nil {
		return nil, res, err
	}

	return NewClientSessionID(res.ID, allowInsecure), res, nil
}

// NewClientPasswordWithSessionId returns a *http.Client first attempting to use
// the supplied `sessionId` with a fallback to `username` and `password`.
func NewClientPasswordWithSessionID(baseURL, username, password, sessionID string, allowInsecure bool) (*http.Client, *AuthResponse, error) {
	sessionID = strings.TrimSpace(sessionID)
	if len(sessionID) > 0 {
		httpClient := NewClientSessionID(sessionID, allowInsecure)
		userURL := urlutil.JoinAbsolute(baseURL, RelPathAPIUserCurrent)
		resp, err := httpClient.Get(userURL)
		if err == nil && resp.StatusCode == 200 {
			return httpClient, nil, nil
		}
	}
	return NewClientPassword(baseURL, username, password, allowInsecure)
}

func NewClientSessionID(sessionID string, allowInsecure bool) *http.Client {
	return authutil.NewClientHeaderQuery(
		http.Header{HeaderMetabaseSession: []string{sessionID}},
		url.Values{},
		allowInsecure)
	/*
		client := &http.Client{}

		header := http.Header{}
		header.Add(HeaderMetabaseSession, sessionID)

		if tlsSkipVerify {
			client = authutil.ClientTLSInsecureSkipVerify(client)
		}

		client.Transport = httputilmore.TransportRequestModifier{
			Transport: client.Transport,
			Header:    header}

		return client
	*/
}

func (cfg *Config) NewClient() (*http.Client, *AuthResponse, error) {
	return NewClient(*cfg)
}

type ConfigEnvOpts struct {
	EnvPaths                 []string
	EnvPathsLoad             bool
	EnvMetabaseBaseURL       string
	EnvMetabaseSessionID     string
	EnvMetabaseUsername      string
	EnvMetabasePassword      string // #nosec G101
	EnvMetabaseTLSSkipVerify string
}

func (opts *ConfigEnvOpts) Defaultify() {
	if len(strings.TrimSpace(opts.EnvMetabaseBaseURL)) == 0 {
		opts.EnvMetabaseBaseURL = EnvMetabaseBaseURL
	}
	if len(strings.TrimSpace(opts.EnvMetabaseUsername)) == 0 {
		opts.EnvMetabaseUsername = EnvMetabaseUsername
	}
	if len(strings.TrimSpace(opts.EnvMetabasePassword)) == 0 {
		opts.EnvMetabasePassword = EnvMetabasePassword
	}
	if len(strings.TrimSpace(opts.EnvMetabaseSessionID)) == 0 {
		opts.EnvMetabaseSessionID = EnvMetabaseSessionID
	}
	if len(strings.TrimSpace(opts.EnvMetabaseTLSSkipVerify)) == 0 {
		opts.EnvMetabaseTLSSkipVerify = EnvMetabaseTLSSkipVerify
	}
}

func (opts *ConfigEnvOpts) LoadEnv() error {
	if opts.EnvPathsLoad && len(opts.EnvPaths) > 0 {
		_, err := config.LoadDotEnv(opts.EnvPaths, -1)
		return err
	}
	return nil
}

func (opts *ConfigEnvOpts) Config() Config {
	return Config{
		BaseURL:       os.Getenv(opts.EnvMetabaseBaseURL),
		Username:      os.Getenv(opts.EnvMetabaseUsername),
		Password:      os.Getenv(opts.EnvMetabasePassword),
		SessionID:     os.Getenv(opts.EnvMetabaseSessionID),
		TLSSkipVerify: stringsutil.ToBool(os.Getenv(opts.EnvMetabaseTLSSkipVerify))}
}

func NewClientEnv(opts *ConfigEnvOpts) (*http.Client, *AuthResponse, *Config, error) {
	if opts == nil {
		opts = &ConfigEnvOpts{}
	} else {
		if opts.EnvPathsLoad && len(opts.EnvPaths) > 0 {
			err := opts.LoadEnv()
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}
	opts.Defaultify()
	cfg := opts.Config()
	if len(strings.TrimSpace(cfg.Username)) == 0 {
		return nil, nil, nil, errors.New("metabase client: no 'username' configured")
	}

	client, authres, err := NewClient(cfg)
	return client, authres, &cfg, err
}

/*
func NewClientEnv(initCfg ConfigEnvOpts) (*http.Client, *AuthResponse, error) {
	if initCfg.LoadEnv && len(initCfg.EnvPath) > 0 {
		err := config.LoadDotEnvSkipEmpty(os.Getenv(initCfg.EnvPath), "./.env")
		if err != nil {
			return nil, nil, err
		}
	}

	initCfg.Defaultify()

	return NewClient(Config{
		BaseURL:       os.Getenv(initCfg.EnvMetabaseBaseUrl),
		Username:      os.Getenv(initCfg.EnvMetabaseUsername),
		Password:      os.Getenv(initCfg.EnvMetabasePassword),
		SessionID:     os.Getenv(initCfg.EnvMetabaseSessionId),
		TLSSkipVerify: initCfg.TlsSkipVerify})
}
*/

// AuthRequest creates an authentiation request that returns a id that is used
// in Metabase API requests. It follows the following curl command:
// curl -v -H "Content-Type: application/json" -d '{"username":"myusername","password":"mypassword"}' -XPOST 'http://example.com/api/session'
func AuthRequest(authURL, username, password string, tlsSkipVerify bool) (*http.Response, error) {
	bodyBytes, err := json.Marshal(authRequest{Username: username, Password: password})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, authURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJSONUtf8)

	client := &http.Client{}

	if tlsSkipVerify { // #nosec G402
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsSkipVerify}}
	}

	return client.Do(req)
}

func BuildURL(server, urlpath string) string {
	if urlutil.IsHTTP(urlpath, true, true) {
		return urlpath
	}
	return urlutil.JoinAbsolute(server, urlpath)
}
