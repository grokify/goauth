package oauth2more

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/grokify/gotilla/net/httputilmore"
	hum "github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/time/timeutil"
	"golang.org/x/oauth2"
)

// RFC7617UserPass base64 encodes a user-id and password per:
// https://tools.ietf.org/html/rfc7617#section-2
func RFC7617UserPass(userid, password string) (string, error) {
	if strings.Index(userid, ":") > -1 {
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
	return BasicPrefix + " " + apiKey, nil
}

// BasicAuthToken provides Basic Authentication support via an oauth2.Token.
func BasicAuthToken(username, password string) (*oauth2.Token, error) {
	basicToken, err := RFC7617UserPass(username, password)
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: basicToken,
		TokenType:   BasicPrefix,
		Expiry:      timeutil.TimeRFC3339Zero()}, nil
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

	client.Transport = hum.TransportWithHeaders{
		Header:    header,
		Transport: client.Transport}
	return client, nil
}
