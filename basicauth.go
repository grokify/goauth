package goauth

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/grokify/mogo/net/httputilmore"
	"github.com/grokify/mogo/time/timeutil"
	"golang.org/x/oauth2"
)

// RFC7617UserPass base64 encodes a user-id and password per:
// https://tools.ietf.org/html/rfc7617#section-2
func RFC7617UserPass(userid, password string) (string, error) {
	if strings.Contains(userid, ":") {
		return "", fmt.Errorf(
			"RFC7617 user-id cannot include a colon (':') [%v]", userid)
	}

	return base64.StdEncoding.EncodeToString(
		[]byte(userid + ":" + password),
	), nil
}

func BasicAuthHeader(userid, password string) (string, error) {
	apiKey, err := RFC7617UserPass(userid, password)
	if err != nil {
		return "", err
	}
	return TokenBasic + " " + apiKey, nil
}

// BasicAuthToken provides Basic Authentication support via an oauth2.Token.
func BasicAuthToken(username, password string) (*oauth2.Token, error) {
	basicToken, err := RFC7617UserPass(username, password)
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: basicToken,
		TokenType:   TokenBasic,
		Expiry:      timeutil.TimeZeroRFC3339()}, nil
}

// NewClientBasicAuth returns a *http.Client given a basic auth
// username and password.
func NewClientBasicAuth(username, password string, tlsInsecureSkipVerify bool) (*http.Client, error) {
	authHeaderVal, err := BasicAuthHeader(username, password)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	header := http.Header{}
	header.Add(httputilmore.HeaderAuthorization, authHeaderVal)

	if tlsInsecureSkipVerify {
		client = ClientTLSInsecureSkipVerify(client)
	}

	client.Transport = httputilmore.TransportWithHeaders{
		Header:    header,
		Transport: client.Transport}
	return client, nil
}

func HandlerFuncWrapBasicAuth(handler http.HandlerFunc, username, password, realm, errmsg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set(httputilmore.HeaderWWWAuthenticate, `Basic realm="`+realm+`"`)
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unauthorized.\n"))
			if err != nil {
				log.Println(err.Error())
			}
			return
		}

		handler(w, r)
	}
}

/*

400	Bad Request	[RFC-ietf-httpbis-semantics, Section 15.5.1]
401	Unauthorized	[RFC-ietf-httpbis-semantics, Section 15.5.2]
402	Payment Required	[RFC-ietf-httpbis-semantics, Section 15.5.3]
403	Forbidden	[RFC-ietf-httpbis-semantics, Section 15.5.4]
404	Not Found	[RFC-ietf-httpbis-semantics, Section 15.5.5]
405	Method Not Allowed	[RFC-ietf-httpbis-semantics, Section 15.5.6]
406	Not Acceptable	[RFC-ietf-httpbis-semantics, Section 15.5.7]
407	Proxy Authentication Required	[RFC-ietf-httpbis-semantics, Section 15.5.8]
408	Request Timeout	[RFC-ietf-httpbis-semantics, Section 15.5.9]
409	Conflict	[RFC-ietf-httpbis-semantics, Section 15.5.10]
410	Gone	[RFC-ietf-httpbis-semantics, Section 15.5.11]
411	Length Required	[RFC-ietf-httpbis-semantics, Section 15.5.12]
412	Precondition Failed	[RFC-ietf-httpbis-semantics, Section 15.5.13]
413	Content Too Large	[RFC-ietf-httpbis-semantics, Section 15.5.14]
414	URI Too Long	[RFC-ietf-httpbis-semantics, Section 15.5.15]
415	Unsupported Media Type	[RFC-ietf-httpbis-semantics, Section 15.5.16]
416	Range Not Satisfiable	[RFC-ietf-httpbis-semantics, Section 15.5.17]
417	Expectation Failed	[RFC-ietf-httpbis-semantics, Section 15.5.18]

*/
