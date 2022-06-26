package credentials

import (
	"net/http"
	"net/url"

	"github.com/grokify/goauth"
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
