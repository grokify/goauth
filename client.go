package goauth

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/mogo/net/httputilmore"
	"github.com/grokify/mogo/time/timeutil"
	"golang.org/x/oauth2"
)

func NewClientPassword(conf oauth2.Config, ctx context.Context, username, password string) (*http.Client, error) {
	token, err := BasicAuthToken(username, password)
	if err != nil {
		return nil, err
	}
	return conf.Client(ctx, token), nil
}

func NewClientPasswordConf(conf oauth2.Config, username, password string) (*http.Client, error) {
	token, err := conf.PasswordCredentialsToken(context.Background(), username, password)
	if err != nil {
		return &http.Client{}, err
	}

	return conf.Client(context.Background(), token), nil
}

func NewClientAuthCode(conf oauth2.Config, authCode string) (*http.Client, error) {
	token, err := conf.Exchange(context.Background(), authCode)
	if err != nil {
		return &http.Client{}, err
	}
	return conf.Client(context.Background(), token), nil
}

func NewClientTokenJSON(ctx context.Context, tokenJSON []byte) (*http.Client, error) {
	token := &oauth2.Token{}
	err := json.Unmarshal(tokenJSON, token)
	if err != nil {
		return nil, err
	}

	oAuthConfig := &oauth2.Config{}

	return oAuthConfig.Client(ctx, token), nil
}

// NewClientHeaderQuery returns a new `*http.Client` that will set headers and query
// string parameters on very request.
func NewClientHeaderQuery(header http.Header, query url.Values, allowInsecure bool) *http.Client {
	client := &http.Client{}

	if allowInsecure {
		client = ClientTLSInsecureSkipVerify(client)
	}

	client.Transport = httputilmore.TransportRequestModifier{
		Header:    header,
		Query:     query,
		Transport: client.Transport}

	return client
}

func NewClientToken(tokenType, tokenValue string, allowInsecure bool) *http.Client {
	return NewClientHeaderQuery(
		http.Header{httputilmore.HeaderAuthorization: []string{tokenType + " " + tokenValue}},
		url.Values{},
		allowInsecure)
}

func NewClientTokenBase64Encode(tokenType, tokenValue string, allowInsecure bool) *http.Client {
	return NewClientToken(
		tokenType,
		base64.StdEncoding.EncodeToString([]byte(tokenValue)),
		allowInsecure)
}

// NewClientAuthzTokenSimple returns a *http.Client given a token type and token string.
func NewClientAuthzTokenSimple(tokenType, accessToken string) *http.Client {
	token := &oauth2.Token{
		AccessToken: strings.TrimSpace(accessToken),
		TokenType:   strings.TrimSpace(tokenType),
		Expiry:      timeutil.TimeZeroRFC3339()}

	oAuthConfig := oauth2.Config{}
	return oAuthConfig.Client(context.Background(), token)
}

func NewClientTokenOAuth2(token *oauth2.Token) *http.Client {
	oAuthConfig := oauth2.Config{}
	return oAuthConfig.Client(context.Background(), token)
}

func NewClientBearerTokenSimpleOrJSON(ctx context.Context, tokenOrJSON []byte) (*http.Client, error) {
	tokenOrJSONString := strings.TrimSpace(string(tokenOrJSON))
	if len(tokenOrJSONString) == 0 {
		return nil, fmt.Errorf("no token [%v]", string(tokenOrJSON))
	} else if strings.Index(tokenOrJSONString, "{") == 0 {
		return NewClientTokenJSON(ctx, tokenOrJSON)
	} else {
		return NewClientAuthzTokenSimple(TokenBearer, tokenOrJSONString), nil
	}
}

func NewClientTLSToken(ctx context.Context, tlsConfig *tls.Config, token *oauth2.Token) *http.Client {
	tlsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig}}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, tlsClient)

	cfg := &oauth2.Config{}

	return cfg.Client(ctx, token)
}

func ClientTLSInsecureSkipVerify(client *http.Client) *http.Client {
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true}} // #nosec G402
	return client
}
