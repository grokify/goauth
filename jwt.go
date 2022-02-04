package goauth

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/httputilmore"
	"golang.org/x/oauth2"
)

func ParseJwtTokenString(tokenString string, secretKey string, claims jwt.Claims) (*jwt.Token, error) {
	// https://stackoverflow.com/questions/41077953/go-language-and-verify-jwt
	if claims == nil {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			return nil, errorsutil.Wrap(err, "ParseTokenString.jwt.Parse")
		}
		return token, nil
	}
	// *jwt.StandardClaims
	// https://stackoverflow.com/questions/45405626/decoding-jwt-token-in-golang
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errorsutil.Wrap(err, "ParseTokenString.jwt.ParseWithClaims")
	}
	return token, nil
}

func NewTokenOAuth2Jwt(tokenURL, clientId, clientSecret, jwtBase64Enc string) (*oauth2.Token, error) {
	body := url.Values{
		ParamGrantType: {GrantTypeJWTBearer},
		ParamAssertion: {jwtBase64Enc}}

	req, err := http.NewRequest(
		http.MethodPost, tokenURL,
		strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppFormUrlEncoded)

	if len(clientId) > 0 || len(clientSecret) > 0 {
		b64Enc, err := RFC7617UserPass(clientId, clientSecret)
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
}
