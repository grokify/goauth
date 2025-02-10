package jwtutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
)

const (
	// claims from https://datatracker.ietf.org/doc/html/rfc7519 .
	JWTClaimAudience   = "aud"
	JWTClaimExpiration = "exp"
	JWTClaimIssuedAt   = "iat"
	JWTClaimIssuer     = "iss"
	JWTClaimJWTID      = "jti"
	JWTClaimNotBefore  = "nbf"
	JWTClaimSubject    = "sub"
)

func ParseJWTString(tokenString string, secretKey string, claims jwt.Claims) (*jwt.Token, error) {
	// https://stackoverflow.com/questions/41077953/go-language-and-verify-jwt
	if claims == nil {
		if token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		}); err != nil {
			return nil, errorsutil.Wrap(err, "ParseTokenString.jwt.Parse")
		} else {
			return token, nil
		}
	}
	// *jwt.StandardClaims
	// https://stackoverflow.com/questions/45405626/decoding-jwt-token-in-golang
	if token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	}); err != nil {
		return nil, errorsutil.Wrap(err, "ParseTokenString.jwt.ParseWithClaims")
	} else {
		return token, nil
	}
}

func NewTokenOAuth2JWT(ctx context.Context, tokenURL, clientID, clientSecret, jwtBase64Enc string) (*oauth2.Token, error) {
	sreq := httpsimple.Request{
		Method:  http.MethodPost,
		URL:     tokenURL,
		Headers: http.Header{},
		Body: url.Values{
			authutil.ParamGrantType: {authutil.GrantTypeJWTBearer},
			authutil.ParamAssertion: {jwtBase64Enc}},
		BodyType: httpsimple.BodyTypeForm,
	}
	if len(clientID) > 0 || len(clientSecret) > 0 {
		if authHeaderVal, err := authutil.BasicAuthHeader(clientID, clientSecret); err != nil {
			return nil, errorsutil.NewErrorWithLocation(err.Error())
		} else {
			sreq.Headers.Add(httputilmore.HeaderAuthorization, authHeaderVal)
		}
	}
	if hreq, err := sreq.HTTPRequest(ctx); err != nil {
		return nil, errorsutil.NewErrorWithLocation(err.Error())
	} else if resp, err := ctxhttp.Do(ctx, &http.Client{}, hreq); err != nil {
		return nil, errorsutil.NewErrorWithLocation(err.Error())
	} else if resp.StatusCode >= 300 {
		msa := map[string]any{
			"func":              "jwtutil.NewTokenOAuth2JWT()",
			"httpResStatusCode": resp.StatusCode,
			"httpReqURL":        sreq.URL,
			"httpReqMethod":     sreq.Method,
		}
		b, err := json.Marshal(msa)
		if err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("tokenURL (httpResStatus: %d) %s", resp.StatusCode, string(b))
	} else {
		return authutil.ParseTokenReader(resp.Body)
	}

	/*
		body := url.Values{
			ParamGrantType: {GrantTypeJWTBearer},
			ParamAssertion: {jwtBase64Enc}}

		req, err := http.NewRequest(
			http.MethodPost, tokenURL,
			strings.NewReader(body.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppFormURLEncoded)

		if len(clientID) > 0 || len(clientSecret) > 0 {
			b64Enc, err := RFC7617UserPass(clientID, clientSecret)
			if err != nil {
				return nil, err
			}
			req.Header.Add(httputilmore.HeaderAuthorization, TokenBasic+" "+b64Enc)
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		return ParseTokenReader(resp.Body)
	*/
}
