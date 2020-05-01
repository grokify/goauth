package oauth2more

import (
	"encoding/json"

	"golang.org/x/oauth2"
)

// ParseToken parses a JSON token and returns an
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
		return nil, err
	}
	return tok.WithExtra(msi), nil
}
