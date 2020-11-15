package metabase

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/encoding/jsonutil"
	hum "github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/net/urlutil"
	"github.com/grokify/gotilla/type/stringsutil"
	om "github.com/grokify/oauth2more"
)

const (
	MetabaseSessionHeader = "X-Metabase-Session"
	RelPathApiDatabase    = "api/database"
	RelPathApiSession     = "api/session"
	RelPathApiUserCurrent = "api/user/current"

	// Example environment variables
	EnvMetabaseBaseURL       = "METABASE_BASE_URL"
	EnvMetabaseUsername      = "METABASE_USERNAME"
	EnvMetabasePassword      = "METABASE_PASSWORD"
	EnvMetabaseSessionID     = "METABASE_SESSION_ID"
	EnvMetabaseTLSSkipVerify = "METABASE_TLS_SKIP_VERIFY"
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
	if len(cfg.SessionID) == 0 ||
		(len(cfg.Username) == 0 || len(cfg.Password) == 0) {
		missing = append(missing, "SessionID or Username/Password")
	}
	if len(missing) > 0 {
		return fmt.Errorf("Config Missing: [%s]", strings.Join(missing, ","))
	}
	return nil
}

func NewClient(cfg Config) (*http.Client, *AuthResponse, error) {
	cfg.SessionID = strings.TrimSpace(cfg.SessionID)
	if len(cfg.SessionID) > 0 {
		httpClient := NewClientSessionId(cfg.SessionID, cfg.TLSSkipVerify)
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
	_, err = jsonutil.UnmarshalIoReader(resp.Body, res)
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

func (cfg *Config) NewClient() (*http.Client, *AuthResponse, error) {
	return NewClient(*cfg)
}

type ConfigEnvOpts struct {
	EnvPaths                 []string
	EnvPathsLoad             bool
	EnvMetabaseBaseURL       string
	EnvMetabaseSessionID     string
	EnvMetabaseUsername      string
	EnvMetabasePassword      string
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
		return config.LoadDotEnvSkipEmpty(opts.EnvPaths...)
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
		return nil, nil, nil, errors.New("Metabase Client: No 'username' configured")
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
