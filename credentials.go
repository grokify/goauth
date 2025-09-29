package goauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	TypeBasic        = "basic"
	TypeHeaderQuery  = "headerquery"
	TypeOAuth2       = "oauth2"
	TypeJWT          = "jwt"
	TypeGCPSA        = "gcpsa" // Google Cloud Platform Service Account
	TypeGoogleOAuth2 = "googleoauth2"
)

var ErrsInclLocation = false

type Credentials struct {
	Service      string                   `json:"service,omitempty"`
	Type         string                   `json:"type,omitempty"`
	Subdomain    string                   `json:"subdomain,omitempty"`
	Basic        *CredentialsBasicAuth    `json:"basic,omitempty"`
	HeaderQuery  *CredentialsHeaderQuery  `json:"headerquery,omitempty"`
	GCPSA        *CredentialsGCP          `json:"gcpsa,omitempty"`
	GoogleOAuth2 *CredentialsGoogleOAuth2 `json:"googleoauth2,omitempty"`
	JWT          *CredentialsJWT          `json:"jwt,omitempty"`
	OAuth2       *CredentialsOAuth2       `json:"oauth2,omitempty"`
	Token        *oauth2.Token            `json:"token,omitempty"`
	Additional   url.Values               `json:"additional,omitempty"`
}

func NewCredentialsFromCLI(inclAccountsOnError bool) (Credentials, error) {
	if opts, err := ParseOptions(); err != nil {
		return Credentials{}, err
	} else {
		return opts.Credentials()
	}
}

func NewCredentialsFromJSON(credsData, accessToken []byte) (Credentials, error) {
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

func NewCredentialsFromSetFile(credentialsSetFilename, accountKey string, inclAccountsOnError bool) (Credentials, error) {
	set, err := ReadFileCredentialsSet(credentialsSetFilename, true)
	if err != nil {
		return Credentials{}, err
	}
	creds, err := set.Get(accountKey)
	if err != nil {
		if inclAccountsOnError {
			return creds, errorsutil.Wrap(err,
				fmt.Sprintf("validAccounts [%s]", strings.Join(set.Accounts(), ",")))
		}
		return creds, err
	}
	return creds, nil
}

func (creds *Credentials) Inflate() error {
	if len(strings.TrimSpace(creds.Service)) > 0 {
		if creds.Type == TypeOAuth2 {
			if creds.OAuth2 == nil {
				return errors.New("type `oauth2` is not set")
			}
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

func NewClient(ctx context.Context, goauthfile, goauthkey string) (*http.Client, error) {
	if creds, err := NewCredentialsFromSetFile(goauthfile, goauthkey, false); err != nil {
		return nil, err
	} else {
		return creds.NewClient(ctx)
	}
}

func (creds *Credentials) NewClient(ctx context.Context) (*http.Client, error) {
	switch creds.Type {
	case TypeBasic:
		if creds.Basic == nil {
			return nil, ErrBasicAuthNotPopulated
		} else {
			return creds.Basic.NewClient()
		}
	case TypeGCPSA:
		if creds.GCPSA == nil {
			return nil, ErrGCPSANotPopulated
		} else {
			return creds.GCPSA.NewClient(ctx)
		}
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

	if creds.OAuth2 != nil && (creds.OAuth2.GrantType == authutil.GrantTypeClientCredentials ||
		strings.Contains(creds.OAuth2.GrantType, TypeJWT)) {
		clt, _, err := creds.OAuth2.NewClient(ctx)
		return clt, err
	}
	if tok, err := creds.NewToken(ctx); err != nil {
		return nil, errorsutil.Wrap(err, "Credentials.NewToken()")
	} else {
		creds.Token = tok
		return authutil.NewClientToken(authutil.TokenBearer, tok.AccessToken, false), nil
	}
}

func (creds *Credentials) NewSimpleClient(ctx context.Context) (*httpsimple.Client, error) {
	if httpClient, err := creds.NewClient(ctx); err != nil {
		return nil, err
	} else {
		return creds.NewSimpleClientHTTP(httpClient)
	}
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

func (creds *Credentials) NewClientCLI(ctx context.Context, oauth2State string) (*http.Client, error) {
	if tok, err := creds.NewTokenCLI(ctx, oauth2State); err != nil {
		return nil, err
	} else {
		creds.Token = tok
		return authutil.NewClientToken(authutil.TokenBearer, tok.AccessToken, false), nil
	}
}

func (creds *Credentials) NewOrExistingValidToken(ctx context.Context) (*oauth2.Token, error) {
	if tok, err := creds.ExistingValidToken(); err == nil && tok != nil {
		return tok, nil
	} else {
		return creds.NewToken(ctx)
	}
}

func (creds *Credentials) ExistingValidToken() (*oauth2.Token, error) {
	if creds.Type == TypeOAuth2 && creds.OAuth2 != nil && creds.OAuth2.Token != nil && creds.OAuth2.Token.Valid() {
		return creds.OAuth2.Token, nil
	} else if creds.Type == TypeGoogleOAuth2 && creds.GoogleOAuth2 != nil && creds.GoogleOAuth2.Token.Valid() {
		return creds.GoogleOAuth2.Token, nil
	} else {
		return nil, nil
	}
}

func (creds *Credentials) NewToken(ctx context.Context) (*oauth2.Token, error) {
	switch creds.Type {
	case TypeOAuth2:
		if creds.OAuth2 == nil {
			return nil, fmt.Errorf("credentials.%s is nil for type `%s`", TypeOAuth2, TypeOAuth2)
		} else {
			return creds.OAuth2.NewToken(ctx)
		}
	case TypeGoogleOAuth2:
		if creds.GoogleOAuth2 == nil {
			return nil, fmt.Errorf("credentials.%s is nil for type `%s`", TypeGoogleOAuth2, TypeGoogleOAuth2)
		} else {
			credsOAuth2 := creds.GoogleOAuth2.CredentialsOAuth2()
			return credsOAuth2.NewToken(ctx)
		}
	default:
		return nil, fmt.Errorf("creds type not supported [%s]", creds.Type)
	}
}

// NewTokenCLI retrieves a token using CLI approach for
// OAuth 2.0 authorization code or password grant.
func (creds *Credentials) NewTokenCLI(ctx context.Context, oauth2State string) (*oauth2.Token, error) {
	if strings.EqualFold(strings.TrimSpace(creds.OAuth2.GrantType), authutil.GrantTypeAuthorizationCode) {
		return NewTokenCLI(ctx, *creds, oauth2State)
	} else {
		return creds.NewToken(ctx)
	}
}
