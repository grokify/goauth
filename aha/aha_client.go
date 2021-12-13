package aha

import (
	"fmt"
	"net/http"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/net/httputilmore"
	"golang.org/x/oauth2"
)

const (
	APIMeURL         = "https://secure.aha.io/api/v1/me"
	AuthURLFormat    = "https://%s.aha.io/oauth/authorize"
	TokenURLFormat   = "https://%s.aha.io/oauth/token"
	AhaAccountHeader = "X-AHA-ACCOUNT"
)

var (
	AhaAccountEnv = "AHA_ACCOUNT"
	AhaApiKeyEnv  = "AHA_API_KEY"
)

func NewEndpoint(subdomain string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf(AuthURLFormat, subdomain),
		TokenURL: fmt.Sprintf(TokenURLFormat, subdomain)}
}

func NewClient(subdomain, token string) *http.Client {
	client := goauth.NewClientAuthzTokenSimple(goauth.TokenBearer, token)

	header := http.Header{}
	header.Add(AhaAccountHeader, subdomain)

	client.Transport = httputilmore.TransportWithHeaders{
		Transport: client.Transport,
		Header:    header}
	return client
}
