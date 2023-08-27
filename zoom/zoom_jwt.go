package zoom

import (
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/grokify/goauth"
	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/endpoints"
	"github.com/grokify/mogo/net/http/httputilmore"
)

const (
	EnvZoomAPIKey           = "ZOOM_API_KEY" // #nosec G101
	EnvZoomAPISecret        = "ZOOM_API_SECRET"
	HeaderUserAgentJWTValue = "Zoom-api-Jwt-Request"
)

func CreateJwtToken(apiKey, apiSecret string, tokenDuration time.Duration) (*jwt.Token, string, error) {
	jwtCreds := goauth.CredentialsJWT{
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
	return authutil.NewClientHeaderQuery(
		map[string][]string{
			httputilmore.HeaderAuthorization: {authutil.TokenBearer + " " + bearerToken},
			httputilmore.HeaderUserAgent:     {HeaderUserAgentJWTValue}},
		map[string][]string{},
		false)
}
