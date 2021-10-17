package zoom

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/grokify/goauth"
	"github.com/grokify/goauth/credentials"
	"github.com/grokify/goauth/endpoints"
	"github.com/grokify/simplego/net/httputilmore"
)

const (
	EnvZoomApiKey           = "ZOOM_API_KEY"
	EnvZoomApiSecret        = "ZOOM_API_SECRET"
	HeaderUserAgentJwtValue = "Zoom-api-Jwt-Request"
)

func CreateJwtToken(apiKey, apiSecret string, tokenDuration time.Duration) (*jwt.Token, string, error) {
	jwtCreds := credentials.JWTCredentials{
		Issuer:        apiKey,
		PrivateKey:    apiSecret,
		SigningMethod: endpoints.ZoomJWTSigningMethod}
	return jwtCreds.StandardToken(tokenDuration)
}

func NewClient(apiKey, apiSecret string, tokenDuration time.Duration) (*http.Client, error) {
	_, jwtString, err := CreateJwtToken(apiKey, apiSecret, tokenDuration)
	if err != nil {
		return nil, err
	}
	return NewClientToken(jwtString), nil
}

func NewClientToken(bearerToken string) *http.Client {
	return goauth.NewClientHeaders(
		map[string][]string{
			httputilmore.HeaderAuthorization: []string{"Bearer " + bearerToken},
			httputilmore.HeaderUserAgent:     []string{HeaderUserAgentJwtValue}},
		false)
}
