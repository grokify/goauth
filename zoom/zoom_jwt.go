package zoom

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/credentials"
	"github.com/grokify/oauth2more/endpoints"
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
	/*
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.StandardClaims{
				Issuer:    apiKey,
				ExpiresAt: time.Now().Add(tokenDuration).Unix()})
		tokenString, err := token.SignedString([]byte(apiSecret))
		return token, tokenString, err
	*/
}

func NewClient(apiKey, apiSecret string, tokenDuration time.Duration) (*http.Client, error) {
	_, jwtString, err := CreateJwtToken(apiKey, apiSecret, tokenDuration)
	if err != nil {
		return nil, err
	}
	return NewClientToken(jwtString), nil
}

func NewClientToken(bearerToken string) *http.Client {
	return oauth2more.NewClientHeaders(
		map[string][]string{
			httputilmore.HeaderAuthorization: []string{"Bearer " + bearerToken},
			httputilmore.HeaderUserAgent:     []string{HeaderUserAgentJwtValue}},
		false)
}
