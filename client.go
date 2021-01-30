package oauth2more

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/grokify/simplego/net/httputilmore"
	"github.com/grokify/simplego/time/timeutil"
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
	token, err := conf.PasswordCredentialsToken(oauth2.NoContext, username, password)
	if err != nil {
		return &http.Client{}, err
	}

	return conf.Client(oauth2.NoContext, token), nil
}

func NewClientAuthCode(conf oauth2.Config, authCode string) (*http.Client, error) {
	token, err := conf.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		return &http.Client{}, err
	}
	return conf.Client(oauth2.NoContext, token), nil
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

func NewClientHeaders(headersMap map[string]string, tlsInsecureSkipVerify bool) *http.Client {
	client := &http.Client{}
	header := httputilmore.NewHeadersMSS(headersMap)

	if tlsInsecureSkipVerify {
		client = ClientTLSInsecureSkipVerify(client)
	}

	client.Transport = httputilmore.TransportWithHeaders{
		Header:    header,
		Transport: client.Transport}

	return client
}

func NewClientToken(tokenType, tokenValue string, tlsInsecureSkipVerify bool) *http.Client {
	client := &http.Client{}

	header := http.Header{}
	header.Add(httputilmore.HeaderAuthorization, tokenType+" "+tokenValue)

	if tlsInsecureSkipVerify {
		client = ClientTLSInsecureSkipVerify(client)
	}

	client.Transport = httputilmore.TransportWithHeaders{
		Header:    header,
		Transport: client.Transport}

	return client
}

func NewClientTokenBase64Encode(tokenType, tokenValue string, tlsInsecureSkipVerify bool) *http.Client {
	return NewClientToken(
		tokenType,
		base64.StdEncoding.EncodeToString([]byte(tokenValue)),
		tlsInsecureSkipVerify)
}

// NewClientAuthzTokenSimple returns a *http.Client given a token type and token string.
func NewClientAuthzTokenSimple(tokenType, accessToken string) *http.Client {
	token := &oauth2.Token{
		AccessToken: strings.TrimSpace(accessToken),
		TokenType:   strings.TrimSpace(tokenType),
		Expiry:      timeutil.TimeZeroRFC3339()}

	oAuthConfig := &oauth2.Config{}

	return oAuthConfig.Client(oauth2.NoContext, token)
}

func NewClientTokenOAuth2(token *oauth2.Token) *http.Client {
	oAuthConfig := &oauth2.Config{}
	return oAuthConfig.Client(oauth2.NoContext, token)
}

func NewClientBearerTokenSimpleOrJson(ctx context.Context, tokenOrJson []byte) (*http.Client, error) {
	tokenOrJsonString := strings.TrimSpace(string(tokenOrJson))
	if len(tokenOrJsonString) == 0 {
		return nil, fmt.Errorf("No token [%v]", string(tokenOrJson))
	} else if strings.Index(tokenOrJsonString, "{") == 0 {
		return NewClientTokenJSON(ctx, tokenOrJson)
	} else {
		return NewClientAuthzTokenSimple(TokenBearer, tokenOrJsonString), nil
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
			InsecureSkipVerify: true}}
	return client
}
