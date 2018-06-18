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
	hum "github.com/grokify/gotilla/net/httputilmore"
)

func CreatePKCECodeVerifier() string {
	verifier := make([]byte, 32)
	rand.Read(verifier)
	return base64.RawURLEncoding.EncodeToString(verifier[:])
}

func CreatePKCEChallengeS256(verifier string) string {
	challenge := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(challenge[:])
}

type PKCEAuthorizationUrlInfo struct {
	Host                string `url:"-"`
	Audience            string `url:"audience"`
	Scope               string `url:"scope"`
	ResponseType        string `url:"response_type"`
	ClientId            string `url:"client_id"`
	CodeChallenge       string `url:"code_challenge"`
	CodeChallengeMethod string `url:"code_challenge_method"`
	RedirectUri         string `url:"redirect_uri"`
}

func (au *PKCEAuthorizationUrlInfo) url() (string, error) {
	baseUrl := fmt.Sprintf("https://%s/authorize", au.Host)
	au.ResponseType = "code"
	au.CodeChallengeMethod = "S256"
	v, err := query.Values(au)
	if err != nil {
		return baseUrl, err
	}
	return baseUrl + "?" + v.Encode(), nil
}

func (au *PKCEAuthorizationUrlInfo) Data() (string, string, string, error) {
	verifier := CreatePKCECodeVerifier()
	challenge := CreatePKCEChallengeS256(verifier)
	au.CodeChallenge = challenge
	myUrl, err := au.url()
	return verifier, challenge, myUrl, err
}

type PKCETokenUrlInfo struct {
	Host         string `json:"-"`
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	CodeVerifier string `json:"code_verifier"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
}

func (tu *PKCETokenUrlInfo) URL() string {
	return fmt.Sprintf("https://%s/oauth/token", tu.Host)
}

func (tu *PKCETokenUrlInfo) Body() ([]byte, error) {
	return json.Marshal(tu)
}

func (tu *PKCETokenUrlInfo) Exchange() (*http.Response, error) {
	body, err := tu.Body()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, tu.URL(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set(hum.HeaderContentType, hum.ContentTypeAppJsonUtf8)
	client := &http.Client{}
	return client.Do(req)
}
