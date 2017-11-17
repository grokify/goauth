package aha

import (
	"fmt"
	"net/http"

	hum "github.com/grokify/gotilla/net/httputilmore"
	ou "github.com/grokify/oauth2util"
	"golang.org/x/oauth2"
)

const (
	AuthURLFormat    = "https://%s.aha.io/oauth/authorize"
	TokenURLFormat   = "https://%s.aha.io/oauth/token"
	AhaAccountHeader = "X-AHA-ACCOUNT"
)

func NewEndpoint(subdomain string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf(AuthURLFormat, subdomain),
		TokenURL: fmt.Sprintf(TokenURLFormat, subdomain)}
}

func NewClient(subdomain, token string) *http.Client {
	client := ou.NewClientAccessToken(token)

	header := http.Header{}
	header.Add(AhaAccountHeader, subdomain)

	client.Transport = hum.TransportWithHeaders{
		Transport: client.Transport,
		Header:    header}
	return client
}
