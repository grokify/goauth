package authutil

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/type/stringsutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func ClientCredentialsToken(ctx context.Context, cfg clientcredentials.Config) (*oauth2.Token, error) {
	basicHeader, err := BasicAuthHeader(cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		return nil, err
	}
	body := url.Values{
		ParamGrantType: []string{GrantTypeClientCredentials},
	}
	scopes := stringsutil.SliceCondenseSpace(cfg.Scopes, true, false)
	if len(scopes) > 0 {
		body.Add(ParamScope, strings.Join(scopes, " "))
	}
	sr := httpsimple.Request{
		Method: http.MethodPost,
		URL:    cfg.TokenURL,
		Headers: http.Header{
			httputilmore.HeaderAuthorization: []string{basicHeader},
			httputilmore.HeaderContentType:   []string{httputilmore.ContentTypeAppFormURLEncodedUtf8},
		},
		Body: body.Encode(),
	}
	resp, err := sr.Do(ctx, nil)
	if err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API respose status code not successfui (%d)", resp.StatusCode)
	} else if tok, err := ParseTokenReader(resp.Body); err != nil {
		return nil, err
	} else {
		return tok, nil
	}
}
