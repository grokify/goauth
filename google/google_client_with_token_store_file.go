package google

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/grokify/oauth2more"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/pkg/errors"

	"golang.org/x/oauth2"
)

type GoogleConfigFileStore struct {
	CredentialsRaw []byte
	Credentials    *Credentials
	OAuthConfig    *oauth2.Config
	Scopes         []string
	TokenPath      string
	UseDefaultDir  bool
	ForceNewToken  bool
	State          string
}

// LoadCredentialsBytes set this after setting Scopes.
func (gc *GoogleConfigFileStore) LoadCredentialsBytes(bytes []byte) error {
	if len(strings.TrimSpace(string(bytes))) == 0 {
		return errors.Wrap(errors.New("No Credentials Provided"), "GoogleConfigFileStore.LoadCredentialsBytes()")
	}
	o2Config, err := ConfigFromBytes(bytes, gc.Scopes)
	if err != nil {
		return errors.Wrap(err, "GoogleConfigFileStore.LoadCredentialsBytes() - ConfigFromBytes")
	}
	gc.CredentialsRaw = bytes
	credsContainer, err := CredentialsContainerFromBytes(bytes)
	if err != nil {
		return errors.Wrap(err, "GoogleConfigFileStore.LoadCredentialsBytes() - CredentialsContainerFromBytes")
	}
	gc.Credentials = credsContainer.Credentials()
	gc.OAuthConfig = o2Config
	return nil
}

// SetDefaultFilepath creates a default filepath for the
// file system based token file.
func (gc *GoogleConfigFileStore) SetDefaultFilepath() error {
	projectID := "PlaceholderProjectId"
	clientID := "PlaceholderClientId"
	creds := gc.Credentials
	if creds == nil || len(strings.TrimSpace(creds.ClientID)) == 0 {
		return errors.New("GoogleConfigFileStore.SetDefaultFilepath() - No Credentials Loaded")
	}
	creds.ProjectID = strings.TrimSpace(creds.ProjectID)
	creds.ClientID = strings.TrimSpace(creds.ClientID)
	if len(creds.ProjectID) > 0 {
		projectID = creds.ProjectID
	}
	if len(creds.ClientID) > 0 {
		clientID = creds.ClientID
	}
	scopesShort := []string{}
	for _, scope := range gc.Scopes {
		leaf, err := urlutil.GetPathLeaf(scope)
		if err != nil {
			return errors.Wrap(err, "GoogleConfigFileStore.SetDefaultFilepath - GetPathLeaf")
		}
		leaf = strings.TrimSpace(leaf)
		if len(leaf) > 0 {
			scopesShort = append(scopesShort, leaf)
		}
	}
	sort.Strings(scopesShort)
	scopesStr := strings.Join(scopesShort, "---")
	filename := fmt.Sprintf(
		`google__project-id--%s__client-id--%s__scopes--%s.json`,
		projectID, clientID, scopesStr)
	gc.TokenPath = filename
	gc.UseDefaultDir = true
	return nil
}

// Client returns a `*http.Client`.
func (gc *GoogleConfigFileStore) Client() (*http.Client, error) {
	return NewClientFileStore(
		gc.CredentialsRaw,
		gc.Scopes,
		gc.TokenPath,
		gc.UseDefaultDir,
		gc.ForceNewToken,
		gc.State)
}

// NewClientFileStore returns a `*http.Client` with Google credentials
// in a token store file. It will use the token file credentials unless
// `forceNewToken` is set to true.
func NewClientFileStore(
	credentials []byte,
	scopes []string,
	tokenPath string,
	useDefaultDir bool,
	forceNewToken bool,
	state string) (*http.Client, error) {

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
	googHttpClient, err := oauth2more.NewClientWebTokenStore(context.Background(), conf, tokenStore, forceNewToken, state)
	if err != nil {
		return nil, err
	}
	if !forceNewToken {
		cu := NewClientUtil(googHttpClient)
		_, err := cu.GetUserinfo()
		if err != nil {
			fmt.Printf("E_GOOGLE_USER_PROFILE_API_ERROR [%v] ... Getting New Token...\n", err.Error())
			googHttpClient, err = oauth2more.NewClientWebTokenStore(context.Background(), conf, tokenStore, true, state)
			if err != nil {
				return nil, err
			}
		}
	}

	return googHttpClient, err
}

// NewClientFileStoreWithDefaults returns a `*http.Client` using file system cache
// for access tokens.
func NewClientFileStoreWithDefaults(googleCredentials []byte, googleScopes []string, forceNewToken bool) (*http.Client, error) {
	gcfs := GoogleConfigFileStore{
		Scopes:        googleScopes,
		ForceNewToken: forceNewToken}
	err := gcfs.LoadCredentialsBytes(googleCredentials)
	if err != nil {
		return nil, errors.Wrap(err, "NewGoogleClient - LoadCredentialsBytes")
	}
	err = gcfs.SetDefaultFilepath()
	if err != nil {
		return nil, errors.Wrap(err, "NewGoogleClient - SetDefaultFilepath")
	}
	return gcfs.Client()
}

// NewClientFileStoreWithDefaultsCliEnv instantiates an `*http.Client` for the
// Google API for use from the command line interface (CLI). It will prompt
// the user to open the browser to auth when necessary.
func NewClientFileStoreWithDefaultsCliEnv(googleCredentialsEnvVar, googleScopesEnvVar string) (*http.Client, error) {
	googleCredentialsEnvVar = strings.TrimSpace(googleCredentialsEnvVar)
	googleScopesEnvVar = strings.TrimSpace(googleScopesEnvVar)
	if len(googleCredentialsEnvVar) == 0 {
		googleCredentialsEnvVar = EnvGoogleAppCredentials
	}
	if len(googleScopesEnvVar) == 0 {
		googleScopesEnvVar = EnvGoogleAppScopes
	}
	return NewClientFileStoreWithDefaults(
		[]byte(os.Getenv(googleCredentialsEnvVar)),
		stringsutil.SplitCondenseSpace(os.Getenv(googleScopesEnvVar), ","),
		false)
}
