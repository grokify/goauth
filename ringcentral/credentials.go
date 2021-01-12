package ringcentral

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/scim"
	"golang.org/x/oauth2"
)

const (
	EnvServerURL    = "RINGCENTRAL_SERVER_URL"
	EnvClientID     = "RINGCENTRAL_CLIENT_ID"
	EnvClientSecret = "RINGCENTRAL_CLIENT_SECRET"
	EnvAppName      = "RINGCENTRAL_APP_NAME"
	EnvAppVersion   = "RINGCENTRAL_APP_VERSION"
	EnvRedirectURL  = "RINGCENTRAL_OAUTH_REDIRECT_URL"
	EnvUsername     = "RINGCENTRAL_USERNAME"
	EnvExtension    = "RINGCENTRAL_EXTENSION"
	EnvPassword     = "RINGCENTRAL_PASSWORD"
)

type Credentials struct {
	Application         ApplicationCredentials `json:"application,omitempty"`
	PasswordCredentials PasswordCredentials    `json:"passwordCredentials,omitempty"`
	Token               *oauth2.Token          `json:"token,omitempty"`
}

func NewCredentialsJSONs(appJson, userJson, accessToken []byte) (Credentials, error) {
	creds := Credentials{}
	if len(appJson) > 1 {
		app := ApplicationCredentials{}
		err := json.Unmarshal(appJson, &app)
		if err != nil {
			return creds, err
		}
		creds.Application = app
	}
	if len(userJson) > 0 {
		user := PasswordCredentials{}
		err := json.Unmarshal(userJson, &user)
		if err != nil {
			return creds, err
		}
		creds.PasswordCredentials = user
	}
	if len(accessToken) > 0 {
		creds.Token = &oauth2.Token{
			AccessToken: string(accessToken)}
	}
	return creds, nil
}

func NewCredentialsJSON(jsonData []byte) (Credentials, error) {
	creds := Credentials{}
	return creds, json.Unmarshal(jsonData, &creds)
}

func NewCredentialsEnv() Credentials {
	return Credentials{
		Application:         NewApplicationCredentialsEnv(),
		PasswordCredentials: NewPasswordCredentialsEnv()}
}

func (creds *Credentials) NewClient() (*http.Client, error) {
	tok, err := creds.NewToken()
	if err != nil {
		return nil, err
	}
	creds.Token = tok
	return oauth2more.NewClientToken(
		oauth2more.TokenBearer, tok.AccessToken, false), nil
	/*return NewClientPassword(creds.Application, creds.PasswordCredentials)*/
}

func (creds *Credentials) NewToken() (*oauth2.Token, error) {
	tok, err := NewTokenPassword(
		creds.Application, creds.PasswordCredentials)
	if err == nil {
		creds.Token = tok
	}
	return tok, err
}

func (creds *Credentials) NewClientUtil() (ClientUtil, error) {
	httpClient, err := creds.NewClient()
	if err != nil {
		return ClientUtil{}, err
	}
	return ClientUtil{
		Client:    httpClient,
		ServerURL: creds.Application.ServerURL}, nil
}

func (creds *Credentials) Me() (RingCentralExtensionInfo, error) {
	cu, err := creds.NewClientUtil()
	if err != nil {
		return RingCentralExtensionInfo{}, err
	}
	return cu.GetUserinfo()
}

func (creds *Credentials) MeScim() (scim.User, error) {
	cu, err := creds.NewClientUtil()
	if err != nil {
		return scim.User{}, err
	}
	return cu.GetSCIMUser()
}

type ApplicationCredentials struct {
	ServerURL       string `json:"serverURL,omitempty"`
	ApplicationID   string `json:"applicationID,omitempty"`
	ClientID        string `json:"clientID,omitempty"`
	ClientSecret    string `json:"clientSecret,omitempty"`
	RedirectURL     string `json:"redirectURL,omitempty"`
	AppName         string `json:"applicationName,omitempty"`
	AppVersion      string `json:"applicationVersion,omitempty"`
	OAuthEndpointID string `json:"oauthEndpointID,omitempty"`
	AccessTokenTTL  int64  `json:"accessTokenTTL,omitempty"`
	RefreshTokenTTL int64  `json:"refreshTokenTTL,omitempty"`
	GrantType       string `json:"grantType,omitempty"`
}

func NewApplicationCredentialsEnv() ApplicationCredentials {
	return ApplicationCredentials{
		ServerURL:    os.Getenv(EnvServerURL),
		ClientID:     os.Getenv(EnvClientID),
		ClientSecret: os.Getenv(EnvClientSecret),
		AppName:      os.Getenv(EnvAppName),
		AppVersion:   os.Getenv(EnvAppVersion)}
}

func (app *ApplicationCredentials) Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Endpoint:     NewEndpoint(app.ServerURL),
		RedirectURL:  app.RedirectURL}
}

func (app *ApplicationCredentials) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	cfg := app.Config()
	return cfg.AuthCodeURL(state, opts...)
}

func (app *ApplicationCredentials) Exchange(code string) (*RcToken, error) {
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

type PasswordCredentials struct {
	GrantType            string `url:"grant_type"`
	AccessTokenTTL       int64  `url:"access_token_ttl"`
	RefreshTokenTTL      int64  `url:"refresh_token_ttl"`
	Username             string `json:"username" url:"username"`
	Extension            string `json:"extension" url:"extension"`
	Password             string `json:"password" url:"password"`
	EndpointId           string `url:"endpoint_id"`
	EngageVoiceAccountId int64  `json:"engageVoiceAccountId"`
}

func NewPasswordCredentialsEnv() PasswordCredentials {
	return PasswordCredentials{
		Username:  os.Getenv(EnvUsername),
		Extension: os.Getenv(EnvExtension),
		Password:  os.Getenv(EnvPassword)}
}

func (pw *PasswordCredentials) URLValues() url.Values {
	v := url.Values{
		"grant_type": {"password"},
		"username":   {pw.Username},
		"password":   {pw.Password}}
	if pw.AccessTokenTTL != 0 {
		v.Set("access_token_ttl", strconv.Itoa(int(pw.AccessTokenTTL)))
	}
	if pw.RefreshTokenTTL != 0 {
		v.Set("refresh_token_ttl", strconv.Itoa(int(pw.RefreshTokenTTL)))
	}
	if len(pw.Extension) > 0 {
		v.Set("extension", pw.Extension)
	}
	if len(pw.EndpointId) > 0 {
		v.Set("endpoint_id", pw.EndpointId)
	}
	return v
}

func (uc *PasswordCredentials) UsernameSimple() string {
	if len(strings.TrimSpace(uc.Extension)) > 0 {
		return strings.Join([]string{uc.Username, uc.Extension}, "*")
	}
	return uc.Username
}
