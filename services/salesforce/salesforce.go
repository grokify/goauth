package salesforce

import (
	"fmt"
	"net/http"
	"net/url"

	ou "github.com/grokify/oauth2util-go"
	"golang.org/x/oauth2"
	"gopkg.in/jeevatkm/go-model.v1"
)

const (
	AuthzURL        = "https://login.salesforce.com/services/oauth2/authorize"
	TokenURL        = "https://login.salesforce.com/services/oauth2/token"
	RevokeURL       = "https://login.salesforce.com/services/oauth2/revoke"
	ServerURLFormat = "https://%v.salesforce.com"
	HostFormat      = "%v.salesforce.com"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  AuthzURL,
	TokenURL: TokenURL,
}

func NewClientPassword(app ou.ApplicationCredentials, user ou.UserCredentials) (*http.Client, error) {
	conf := oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Endpoint:     Endpoint}

	if model.IsZero(app.Endpoint) {
		conf.Endpoint = Endpoint
	} else {
		conf.Endpoint = app.Endpoint
	}

	return ou.NewClientPasswordConf(conf, user.Username, user.Password)
}

type URLBuilder struct {
	BaseURL url.URL
}

func NewURLBuilder(instanceName string) URLBuilder {
	u := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf(HostFormat, instanceName),
	}
	return URLBuilder{BaseURL: u}
}

func (b *URLBuilder) Build(path string) url.URL {
	u := b.BaseURL
	u.Path = path
	return u
}
