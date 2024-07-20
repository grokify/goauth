package jwtutil

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateJWTHS256SignedString(key []byte, claims map[string]any) (string, error) {
	claims2 := jwt.MapClaims(claims)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&claims2)
	return token.SignedString(key)
}
