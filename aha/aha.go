package aha

import (
	"fmt"

	"golang.org/x/oauth2"
)

const (
	AuthURLFormat  = "https://%s.aha.io/oauth/authorize"
	TokenURLFormat = "https://%s.aha.io/oauth/token"
)

func NewEndpoint(subdomain string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf(AuthURLFormat, subdomain),
		TokenURL: fmt.Sprintf(TokenURLFormat, subdomain)}
}
