package hubspot

import (
	"net/http"
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/endpoints"
	"github.com/grokify/gohttp/httpsimple"
)

const (
	APIKeyQueryParameter = "hapikey"
)

func NewClientAPIKey(apiKey string) *http.Client {
	apiKey = strings.TrimSpace(apiKey)

	return goauth.NewClientHeadersQuery(
		http.Header{},
		map[string][]string{APIKeyQueryParameter: []string{apiKey}},
		false)
}

func NewSimpleClientAPIKey(apiKey string) httpsimple.SimpleClient {
	return httpsimple.SimpleClient{
		BaseURL:    endpoints.HubspotServerURL,
		HTTPClient: NewClientAPIKey(apiKey)}
}
