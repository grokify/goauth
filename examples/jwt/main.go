package main

import (
	"fmt"
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/grokify/goauth/authutil/jwtutil"
)

// Answer for https://stackoverflow.com/a/61284284/1908967

func createJWTHS256SignedString(secretKey string, data map[string]any) (string, error) {
	claims := map[string]any{
		"iss":  "issuer",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"data": data,
	}
	return jwtutil.CreateHS256SignedString([]byte(secretKey), claims)
}

func parseJWTSubClaimName(tokenString, secretKey, field string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)

	data := claims["data"].(map[string]any)
	return data[field].(string), nil
}

func main() {
	secretKey := "foobar"

	tokenString, err := createJWTHS256SignedString(
		secretKey,
		map[string]any{
			"id": "123", "name": "JohnDoe"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tokenString)

	key := "name"
	value, err := parseJWTSubClaimName(tokenString, secretKey, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("KEY [%s] VAL [%s]\n", key, value)

	fmt.Println("DONE")
}
