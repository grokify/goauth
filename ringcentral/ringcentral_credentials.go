package ringcentral

import (
	"github.com/grokify/goauth"
	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/endpoints"
)

// CredentialsJWTBearer implements RingCentral's JWT Bearer Flow: https://developers.ringcentral.com/guide/authentication/jwt/quick-start
func CredentialsJWTBearer(clientID, clientSecret, jwt string, sandbox bool) goauth.Credentials {
	creds := goauth.Credentials{
		Type:    goauth.TypeOAuth2,
		Service: endpoints.ServiceRingcentral,
		OAuth2: &goauth.CredentialsOAuth2{
			GrantType:    authutil.GrantTypeJWTBearer,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			JWT:          jwt,
		},
	}
	if sandbox {
		creds.Service = endpoints.ServiceRingcentralSandbox
	}
	if err := creds.Inflate(); err != nil {
		panic(err)
	} else {
		return creds
	}
}
