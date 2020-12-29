package zoom

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/grokify/oauth2more"
	"github.com/grokify/simplego/net/httputilmore"
)

const (
	EnvZoomApiKey           = "ZOOM_API_KEY"
	EnvZoomApiSecret        = "ZOOM_API_SECRET"
	HeaderUserAgentJwtValue = "Zoom-api-Jwt-Request"
)

func CreateJwtToken(apiKey, apiSecret string, tokenDur time.Duration) (*jwt.Token, string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.StandardClaims{
			Issuer:    apiKey,
			ExpiresAt: time.Now().Add(tokenDur).Unix()})
	tokenString, err := token.SignedString([]byte(apiSecret))
	return token, tokenString, err
}

func NewClientToken(bearerToken string) *http.Client {
	return oauth2more.NewClientHeaders(
		map[string]string{
			httputilmore.HeaderAuthorization: "Bearer " + bearerToken,
			httputilmore.HeaderUserAgent:     HeaderUserAgentJwtValue},
		false,
	)
}
