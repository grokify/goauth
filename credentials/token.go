package credentials

import (
	"fmt"
	"strings"
	"time"

	"github.com/grokify/oauth2more"
	"golang.org/x/oauth2"
)

func NewTokenCli(creds Credentials, state string) (token *oauth2.Token, err error) {
	if creds.Application.IsGrantType(oauth2more.GrantTypeAuthorizationCode) {
		state = strings.TrimSpace(state)
		if len(state) == 0 {
			state = "oauth2more-" + time.Now().UTC().Format(time.RFC3339)
		}
		fmt.Printf("OAuth State [%s]\n", state)
		cfg := creds.Application.Config()
		token, err = oauth2more.NewTokenCliFromWeb(&cfg, state)
		if err != nil {
			return token, err
		}
	} else {
		token, err = creds.NewToken()
		if err != nil {
			return token, err
		}
	}
	token.Expiry = token.Expiry.UTC()
	return token, nil
}
