package zoom

import (
	"net/url"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/endpoints"
)

const (
	EnvZoomClientID      = "ZOOM_CLIENT_ID"
	EnvZoomCLientSecret  = "ZOOM_CLIENT_SECRET"
	EnvZoomApplicationID = "ZOOM_APPLICATION_ID"

	tokenBodyParamAccountID = "account_id"
)

// ServerToServerOAuth2Credentials implements Zoom's Server-to-Server OAuth 2.0 flow
// described here: https://developers.zoom.us/docs/internal-apps/s2s-oauth/ .
func ServerToServerOAuth2Credentials(clientID, clientSecret, accountID string) (goauth.Credentials, error) {
	creds := goauth.Credentials{
		Type:    goauth.TypeOAuth2,
		Service: endpoints.ServiceZoom,
		OAuth2: &goauth.CredentialsOAuth2{
			GrantType:     authutil.GrantTypeAccountCredentials,
			ClientID:      clientID,
			ClientSecret:  clientSecret,
			TokenBodyOpts: url.Values{},
		},
	}
	if accountID != "" {
		creds.OAuth2.TokenBodyOpts.Add(tokenBodyParamAccountID, accountID)
	}
	err := creds.Inflate()
	return creds, err
}
