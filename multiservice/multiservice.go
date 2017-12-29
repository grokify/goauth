package multiservice

import (
	"fmt"
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

type AppConfigs struct {
	ConfigsMap map[string]*oauth2.Config
}

func NewAppConfigs() *AppConfigs {
	return &AppConfigs{ConfigsMap: map[string]*oauth2.Config{}}
}

func (cfgs *AppConfigs) AddAppConfigWrapperBytes(key string, val []byte) error {
	acw, err := oauth2more.NewAppCredentialsWrapperFromBytes(val)
	if err != nil {
		return err
	}
	return cfgs.AddAppConfigWrapper(key, acw)
}

func (cfgs *AppConfigs) AddAppConfigWrapper(key string, acw oauth2more.AppCredentialsWrapper) error {
	cfg, err := acw.Config()
	if err != nil {
		return err
	}
	cfgs.ConfigsMap[key] = cfg
	return nil
}

func (cfgs *AppConfigs) Get(key string) (*oauth2.Config, error) {
	if cfg, ok := cfgs.ConfigsMap[key]; ok {
		return cfg, nil
	}
	return nil, fmt.Errorf("AppConfig not found for %v", key)
}

func (cfgs *AppConfigs) MustGet(key string) *oauth2.Config {
	c, err := cfgs.Get(key)
	if err != nil {
		panic(err)
	}
	return c
}

// EnvOAuth2ConfigMap returns a map of *oauth2.Config from environment
// variables in AppCredentialsWrapper format.
func EnvOAuth2ConfigMap(env []osutil.EnvVar, prefix string) (*AppConfigs, error) {
	cfgs := NewAppConfigs()

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

func GetClientUtilForServiceType(svcType string) (oauth2more.OAuth2Util, error) {
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
