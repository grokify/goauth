package credentials

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/grokify/simplego/net/urlutil"

	"golang.org/x/oauth2"
)

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
}

func (app *ApplicationCredentials) Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Endpoint:     app.OAuth2Endpoint,
		RedirectURL:  app.RedirectURL}
}

/*
func (app *ApplicationCredentials) NewEndpoint() oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  app.OAuth2AuthzURL,
		TokenURL: app.OAuth2TokenURL}
}*/

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
		//params.Set("accessTokenTtl", strconv.Itoa(int(app.AccessTokenTTL)))
	}
	if app.RefreshTokenTTL > 0 {
		authCodeOptions = append(authCodeOptions,
			oauth2.SetAuthURLParam("refreshTokenTtl", strconv.Itoa(int(app.RefreshTokenTTL))))
		//params.Set("refreshTokenTtl", strconv.Itoa(int(app.RefreshTokenTTL)))
	}

	return cfg.Exchange(context.Background(), code)
	/*
		params := url.Values{}
		params.Set("grant_type", "authorization_code")
		params.Set("code", code)
		params.Set("redirect_uri", app.RedirectURL)
		if len(app.OAuthEndpointID) > 0 {
			params.Set("endpoint_id", app.OAuthEndpointID)
		}
		if app.AccessTokenTTL > 0 {
			params.Set("accessTokenTtl", strconv.Itoa(int(app.AccessTokenTTL)))
		}
		if app.RefreshTokenTTL > 0 {
			params.Set("refreshTokenTtl", strconv.Itoa(int(app.RefreshTokenTTL)))
		}
		return RetrieveRcToken(app.Config(), params)
	*/
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
	if strings.TrimSpace(grantType) == strings.TrimSpace(app.GrantType) {
		return true
	}
	return false
}

func (app *ApplicationCredentials) InflateURL(apiUrlPath string) string {
	return urlutil.JoinAbsolute(app.ServerURL, apiUrlPath)
}
