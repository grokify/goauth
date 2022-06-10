package credentials

import (
	"net/http"
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/net/httputilmore"
)

type CredentialsBasicAuth struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Encoded       string `json:"encoded,omitempty"`
	ServerURL     string `json:"serverURL,omitempty"`
	AllowInsecure bool   `json:"allowInsecure,omitempty"`
}

func (c *CredentialsBasicAuth) NewClient() (*http.Client, error) {
	if len(c.Encoded) > 0 {
		if strings.Index(strings.ToLower(strings.TrimSpace(c.Encoded)), "basic ") == 0 {
			hdr := http.Header{}
			hdr.Add(httputilmore.HeaderAuthorization, c.Encoded)
			return goauth.NewClientHeaderQuery(hdr, map[string][]string{}, c.AllowInsecure), nil
		}
		return goauth.NewClientToken(goauth.TokenBasic, c.Encoded, c.AllowInsecure), nil
	} else if len(c.Username) > 0 || len(c.Password) > 0 {
		return goauth.NewClientBasicAuth(c.Username, c.Password, c.AllowInsecure)
	}
	return &http.Client{}, nil
}