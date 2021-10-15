package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/grokify/oauth2more"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// ApplicationCredentials supports OAuth 2.0 authorization_code, password,
// and client_credentials grant flows.
type ApplicationCredentials struct {
	ServerURL       string          `json:"serverURL,omitempty"`
	ApplicationID   string          `json:"applicationID,omitempty"`
	ClientID        string          `json:"clientID,omitempty"`
	ClientSecret    string          `json:"clientSecret,omitempty"`
	OAuth2Endpoint  oauth2.Endpoint `json:"oauth2Endpoint,omitempty"`
	RedirectURL     string          `json:"redirectURL,omitempty"`
	AppName         string          `json:"applicationName,omitempty"`
	AppVersion      string          `json:"applicationVersion,omitempty"`
	OAuthEndpointID string          `json:"oauthEndpointID,omitempty"`
	AccessTokenTTL  int64           `json:"accessTokenTTL,omitempty"`
	RefreshTokenTTL int64           `json:"refreshTokenTTL,omitempty"`
	GrantType       string          `json:"grantType,omitempty"`
	Scopes          []string        `json:"scopes,omitempty"`
}

func (app *ApplicationCredentials) Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Endpoint:     app.OAuth2Endpoint,
		RedirectURL:  app.RedirectURL,
		Scopes:       app.Scopes}
}

func (app *ApplicationCredentials) ConfigClientCredentials() clientcredentials.Config {
	return clientcredentials.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		TokenURL:     app.OAuth2Endpoint.TokenURL,
		Scopes:       app.Scopes,
		AuthStyle:    oauth2.AuthStyleAutoDetect}

}

func (app *ApplicationCredentials) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	cfg := app.Config()
	return cfg.AuthCodeURL(state, opts...)
}

func (app *ApplicationCredentials) Exchange(code string) (*oauth2.Token, error) {
	cfg := app.Config()
	authCodeOptions := []oauth2.AuthCodeOption{}

	if len(app.OAuthEndpointID) > 0 {
		authCodeOptions = append(authCodeOptions,
			oauth2.SetAuthURLParam("endpoint_id", app.OAuthEndpointID))
	}
	if app.AccessTokenTTL > 0 {
		authCodeOptions = append(authCodeOptions,
			oauth2.SetAuthURLParam("accessTokenTtl", strconv.Itoa(int(app.AccessTokenTTL))))
	}
	if app.RefreshTokenTTL > 0 {
		authCodeOptions = append(authCodeOptions,
			oauth2.SetAuthURLParam("refreshTokenTtl", strconv.Itoa(int(app.RefreshTokenTTL))))
	}

	return cfg.Exchange(context.Background(), code, authCodeOptions...)
}

func (ac *ApplicationCredentials) AppNameAndVersion() string {
	parts := []string{}
	ac.AppName = strings.TrimSpace(ac.AppName)
	ac.AppVersion = strings.TrimSpace(ac.AppVersion)
	if len(ac.AppName) > 0 {
		parts = append(parts, ac.AppName)
	}
	if len(ac.AppVersion) > 0 {
		parts = append(parts, fmt.Sprintf("v%v", ac.AppVersion))
	}
	return strings.Join(parts, "-")
}

func (app *ApplicationCredentials) IsGrantType(grantType string) bool {
	return strings.EqualFold(
		strings.TrimSpace(grantType),
		strings.TrimSpace(app.GrantType))
}

func (app *ApplicationCredentials) InflateURL(apiUrlPath string) string {
	return urlutil.JoinAbsolute(app.ServerURL, apiUrlPath)
}

// NewClient returns a `*http.Client` for applications using `client_credentials`
// grant. The client can be modified using context, e.g. ignoring bad certs or otherwise.
func (app *ApplicationCredentials) NewClient(ctx context.Context) (*http.Client, error) {
	if app.GrantType != oauth2more.GrantTypeClientCredentials {
		return nil, errors.New("grant type is not client_credentials")
	}
	config := app.ConfigClientCredentials()
	return config.Client(ctx), nil
}

// NewToken retrieves an `*oauth2.Token` when the requisite information is available.
// Note this uses `clientcredentials.Config.Token()` which doesn't always work. In
// This situation, use `oauth2more.TokenClientCredentials()` as an alternative.
func (app *ApplicationCredentials) NewToken(ctx context.Context) (*oauth2.Token, error) {
	if app.GrantType != oauth2more.GrantTypeClientCredentials {
		return nil, errors.New("grant type is not client_credentials")
	}
	config := app.ConfigClientCredentials()
	return config.Token(ctx)
}
