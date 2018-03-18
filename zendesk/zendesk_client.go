package zendesk

import (
	"context"
	"net/http"

	"github.com/grokify/oauth2more"
	"golang.org/x/oauth2"
)

var (
	EnvZendeskUsername  = "ZENDESK_USERNAME"
	EnvZendeskPassword  = "ZENDESK_PASSWORD"
	EnvZendeskSubdomain = "ZENDESK_SUBDOMAIN"
)

func NewClient(ctx context.Context, subdomain, username, password string) (*http.Client, error) {
	token, err := oauth2more.BasicAuthToken(username, password)
	if err != nil {
		return nil, err
	}

	cfg := oauth2.Config{}
	return cfg.Client(ctx, token), nil
}
