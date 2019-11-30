package google

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	o2g "golang.org/x/oauth2/google"
)

var (
	ClientSecretEnv         = "GOOGLE_APP_CLIENT_SECRET"
	EnvGoogleAppCredentials = "GOOGLE_APP_CREDENTIALS"
)

func ConfigFromFile(file string, scopes []string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(file) // Google client_secret.json
	if err != nil {
		return &oauth2.Config{},
			errors.Wrap(err, fmt.Sprintf("Unable to read client secret file: %v", err))
	}
	return o2g.ConfigFromJSON(b, scopes...)
}

func ConfigFromEnv(envVar string, scopes []string) (*oauth2.Config, error) {
	return o2g.ConfigFromJSON([]byte(os.Getenv(envVar)), scopes...)
}

// ConfigFromBytes returns an *oauth2.Config given a byte array
// containing the Google client_secret.json data.
func ConfigFromBytes(configJson []byte, scopes []string) (*oauth2.Config, error) {
	if len(strings.TrimSpace(string(configJson))) == 0 {
		return nil, errors.Wrap(errors.New("No Credentials Provided"), "oauth2more/google.ConfigFromBytes()")
	}

	return o2g.ConfigFromJSON(configJson, scopes...)
}
