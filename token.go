package oauth2more

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/simplego/net/httputilmore"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
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
	//authURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	authURL := cfg.AuthCodeURL(state)
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

// TokenClientCredentials is an alternative to `clientcredentials.Config.Token()`
// which does not work for some APIs. More investigation is needed but it appears
// the issue is encoding the HTTP request body. The approach here uses `&` in the
// URL encoded values.
func TokenClientCredentials(cfg clientcredentials.Config) (*oauth2.Token, error) {
	body := url.Values{}
	body.Add("grant_type", GrantTypeClientCredentials)
	for _, scope := range cfg.Scopes {
		body.Add("scope", scope)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		cfg.TokenURL,
		strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	b64, err := RFC7617UserPass(cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		return nil, err
	}
	req.Header.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppFormUrlEncoded)
	req.Header.Add(httputilmore.HeaderAuthorization, "Basic "+b64)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	return tok, json.Unmarshal(data, tok)
}
