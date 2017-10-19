package oauth2util

import (
	"net/http"

	"github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/oauth2util-go/scimutil"
	"golang.org/x/oauth2"
)

type ApplicationCredentials struct {
	ServerURL    string
	ClientID     string
	ClientSecret string
	Endpoint     oauth2.Endpoint
}

type UserCredentials struct {
	Username string
	Password string
}

type OAuth2Util interface {
	SetClient(*http.Client)
	GetSCIMUser() (scimutil.User, error)
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
