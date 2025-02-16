package goauth

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/goauth/authutil/jwtutil"
	"github.com/grokify/mogo/crypto/randutil"
	"github.com/grokify/mogo/encoding/basex"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/type/stringsutil"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// CredentialsOAuth2 supports OAuth 2.0 authorization_code, password, and client_credentials grant flows.
type CredentialsOAuth2 struct {
	ServerURL            string              `json:"serverURL,omitempty"`
	ApplicationID        string              `json:"applicationID,omitempty"`
	ClientID             string              `json:"clientID,omitempty"`
	ClientSecret         string              `json:"clientSecret,omitempty"`
	Endpoint             oauth2.Endpoint     `json:"endpoint,omitempty"`
	RedirectURL          string              `json:"redirectURL,omitempty"`
	OAuthEndpointID      string              `json:"oauthEndpointID,omitempty"`
	Scopes               []string            `json:"scope,omitempty"`
	GrantType            string              `json:"grantType,omitempty"`
	PKCE                 bool                `json:"pkce"`
	Username             string              `json:"username,omitempty"`
	Password             string              `json:"password,omitempty"`
	JWT                  string              `json:"jwt,omitempty"`
	Token                *oauth2.Token       `json:"token,omitempty"`
	AuthCodeOpts         map[string][]string `json:"authCodeOpts,omitempty"`
	AuthCodeExchangeOpts map[string][]string `json:"authCodeExchangeOpts,omitempty"`
	TokenBodyOpts        url.Values          `json:"tokenBodyOpts,omitempty"`
	Metadata             map[string]string   `json:"metadata,omitempty"`
}

func ParseCredentialsOAuth2(b []byte) (CredentialsOAuth2, error) {
	creds := CredentialsOAuth2{}
	return creds, json.Unmarshal(b, &creds)
}

// MarshalJSON returns JSON. It is useful for exporting creating configs to be parsed.
func (oc *CredentialsOAuth2) MarshalJSON(prefix, indent string) ([]byte, error) {
	return jsonutil.MarshalSimple(*oc, prefix, indent)
}

func (oc *CredentialsOAuth2) Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     oc.ClientID,
		ClientSecret: oc.ClientSecret,
		Endpoint:     oc.Endpoint,
		RedirectURL:  oc.RedirectURL,
		Scopes:       oc.Scopes}
}

func (oc *CredentialsOAuth2) ConfigClientCredentials() clientcredentials.Config {
	return clientcredentials.Config{
		ClientID:       oc.ClientID,
		ClientSecret:   oc.ClientSecret,
		TokenURL:       oc.Endpoint.TokenURL,
		Scopes:         oc.Scopes,
		EndpointParams: oc.TokenBodyOpts,
		AuthStyle:      oauth2.AuthStyleAutoDetect}
}

type AuthCodeOptions []oauth2.AuthCodeOption

func (opts *AuthCodeOptions) Add(k, v string) {
	*opts = append(*opts, oauth2.SetAuthURLParam(k, v))
}

func (opts *AuthCodeOptions) AddMap(m map[string][]string) {
	for k, vs := range m {
		for _, v := range vs {
			opts.Add(k, v)
		}
	}
}

func (oc *CredentialsOAuth2) AuthCodeURL(state string, opts map[string][]string) string {
	authCodeOptions := AuthCodeOptions{}
	authCodeOptions.AddMap(oc.AuthCodeOpts)
	authCodeOptions.AddMap(opts)
	cfg := oc.Config()
	return cfg.AuthCodeURL(state, authCodeOptions...)
}

func (oc *CredentialsOAuth2) BasicAuthHeader() (string, error) {
	return authutil.BasicAuthHeader(oc.ClientID, oc.ClientSecret)
}

func (oc *CredentialsOAuth2) Exchange(ctx context.Context, code string, opts map[string][]string) (*oauth2.Token, error) {
	/*
		authCodeOptions := []oauth2.AuthCodeOption{}
		if len(oc.OAuthEndpointID) > 0 {
			authCodeOptions = append(authCodeOptions,
				oauth2.SetAuthURLParam("endpoint_id", oc.OAuthEndpointID))
		}
		if oc.AccessTokenTTL > 0 {
			authCodeOptions = append(authCodeOptions,
				oauth2.SetAuthURLParam("accessTokenTtl", strconv.Itoa(int(oc.AccessTokenTTL))))
		}
		if oc.RefreshTokenTTL > 0 {
			authCodeOptions = append(authCodeOptions,
				oauth2.SetAuthURLParam("refreshTokenTtl", strconv.Itoa(int(oc.RefreshTokenTTL))))
		}
	*/
	authCodeOptions := AuthCodeOptions{}
	authCodeOptions.AddMap(oc.AuthCodeExchangeOpts)
	authCodeOptions.AddMap(opts)
	cfg := oc.Config()
	return cfg.Exchange(ctx, code, authCodeOptions...)
}

func (oc *CredentialsOAuth2) IsGrantType(grantType string) bool {
	return strings.EqualFold(
		strings.TrimSpace(grantType),
		strings.TrimSpace(oc.GrantType))
}

func (oc *CredentialsOAuth2) InflateURL(apiURLPath string) string {
	return urlutil.JoinAbsolute(oc.ServerURL, apiURLPath)
}

func (oc *CredentialsOAuth2) NewClient(ctx context.Context) (*http.Client, *oauth2.Token, error) {
	if tok, err := oc.NewToken(ctx); err != nil {
		return nil, tok, err
	} else {
		oc.Token = tok
		config := oc.Config()
		return config.Client(ctx, tok), tok, nil
	}
}

// NewToken retrieves an `*oauth2.Token` when the requisite information is available.
// Note this uses `clientcredentials.Config.Token()` which doesn't always work. In
// This situation, use `authutil.TokenClientCredentials()` as an alternative. Note: authorization
// code is only supported for CLI testing purposes. In a production application, it should be
// done in a multi-step process to redirect the user to the authorization URL, retrieve the
// auth code and then `Exchange` it for a token. The `state` value is currently a randomly generated
// string as this should be used for testing purposes only.
func (oc *CredentialsOAuth2) NewToken(ctx context.Context) (*oauth2.Token, error) {
	if oc.Token != nil && len(strings.TrimSpace(oc.Token.AccessToken)) > 0 {
		return oc.Token, nil
	} else if strings.Contains(strings.ToLower(oc.GrantType), "jwt") {
		return jwtutil.NewTokenOAuth2JWT(ctx, oc.Endpoint.TokenURL, oc.ClientID, oc.ClientSecret, oc.JWT)
	} else if oc.IsGrantType(authutil.GrantTypeAccountCredentials) {
		return authutil.NewTokenAccountCredentials(ctx, oc.Endpoint.TokenURL, oc.ClientID, oc.ClientSecret, oc.TokenBodyOpts)
	} else if oc.IsGrantType(authutil.GrantTypeClientCredentials) {
		config := oc.ConfigClientCredentials()
		return authutil.ClientCredentialsToken(ctx, config)
		// return config.Token(ctx)
	} else if oc.IsGrantType(authutil.GrantTypePassword) {
		// cfg := oc.Config()
		// return cfg.PasswordCredentialsToken(ctx, oc.Username, oc.Password)
		return oc.NewTokenPasswordCredentials(ctx) // supports custom request params
	} else if oc.IsGrantType(authutil.GrantTypeAuthorizationCode) {
		state := randutil.RandString(basex.AlphabetBase62, 12)
		authURL := oc.AuthCodeURL(state, map[string][]string{})
		fmt.Printf("Authorization URL: %s\n\n", authURL)
		fmt.Printf("Authorization URL State: %s\n\n", state)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Authorization Code:")
		authCode, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return oc.Exchange(ctx, authCode, map[string][]string{})
	} else {
		return nil, fmt.Errorf("grant type [%s] is not supported in CredentialsOAuth2.NewToken()", oc.GrantType)
	}
}

func (oc *CredentialsOAuth2) newTokenPasswordCredentialsRequest() (*httpsimple.Request, error) {
	cfg := oc.Config()
	basicHeaderVal, err := authutil.BasicAuthHeader(cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		return nil, err
	}
	req := httpsimple.Request{
		Method: http.MethodPost,
		URL:    cfg.Endpoint.TokenURL,
		Headers: http.Header{
			httputilmore.HeaderAuthorization: []string{basicHeaderVal},
		},
		BodyType: httpsimple.BodyTypeForm,
	}
	body := url.Values{}
	if len(oc.TokenBodyOpts) > 0 {
		for k, vals := range oc.TokenBodyOpts {
			for _, v := range vals {
				body.Add(k, v)
			}
		}
	}
	body.Add(authutil.ParamGrantType, authutil.GrantTypePassword)
	body.Add(authutil.ParamUsername, oc.Username)
	body.Add(authutil.ParamPassword, oc.Password)
	if len(oc.Scopes) > 0 {
		body.Add(authutil.ParamScope, strings.Join(stringsutil.SliceCondenseSpace(oc.Scopes, true, false), ","))
	}
	req.Body = body
	return &req, nil
}

// NewTokenPasswordCredentials provides fine-grained token request.
func (oc *CredentialsOAuth2) NewTokenPasswordCredentials(ctx context.Context) (*oauth2.Token, error) {
	if sreq, err := oc.newTokenPasswordCredentialsRequest(); err != nil {
		return nil, err
	} else if hreq, err := sreq.HTTPRequest(ctx); err != nil {
		return nil, err
	} else if resp, err := ctxhttp.Do(ctx, &http.Client{}, hreq); err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("receted status code (%d)", resp.StatusCode)
	} else {
		return authutil.ParseTokenReader(resp.Body)
	}
}

func (oc *CredentialsOAuth2) NewSimpleClient(ctx context.Context) (*httpsimple.Client, error) {
	if c, _, err := oc.NewClient(ctx); err != nil {
		return nil, err
	} else {
		return &httpsimple.Client{
			BaseURL:    oc.ServerURL,
			HTTPClient: c}, nil
	}
}

/*
func NewTokenOAuth2(credsOA2 credentials.CredentialsOAuth2) (*oauth2.Token, error) {
	if credsOA2.GrantType == authutil.GrantTypeAuthorizationCode {
		authURL := credsOA2.AuthCodeURL("abc", map[string][]string{})
		fmt.Println(authURL)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Authorization Code:")
		authCode, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return credsOA2.Exchange(context.Background(), authCode, map[string][]string{})
	} else if credsOA2.GrantType == authutil.GrantTypeClientCredentials {
		return credsOA2.NewToken(context.Background()) // direct `creds.NewToken()` doesn't work
	}
	return nil, fmt.Errorf("grant type not supported (%s)", credsOA2.GrantType)
}
*/

func (oc *CredentialsOAuth2) RefreshToken(ctx context.Context, tok *oauth2.Token) (*oauth2.Token, []byte, error) {
	if tok == nil {
		return nil, []byte{}, errors.New("token not supplied")
	} else {
		return oc.RefreshTokenSimple(ctx, tok.RefreshToken)
	}
}

func (oc *CredentialsOAuth2) RefreshTokenSimple(ctx context.Context, refreshToken string) (*oauth2.Token, []byte, error) {
	basicAuthHeader, err := oc.BasicAuthHeader()
	if err != nil {
		return nil, []byte{}, err
	}
	body := url.Values{}
	body.Add(authutil.ParamRefreshToken, refreshToken)
	body.Add(authutil.ParamGrantType, authutil.GrantTypeRefreshToken)
	if len(oc.Scopes) > 0 {
		body.Add(authutil.ParamScope, strings.Join(oc.Scopes, " "))
	}

	sr := httpsimple.Request{
		Method: http.MethodPost,
		URL:    oc.Endpoint.TokenURL,
		Headers: map[string][]string{
			httputilmore.HeaderContentType:   {httputilmore.ContentTypeAppFormURLEncoded},
			httputilmore.HeaderAuthorization: {basicAuthHeader},
		},
		Body: []byte(body.Encode()),
	}

	if resp, err := sr.Do(ctx); err != nil {
		return nil, []byte{}, err
	} else if tokBody, err := io.ReadAll(resp.Body); err != nil {
		return nil, tokBody, err
	} else if resp.StatusCode >= 300 {
		return nil, tokBody, fmt.Errorf("status code (%d)", resp.StatusCode)
	} else {
		tok, err := authutil.ParseToken(tokBody)
		return tok, tokBody, err
	}
	/*
		oaTok, err := goauth.ParseOAuth2Token(tokBody)
		if err != nil {
			return nil, tokBody, err
		}
		tok := oaTok.Token()
		return tok, tokBody, nil
	*/
}

func (oc *CredentialsOAuth2) PasswordRequestBody() url.Values {
	body := url.Values{
		authutil.ParamGrantType: {authutil.GrantTypePassword},
		authutil.ParamUsername:  {oc.Username},
		authutil.ParamPassword:  {oc.Password}}
	if len(oc.TokenBodyOpts) > 0 {
		for k, vals := range oc.TokenBodyOpts {
			for _, v := range vals {
				body.Set(k, v)
			}
		}
	}
	return body
}

func NewCredentialsOAuth2Env(envPrefix string) CredentialsOAuth2 {
	creds := CredentialsOAuth2{
		ClientID:     os.Getenv(envPrefix + "CLIENT_ID"),
		ClientSecret: os.Getenv(envPrefix + "CLIENT_SECRET"),
		ServerURL:    os.Getenv(envPrefix + "SERVER_URL"),
		Username:     os.Getenv(envPrefix + "USERNAME"),
		Password:     os.Getenv(envPrefix + "PASSWORD")}
	if len(strings.TrimSpace(creds.Username)) > 0 {
		creds.GrantType = authutil.GrantTypePassword
	}
	return creds
}
