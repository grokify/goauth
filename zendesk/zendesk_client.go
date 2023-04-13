package zendesk

import (
	"context"
	"net/http"

	"github.com/grokify/goauth/authutil"
	"golang.org/x/oauth2"
)

var (
	EnvZendeskUsername  = "ZENDESK_USERNAME"
	EnvZendeskPassword  = "ZENDESK_PASSWORD"
	EnvZendeskSubdomain = "ZENDESK_SUBDOMAIN"
)

// NewClientPassword creates a new http.Client using basic authentication.
func NewClientPassword(ctx context.Context, emailAddress, password string) (*http.Client, error) {
	token, err := authutil.BasicAuthToken(emailAddress, password)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{}
	return conf.Client(ctx, token), nil
}

// NewClientToken creates a new http.Client using the Zendesk API token
// as specified here: https://developer.zendesk.com/rest_api/docs/core/introduction
func NewClientToken(ctx context.Context, emailAddress, apiToken string) (*http.Client, error) {
	token, err := authutil.BasicAuthToken(emailAddress+"/token", apiToken)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{}
	return conf.Client(ctx, token), nil
}
