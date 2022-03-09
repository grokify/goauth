package multiservice

import (
	"fmt"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

type ConfigMoreSet struct {
	ConfigMoreMap map[string]*O2ConfigMore
}

func NewConfigMoreSet() *ConfigMoreSet {
	return &ConfigMoreSet{ConfigMoreMap: map[string]*O2ConfigMore{}}
}

func (cfgs *ConfigMoreSet) AddConfigMoreJSON(key string, val []byte) error {
	key = strings.TrimSpace(key)
	cfg, err := NewO2ConfigMoreFromJSON(val)
	if err != nil {
		return err
	}
	cfgs.ConfigMoreMap[key] = cfg
	return nil
}

func (cfgs *ConfigMoreSet) Has(key string) bool {
	if _, ok := cfgs.ConfigMoreMap[key]; ok {
		return true
	}
	return false
}

func (cfgs *ConfigMoreSet) Get(key string) (*O2ConfigMore, error) {
	if cfg, ok := cfgs.ConfigMoreMap[key]; ok {
		return cfg, nil
	}
	return nil, fmt.Errorf("AppConfig not found for %v", key)
}

func (cfgs *ConfigMoreSet) MustGet(key string) *O2ConfigMore {
	c, err := cfgs.Get(key)
	if err != nil {
		panic(err)
	}
	return c
}

func (cfgs *ConfigMoreSet) Slugs() []string {
	slugs := []string{}
	for slug := range cfgs.ConfigMoreMap {
		slugs = append(slugs, slug)
	}
	return slugs
}

func (cfgs *ConfigMoreSet) ClientURLsMap() map[string]AppURLs {
	apps := map[string]AppURLs{}
	for slug, cfg := range cfgs.ConfigMoreMap {
		apps[slug] = AppURLs{
			AuthURL:     cfg.AuthURI,
			RedirectURL: stringsutil.SliceIndexOrEmpty(cfg.RedirectURIs, 0),
		}
	}
	return apps
}
