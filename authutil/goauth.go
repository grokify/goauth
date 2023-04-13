package authutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/grokify/goauth/scim"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

const (
	VERSION = "0.10"
	PATH    = "github.com/grokify/goauth"
)

type AuthorizationType int

const (
	Anonymous AuthorizationType = iota
	Basic
	Bearer
	Digest
	NTLM
	Negotiate
	OAuth
)

var authorizationTypes = [...]string{
	"Anonymous",
	"Basic",
	"Bearer",
	"Digest",
	"NTLM",
	"Negotiate",
	"OAuth",
}

// String returns the English name of the authorizationTypes ("Basic", "Bearer", ...).
func (a AuthorizationType) String() string {
	if Basic <= a && a <= OAuth {
		return authorizationTypes[a]
	}
	buf := make([]byte, 20)
	n := fmtInt(buf, uint64(a))
	return "%!AuthorizationType(" + string(buf[n:]) + ")"
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

func PathVersion() string {
	return fmt.Sprintf("%v-v%v", PATH, VERSION)
}

type ServiceType int

const (
	Google ServiceType = iota
	Facebook
	RingCentral
	Aha
)

/*
// ApplicationCredentials represents information for an app.

	type ApplicationCredentials struct {
		ServerURL    string
		ClientID     string
		ClientSecret string
		Endpoint     oauth2.Endpoint
	}
*/
type AppCredentials struct {
	Service      string   `json:"service,omitempty"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIs []string `json:"redirect_uris"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
	Scopes       []string `json:"scopes"`
}

func (ac *AppCredentials) Defaultify() {
	switch ac.Service {
	case "facebook":
		if len(ac.AuthURI) == 0 || len(ac.TokenURI) == 0 {
			endpoint := facebook.Endpoint
			if len(ac.AuthURI) == 0 {
				ac.AuthURI = endpoint.AuthURL
			}
			if len(ac.TokenURI) == 0 {
				ac.TokenURI = endpoint.TokenURL
			}
		}
	}
}

type AppCredentialsWrapper struct {
	Web       *AppCredentials `json:"web"`
	Installed *AppCredentials `json:"installed"`
}

func (w *AppCredentialsWrapper) Config() (*oauth2.Config, error) {
	var c *AppCredentials
	if w.Web != nil {
		c = w.Web
	} else if w.Installed != nil {
		c = w.Installed
	} else {
		return nil, errors.New("no OAuth2 config info")
	}
	c.Defaultify()
	return c.Config(), nil
}

func NewAppCredentialsWrapperFromBytes(data []byte) (AppCredentialsWrapper, error) {
	var acw AppCredentialsWrapper
	err := json.Unmarshal(data, &acw)
	if err != nil {
		panic(err)
	}
	return acw, err
}

func (ac *AppCredentials) Config() *oauth2.Config {
	cfg := &oauth2.Config{
		ClientID:     ac.ClientID,
		ClientSecret: ac.ClientSecret,
		Scopes:       ac.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  ac.AuthURI,
			TokenURL: ac.TokenURI}}

	if len(ac.RedirectURIs) > 0 {
		cfg.RedirectURL = ac.RedirectURIs[0]
	}
	return cfg
}

// UserCredentials represents a user's credentials.
type UserCredentials struct {
	Username string
	Password string
}

type OAuth2Util interface {
	SetClient(*http.Client)
	GetSCIMUser() (scim.User, error)
}
