package credentials

import (
	"context"
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

// OAuth2Credentials supports OAuth 2.0 authorization_code, password,
// and client_credentials grant flows.
type OAuth2Credentials struct {
	ServerURL       string              `json:"serverURL,omitempty"`
	ApplicationID   string              `json:"applicationID,omitempty"`
	ClientID        string              `json:"clientID,omitempty"`
	ClientSecret    string              `json:"clientSecret,omitempty"`
	OAuth2Endpoint  oauth2.Endpoint     `json:"oauth2Endpoint,omitempty"`
	RedirectURL     string              `json:"redirectURL,omitempty"`
	AppName         string              `json:"applicationName,omitempty"`
	AppVersion      string              `json:"applicationVersion,omitempty"`
	OAuthEndpointID string              `json:"oauthEndpointID,omitempty"`
	AccessTokenTTL  int64               `json:"accessTokenTTL,omitempty"`
	RefreshTokenTTL int64               `json:"refreshTokenTTL,omitempty"`
	GrantType       string              `json:"grantType,omitempty"`
	PKCE            bool                `json:"pkce"`
	Username        string              `json:"username,omitempty"`
	Password        string              `json:"password,omitempty"`
	JWT             string              `json:"jwt,omitempty"`
	OtherParams     map[string][]string `json:"otherParams,omitempty"`
	Scopes          []string            `json:"scopes,omitempty"`
}

func (oc *OAuth2Credentials) Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     oc.ClientID,
		ClientSecret: oc.ClientSecret,
		Endpoint:     oc.OAuth2Endpoint,
		RedirectURL:  oc.RedirectURL,
		Scopes:       oc.Scopes}
}

func (oc *OAuth2Credentials) ConfigClientCredentials() clientcredentials.Config {
	return clientcredentials.Config{
		ClientID:     oc.ClientID,
		ClientSecret: oc.ClientSecret,
		TokenURL:     oc.OAuth2Endpoint.TokenURL,
		Scopes:       oc.Scopes,
		AuthStyle:    oauth2.AuthStyleAutoDetect}
}

// func (oc *OAuth2Credentials) AuthCodeURL(state string, opts url.Values) string {
func (oc *OAuth2Credentials) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	cfg := oc.Config()
	return cfg.AuthCodeURL(state, opts...)
}

func (oc *OAuth2Credentials) Exchange(code string) (*oauth2.Token, error) {
	cfg := oc.Config()
	authCodeOptions := []oauth2.AuthCodeOption{}

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

	return cfg.Exchange(context.Background(), code, authCodeOptions...)
}

func (oc *OAuth2Credentials) AppNameAndVersion() string {
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

func (oc *OAuth2Credentials) IsGrantType(grantType string) bool {
	return strings.EqualFold(
		strings.TrimSpace(grantType),
		strings.TrimSpace(oc.GrantType))
}

func (oc *OAuth2Credentials) InflateURL(apiUrlPath string) string {
	return urlutil.JoinAbsolute(oc.ServerURL, apiUrlPath)
}

// NewClient returns a `*http.Client` for applications using `client_credentials`
// grant. The client can be modified using context, e.g. ignoring bad certs or otherwise.
func (oc *OAuth2Credentials) NewClient(ctx context.Context) (*http.Client, error) {
	if strings.Contains(strings.ToLower(oc.GrantType), "jwt") {
		tok, err := oc.NewToken(ctx)
		if err != nil {
			return nil, err
		}
		return goauth.NewClientToken(
			goauth.TokenBearer, tok.AccessToken, false), nil
	} else if oc.GrantType == goauth.GrantTypeClientCredentials {
		config := oc.ConfigClientCredentials()
		return config.Client(ctx), nil
	}
	return nil, fmt.Errorf("grant type is not client_credentials or jwt-bearer [%s]", oc.GrantType)
}

// NewToken retrieves an `*oauth2.Token` when the requisite information is available.
// Note this uses `clientcredentials.Config.Token()` which doesn't always work. In
// This situation, use `goauth.TokenClientCredentials()` as an alternative.
func (oc *OAuth2Credentials) NewToken(ctx context.Context) (*oauth2.Token, error) {
	if strings.Contains(strings.ToLower(oc.GrantType), "jwt") {
		return goauth.NewTokenOAuth2Jwt(oc.OAuth2Endpoint.TokenURL,
			oc.ClientID, oc.ClientSecret, oc.JWT)
	} else if oc.GrantType == goauth.GrantTypeClientCredentials {
		config := oc.ConfigClientCredentials()
		return config.Token(ctx)
	}
	return nil, fmt.Errorf("grant type is not client_credentials or jwt-bearer [%s]", oc.GrantType)
}

func (oc *OAuth2Credentials) PasswordRequestBody() url.Values {
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
	if len(oc.OtherParams) > 0 {
		for k, vals := range oc.OtherParams {
			for _, v := range vals {
				body.Set(k, v)
			}
		}
	}
	return body
}

func NewOAuth2CredentialsEnv(envPrefix string) OAuth2Credentials {
	creds := OAuth2Credentials{
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
