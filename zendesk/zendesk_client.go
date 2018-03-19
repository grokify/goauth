package zendesk

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"

	om "github.com/grokify/oauth2more"
)

var (
	EnvZendeskUsername  = "ZENDESK_USERNAME"
	EnvZendeskPassword  = "ZENDESK_PASSWORD"
	EnvZendeskSubdomain = "ZENDESK_SUBDOMAIN"
)

func NewClientPassword(ctx context.Context, username, password string) (*http.Client, error) {
	token, err := om.BasicAuthToken(username, password)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{}
	return conf.Client(ctx, token), nil
}
