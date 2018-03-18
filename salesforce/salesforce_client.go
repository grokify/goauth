package salesforce

import (
	"net/http"
	"os"
	"strings"

	om "github.com/grokify/oauth2more"
	"golang.org/x/oauth2"
	"gopkg.in/jeevatkm/go-model.v1"
)

const (
	AuthzURL        = "https://login.salesforce.com/services/oauth2/authorize"
	TokenURL        = "https://login.salesforce.com/services/oauth2/token"
	RevokeURL       = "https://login.salesforce.com/services/oauth2/revoke"
	ServerURLFormat = "https://%v.salesforce.com"
	HostFormat      = "%v.salesforce.com"
	TestServerURL   = "https://test.salesforce.com"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  AuthzURL,
	TokenURL: TokenURL,
}

func NewClientPassword(app om.ApplicationCredentials, user om.UserCredentials) (*http.Client, error) {
	conf := oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret}

	//conf.Endpoint = app.Endpoint

	if 1 == 0 {
		if len(strings.TrimSpace(app.Endpoint.AuthURL)) == 0 {
			conf.Endpoint = Endpoint
		} else {
			conf.Endpoint = app.Endpoint
		}
	}
	if 1 == 1 {
		if model.IsZero(app.Endpoint) {
			conf.Endpoint = Endpoint
		} else {
			conf.Endpoint = app.Endpoint
		}
	}
	return om.NewClientPasswordConf(conf, user.Username, user.Password)
}

func NewClientPasswordSalesforceEnv() (*http.Client, error) {
	return NewClientPassword(
		om.ApplicationCredentials{
			ClientID:     os.Getenv("SALESFORCE_CLIENT_ID"),
			ClientSecret: os.Getenv("SALESFORCE_CLIENT_SECRET")},
		om.UserCredentials{
			Username: os.Getenv("SALESFORCE_USERNAME"),
			Password: strings.Join([]string{
				os.Getenv("SALESFORCE_PASSWORD"),
				os.Getenv("SALESFORCE_SECURITY_TOKEN"),
			}, "")})
}
