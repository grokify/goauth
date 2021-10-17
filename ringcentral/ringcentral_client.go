package ringcentral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	om "github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/credentials"
	hum "github.com/grokify/simplego/net/httputilmore"
	"golang.org/x/oauth2"
)

func NewTokenPassword(oc credentials.OAuth2Credentials) (*oauth2.Token, error) {
	return RetrieveToken(
		oauth2.Config{
			ClientID:     oc.ClientID,
			ClientSecret: oc.ClientSecret,
			Endpoint:     oc.OAuth2Endpoint},
		oc.PasswordRequestBody())
}

// NewClientPassword uses dedicated password grant handling.
func NewClientPassword(oc credentials.OAuth2Credentials) (*http.Client, error) {
	c := oc.Config()
	token, err := RetrieveToken(c, oc.PasswordRequestBody())
	if err != nil {
		return nil, err
	}

	httpClient := c.Client(oauth2.NoContext, token)

	header := getClientHeader(oc)
	if len(header) > 0 {
		httpClient.Transport = hum.TransportWithHeaders{
			Transport: httpClient.Transport,
			Header:    header}
	}
	return httpClient, nil
}

// NewClientPasswordSimple uses OAuth2 package password grant handling.
func NewClientPasswordSimple(oc credentials.OAuth2Credentials) (*http.Client, error) {
	httpClient, err := om.NewClientPasswordConf(
		oauth2.Config{
			ClientID:     oc.ClientID,
			ClientSecret: oc.ClientSecret,
			Endpoint:     oc.OAuth2Endpoint},
		oc.Username,
		oc.Password)
	if err != nil {
		return nil, err
	}

	header := getClientHeader(oc)
	if len(header) > 0 {
		httpClient.Transport = hum.TransportWithHeaders{
			Transport: httpClient.Transport,
			Header:    header}
	}
	return httpClient, nil
}

func getClientHeader(oc credentials.OAuth2Credentials) http.Header {
	userAgentParts := []string{om.PathVersion()}
	if len(oc.AppNameAndVersion()) > 0 {
		userAgentParts = append([]string{oc.AppNameAndVersion()}, userAgentParts...)
	}
	userAgent := strings.TrimSpace(strings.Join(userAgentParts, "; "))

	header := http.Header{}
	if len(userAgent) > 0 {
		header.Add(hum.HeaderUserAgent, userAgent)
		header.Add("X-User-Agent", userAgent)
	}
	return header
}

func RetrieveToken(cfg oauth2.Config, params url.Values) (*oauth2.Token, error) {
	rcToken, err := RetrieveRcToken(cfg, params)
	if err != nil {
		return nil, err
	}
	return rcToken.OAuth2Token()
}

func RetrieveRcToken(cfg oauth2.Config, params url.Values) (*RcToken, error) {
	r, err := http.NewRequest(
		http.MethodPost,
		cfg.Endpoint.TokenURL,
		strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	basicAuthHeader, err := om.BasicAuthHeader(cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		return nil, err
	}

	r.Header.Add(hum.HeaderAuthorization, basicAuthHeader)
	r.Header.Add(hum.HeaderContentType, hum.ContentTypeAppFormUrlEncoded)
	r.Header.Add(hum.HeaderContentLength, strconv.Itoa(len(params.Encode())))

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("RingCentral API Response Status [%v][%s]", resp.StatusCode, string(bytes))
	}

	rcToken := &RcToken{}
	err = json.Unmarshal(bytes, rcToken)
	if err != nil {
		return nil, err
	}
	err = rcToken.Inflate()
	return rcToken, err
}

type RcToken struct {
	AccessToken           string    `json:"access_token,omitempty"`
	TokenType             string    `json:"token_type,omitempty"`
	Scope                 string    `json:"scope,omitempty"`
	ExpiresIn             int64     `json:"expires_in,omitempty"`
	RefreshToken          string    `json:"refresh_token,omitempty"`
	RefreshTokenExpiresIn int64     `json:"refresh_token_expires_in,omitempty"`
	OwnerID               string    `json:"owner_id,omitempty"`
	EndpointID            string    `json:"endpoint_id,omitempty"`
	Expiry                time.Time `json:"expiry,omitempty"`
	RefreshTokenExpiry    time.Time `json:"refresh_token_expiry,omitempty"`
	inflated              bool      `json:"inflated"`
}

func (rcTok *RcToken) Inflate() error {
	now := time.Now()
	if (rcTok.ExpiresIn) > 0 {
		expiresIn, err := time.ParseDuration(fmt.Sprintf("%vs", rcTok.ExpiresIn))
		if err != nil {
			return err
		}
		rcTok.Expiry = now.Add(expiresIn)
	}
	if (rcTok.RefreshTokenExpiresIn) > 0 {
		expiresIn, err := time.ParseDuration(fmt.Sprintf("%vs", rcTok.RefreshTokenExpiresIn))
		if err != nil {
			return err
		}
		rcTok.RefreshTokenExpiry = now.Add(expiresIn)
	}

	rcTok.inflated = true
	return nil
}

func (rcTok *RcToken) OAuth2Token() (*oauth2.Token, error) {
	if !rcTok.inflated {
		err := rcTok.Inflate()
		if err != nil {
			return nil, err
		}
	}

	tok := &oauth2.Token{
		AccessToken:  rcTok.AccessToken,
		TokenType:    rcTok.TokenType,
		RefreshToken: rcTok.RefreshToken,
		Expiry:       rcTok.Expiry}

	return tok, nil
}
