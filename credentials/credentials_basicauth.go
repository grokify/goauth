package credentials

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/gohttp/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
)

type CredentialsBasicAuth struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Encoded       string `json:"encoded,omitempty"`
	ServerURL     string `json:"serverURL,omitempty"`
	AllowInsecure bool   `json:"allowInsecure,omitempty"`
}

func (c *CredentialsBasicAuth) NewClient() (*http.Client, error) {
	if len(strings.TrimSpace(c.Encoded)) > 0 {
		if strings.Index(strings.ToLower(strings.TrimSpace(c.Encoded)), "basic ") == 0 {
			return goauth.NewClientHeaderQuery(
				http.Header{httputilmore.HeaderAuthorization: []string{c.Encoded}},
				url.Values{},
				c.AllowInsecure), nil
		}
		return goauth.NewClientToken(goauth.TokenBasic, c.Encoded, c.AllowInsecure), nil
	} else if len(c.Username) > 0 || len(c.Password) > 0 {
		return goauth.NewClientBasicAuth(c.Username, c.Password, c.AllowInsecure)
	}
	return &http.Client{}, nil
}

func (c *CredentialsBasicAuth) NewSimpleClient() (httpsimple.SimpleClient, error) {
	hclient, err := c.NewClient()
	if err != nil {
		return httpsimple.SimpleClient{}, err
	}
	return httpsimple.SimpleClient{
		HTTPClient: hclient,
		BaseURL:    c.ServerURL}, nil
}
