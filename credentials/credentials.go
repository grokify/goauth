package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/endpoints"
	"github.com/grokify/gohttp/httpsimple"
	"github.com/grokify/mogo/errors/errorsutil"
	"golang.org/x/oauth2"
)

const (
	TypeOAuth2 = "oauth2"
	TypeJWT    = "jwt"
)

type Credentials struct {
	Service   string            `json:"service,omitempty"`
	Type      string            `json:"type,omitempty"`
	Subdomain string            `json:"subdomain,omitempty"`
	OAuth2    CredentialsOAuth2 `json:"oauth2,omitempty"`
	JWT       CredentialsJWT    `json:"jwt,omitempty"`
	Token     *oauth2.Token     `json:"token,omitempty"`
}

func NewCredentialsJSON(credsData, accessToken []byte) (Credentials, error) {
	var creds Credentials
	err := json.Unmarshal(credsData, &creds)
	if err != nil {
		return creds, err
	}
	err = creds.Inflate()
	if err != nil {
		return creds, err
	}
	if len(accessToken) > 0 {
		creds.Token = &oauth2.Token{
			AccessToken: string(accessToken)}
	}
	return creds, nil
}

func (creds *Credentials) Inflate() error {
	if len(strings.TrimSpace(creds.Service)) > 0 {
		ep, svcURL, err := endpoints.NewEndpoint(creds.Service, creds.Subdomain)
		if err != nil {
			return err
		}
		if creds.OAuth2.Endpoint == (oauth2.Endpoint{}) {
			creds.OAuth2.Endpoint = ep
		}
		if len(strings.TrimSpace(creds.OAuth2.ServerURL)) == 0 {
			creds.OAuth2.ServerURL = svcURL
		}
	}
	return nil
}

func (creds *Credentials) NewClient(ctx context.Context) (*http.Client, error) {
	if creds.Type == TypeJWT {
		return nil, errors.New("NewClient() does not support jwt")
	}
	if creds.Token != nil {
		return goauth.NewClientToken(goauth.TokenBearer, creds.Token.AccessToken, false), nil
	}
	if creds.OAuth2.GrantType == goauth.GrantTypeClientCredentials ||
		strings.Contains(creds.OAuth2.GrantType, "jwt") {
		return creds.OAuth2.NewClient(ctx)
	}
	tok, err := creds.NewToken()
	if err != nil {
		return nil, errorsutil.Wrap(err, "Credentials.NewToken()")
	}
	creds.Token = tok
	return goauth.NewClientToken(goauth.TokenBearer, tok.AccessToken, false), nil
}

func (creds *Credentials) NewSimpleClient(httpClient *http.Client) (*httpsimple.SimpleClient, error) {
	return &httpsimple.SimpleClient{
		BaseURL:    creds.OAuth2.ServerURL,
		HTTPClient: httpClient}, nil
}

func (creds *Credentials) NewClientCli(oauth2State string) (*http.Client, error) {
	tok, err := creds.NewTokenCli(oauth2State)
	if err != nil {
		return nil, err
	}
	creds.Token = tok
	return goauth.NewClientToken(
		goauth.TokenBearer, tok.AccessToken, false), nil
}

func (creds *Credentials) NewToken() (*oauth2.Token, error) {
	cfg := creds.OAuth2.Config()
	return cfg.PasswordCredentialsToken(
		context.Background(),
		creds.OAuth2.Username,
		creds.OAuth2.Password)
}

// NewTokenCli retrieves a token using CLI approach for
// OAuth 2.0 authorization code or password grant.
func (creds *Credentials) NewTokenCli(oauth2State string) (*oauth2.Token, error) {
	if strings.EqualFold(strings.TrimSpace(creds.OAuth2.GrantType), goauth.GrantTypeAuthorizationCode) {
		return NewTokenCli(*creds, oauth2State)
	}
	return creds.NewToken()
}
