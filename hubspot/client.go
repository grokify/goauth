package hubspot

import (
	"net/http"
	"strings"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/endpoints"
	"github.com/grokify/mogo/net/http/httpsimple"
)

const (
	APIKeyQueryParameter = "hapikey"
)

func NewClientAPIKey(apiKey string) *http.Client {
	return authutil.NewClientHeaderQuery(
		http.Header{},
		map[string][]string{APIKeyQueryParameter: {strings.TrimSpace(apiKey)}},
		false)
}

func NewSimpleClientAPIKey(apiKey string) httpsimple.Client {
	return httpsimple.Client{
		BaseURL:    endpoints.HubspotServerURL,
		HTTPClient: NewClientAPIKey(apiKey)}
}
