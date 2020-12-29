package google

import (
	"io/ioutil"

	"github.com/grokify/simplego/encoding/jsonutil"
	json "github.com/pquerna/ffjson/ffjson"
	"golang.org/x/oauth2"
	o2g "golang.org/x/oauth2/google"
)

type CredentialsContainer struct {
	Web    *Credentials `json:"web,omitempty"`
	Raw    []byte       `json:"-"`
	Scopes []string     `json:"scopes,omitempty"` // optional for self-contained app credentials
}

func (cc *CredentialsContainer) OAuth2Config(scopes ...string) (*oauth2.Config, error) {
	return o2g.ConfigFromJSON(jsonutil.MustMarshalSimple(cc, "", ""), scopes...)
}

func (cc *CredentialsContainer) Credentials() *Credentials {
	return cc.Web
}

type Credentials struct {
	Type                    string   `json:"type,omitempty"`
	ProjectID               string   `json:"project_id,omitempty"`
	PrivateKeyID            string   `json:"private_key_id,omitempty"`
	PrivateKey              string   `json:"private_key,omitempty"`
	ClientEmail             string   `json:"client_email,omitempty"`
	ClientID                string   `json:"client_id,omitempty"`
	ClientSecret            string   `json:"client_secret,omitempty"`
	AuthURI                 string   `json:"auth_uri,omitempty"`
	TokenURI                string   `json:"token_uri,omitempty"`
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url,omitempty"`
	ClientX509CertURL       string   `json:"client_x509_cert_url,omitempty"`
	RedirectURIs            []string `json:"redirect_uris,omitempty"`
}

func CredentialsContainerFromBytes(bytes []byte) (CredentialsContainer, error) {
	creds := CredentialsContainer{Raw: bytes}
	return creds, json.Unmarshal(bytes, &creds)
}

func CredentialsContainerFromFile(file string) (CredentialsContainer, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return CredentialsContainer{}, err
	}
	return CredentialsContainerFromBytes(bytes)
}

func CredentialsFromFile(file string) (Credentials, error) {
	c := Credentials{}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}
	return c, json.Unmarshal(bytes, &c)
}
