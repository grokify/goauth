package oauth2util

import (
	"fmt"
	"net/http"

	b64 "github.com/grokify/gotilla/encoding/base64"
	"github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/oauth2util/scimutil"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// ApplicationCredentials represents information for an app.
type ApplicationCredentials struct {
	ServerURL    string
	ClientID     string
	ClientSecret string
	Endpoint     oauth2.Endpoint
}

// UserCredentials represents a user's credentials.
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

func BasicAuthToken(username, password string) (*oauth2.Token, error) {
	basicToken, err := b64.RFC7617UserPass(username, password)
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: basicToken,
		TokenType:   "Basic",
		Expiry:      timeutil.TimeRFC3339Zero()}, nil
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
