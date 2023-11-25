package goauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/endpoints"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"golang.org/x/oauth2"
)

const (
	TypeBasic       = "basic"
	TypeHeaderQuery = "headerquery"
	TypeOAuth2      = "oauth2"
	TypeJWT         = "jwt"
	TypeGCPSA       = "gcpsa" // Google Cloud Platform Service Account
)

type Credentials struct {
	Service     string                  `json:"service,omitempty"`
	Type        string                  `json:"type,omitempty"`
	Subdomain   string                  `json:"subdomain,omitempty"`
	Basic       *CredentialsBasicAuth   `json:"basic,omitempty"`
	HeaderQuery *CredentialsHeaderQuery `json:"headerquery,omitempty"`
	GCPSA       *CredentialsGCP         `json:"gcpsa,omitempty"`
	JWT         *CredentialsJWT         `json:"jwt,omitempty"`
	OAuth2      *CredentialsOAuth2      `json:"oauth2,omitempty"`
	Token       *oauth2.Token           `json:"token,omitempty"`
	Additional  url.Values              `json:"additional,omitempty"`
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
		if creds.Type == TypeOAuth2 {
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
	}
	return nil
}

var (
	ErrBasicAuthNotPopulated   = errors.New("basic auth is not populated")
	ErrHeaderQueryNotPopulated = errors.New("header query is not populated")
	ErrJWTNotPopulated         = errors.New("jwt is not populated")
	ErrJWTNotSupported         = errors.New("jwt is not supported for function")
	ErrOAuth2NotPopulated      = errors.New("oauth2 is not populated")
	ErrTypeNotSupported        = errors.New("credentials type not supported")
	ErrGCPSANotPopulated       = errors.New("gcp service account credentials are not populated")
)

func (creds *Credentials) NewClient(ctx context.Context) (*http.Client, error) {
	switch creds.Type {
	case TypeBasic:
		if creds.Basic == nil {
			return nil, ErrBasicAuthNotPopulated
		}
		return creds.Basic.NewClient()
	case TypeGCPSA:
		if creds.GCPSA == nil {
			return nil, ErrGCPSANotPopulated
		}
		return creds.GCPSA.NewClient(ctx)
	case TypeHeaderQuery:
		if creds.HeaderQuery == nil {
			return nil, ErrHeaderQueryNotPopulated
		}
		return creds.HeaderQuery.NewClient(), nil
	case TypeJWT:
		return nil, ErrJWTNotSupported
	}
	if creds.Token != nil {
		return authutil.NewClientToken(authutil.TokenBearer, creds.Token.AccessToken, false), nil
	}
	if creds.OAuth2.GrantType == authutil.GrantTypeClientCredentials ||
		strings.Contains(creds.OAuth2.GrantType, TypeJWT) {
		clt, _, err := creds.OAuth2.NewClient(ctx)
		return clt, err
	}
	tok, err := creds.NewToken()
	if err != nil {
		return nil, errorsutil.Wrap(err, "Credentials.NewToken()")
	}
	creds.Token = tok
	return authutil.NewClientToken(authutil.TokenBearer, tok.AccessToken, false), nil
}

func (creds *Credentials) NewSimpleClient(ctx context.Context) (*httpsimple.Client, error) {
	httpClient, err := creds.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return creds.NewSimpleClientHTTP(httpClient)
}

func (creds *Credentials) NewSimpleClientHTTP(httpClient *http.Client) (*httpsimple.Client, error) {
	switch creds.Type {
	case TypeJWT:
		return nil, ErrJWTNotSupported
	case TypeBasic:
		return &httpsimple.Client{
			BaseURL:    creds.Basic.ServerURL,
			HTTPClient: httpClient}, nil
	case TypeHeaderQuery:
		return &httpsimple.Client{
			BaseURL:    creds.HeaderQuery.ServerURL,
			HTTPClient: httpClient}, nil
	case TypeOAuth2:
		return &httpsimple.Client{
			BaseURL:    creds.OAuth2.ServerURL,
			HTTPClient: httpClient}, nil
	default:
		return nil, ErrTypeNotSupported
	}
}

func (creds *Credentials) NewClientCLI(oauth2State string) (*http.Client, error) {
	tok, err := creds.NewTokenCLI(oauth2State)
	if err != nil {
		return nil, err
	}
	creds.Token = tok
	return authutil.NewClientToken(
		authutil.TokenBearer, tok.AccessToken, false), nil
}

func (creds *Credentials) NewToken() (*oauth2.Token, error) {
	cfg := creds.OAuth2.Config()
	return cfg.PasswordCredentialsToken(
		context.Background(),
		creds.OAuth2.Username,
		creds.OAuth2.Password)
}

// NewTokenCLI retrieves a token using CLI approach for
// OAuth 2.0 authorization code or password grant.
func (creds *Credentials) NewTokenCLI(oauth2State string) (*oauth2.Token, error) {
	if strings.EqualFold(strings.TrimSpace(creds.OAuth2.GrantType), authutil.GrantTypeAuthorizationCode) {
		return NewTokenCLI(*creds, oauth2State)
	}
	return creds.NewToken()
}
