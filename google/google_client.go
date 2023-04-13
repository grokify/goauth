package google

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func ClientFromFile(ctx context.Context, filepath string, scopes []string, tok *oauth2.Token) (*http.Client, error) {
	conf, err := ConfigFromFile(filepath, scopes)
	if err != nil {
		return &http.Client{}, errorsutil.Wrap(err, fmt.Sprintf("Unable to read app config file: %v", filepath))
	}

	return conf.Client(ctx, tok), nil
}

type ClientOauthCliTokenStoreConfig struct {
	Context       context.Context
	AppConfig     []byte
	Scopes        []string
	TokenFile     string
	ForceNewToken bool
	State         string
}

func NewClientOauthCliTokenStore(cfg ClientOauthCliTokenStoreConfig) (*http.Client, error) {
	conf, err := ConfigFromBytes(cfg.AppConfig, cfg.Scopes)
	if err != nil {
		return nil, err
	}

	tokenStore, err := authutil.NewTokenStoreFileDefault(cfg.TokenFile, true, 0700)
	if err != nil {
		return nil, err
	}

	return authutil.NewClientWebTokenStore(cfg.Context, conf, tokenStore, cfg.ForceNewToken, cfg.State)
}

func NewClientSvcAccountFromFile(ctx context.Context, svcAccountConfigFile string, scopes ...string) (*http.Client, error) {
	svcAccountConfig, err := os.ReadFile(svcAccountConfigFile)
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
