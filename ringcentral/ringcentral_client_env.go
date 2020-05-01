package ringcentral

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	om "github.com/grokify/oauth2more"
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
		return om.NewClientBearerTokenSimple(token), nil
	}

	envPassword := strings.TrimSpace(envPrefix + "PASSWORD")
	password := os.Getenv(envPassword)
	if len(password) > 0 {
		return NewClientPassword(
			ApplicationCredentials{
				ClientID:     os.Getenv(envPrefix + "CLIENT_ID"),
				ClientSecret: os.Getenv(envPrefix + "CLIENT_SECRET"),
				ServerURL:    os.Getenv(envPrefix + "SERVER_URL")},
			PasswordCredentials{
				Username:  os.Getenv(envPrefix + "USERNAME"),
				Extension: os.Getenv(envPrefix + "EXTENSION"),
				Password:  os.Getenv(envPrefix + "PASSWORD")})
	}

	return nil, fmt.Errorf("Cannot load client from ENV for prefix [%v]", envPassword)
}
