package salesforce

import (
	"net/http"
	"os"
	"strings"

	"github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/credentials"
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
	TokenURL: TokenURL}

func NewClientPassword(oc credentials.OAuth2Credentials) (*http.Client, error) {
	conf := oauth2.Config{
		ClientID:     oc.ClientID,
		ClientSecret: oc.ClientSecret}

	if 1 == 0 {
		if len(strings.TrimSpace(oc.OAuth2Endpoint.AuthURL)) == 0 {
			conf.Endpoint = Endpoint
		} else {
			conf.Endpoint = oc.OAuth2Endpoint
		}
	}
	if 1 == 1 {
		if model.IsZero(oc.OAuth2Endpoint) {
			conf.Endpoint = Endpoint
		} else {
			conf.Endpoint = oc.OAuth2Endpoint
		}
	}
	return oauth2more.NewClientPasswordConf(conf, oc.Username, oc.Password)
}

func NewClientPasswordSalesforceEnv() (*http.Client, error) {
	return NewClientPassword(
		credentials.OAuth2Credentials{
			ClientID:     os.Getenv("SALESFORCE_CLIENT_ID"),
			ClientSecret: os.Getenv("SALESFORCE_CLIENT_SECRET"),
			Username:     os.Getenv("SALESFORCE_USERNAME"),
			Password: strings.Join([]string{
				os.Getenv("SALESFORCE_PASSWORD"),
				os.Getenv("SALESFORCE_SECURITY_TOKEN"),
			}, "")})
}
