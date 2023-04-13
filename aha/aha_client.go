package aha

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"golang.org/x/oauth2"
)

const (
	APIMeURL         = "https://secure.aha.io/api/v1/me"
	AuthURLFormat    = "https://%s.aha.io/oauth/authorize"
	TokenURLFormat   = "https://%s.aha.io/oauth/token" // #nosec G101
	AhaAccountHeader = "X-AHA-ACCOUNT"
)

var (
	AhaAccountEnv = "AHA_ACCOUNT"
	AhaAPIKeyEnv  = "AHA_API_KEY" // #nosec G101
)

func NewEndpoint(subdomain string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf(AuthURLFormat, subdomain),
		TokenURL: fmt.Sprintf(TokenURLFormat, subdomain)}
}

func NewClient(subdomain, token string) *http.Client {
	return authutil.NewClientHeaderQuery(
		http.Header{
			httputilmore.HeaderAuthorization: []string{authutil.TokenBearer + " " + token},
			AhaAccountHeader:                 []string{subdomain}},
		url.Values{},
		false)
}

/*
func NewClient(subdomain, token string) *http.Client {
	client := authutil.NewClientAuthzTokenSimple(authutil.TokenBearer, token)

	header := http.Header{}
	header.Add(AhaAccountHeader, subdomain)

	client.Transport = httputilmore.TransportRequestModifier{
		Transport: client.Transport,
		Header:    header}
	return client
}
*/
