package ringcentral

import (
	"github.com/caarlos0/env"
	"golang.org/x/oauth2"
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
