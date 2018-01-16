package ringcentral

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	hum "github.com/grokify/gotilla/net/httputilmore"
	ou "github.com/grokify/oauth2more"
	"golang.org/x/oauth2"
)

var (
	EnvServerURL    = "RINGCENTRAL_SERVER_URL"
	EnvClientID     = "RINGCENTRAL_CLIENT_ID"
	EnvClientSecret = "RINGCENTRAL_CLIENT_SECRET"
	EnvAppName      = "RINGCENTRAL_APP_NAME"
	EnvAppVersion   = "RINGCENTRAL_APP_VERSION"
	EnvUsername     = "RINGCENTRAL_USERNAME"
	EnvExtension    = "RINGCENTRAL_EXTENSION"
	EnvPassword     = "RINGCENTRAL_PASSWORD"
)

type ApplicationCredentials struct {
	ServerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AppName      string
	AppVersion   string
}

func (ac *ApplicationCredentials) AppNameAndVersion() string {
	parts := []string{}
	ac.AppName = strings.TrimSpace(ac.AppName)
	ac.AppVersion = strings.TrimSpace(ac.AppVersion)
	if len(ac.AppName) > 0 {
		parts = append(parts, ac.AppName)
	}
	if len(ac.AppVersion) > 0 {
		parts = append(parts, fmt.Sprintf("v%v", ac.AppVersion))
	}
	return strings.Join(parts, "-")
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
	httpClient, err := ou.NewClientPasswordConf(
		oauth2.Config{
			ClientID:     app.ClientID,
			ClientSecret: app.ClientSecret,
			Endpoint:     NewEndpoint(app.ServerURL)},
		user.UsernameSimple(),
		user.Password)
	if err != nil {
		return nil, err
	}

	header := getClientHeader(app)
	if len(header) > 0 {
		httpClient.Transport = hum.TransportWithHeaders{
			Transport: httpClient.Transport,
			Header:    header}
	}

	return httpClient, nil
}

func getClientHeader(app ApplicationCredentials) http.Header {
	userAgentParts := []string{ou.PathVersion()}
	if len(app.AppNameAndVersion()) > 0 {
		userAgentParts = append([]string{app.AppNameAndVersion()}, userAgentParts...)
	}
	userAgent := strings.TrimSpace(strings.Join(userAgentParts, "; "))

	header := http.Header{}
	if len(userAgent) > 0 {
		header.Add("User-Agent", userAgent)
		header.Add("X-User-Agent", userAgent)
	}
	return header
}

func NewClientPasswordEnv() (*http.Client, error) {
	return NewClientPassword(
		NewApplicationCredentialsEnv(),
		NewUserCredentialsEnv(),
	)
}

func NewApplicationCredentialsEnv() ApplicationCredentials {
	return ApplicationCredentials{
		ServerURL:    os.Getenv(EnvServerURL),
		ClientID:     os.Getenv(EnvClientID),
		ClientSecret: os.Getenv(EnvClientSecret),
		AppName:      os.Getenv(EnvAppName),
		AppVersion:   os.Getenv(EnvAppVersion)}
}

func NewUserCredentialsEnv() UserCredentials {
	return UserCredentials{
		Username:  os.Getenv(EnvUsername),
		Extension: os.Getenv(EnvExtension),
		Password:  os.Getenv(EnvPassword)}
}
