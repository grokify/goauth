package ringcentral

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/caarlos0/env"
	"golang.org/x/oauth2"

	om "github.com/grokify/oauth2more"
)

// ApplicationConfigEnv returns a struct designed to be used to
// read values from the environment.
type ApplicationConfigEnv struct {
	ClientID     string `env:"RINGCENTRAL_CLIENT_ID"`
	ClientSecret string `env:"RINGCENTRAL_CLIENT_SECRET"`
	ServerURL    string `env:"RINGCENTRAL_SERVER_URL" envDefault:"https://platform.ringcentral.com"`
	AccessToken  string `env:"RINGCENTRAL_ACCESS_TOKEN"`
	Username     string `env:"RINGCENTRAL_USERNAME"`
	Extension    string `env:"RINGCENTRAL_EXTENSION"`
	Password     string `env:"RINGCENTRAL_PASSWORD"`
}

// NewApplicationConfigEnv returns a new ApplicationConfigEnv
// populated with values from the environment.
func NewApplicationConfigEnv() (ApplicationConfigEnv, error) {
	cfg := ApplicationConfigEnv{}
	return cfg, env.Parse(&cfg)
}

// ApplicationCredentials returns a ApplicationCredentials struct.
func (cfg *ApplicationConfigEnv) ApplicationCredentials() ApplicationCredentials {
	return ApplicationCredentials{
		ServerURL:    cfg.ServerURL,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret}
}

// PasswordCredentials returns a PasswordCredentials struct.
func (cfg *ApplicationConfigEnv) PasswordCredentials() PasswordCredentials {
	return PasswordCredentials{
		Username:  cfg.Username,
		Extension: cfg.Extension,
		Password:  cfg.Password}
}

// LoadToken loads and returns an OAuth token.
func (cfg *ApplicationConfigEnv) LoadToken() (*oauth2.Token, error) {
	tok, err := NewTokenPassword(
		cfg.ApplicationCredentials(),
		cfg.PasswordCredentials())
	if err == nil {
		cfg.AccessToken = tok.AccessToken
	}
	return tok, err
}

func NewHttpClientEnvFlexStatic(envPrefix string) (*http.Client, error) {
	envPrefix = strings.TrimSpace(envPrefix)
	if len(envPrefix) == 0 {
		envPrefix = "RINGCENTRAL_"
	}

	envToken := strings.TrimSpace(envPrefix + "TOKEN")
	token := os.Getenv(envToken)
	if len(token) > 0 {
		return om.NewClientBearerTokenSimple(token), nil
	}

	envPassword := strings.TrimSpace(envPrefix + "PASSWORD")
	password := os.Getenv(envPassword)
	if len(password) > 0 {
		return NewClientPassword(
			ApplicationCredentials{
				ClientID:     os.Getenv(envPrefix + "CLIENT_ID"),
				ClientSecret: os.Getenv(envPrefix + "CLIENT_SECRET"),
				ServerURL:    os.Getenv(envPrefix + "SERVER_URL")},
			PasswordCredentials{
				Username:  os.Getenv(envPrefix + "USERNAME"),
				Extension: os.Getenv(envPrefix + "EXTENSION"),
				Password:  os.Getenv(envPrefix + "PASSWORD")})
	}

	return nil, fmt.Errorf("Cannot load client from ENV for prefix [%v]", envPassword)
}
