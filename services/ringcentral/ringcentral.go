package ringcentral

import (
	"golang.org/x/oauth2"
)

// Endpoint is RingCentral's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://www.facebook.com/dialog/oauth",
	TokenURL: "https://graph.facebook.com/oauth/access_token",
}
