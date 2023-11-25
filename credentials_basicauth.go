package goauth

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
)

type CredentialsBasicAuth struct {
	Username      string            `json:"username,omitempty"`
	Password      string            `json:"password,omitempty"`
	Encoded       string            `json:"encoded,omitempty"`
	ServerURL     string            `json:"serverURL,omitempty"`
	AllowInsecure bool              `json:"allowInsecure,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

func (c *CredentialsBasicAuth) NewClient() (*http.Client, error) {
	if len(strings.TrimSpace(c.Encoded)) > 0 {
		if strings.Index(strings.ToLower(strings.TrimSpace(c.Encoded)), "basic ") == 0 {
			return authutil.NewClientHeaderQuery(
				http.Header{httputilmore.HeaderAuthorization: []string{c.Encoded}},
				url.Values{},
				c.AllowInsecure), nil
		}
		return authutil.NewClientToken(authutil.TokenBasic, c.Encoded, c.AllowInsecure), nil
	} else if len(c.Username) > 0 || len(c.Password) > 0 {
		return authutil.NewClientBasicAuth(c.Username, c.Password, c.AllowInsecure)
	}
	return &http.Client{}, nil
}

func (c *CredentialsBasicAuth) NewSimpleClient() (httpsimple.Client, error) {
	hclient, err := c.NewClient()
	if err != nil {
		return httpsimple.Client{}, err
	}
	return httpsimple.Client{
		HTTPClient: hclient,
		BaseURL:    c.ServerURL}, nil
}
