package authutil

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
	"golang.org/x/oauth2"
)

// NewTokenAccountCredentials is to support Zoom API's `application_credentials` OAuth 2.0 grant type.
// It is unknown if anyone else uses this at the time of this writing. This grant type is described here:
// https://developers.zoom.us/docs/internal-apps/s2s-oauth/ .
func NewTokenAccountCredentials(ctx context.Context, tokenEndpoint, clientID, clientSecret string, bodyOpts url.Values) (*oauth2.Token, error) {
	basicAuthHeaderValue, err := BasicAuthHeader(clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	bodyOpts.Add(ParamGrantType, GrantTypeAccountCredentials)
	req := httpsimple.Request{
		Method:   http.MethodPost,
		URL:      tokenEndpoint,
		Headers:  map[string][]string{httputilmore.HeaderAuthorization: {basicAuthHeaderValue}},
		Body:     bodyOpts,
		BodyType: httpsimple.BodyTypeForm}
	if resp, err := httpsimple.Do(req); err != nil {
		return nil, err
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		tok := &oauth2.Token{}
		return tok, json.Unmarshal(b, tok)
	}
}
