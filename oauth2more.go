package oauth2more

import (
	"context"
	"crypto/tls"
	"encoding/json"
	errr "errors"
	"fmt"
	"net/http"

	"github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/oauth2more/scim"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type ServiceType int

const (
	Google ServiceType = iota
	Facebook
	RingCentral
	Aha
)

// ApplicationCredentials represents information for an app.
type ApplicationCredentials struct {
	ServerURL    string
	ClientID     string
	ClientSecret string
	Endpoint     oauth2.Endpoint
}

type AppCredentials struct {
	Service      string   `json:"service,omitempty"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIs []string `json:"redirect_uris"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
	Scopes       []string `json:"scopes"`
}

func (ac *AppCredentials) Defaultify() {
	switch ac.Service {
	case "facebook":
		if len(ac.AuthURI) == 0 || len(ac.TokenURI) == 0 {
			endpoint := facebook.Endpoint
			if len(ac.AuthURI) == 0 {
				ac.AuthURI = endpoint.AuthURL
			}
			if len(ac.TokenURI) == 0 {
				ac.TokenURI = endpoint.TokenURL
			}
		}
	}
}

type AppCredentialsWrapper struct {
	Web       *AppCredentials `json:"web"`
	Installed *AppCredentials `json:"installed"`
}

func (w *AppCredentialsWrapper) Config() (*oauth2.Config, error) {
	var c *AppCredentials
	if w.Web != nil {
		c = w.Web
	} else if w.Installed != nil {
		c = w.Installed
	} else {
		return nil, errr.New("No OAuth2 config info")
	}
	c.Defaultify()
	return c.Config(), nil
}

func NewAppCredentialsWrapperFromBytes(data []byte) (AppCredentialsWrapper, error) {
	var acw AppCredentialsWrapper
	err := json.Unmarshal(data, &acw)
	if err != nil {
		panic(err)
	}
	return acw, err
}

func (c *AppCredentials) Config() *oauth2.Config {
	cfg := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Scopes:       c.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AuthURI,
			TokenURL: c.TokenURI,
		},
	}
	if len(c.RedirectURIs) > 0 {
		cfg.RedirectURL = c.RedirectURIs[0]
	}
	return cfg
}

// UserCredentials represents a user's credentials.
type UserCredentials struct {
	Username string
	Password string
}

type OAuth2Util interface {
	SetClient(*http.Client)
	GetSCIMUser() (scim.User, error)
}

func NewClientPasswordConf(conf oauth2.Config, username, password string) (*http.Client, error) {
	token, err := conf.PasswordCredentialsToken(oauth2.NoContext, username, password)
	if err != nil {
		return &http.Client{}, err
	}

	return conf.Client(oauth2.NoContext, token), nil
}

func NewClientAuthCode(conf oauth2.Config, authCode string) (*http.Client, error) {
	token, err := conf.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		return &http.Client{}, err
	}
	return conf.Client(oauth2.NoContext, token), nil
}

func NewClientAccessToken(accessToken string) *http.Client {
	token := &oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		Expiry:      timeutil.TimeRFC3339Zero()}

	oAuthConfig := &oauth2.Config{}

	return oAuthConfig.Client(oauth2.NoContext, token)
}

func NewTokenFromWeb(cfg *oauth2.Config) (*oauth2.Token, error) {
	authURL := cfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to this link in your browser then type the auth code: \n%v\n", authURL)

	code := ""
	if _, err := fmt.Scan(&code); err != nil {
		return &oauth2.Token{}, errors.Wrap(err, "Unable to read auth code")
	}

	tok, err := cfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		return tok, errors.Wrap(err, "Unable to retrieve token from web")
	}
	return tok, nil
}

func NewClientTLSToken(ctx context.Context, tlsConfig *tls.Config, token *oauth2.Token) *http.Client {
	tlsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig}}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, tlsClient)

	cfg := &oauth2.Config{}

	return cfg.Client(ctx, token)
}
