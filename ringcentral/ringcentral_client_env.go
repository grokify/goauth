package ringcentral

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/authutil"
)

func NewHTTPClientEnvFlexStatic(envPrefix string) (*http.Client, error) {
	envPrefix = strings.TrimSpace(envPrefix)
	if len(envPrefix) == 0 {
		envPrefix = "RINGCENTRAL_"
	}

	envToken := strings.TrimSpace(envPrefix + "TOKEN")
	//token := config.JoinEnvNumbered(envToken, "", 2, true)
	token := os.Getenv(envToken)
	if len(token) > 0 {
		return authutil.NewClientAuthzTokenSimple(authutil.TokenBearer, token), nil
	}

	envPassword := strings.TrimSpace(envPrefix + "PASSWORD")
	password := os.Getenv(envPassword)
	if len(password) > 0 {
		return NewClientPassword(goauth.NewCredentialsOAuth2Env(envPrefix))
	}

	return nil, fmt.Errorf("cannot load client from ENV for prefix [%v]", envPassword)
}
