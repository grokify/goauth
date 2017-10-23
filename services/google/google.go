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

func NewClientFromFile(ctx context.Context, filepath string, scopes []string, tok *oauth2.Token) (*http.Client, error) {
	conf, err := NewConfigFromFile(filepath, scopes)
	if err != nil {
		return &http.Client{}, errors.Wrap(err, fmt.Sprintf("Unable to read app config file: %v", filepath))
	}

	return conf.Client(ctx, tok), nil
}

func NewConfigFromFile(file string, scopes []string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(file) // Google client_secret.json
	if err != nil {
		return &oauth2.Config{},
			errors.Wrap(err, fmt.Sprintf("Unable to read client secret file: %v", err))
	}
	return o2g.ConfigFromJSON(b, scopes...)
}

func NewTokenFromWeb(cfg *oauth2.Config) (*oauth2.Token, error) {
	authURL := cfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to this link in your browser then type the auth code: \n%v\n", authURL)

	code := ""
	if _, err := fmt.Scan(&code); err != nil {
		return &oauth2.Token{}, errors.Wrap(err, "Unable to read auth code")
	}

	tok, err := cfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		return tok, errors.Wrap(err, "Unable to retrieve token from web")
	}
	return tok, nil
}
