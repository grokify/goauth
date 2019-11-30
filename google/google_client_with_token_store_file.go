package google

import (
	"net/http"
	"strings"

	"github.com/grokify/oauth2more"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type GoogleConfigFileStore struct {
	Credentials   []byte
	OAuthConfig   *oauth2.Config
	Scopes        []string
	TokenPath     string
	UseDefaultDir bool
	ForceNewToken bool
}

// LoadCredentialsBytes set this after setting Scopes.
func (gc *GoogleConfigFileStore) LoadCredentialsBytes(bytes []byte) error {
	if len(strings.TrimSpace(string(bytes))) == 0 {
		return errors.Wrap(errors.New("No Credentials Provided"), "GoogleConfigFileStore.LoadCredentialsBytes()")
	}
	o2Config, err := ConfigFromBytes(bytes, gc.Scopes)
	if err != nil {
		return errors.Wrap(err, "GoogleConfigFileStore.LoadCredentialsBytes()")
	}
	gc.Credentials = bytes
	gc.OAuthConfig = o2Config
	return nil
}

// Client returns a `*http.Client`.
func (gc *GoogleConfigFileStore) Client() (*http.Client, error) {
	return NewClientFileStore(
		gc.Credentials,
		gc.Scopes,
		gc.TokenPath,
		gc.UseDefaultDir,
		gc.ForceNewToken)
}

// NewClientFileStore returns a `*http.Client` with Google credentials
// in a token store file. It will use the token file credentails unless
// `forceNewToken` is set to true.
func NewClientFileStore(
	credentials []byte,
	scopes []string,
	tokenPath string,
	useDefaultDir bool,
	forceNewToken bool) (*http.Client, error) {

	if len(strings.TrimSpace(string(credentials))) == 0 {
		return nil, errors.Wrap(errors.New("No Credentials Provided"), "GoogleConfigFileStore.LoadCredentialsBytes()")
	}

	conf, err := ConfigFromBytes(credentials, scopes)
	if err != nil {
		return nil, err
	}
	tokenStore, err := oauth2more.NewTokenStoreFileDefault(tokenPath, useDefaultDir, 0700)
	if err != nil {
		return nil, err
	}
	return oauth2more.NewClientWebTokenStore(context.Background(), conf, tokenStore, forceNewToken)
}
