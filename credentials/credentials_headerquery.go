package credentials

import (
	"net/http"
	"net/url"

	"github.com/grokify/goauth"
	"github.com/grokify/gohttp/httpsimple"
)

type CredentialsHeaderQuery struct {
	ServerURL     string      `json:"serverURL,omitempty"`
	Header        http.Header `json:"header,omitempty"`
	Query         url.Values  `json:"query,omitempty"`
	AllowInsecure bool        `json:"allowInsecure,omitempty"`
}

func (c *CredentialsHeaderQuery) NewClient() *http.Client {
	return goauth.NewClientHeaderQuery(c.Header, c.Query, c.AllowInsecure)
}

func (c *CredentialsHeaderQuery) NewSimpleClient() httpsimple.SimpleClient {
	return httpsimple.SimpleClient{
		HTTPClient: c.NewClient(),
		BaseURL:    c.ServerURL}
}