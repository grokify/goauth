package google

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	o2g "golang.org/x/oauth2/google"
)

func ClientFromFile(ctx context.Context, filepath string, scopes []string, tok *oauth2.Token) (*http.Client, error) {
	conf, err := ConfigFromFile(filepath, scopes)
	if err != nil {
		return &http.Client{}, errors.Wrap(err, fmt.Sprintf("Unable to read app config file: %v", filepath))
	}

	return conf.Client(ctx, tok), nil
}

func ConfigFromFile(file string, scopes []string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(file) // Google client_secret.json
	if err != nil {
		return &oauth2.Config{},
			errors.Wrap(err, fmt.Sprintf("Unable to read client secret file: %v", err))
	}
	return o2g.ConfigFromJSON(b, scopes...)
}
