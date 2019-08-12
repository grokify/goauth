package multiservice

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/os/osutil"
	"github.com/grokify/gotilla/type/stringsutil"

	"github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/aha"
	"github.com/grokify/oauth2more/facebook"
	"github.com/grokify/oauth2more/google"
	"github.com/grokify/oauth2more/multiservice/common"
	"github.com/grokify/oauth2more/multiservice/tokenset_memory"
	"github.com/grokify/oauth2more/ringcentral"
)

type OAuth2Manager struct {
	ConfigSet *ConfigSet
	TokenSet  common.TokenSet
}

func NewOAuth2Manager() *OAuth2Manager {
	return &OAuth2Manager{
		ConfigSet: NewConfigSet(),
		TokenSet:  memory.NewTokenSet(),
	}
}

func (cb *OAuth2Manager) GetClient(ctx context.Context, serviceKey string) (*http.Client, error) {
	if cb.ConfigSet == nil {
		return nil, fmt.Errorf("OAuth2Manager.ConfigSet == nil")
	}
	if cb.TokenSet == nil {
		return nil, fmt.Errorf("OAuth2Manager.TokenSet == nil")
	}

	cfgMore, err := cb.ConfigSet.Get(serviceKey)
	if err != nil {
		return nil, err
	}
	cfg := cfgMore.Config()
	tok, err := cb.TokenSet.GetToken(serviceKey)
	if err != nil {
		return nil, err
	}
	return cfg.Client(ctx, tok), nil
}

/*
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
*/
/*
type TokenInfo struct {
	ServiceKey  string
	ServiceType string
	Token       *oauth2.Token
}
*/
type ConfigSet struct {
	ConfigsMap map[string]*O2ConfigMore
}

func NewConfigSet() *ConfigSet {
	return &ConfigSet{ConfigsMap: map[string]*O2ConfigMore{}}
}

func (cfgs *ConfigSet) AddConfigMoreJson(key string, val []byte) error {
	key = strings.TrimSpace(key)
	cfg, err := NewO2ConfigMoreFromJSON(val)
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

func (cfgs *ConfigSet) Get(key string) (*O2ConfigMore, error) {
	if cfg, ok := cfgs.ConfigsMap[key]; ok {
		return cfg, nil
	}
	return nil, fmt.Errorf("AppConfig not found for %v", key)
}

func (cfgs *ConfigSet) MustGet(key string) *O2ConfigMore {
	c, err := cfgs.Get(key)
	if err != nil {
		panic(err)
	}
	return c
}

func (cfgs *ConfigSet) Slugs() []string {
	slugs := []string{}
	for slug := range cfgs.ConfigsMap {
		slugs = append(slugs, slug)
	}
	return slugs
}

func (cfgs *ConfigSet) ClientURLsMap() map[string]AppURLs {
	apps := map[string]AppURLs{}
	for slug, cfg := range cfgs.ConfigsMap {
		apps[slug] = AppURLs{
			AuthURL:     cfg.AuthUri,
			RedirectURL: stringsutil.SliceIndexOrEmpty(cfg.RedirectUris, 0),
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
			key := m[1]
			err := cfgs.AddConfigMoreJson(key, []byte(val))
			if err != nil {
				return nil, err
			}
		}
	}
	return cfgs, nil
}

func NewClientUtilForProviderType(providerType OAuth2Provider) (oauth2more.OAuth2Util, error) {
	switch provider {
	case Aha:
		return &aha.ClientUtil{}, nil
	case Facebook:
		return &facebook.ClientUtil{}, nil
	case Google:
		return &google.ClientUtil{}, nil
	case RingCentral:
		return &ringcentral.ClientUtil{}, nil
	default:
		return nil, fmt.Errorf("Cannot find ClientUtil for provider type [%s]", providerType)
	}
}
