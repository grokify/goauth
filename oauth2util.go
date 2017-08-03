package oauth2util

import (
	"net/http"

	"github.com/grokify/oauth2-util-go/scimutil"
)

type OAuth2Util interface {
	SetClient(*http.Client)
	GetSCIMUser() (scimutil.User, error)
}
