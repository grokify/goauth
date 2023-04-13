package salesforce

import (
	"net/http"
	"os"
	"strings"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/credentials"
	"github.com/grokify/goauth/endpoints"
	"golang.org/x/oauth2"
	model "gopkg.in/jeevatkm/go-model.v1"
)

const (
	AuthzURL        = endpoints.SalesforceAuthzURL
	TokenURL        = endpoints.SalesforceTokenURL // #nosec G101
	RevokeURL       = endpoints.SalesforceRevokeURL
	ServerURLFormat = "https://%v.salesforce.com"
	HostFormat      = "%v.salesforce.com"
	TestServerURL   = "https://test.salesforce.com"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  AuthzURL,
	TokenURL: TokenURL}

func NewClientPassword(oc credentials.CredentialsOAuth2) (*http.Client, error) {
	conf := oauth2.Config{
		ClientID:     oc.ClientID,
		ClientSecret: oc.ClientSecret}

	if 1 == 0 {
		if len(strings.TrimSpace(oc.Endpoint.AuthURL)) == 0 {
			conf.Endpoint = Endpoint
		} else {
			conf.Endpoint = oc.Endpoint
		}
	}

	if model.IsZero(oc.Endpoint) {
		conf.Endpoint = Endpoint
	} else {
		conf.Endpoint = oc.Endpoint
	}

	return authutil.NewClientPasswordConf(conf, oc.Username, oc.Password)
}

func NewClientPasswordSalesforceEnv() (*http.Client, error) {
	return NewClientPassword(
		credentials.CredentialsOAuth2{
			ClientID:     os.Getenv("SALESFORCE_CLIENT_ID"),
			ClientSecret: os.Getenv("SALESFORCE_CLIENT_SECRET"),
			Username:     os.Getenv("SALESFORCE_USERNAME"),
			Password: strings.Join([]string{
				os.Getenv("SALESFORCE_PASSWORD"),
				os.Getenv("SALESFORCE_SECURITY_TOKEN"),
			}, "")})
}
