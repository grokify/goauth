package google

import (
	"fmt"
	"io/ioutil"
	"net/http"

	om "github.com/grokify/oauth2more"
	"github.com/pkg/errors"
	json "github.com/pquerna/ffjson/ffjson"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func ClientFromFile(ctx context.Context, filepath string, scopes []string, tok *oauth2.Token) (*http.Client, error) {
	conf, err := ConfigFromFile(filepath, scopes)
	if err != nil {
		return &http.Client{}, errors.Wrap(err, fmt.Sprintf("Unable to read app config file: %v", filepath))
	}

	return conf.Client(ctx, tok), nil
}

type ClientOauthCliTokenStoreConfig struct {
	Context       context.Context
	AppConfig     []byte
	Scopes        []string
	TokenFile     string
	ForceNewToken bool
}

func NewClientOauthCliTokenStore(cfg ClientOauthCliTokenStoreConfig) (*http.Client, error) {
	conf, err := ConfigFromBytes(cfg.AppConfig, cfg.Scopes)
	if err != nil {
		return nil, err
	}

	tokenStore, err := om.NewTokenStoreFileDefault(cfg.TokenFile, true, 0700)
	if err != nil {
		return nil, err
	}

	return om.NewClientWebTokenStore(cfg.Context, conf, tokenStore, cfg.ForceNewToken)
}

func NewClientSvcAccountFromFile(ctx context.Context, svcAccountConfigFile string, scopes ...string) (*http.Client, error) {
	svcAccountConfig, err := ioutil.ReadFile(svcAccountConfigFile)
	if err != nil {
		return nil, err
	}
	return NewClientFromJWTJSON(ctx, svcAccountConfig, scopes...)
}

func NewClientFromJWTJSON(ctx context.Context, svcAccountConfig []byte, scopes ...string) (*http.Client, error) {
	jwtConf, err := google.JWTConfigFromJSON(svcAccountConfig, scopes...)
	if err != nil {
		return nil, err
	}
	return jwtConf.Client(ctx), nil
}

type Credentials struct {
	Type                    string `json:"type,omitempty"`
	ProjectID               string `json:"project_id,omitempty"`
	PrivateKeyID            string `json:"private_key_id,omitempty"`
	PrivateKey              string `json:"private_key,omitempty"`
	ClientEmail             string `json:"client_email,omitempty"`
	ClientID                string `json:"client_id,omitempty"`
	AuthURI                 string `json:"auth_uri,omitempty"`
	TokenURI                string `json:"token_uri,omitempty"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url,omitempty"`
	ClientX509CertURL       string `json:"client_x509_cert_url,omitempty"`
}

func NewCredentialsFromFile(file string) (Credentials, error) {
	c := Credentials{}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}
	return c, json.Unmarshal(bytes, &c)
}
