package oauth2util

import (
	"net/http"
	"time"

	"github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/oauth2util-go/scimutil"
	"golang.org/x/oauth2"
)

type OAuth2Util interface {
	SetClient(*http.Client)
	GetSCIMUser() (scimutil.User, error)
}

func NewClientAccessToken(accessToken string) *http.Client {
	t0, _ := time.Parse(time.RFC3339, timeutil.RFC3339Zero)

	token := &oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		Expiry:      t0}

	oAuthConfig := &oauth2.Config{}

	return oAuthConfig.Client(oauth2.NoContext, token)
}
