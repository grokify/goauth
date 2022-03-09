package credentials

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	SigningMethodES256 = "ES256"
	SigningMethodES384 = "ES384"
	SigningMethodES512 = "ES512"
	SigningMethodHS256 = "HS256"
	SigningMethodHS384 = "HS384"
	SigningMethodHS512 = "HS512"
)

type CredentialsJWT struct {
	Issuer        string `json:"issuer,omitempty"`
	PrivateKey    string `json:"privateKey,omitempty"`
	SigningMethod string `json:"signingMethod,omitempty"`
}

func (jc *CredentialsJWT) StandardToken(tokenDuration time.Duration) (*jwt.Token, string, error) {
	stdClaims := jwt.RegisteredClaims{}
	if len(jc.Issuer) > 0 {
		stdClaims.Issuer = jc.Issuer
	}
	if tokenDuration > 0 {
		stdClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(tokenDuration))
	}
	token := jwt.NewWithClaims(
		jwt.GetSigningMethod(strings.ToUpper(strings.TrimSpace(jc.SigningMethod))),
		stdClaims)
	tokenString, err := token.SignedString([]byte(jc.PrivateKey))
	return token, tokenString, err
}
