package goauth

import (
	"net/http"
	"net/url"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/net/http/httpsimple"
)

type CredentialsHeaderQuery struct {
	ServerURL     string      `json:"serverURL,omitempty"`
	Header        http.Header `json:"header,omitempty"`
	Query         url.Values  `json:"query,omitempty"`
	AllowInsecure bool        `json:"allowInsecure,omitempty"`
}

func (c *CredentialsHeaderQuery) NewClient() *http.Client {
	return authutil.NewClientHeaderQuery(c.Header, c.Query, c.AllowInsecure)
}

func (c *CredentialsHeaderQuery) NewSimpleClient() httpsimple.SimpleClient {
	return httpsimple.SimpleClient{
		HTTPClient: c.NewClient(),
		BaseURL:    c.ServerURL}
}
