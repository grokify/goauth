package ringcentral

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/scim"
	"github.com/grokify/simplego/net/http/httpsimple"
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
}

func (creds *Credentials) NewSimpleClient() (*httpsimple.SimpleClient, error) {
	httpclient, err := creds.NewClient()
	if err != nil {
		return nil, err
	}
	sc := &httpsimple.SimpleClient{
		BaseURL:    creds.Application.ServerURL,
		HTTPClient: httpclient}
	return sc, nil
}

func (creds *Credentials) NewClientCli(oauth2State string) (*http.Client, error) {
	tok, err := creds.NewTokenCli(oauth2State)
	if err != nil {
		return nil, err
	}
	creds.Token = tok
	return oauth2more.NewClientToken(
		oauth2more.TokenBearer, tok.AccessToken, false), nil
}

func (creds *Credentials) NewToken() (*oauth2.Token, error) {
	tok, err := NewTokenPassword(
		creds.Application, creds.PasswordCredentials)
	if err == nil {
		creds.Token = tok
	}
	return tok, err
}

// NewTokenCli retrieves a token using CLI approach for
// OAuth 2.0 authorization code or password grant.
func (creds *Credentials) NewTokenCli(oauth2State string) (*oauth2.Token, error) {
	if strings.ToLower(strings.TrimSpace(creds.Application.GrantType)) == "code" {
		return NewTokenCli(*creds, oauth2State)
	}
	tok, err := NewTokenPassword(
		creds.Application, creds.PasswordCredentials)
	if err == nil {
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
