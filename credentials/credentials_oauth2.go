package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/net/urlutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// CredentialsOAuth2 supports OAuth 2.0 authorization_code, password, and client_credentials grant flows.
type CredentialsOAuth2 struct {
	ServerURL            string              `json:"serverURL,omitempty"`
	ApplicationID        string              `json:"applicationID,omitempty"`
	ClientID             string              `json:"clientID,omitempty"`
	ClientSecret         string              `json:"clientSecret,omitempty"`
	Endpoint             oauth2.Endpoint     `json:"endpoint,omitempty"`
	RedirectURL          string              `json:"redirectURL,omitempty"`
	AppName              string              `json:"applicationName,omitempty"`
	AppVersion           string              `json:"applicationVersion,omitempty"`
	OAuthEndpointID      string              `json:"oauthEndpointID,omitempty"`
	AccessTokenTTL       int64               `json:"accessTokenTTL,omitempty"`
	RefreshTokenTTL      int64               `json:"refreshTokenTTL,omitempty"`
	GrantType            string              `json:"grantType,omitempty"`
	PKCE                 bool                `json:"pkce"`
	Username             string              `json:"username,omitempty"`
	Password             string              `json:"password,omitempty"`
	JWT                  string              `json:"jwt,omitempty"`
	Token                *oauth2.Token       `json:"token,omitempty"`
	Scopes               []string            `json:"scopes,omitempty"`
	AuthCodeOpts         map[string][]string `json:"authCodeOpts,omitempty"`
	AuthCodeExchangeOpts map[string][]string `json:"authCodeExchangeOpts,omitempty"`
	PasswordOpts         map[string][]string `json:"passwordOpts,omitempty"`
}

func ParseCredentialsOAuth2(b []byte) (CredentialsOAuth2, error) {
	creds := CredentialsOAuth2{}
	return creds, json.Unmarshal(b, &creds)
}

// MarshalJSON returns JSON. It is useful for exporting creating configs to be parsed.
func (oc *CredentialsOAuth2) MarshalJSON() ([]byte, error) {
	return json.Marshal(*oc)
}

func (oc *CredentialsOAuth2) Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     oc.ClientID,
		ClientSecret: oc.ClientSecret,
		Endpoint:     oc.Endpoint,
		RedirectURL:  oc.RedirectURL,
		Scopes:       oc.Scopes}
}

func (oc *CredentialsOAuth2) ConfigClientCredentials() clientcredentials.Config {
	return clientcredentials.Config{
		ClientID:     oc.ClientID,
		ClientSecret: oc.ClientSecret,
		TokenURL:     oc.Endpoint.TokenURL,
		Scopes:       oc.Scopes,
		AuthStyle:    oauth2.AuthStyleAutoDetect}
}

type AuthCodeOptions []oauth2.AuthCodeOption

func (opts *AuthCodeOptions) Add(k, v string) {
	*opts = append(*opts, oauth2.SetAuthURLParam(k, v))
}

func (opts *AuthCodeOptions) AddMap(m map[string][]string) {
	for k, vs := range m {
		for _, v := range vs {
			opts.Add(k, v)
		}
	}
}

func (oc *CredentialsOAuth2) AuthCodeURL(state string, opts map[string][]string) string {
	cfg := oc.Config()
	authCodeOptions := AuthCodeOptions{}
	authCodeOptions.AddMap(oc.AuthCodeOpts)
	authCodeOptions.AddMap(opts)
	return cfg.AuthCodeURL(state, authCodeOptions...)
}

func (oc *CredentialsOAuth2) Exchange(ctx context.Context, code string, opts map[string][]string) (*oauth2.Token, error) {
	cfg := oc.Config()
	authCodeOptions := AuthCodeOptions{}
	authCodeOptions.AddMap(oc.AuthCodeExchangeOpts)
	authCodeOptions.AddMap(opts)
	/*
		authCodeOptions := []oauth2.AuthCodeOption{}
		for k, vs := range oc.AuthCodeExchangeOpts {
			for _, v := range vs {
				authCodeOptions = append(authCodeOptions, oauth2.SetAuthURLParam(k, v))
			}
		}

		if len(oc.OAuthEndpointID) > 0 {
			authCodeOptions = append(authCodeOptions,
				oauth2.SetAuthURLParam("endpoint_id", oc.OAuthEndpointID))
		}
		if oc.AccessTokenTTL > 0 {
			authCodeOptions = append(authCodeOptions,
				oauth2.SetAuthURLParam("accessTokenTtl", strconv.Itoa(int(oc.AccessTokenTTL))))
		}
		if oc.RefreshTokenTTL > 0 {
			authCodeOptions = append(authCodeOptions,
				oauth2.SetAuthURLParam("refreshTokenTtl", strconv.Itoa(int(oc.RefreshTokenTTL))))
		}
		for k, vs := range opts {
			for _, v := range vs {
				authCodeOptions = append(authCodeOptions, oauth2.SetAuthURLParam(k, v))
			}
		}
	*/
	return cfg.Exchange(ctx, code, authCodeOptions...)
}

func (oc *CredentialsOAuth2) AppNameAndVersion() string {
	parts := []string{}
	oc.AppName = strings.TrimSpace(oc.AppName)
	oc.AppVersion = strings.TrimSpace(oc.AppVersion)
	if len(oc.AppName) > 0 {
		parts = append(parts, oc.AppName)
	}
	if len(oc.AppVersion) > 0 {
		parts = append(parts, fmt.Sprintf("v%v", oc.AppVersion))
	}
	return strings.Join(parts, "-")
}

func (oc *CredentialsOAuth2) IsGrantType(grantType string) bool {
	return strings.EqualFold(
		strings.TrimSpace(grantType),
		strings.TrimSpace(oc.GrantType))
}

func (oc *CredentialsOAuth2) InflateURL(apiURLPath string) string {
	return urlutil.JoinAbsolute(oc.ServerURL, apiURLPath)
}

// NewClient returns a `*http.Client` for applications using `client_credentials`
// grant. The client can be modified using context, e.g. ignoring bad certs or otherwise.
func (oc *CredentialsOAuth2) NewClient(ctx context.Context) (*http.Client, error) {
	if oc.Token != nil && len(strings.TrimSpace(oc.Token.AccessToken)) > 0 {
		config := oc.Config()
		return config.Client(ctx, oc.Token), nil
	} else if strings.Contains(strings.ToLower(oc.GrantType), "jwt") ||
		oc.IsGrantType(goauth.GrantTypePassword) {
		tok, err := oc.NewToken(ctx)
		if err != nil {
			return nil, err
		}
		config := oc.Config()
		return config.Client(ctx, tok), nil
	} else if oc.IsGrantType(goauth.GrantTypeClientCredentials) {
		config := oc.ConfigClientCredentials()
		return config.Client(ctx), nil
	}
	return nil, fmt.Errorf("grant type is not supported in CredentialsOAuth2.NewClient() [%s]", oc.GrantType)
}

// NewToken retrieves an `*oauth2.Token` when the requisite information is available.
// Note this uses `clientcredentials.Config.Token()` which doesn't always work. In
// This situation, use `goauth.TokenClientCredentials()` as an alternative.
func (oc *CredentialsOAuth2) NewToken(ctx context.Context) (*oauth2.Token, error) {
	if oc.Token != nil && len(strings.TrimSpace(oc.Token.AccessToken)) > 0 {
		return oc.Token, nil
	} else if strings.Contains(strings.ToLower(oc.GrantType), "jwt") {
		return goauth.NewTokenOAuth2JWT(oc.Endpoint.TokenURL, oc.ClientID, oc.ClientSecret, oc.JWT)
	} else if oc.IsGrantType(goauth.GrantTypeClientCredentials) {
		config := oc.ConfigClientCredentials()
		return config.Token(ctx)
	} else if oc.IsGrantType(goauth.GrantTypePassword) {
		cfg := oc.Config()
		return cfg.PasswordCredentialsToken(ctx, oc.Username, oc.Password)
	}
	return nil, fmt.Errorf("grant type is not supported in CredentialsOAuth2.NewToken() [%s]", oc.GrantType)
}

func (oc *CredentialsOAuth2) PasswordRequestBody() url.Values {
	body := url.Values{
		goauth.ParamGrantType: {goauth.GrantTypePassword},
		goauth.ParamUsername:  {oc.Username},
		goauth.ParamPassword:  {oc.Password}}
	if oc.AccessTokenTTL != 0 {
		body.Set("access_token_ttl", strconv.Itoa(int(oc.AccessTokenTTL)))
	}
	if oc.RefreshTokenTTL != 0 {
		body.Set("refresh_token_ttl", strconv.Itoa(int(oc.RefreshTokenTTL)))
	}
	if len(oc.PasswordOpts) > 0 {
		for k, vals := range oc.PasswordOpts {
			for _, v := range vals {
				body.Set(k, v)
			}
		}
	}
	return body
}

func NewCredentialsOAuth2Env(envPrefix string) CredentialsOAuth2 {
	creds := CredentialsOAuth2{
		ClientID:     os.Getenv(envPrefix + "CLIENT_ID"),
		ClientSecret: os.Getenv(envPrefix + "CLIENT_SECRET"),
		ServerURL:    os.Getenv(envPrefix + "SERVER_URL"),
		Username:     os.Getenv(envPrefix + "USERNAME"),
		Password:     os.Getenv(envPrefix + "PASSWORD")}
	if len(strings.TrimSpace(creds.Username)) > 0 {
		creds.GrantType = goauth.GrantTypePassword
	}
	return creds
}
