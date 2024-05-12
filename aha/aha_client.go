package aha

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/endpoints"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/urlutil"
	"golang.org/x/oauth2"
)

const (
	APIMeURL         = "https://secure.aha.io/api/v1/me" // no longer supported by Aha. Must use sub-domain
	BaseURLFormat    = "https://%s.aha.io/"
	APIMeURLFormat   = "https://%s.aha.io/api/v1/me"
	APIMeURLPath     = "/api/v1/me"
	AuthURLFormat    = "https://%s.aha.io/oauth/authorize"
	TokenURLFormat   = "https://%s.aha.io/oauth/token" // #nosec G101
	AhaAccountHeader = "X-AHA-ACCOUNT"
)

var (
	AhaAccountEnv = "AHA_ACCOUNT"
	AhaServerURL  = "AHA_SERVER_URL"
	AhaAPIKeyEnv  = "AHA_API_KEY" // #nosec G101
)

func NewEndpoint(subdomain string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf(AuthURLFormat, subdomain),
		TokenURL: fmt.Sprintf(TokenURLFormat, subdomain)}
}

func NewClient(subdomain, token string) (*http.Client, error) {
	creds := BuildCredentials(subdomain, token)
	client, _, err := creds.OAuth2.NewClient(context.Background())
	return client, err
	/*
		return authutil.NewClientHeaderQuery(
			http.Header{
				httputilmore.HeaderAuthorization: []string{authutil.TokenBearer + " " + token},
				AhaAccountHeader:                 []string{subdomain}},
			url.Values{},
			false)
	*/
}

func NewSimpleClient(subdomain, token string) (*httpsimple.Client, error) {
	creds := BuildCredentials(subdomain, token)
	return creds.OAuth2.NewSimpleClient(context.Background())
	/*
		return authutil.NewClientHeaderQuery(
			http.Header{
				httputilmore.HeaderAuthorization: []string{authutil.TokenBearer + " " + token},
				AhaAccountHeader:                 []string{subdomain}},
			url.Values{},
			false)
	*/
}

func BuildCredentials(subdomainOrBaseURL, token string) goauth.Credentials {
	if !urlutil.IsHTTP(subdomainOrBaseURL, true, true) {
		subdomainOrBaseURL = fmt.Sprintf(BaseURLFormat, subdomainOrBaseURL)
	}
	return goauth.Credentials{
		Type:    goauth.TypeOAuth2,
		Service: endpoints.ServiceAha,
		OAuth2: &goauth.CredentialsOAuth2{
			ServerURL: subdomainOrBaseURL,
			Token: &oauth2.Token{
				AccessToken: token,
			},
		},
	}
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
