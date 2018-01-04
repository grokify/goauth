package ringcentral

import (
	"net/http"
	"os"
	"strings"

	ou "github.com/grokify/oauth2more"
	"golang.org/x/oauth2"
)

var (
	EnvServerURL    = "RINGCENTRAL_SERVER_URL"
	EnvClientID     = "RINGCENTRAL_CLIENT_ID"
	EnvClientSecret = "RINGCENTRAL_CLIENT_SECRET"
	EnvUsername     = "RINGCENTRAL_USERNAME"
	EnvExtension    = "RINGCENTRAL_EXTENSION"
	EnvPassword     = "RINGCENTRAL_PASSWORD"
)

type ApplicationCredentials struct {
	ServerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (app *ApplicationCredentials) Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Endpoint:     NewEndpoint(app.ServerURL),
		RedirectURL:  app.RedirectURL}
}

type UserCredentials struct {
	Username  string
	Extension string
	Password  string
}

func (uc *UserCredentials) UsernameSimple() string {
	if len(strings.TrimSpace(uc.Extension)) > 0 {
		return strings.Join([]string{uc.Username, uc.Extension}, "*")
	}
	return uc.Username
}

func NewClientPassword(app ApplicationCredentials, user UserCredentials) (*http.Client, error) {
	return ou.NewClientPasswordConf(
		oauth2.Config{
			ClientID:     app.ClientID,
			ClientSecret: app.ClientSecret,
			Endpoint:     NewEndpoint(app.ServerURL)},
		user.UsernameSimple(),
		user.Password)
}

func NewClientPasswordEnv() (*http.Client, error) {
	return NewClientPassword(
		ApplicationCredentials{
			ServerURL:    os.Getenv(EnvServerURL),
			ClientID:     os.Getenv(EnvClientID),
			ClientSecret: os.Getenv(EnvClientSecret),
		},
		UserCredentials{
			Username:  os.Getenv(EnvUsername),
			Password:  os.Getenv(EnvPassword),
			Extension: os.Getenv(EnvExtension),
		},
	)
}
