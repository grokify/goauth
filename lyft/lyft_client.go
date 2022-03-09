package lyft

import (
	"context"
	"net/http"

	"github.com/grokify/goauth/endpoints"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	// OAuth 2.0 Scopes
	Offline      = "offline"
	Profile      = "profile"
	Public       = "public"
	RidesRead    = "rides.read"
	RidesRequest = "rides.request"
)

func Endpoint() oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  endpoints.LyftAuthzURL,
		TokenURL: endpoints.LyftTokenURL}
}

func NewClientCredentials(ctx context.Context, clientID, clientSecret string, scopes []string) *http.Client {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     endpoints.LyftTokenURL,
		Scopes:       scopes}

	return config.Client(ctx)
}
