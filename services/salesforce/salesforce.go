package salesforce

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/grokify/gotilla/net/httputilmore"
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
		ClientSecret: app.ClientSecret}

	if model.IsZero(app.Endpoint) {
		conf.Endpoint = Endpoint
	} else {
		conf.Endpoint = app.Endpoint
	}

	return ou.NewClientPasswordConf(conf, user.Username, user.Password)
}

func NewClientPasswordSalesforceEnv() (*http.Client, error) {
	return NewClientPassword(
		ou.ApplicationCredentials{
			ClientID:     os.Getenv("SALESFORCE_CLIENT_ID"),
			ClientSecret: os.Getenv("SALESFORCE_CLIENT_SECRET")},
		ou.UserCredentials{
			Username: os.Getenv("SALESFORCE_USERNAME"),
			Password: fmt.Sprintf("%v%v",
				os.Getenv("SALESFORCE_PASSWORD"),
				os.Getenv("SALESFORCE_SECURITY_KEY"))})
}

type SalesforceClient struct {
	ClientMore httputilmore.ClientMore
	URLBuilder URLBuilder
}

func NewSalesforceClientEnv() (SalesforceClient, error) {
	sc := SalesforceClient{
		URLBuilder: NewURLBuilder(os.Getenv("SALESFORCE_INSTANCE_NAME")),
	}
	client, err := NewClientPasswordSalesforceEnv()
	if err != nil {
		return sc, err
	}
	sc.ClientMore = httputilmore.ClientMore{Client: client}
	return sc, nil
}

func (sc *SalesforceClient) GetServicesData() (*http.Response, error) {
	apiURL := sc.URLBuilder.Build("services/data")
	return sc.ClientMore.Client.Get(apiURL.String())
}

func (sc *SalesforceClient) CreateContact(contact interface{}) (*http.Response, error) {
	apiURL := sc.URLBuilder.Build("/services/data/v40.0/sobjects/Contact/")
	return sc.ClientMore.PostToJSON(apiURL.String(), contact)
}

type URLBuilder struct {
	BaseURL url.URL
}

func NewURLBuilder(instanceName string) URLBuilder {
	return URLBuilder{BaseURL: url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf(HostFormat, instanceName)}}
}

func (b *URLBuilder) Build(path string) url.URL {
	u := b.BaseURL
	u.Path = path
	return u
}
