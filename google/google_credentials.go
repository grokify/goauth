package google

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/grokify/mogo/encoding/jsonutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type CredentialsContainer struct {
	Web    *Credentials `json:"web,omitempty"`
	Raw    []byte       `json:"-"`
	Scopes []string     `json:"scopes,omitempty"` // optional for self-contained app credentials
}

func (cc *CredentialsContainer) OAuth2Config(scopes ...string) (*oauth2.Config, error) {
	return google.ConfigFromJSON(jsonutil.MustMarshalSimple(cc, "", ""), scopes...)
}

func (cc *CredentialsContainer) Credentials() *Credentials {
	return cc.Web
}

// Credentials represents a full GCP Service Account Key file. A simplified version is available in
// https://github.com/hashicorp/go-gcp-common/blob/main/gcputil/credentials.go#L44-L51
type Credentials struct {
	Type                    string   `json:"type,omitempty"`
	ClientEmail             string   `json:"client_email" structs:"client_email" mapstructure:"client_email"`
	ClientID                string   `json:"client_id" structs:"client_id" mapstructure:"client_id"`
	ClientSecret            string   `json:"client_secret,omitempty"`
	ProjectID               string   `json:"project_id" structs:"project_id" mapstructure:"project_id"`
	PrivateKey              string   `json:"private_key" structs:"private_key" mapstructure:"private_key"`
	PrivateKeyID            string   `json:"private_key_id" structs:"private_key_id" mapstructure:"private_key_id"`
	AuthURI                 string   `json:"auth_uri,omitempty"`
	TokenURI                string   `json:"token_uri,omitempty"`
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url,omitempty"`
	ClientX509CertURL       string   `json:"client_x509_cert_url,omitempty"`
	RedirectURIs            []string `json:"redirect_uris,omitempty"`
}

func (c Credentials) NewClient(ctx context.Context, scopes []string) (*http.Client, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	jwtConf, err := google.JWTConfigFromJSON(b, scopes...)
	if err != nil {
		return nil, err
	}
	return jwtConf.Client(ctx), nil
}

func ReadCredentialsFile(name string) (*Credentials, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	c := &Credentials{}
	err = json.Unmarshal(b, c)
	return c, err
}

func CredentialsContainerFromBytes(bytes []byte) (CredentialsContainer, error) {
	creds := CredentialsContainer{Raw: bytes}
	return creds, json.Unmarshal(bytes, &creds)
}

func CredentialsContainerFromFile(file string) (CredentialsContainer, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return CredentialsContainer{}, err
	}
	return CredentialsContainerFromBytes(bytes)
}

func CredentialsFromFile(file string) (Credentials, error) {
	c := Credentials{}
	bytes, err := os.ReadFile(file)
	if err != nil {
		return c, err
	}
	return c, json.Unmarshal(bytes, &c)
}
