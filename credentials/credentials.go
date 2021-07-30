package credentials

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/grokify/oauth2more"
	"github.com/grokify/simplego/net/http/httpsimple"
	"golang.org/x/oauth2"
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
	cfg := creds.Application.Config()
	return cfg.PasswordCredentialsToken(
		context.Background(),
		creds.PasswordCredentials.Username,
		creds.PasswordCredentials.Password)
	/*
		tok, err := NewTokenPassword(
			creds.Application, creds.PasswordCredentials)
		if err == nil {
			creds.Token = tok
		}
		return tok, err*/
}

// NewTokenCli retrieves a token using CLI approach for
// OAuth 2.0 authorization code or password grant.
func (creds *Credentials) NewTokenCli(oauth2State string) (*oauth2.Token, error) {
	if strings.ToLower(strings.TrimSpace(creds.Application.GrantType)) == "code" {
		return NewTokenCli(*creds, oauth2State)
	}
	return creds.NewToken()
}
