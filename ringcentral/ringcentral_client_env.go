package ringcentral

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/credentials"
)

func NewHttpClientEnvFlexStatic(envPrefix string) (*http.Client, error) {
	envPrefix = strings.TrimSpace(envPrefix)
	if len(envPrefix) == 0 {
		envPrefix = "RINGCENTRAL_"
	}

	envToken := strings.TrimSpace(envPrefix + "TOKEN")
	//token := config.JoinEnvNumbered(envToken, "", 2, true)
	token := os.Getenv(envToken)
	if len(token) > 0 {
		return goauth.NewClientAuthzTokenSimple(goauth.TokenBearer, token), nil
	}

	envPassword := strings.TrimSpace(envPrefix + "PASSWORD")
	password := os.Getenv(envPassword)
	if len(password) > 0 {
		return NewClientPassword(
			credentials.OAuth2Credentials{
				ClientID:     os.Getenv(envPrefix + "CLIENT_ID"),
				ClientSecret: os.Getenv(envPrefix + "CLIENT_SECRET"),
				ServerURL:    os.Getenv(envPrefix + "SERVER_URL"),
				Username:     os.Getenv(envPrefix + "USERNAME"),
				Password:     os.Getenv(envPrefix + "PASSWORD")})
	}

	return nil, fmt.Errorf("Cannot load client from ENV for prefix [%v]", envPassword)
}
