package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Answer for https://stackoverflow.com/a/61284284/1908967

func createJWT(secretKey string, data map[string]interface{}) (string, error) {
	claims := &jwt.MapClaims{
		"iss": "issuer",
		"exp": time.Now().Add(time.Hour).Unix(),
		"data": map[string]string{
			"id":   "123",
			"name": "JohnDoe",
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims)
	return token.SignedString([]byte(secretKey))
}

func parseJWTSubClaimName(tokenString, secretKey, field string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)

	data := claims["data"].(map[string]interface{})
	return data[field].(string), nil
}

func main() {
	secretKey := "foobar"

	tokenString, err := createJWT(
		secretKey,
		map[string]interface{}{
			"id": "123", "name": "JohnDoe"})

	fmt.Println(tokenString)

	key := "name"
	value, err := parseJWTSubClaimName(tokenString, secretKey, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("KEY [%s] VAL [%s]\n", key, value)

	fmt.Println("DONE")
}
