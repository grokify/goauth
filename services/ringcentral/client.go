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

func NewClientPassword(app ApplicationCredentials, user UserCredentials) (*http.Client, error) {
	cfg := oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Endpoint:     NewEndpoint(app.ServerURL)}

	token, err := cfg.PasswordCredentialsToken(
		oauth2.NoContext,
		user.Username,
		user.Password)

	if err != nil {
		return &http.Client{}, err
	}

	return cfg.Client(oauth2.NoContext, token), nil
}
