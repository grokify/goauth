package main

import (
	"errors"
	"fmt"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
)

func main() {
	str, err := getJWT()
	logutil.FatalErr(err)

	clm, err := ParseJWTWithoutVerify(str)
	logutil.FatalErr(err)
	if err := fmtutil.PrintJSON(clm); err != nil {
		logutil.FatalErr(err)
	}

	fmt.Println("DONE")
}

const EnvJWTParse = "JWT_PARSE"

func getJWT() (string, error) {
	if tok := os.Getenv(EnvJWTParse); tok != "" {
		return tok, nil
	}
	return "", errors.New("env var " + EnvJWTParse + " is not set")
}

func ParseJWTWithoutVerify(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid JWT claims format")
}
