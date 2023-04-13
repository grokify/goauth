package authutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/time/timeutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func ParseTokenReader(r io.Reader) (*oauth2.Token, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseToken(data)
}

// ParseToken parses a OAuth 2 token and returns an `*oauth2.Token` with custom properties.
func ParseToken(rawToken []byte) (*oauth2.Token, error) {
	tok := &oauth2.Token{}
	err := json.Unmarshal(rawToken, tok)
	if err != nil {
		return tok, err
	}
	msi := map[string]any{}
	err = json.Unmarshal(rawToken, &msi)
	if err != nil {
		return tok, err
	}
	tok = tok.WithExtra(msi)
	// convert `expires_in` to `Expiry` with 1 minute leeway.
	if timeutil.NewTimeMore(tok.Expiry, 0).IsZeroAny() {
		expiresIn := tok.Extra(OAuth2TokenPropExpiresIn)
		if expiresIn != nil {
			if expiresInFloat, ok := expiresIn.(float64); ok {
				if expiresInFloat > 60 { // subtract 1 minute for defensive handling
					expiresInFloat -= 60
				}
				if expiresInFloat > 0 {
					tok.Expiry = time.Now().UTC().Add(time.Duration(expiresInFloat) * time.Second)
				}
			}
		}
	}
	return tok, nil
}

// NewTokenCLIFromWeb enables a CLI app with no UI to generate
// a OAuth2 AuthURL which is copy and pasted into a web browser to
// return an an OAuth 2 authorization code and state, where the
// authorization code is entered on the command line.
func NewTokenCLIFromWeb(cfg *oauth2.Config, state string) (*oauth2.Token, error) {
	//authURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	authURL := cfg.AuthCodeURL(state)
	fmt.Printf("Go to this link in your browser then type in the auth code from the webpage and click `return` to continue: \n%v\n", authURL)

	code := ""
	if _, err := fmt.Scan(&code); err != nil {
		return nil, errorsutil.Wrap(err, "Unable to read auth code")
	}

	tok, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return tok, errorsutil.Wrap(err, "Unable to retrieve token from web")
	}
	return tok, nil
}

// TokenClientCredentials is an alternative to `clientcredentials.Config.Token()`
// which does not work for some APIs. More investigation is needed but it appears
// the issue is encoding the HTTP request body. The approach here uses `&` in the
// URL encoded values.
func TokenClientCredentials(cfg clientcredentials.Config) (*oauth2.Token, error) {
	body := url.Values{}
	body.Add(ParamGrantType, GrantTypeClientCredentials)
	for _, scope := range cfg.Scopes {
		body.Add(ParamScope, scope)
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
	req.Header.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppFormURLEncoded)
	req.Header.Add(httputilmore.HeaderAuthorization, TokenBasic+" "+b64)

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

/*
// OAuth2Token is a bridge struct to `oauth2.Token` since the RFC-6749 uses `expires_in` and
// golang uses `expiry`.
type OAuth2Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (ot OAuth2Token) Token() *oauth2.Token {
	expiresIn := ot.ExpiresIn
	if expiresIn > 60 { // subtract 1 minute for defensive handling
		expiresIn -= 60
	}
	return &oauth2.Token{
		AccessToken:  ot.AccessToken,
		RefreshToken: ot.RefreshToken,
		TokenType:    ot.TokenType,
		Expiry:       time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
	}
}

func ParseOAuth2Token(rawToken []byte) (*OAuth2Token, error) {
	oTok := &OAuth2Token{}
	err := json.Unmarshal(rawToken, oTok)
	return oTok, err
}
*/

/*
   "access_token":"2YotnFZFEjr1zCsicMWpAA",
   "token_type":"example",
   "expires_in":3600,
   "refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
*/
