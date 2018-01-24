package multiservice

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/os/osutil"
	"golang.org/x/oauth2"

	"github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/aha"
	"github.com/grokify/oauth2more/facebook"
	"github.com/grokify/oauth2more/google"
	"github.com/grokify/oauth2more/ringcentral"
)

type OAuth2Manager struct {
	ConfigSet *ConfigSet
	TokenSet  *TokenSet
}

func NewOAuth2Manager() *OAuth2Manager {
	return &OAuth2Manager{
		ConfigSet: NewConfigSet(),
		TokenSet:  NewTokenSet(),
	}
}

func (cb *OAuth2Manager) GetClient(ctx context.Context, serviceKey string) (*http.Client, error) {
	if cb.ConfigSet == nil {
		return nil, fmt.Errorf("OAuth2Manager.ConfigSet == nil")
	}
	if cb.TokenSet == nil {
		return nil, fmt.Errorf("OAuth2Manager.TokenSet == nil")
	}

	cfg, err := cb.ConfigSet.Get(serviceKey)
	if err != nil {
		return nil, err
	}
	tok, err := cb.TokenSet.Get(serviceKey)
	if err != nil {
		return nil, err
	}
	return cfg.Client(ctx, tok), nil
}

type TokenSet struct {
	TokenMap map[string]*TokenInfo
}

func NewTokenSet() *TokenSet {
	return &TokenSet{TokenMap: map[string]*TokenInfo{}}
}

func (toks *TokenSet) Get(key string) (*oauth2.Token, error) {
	if tok, ok := toks.TokenMap[key]; ok {
		return tok.Token, nil
	}
	return nil, fmt.Errorf("AppConfig not found for %v", key)
}

type TokenInfo struct {
	ServiceKey  string
	ServiceType string
	Token       *oauth2.Token
}

type ConfigSet struct {
	ConfigsMap map[string]*oauth2.Config
}

func NewConfigSet() *ConfigSet {
	return &ConfigSet{ConfigsMap: map[string]*oauth2.Config{}}
}

func (cfgs *ConfigSet) AddAppConfigWrapperBytes(key string, val []byte) error {
	acw, err := oauth2more.NewAppCredentialsWrapperFromBytes(val)
	if err != nil {
		return err
	}
	return cfgs.AddAppConfigWrapper(key, acw)
}

func (cfgs *ConfigSet) AddAppConfigWrapper(key string, acw oauth2more.AppCredentialsWrapper) error {
	cfg, err := acw.Config()
	if err != nil {
		return err
	}
	cfgs.ConfigsMap[key] = cfg
	return nil
}

func (cfgs *ConfigSet) Has(key string) bool {
	if _, ok := cfgs.ConfigsMap[key]; ok {
		return true
	}
	return false
}

func (cfgs *ConfigSet) Get(key string) (*oauth2.Config, error) {
	if cfg, ok := cfgs.ConfigsMap[key]; ok {
		return cfg, nil
	}
	return nil, fmt.Errorf("AppConfig not found for %v", key)
}

func (cfgs *ConfigSet) MustGet(key string) *oauth2.Config {
	c, err := cfgs.Get(key)
	if err != nil {
		panic(err)
	}
	return c
}

func (cfgs *ConfigSet) Slugs() []string {
	slugs := []string{}
	for slug, _ := range cfgs.ConfigsMap {
		slugs = append(slugs, slug)
	}
	return slugs
}

func (cfgs *ConfigSet) ClientURLsMap() map[string]AppURLs {
	apps := map[string]AppURLs{}
	for slug, cfg := range cfgs.ConfigsMap {
		apps[slug] = AppURLs{
			AuthURL:     cfg.Endpoint.AuthURL,
			RedirectURL: cfg.RedirectURL,
		}
	}
	return apps
}

type AppURLs struct {
	AuthURL     string `json:"authUrl,omitempty"`
	TokenURL    string `json:"tokenUrl,omitempty"`
	RedirectURL string `json:"redirectUrl,omitempty"`
}

// EnvOAuth2ConfigMap returns a map of *oauth2.Config from environment
// variables in AppCredentialsWrapper format.
func EnvOAuth2ConfigMap(env []osutil.EnvVar, prefix string) (*ConfigSet, error) {
	cfgs := NewConfigSet()

	rx, err := regexp.Compile(fmt.Sprintf(`^%v(.*)`, prefix))
	if err != nil {
		return nil, err
	}

	for _, pair := range env {
		key := strings.TrimSpace(pair.Key)
		val := pair.Value
		m := rx.FindStringSubmatch(key)
		if len(m) > 0 {
			fmt.Println(val)
			key := m[1]
			err := cfgs.AddAppConfigWrapperBytes(key, []byte(val))
			if err != nil {
				return nil, err
			}
		}
	}
	return cfgs, nil
}

func NewClientUtilForServiceType(svcType string) (oauth2more.OAuth2Util, error) {
	switch strings.ToLower(strings.TrimSpace(svcType)) {
	case "aha":
		return &aha.ClientUtil{}, nil
	case "facebook":
		return &facebook.ClientUtil{}, nil
	case "google":
		return &google.ClientUtil{}, nil
	case "ringcentral":
		return &ringcentral.ClientUtil{}, nil
	default:
		return nil, fmt.Errorf("Cannot find ClientUtil for service type %v", svcType)
	}
}
