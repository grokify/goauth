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
	if conf, err := ConfigFromFile(filepath, scopes); err != nil {
		return &http.Client{}, errorsutil.Wrap(err, fmt.Sprintf("Unable to read app config file: %v", filepath))
	} else {
		return conf.Client(ctx, tok), nil
	}
}

type ClientOAuthCLITokenStoreConfig struct {
	Context       context.Context
	AppConfig     []byte
	Scopes        []string
	TokenFile     string
	ForceNewToken bool
	State         string
}

func NewClientOAuthCLITokenStore(cfg ClientOAuthCLITokenStoreConfig) (*http.Client, error) {
	if conf, err := ConfigFromBytes(cfg.AppConfig, cfg.Scopes); err != nil {
		return nil, err
	} else if tokenStore, err := authutil.NewTokenStoreFileDefault(cfg.TokenFile, true, 0600); err != nil {
		return nil, err
	} else {
		return authutil.NewClientWebTokenStore(cfg.Context, conf, tokenStore, cfg.ForceNewToken, cfg.State)
	}
}

func NewClientSvcAccountFromFile(ctx context.Context, svcAccountConfigFile string, scopes ...string) (*http.Client, error) {
	if svcAccountConfig, err := os.ReadFile(svcAccountConfigFile); err != nil {
		return nil, err
	} else {
		return NewClientFromJWTJSON(ctx, svcAccountConfig, scopes...)
	}
}

func NewClientFromJWTJSON(ctx context.Context, svcAccountConfig []byte, scopes ...string) (*http.Client, error) {
	if jwtConf, err := google.JWTConfigFromJSON(svcAccountConfig, scopes...); err != nil {
		return nil, err
	} else {
		return jwtConf.Client(ctx), nil
	}
}
