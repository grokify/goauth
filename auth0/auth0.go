// auth0 contains a Go implementation of Auth0's PKCE support:
// https://auth0.com/docs/api-auth/tutorials/authorization-code-grant-pkce
package auth0

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/grokify/mogo/net/http/httputilmore"
)

func CreatePKCECodeVerifier() (string, error) {
	verifier := make([]byte, 32)
	_, err := rand.Read(verifier)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(verifier[:]), nil
}

func CreatePKCEChallengeS256(verifier string) string {
	challenge := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(challenge[:])
}

type PKCEAuthorizationURLInfo struct {
	Host                string `url:"-"`
	Audience            string `url:"audience"`
	Scope               string `url:"scope"`
	ResponseType        string `url:"response_type"`
	ClientID            string `url:"client_id"`
	CodeChallenge       string `url:"code_challenge"`
	CodeChallengeMethod string `url:"code_challenge_method"`
	RedirectURI         string `url:"redirect_uri"`
}

func (au *PKCEAuthorizationURLInfo) url() (string, error) {
	baseURL := fmt.Sprintf("https://%s/authorize", au.Host)
	au.ResponseType = "code"
	au.CodeChallengeMethod = "S256"
	v, err := query.Values(au)
	if err != nil {
		return baseURL, err
	}
	return baseURL + "?" + v.Encode(), nil
}

func (au *PKCEAuthorizationURLInfo) Data() (string, string, string, error) {
	verifier, err := CreatePKCECodeVerifier()
	if err != nil {
		return "", "", "", err
	}
	challenge := CreatePKCEChallengeS256(verifier)
	au.CodeChallenge = challenge
	myURL, err := au.url()
	return verifier, challenge, myURL, err
}

type PKCETokenURLInfo struct {
	Host         string `json:"-"`
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	CodeVerifier string `json:"code_verifier"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

func (tu *PKCETokenURLInfo) URL() string {
	return fmt.Sprintf("https://%s/oauth/token", tu.Host)
}

func (tu *PKCETokenURLInfo) Body() ([]byte, error) {
	return json.Marshal(tu)
}

func (tu *PKCETokenURLInfo) Exchange() (*http.Response, error) {
	body, err := tu.Body()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, tu.URL(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJSONUtf8)
	client := &http.Client{}
	return client.Do(req)
}
