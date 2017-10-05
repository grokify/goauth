package ringcentral

import (
	"golang.org/x/oauth2"
	"net/http"
)

var (
	EnvServerURL    = "RC_SERVER_URL"
	EnvClientID     = "RC_CLIENT_ID"
	EnvClientSecret = "RC_CLIENT_SECRET"
	EnvUsername     = "RC_USER_USERNAME"
	EnvExtension    = "RC_USER_EXTENSION"
	EnvPassword     = "RC_USER_PASSWORD"
)

type ApplicationCredentials struct {
	ServerURL    string
	ClientID     string
	ClientSecret string
}

type UserCredentials struct {
	Username  string
	Extension string
	Password  string
}

func NewClientPasswordConf(conf oauth2.Config, username, password string) (*http.Client, error) {
	token, err := conf.PasswordCredentialsToken(oauth2.NoContext, username, password)

	if err != nil {
		return &http.Client{}, err
	}

	return conf.Client(oauth2.NoContext, token), nil
}

func NewClientPassword(app ApplicationCredentials, user UserCredentials) (*http.Client, error) {
	conf := oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Endpoint:     NewEndpoint(app.ServerURL)}

	return NewClientPasswordConf(conf, user.Username, user.Password)
}
