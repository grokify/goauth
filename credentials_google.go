package goauth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/grokify/goauth/google"
)

// CredentialsOAuth2 supports OAuth 2.0 authorization_code, password, and client_credentials grant flows.
type CredentialsGCP struct {
	GCPCredentials google.Credentials `json:"gcpCredentials,omitempty"`
	Scopes         []string           `json:"scopes,omitempty"`
}

// NewClient returns a `*http.Client` and `error`.
func (cg *CredentialsGCP) NewClient(ctx context.Context) (*http.Client, error) {
	hclient, err := cg.GCPCredentials.NewClient(ctx, cg.Scopes)
	return hclient, err
}

func CredentialsGCPReadFile(name string) (*CredentialsGCP, error) {
	if b, err := os.ReadFile(name); err != nil {
		return nil, err
	} else {
		var c *CredentialsGCP
		return c, json.Unmarshal(b, c)
	}
}
