package oauth2more

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// ParseToken parses a OAuth 2 token and returns an
// `*oauth2.Token` with custom properties.
func ParseToken(rawToken []byte) (*oauth2.Token, error) {
	tok := &oauth2.Token{}
	err := json.Unmarshal([]byte(rawToken), tok)
	if err != nil {
		return tok, err
	}
	msi := map[string]interface{}{}
	err = json.Unmarshal(rawToken, &msi)
	if err != nil {
		return tok, err
	}
	return tok.WithExtra(msi), nil
}

// NewTokenCliFromWeb enables a CLI app with no UI to generate
// a OAuth2 AuthURL which is copy and pasted into a web browser to
// return an an OAuth 2 authorization code and state, where the
// authorization code is entered on the command line.
func NewTokenCliFromWeb(cfg *oauth2.Config, state string) (*oauth2.Token, error) {
	authURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Go to this link in your browser then type in the auth code from the webpage and click `return` to continue: \n%v\n", authURL)

	code := ""
	if _, err := fmt.Scan(&code); err != nil {
		return nil, errors.Wrap(err, "Unable to read auth code")
	}

	tok, err := cfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		return tok, errors.Wrap(err, "Unable to retrieve token from web")
	}
	return tok, nil
}
