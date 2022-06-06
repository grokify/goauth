package credentials

import (
	"net/http"
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/net/httputilmore"
)

type CredentialsBasicAuth struct {
	Username      string
	Password      string
	Encoded       string
	BaseURL       string
	AllowInsecure bool
}

func (c *CredentialsBasicAuth) NewClient() (*http.Client, error) {
	if len(c.Encoded) > 0 {
		if strings.Index(strings.ToLower(strings.TrimSpace(c.Encoded)), "basic ") == 0 {
			hdr := http.Header{}
			hdr.Add(httputilmore.HeaderAuthorization, c.Encoded)
			return goauth.NewClientHeaders(hdr, c.AllowInsecure), nil
		}
		return goauth.NewClientToken(goauth.TokenBasic, c.Encoded, c.AllowInsecure), nil
	} else if len(c.Username) > 0 || len(c.Password) > 0 {
		return goauth.NewClientBasicAuth(c.Username, c.Password, c.AllowInsecure)
	}
	return &http.Client{}, nil
}
