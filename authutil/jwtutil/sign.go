package jwtutil

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateHS256SignedString(key []byte, claims map[string]any) (string, error) {
	jmc := jwt.MapClaims(claims)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&jmc)
	return token.SignedString(key)
}
