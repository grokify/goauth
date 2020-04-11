package oauth2more

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func ParseJwtTokenString(tokenString string, secretKey string, claims *jwt.MapClaims) (*jwt.Token, jwt.MapClaims, error) {
	// https://stackoverflow.com/questions/41077953/go-language-and-verify-jwt
	if claims == nil {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			return nil, jwt.MapClaims{}, errors.Wrap(err, "ParseTokenString.jwt.Parse")
		}
		claims := token.Claims.(jwt.MapClaims)
		return token, claims, nil
	}
	// *jwt.StandardClaims
	// https://stackoverflow.com/questions/45405626/decoding-jwt-token-in-golang
	token, err := jwt.ParseWithClaims(tokenString, *claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, jwt.MapClaims{}, errors.Wrap(err, "ParseTokenString.jwt.ParseWithClaims")
	}
	tokenClaims := token.Claims.(jwt.MapClaims)
	return token, tokenClaims, nil
}
