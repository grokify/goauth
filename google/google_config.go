package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/stringsutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func ConfigFromFile(file string, scopes []string) (*oauth2.Config, error) {
	b, err := os.ReadFile(file) // Google client_secret.json
	if err != nil {
		return &oauth2.Config{},
			errorsutil.Wrap(err, fmt.Sprintf("unable to read client secret file [%v]", err))
	}
	return google.ConfigFromJSON(b, scopes...)
}

func ConfigFromEnv(envVar string, scopes []string) (*oauth2.Config, error) {
	envVar = strings.TrimSpace(envVar)
	if len(envVar) == 0 {
		envVar = EnvGoogleAppCredentials
	}
	if len(scopes) == 0 {
		scopesString := os.Getenv(EnvGoogleAppScopes)
		scopes = stringsutil.SplitCondenseSpace(scopesString, ",")
	}
	return google.ConfigFromJSON([]byte(os.Getenv(envVar)), scopes...)
}

// ConfigFromBytes returns an *oauth2.Config given a byte array
// containing the Google client_secret.json data.
func ConfigFromBytes(configJSON []byte, scopes []string) (*oauth2.Config, error) {
	if len(strings.TrimSpace(string(configJSON))) == 0 {
		return nil, errorsutil.Wrap(errors.New("no credentials provided"), "goauth/google.ConfigFromBytes()")
	}

	if len(scopes) == 0 {
		cc := CredentialsContainer{}
		err := json.Unmarshal(configJSON, &cc)
		if err != nil {
			return nil, errorsutil.Wrap(err, "func ConfigFromBytes")
		}
		if len(cc.Scopes) > 0 {
			scopes = append(scopes, cc.Scopes...)
		}
	}

	return google.ConfigFromJSON(configJSON, scopes...)
}
