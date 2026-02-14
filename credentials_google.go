package goauth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"slices"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/google"
	"golang.org/x/oauth2"
)

// CredentialsGCP supports OAuth 2.0 authorization_code, password, and client_credentials grant flows.
type CredentialsGCP struct {
	GCPCredentials google.Credentials `json:"gcpCredentials,omitempty"`
	Scopes         []string           `json:"scopes,omitempty"`
}

// NewClient returns a `*http.Client` and `error`.
func (cg *CredentialsGCP) NewClient(ctx context.Context) (*http.Client, error) {
	return cg.GCPCredentials.NewClient(ctx, cg.Scopes)
}

func CredentialsGCPReadFile(name string) (*CredentialsGCP, error) {
	if b, err := os.ReadFile(name); err != nil {
		return nil, err
	} else {
		var c *CredentialsGCP
		return c, json.Unmarshal(b, c)
	}
}

type CredentialsGoogleOAuth2 struct {
	GoogleWebCredentials google.Credentials `json:"web,omitempty"` // "web"
	Scopes               []string           `json:"scopes,omitempty"`
	Token                *oauth2.Token      `json:"token,omitempty"`
}

func (cgo CredentialsGoogleOAuth2) CredentialsOAuth2() CredentialsOAuth2 {
	gcreds := cgo.GoogleWebCredentials
	coauth2 := CredentialsOAuth2{
		GrantType:    authutil.GrantTypeAuthorizationCode,
		ClientID:     gcreds.ClientID,
		ClientSecret: gcreds.ClientSecret,
		Endpoint:     gcreds.OAuth2Endpoint(),
		Scopes:       slices.Clone(cgo.Scopes)}
	coauth2.RedirectURL = gcreds.FirstRedirectURI()
	return coauth2
}
