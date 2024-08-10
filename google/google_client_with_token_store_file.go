package google

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/type/stringsutil"
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
		return errorsutil.Wrap(errors.New("no credentials provided"), "err GoogleConfigFileStore.LoadCredentialsBytes()")
	}
	o2Config, err := ConfigFromBytes(bytes, gc.Scopes)
	if err != nil {
		return errorsutil.Wrap(err, "err GoogleConfigFileStore.LoadCredentialsBytes() - ConfigFromBytes")
	}
	gc.CredentialsRaw = bytes
	credsContainer, err := CredentialsContainerFromBytes(bytes)
	if err != nil {
		return errorsutil.Wrap(err, "err GoogleConfigFileStore.LoadCredentialsBytes() - CredentialsContainerFromBytes")
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
		return errors.New("err GoogleConfigFileStore.SetDefaultFilepath() - No Credentials Loaded")
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
			return errorsutil.Wrap(err, "err GoogleConfigFileStore.SetDefaultFilepath - GetPathLeaf")
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
func (gc *GoogleConfigFileStore) Client(ctx context.Context) (*http.Client, error) {
	return NewClientFileStore(
		ctx,
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
	ctx context.Context,
	credentials []byte,
	scopes []string,
	tokenPath string,
	useDefaultDir bool,
	forceNewToken bool,
	state string) (*http.Client, error) {
	if len(strings.TrimSpace(string(credentials))) == 0 {
		return nil, errorsutil.Wrap(errors.New("no credentials provided"), "err GoogleConfigFileStore.LoadCredentialsBytes()")
	}

	conf, err := ConfigFromBytes(credentials, scopes)
	if err != nil {
		return nil, err
	}
	tokenStoreFile, err := authutil.NewTokenStoreFileDefault(tokenPath, useDefaultDir, 0700)
	if err != nil {
		return nil, err
	}
	googHTTPClient, err := authutil.NewClientWebTokenStore(ctx, conf, tokenStoreFile, forceNewToken, state)
	if err != nil {
		return nil, err
	}
	if !forceNewToken {
		cu := NewClientUtil(googHTTPClient)
		_, err := cu.GetUserinfo(ctx)
		if err != nil {
			fmt.Printf("error for Google user profile API [%v] ... Getting New Token", err.Error())
			googHTTPClient, err = authutil.NewClientWebTokenStore(ctx, conf, tokenStoreFile, true, state)
			if err != nil {
				return nil, err
			}
		}
	}

	return googHTTPClient, err
}

// NewClientFileStoreWithDefaults returns a `*http.Client` using file system cache
// for access tokens.
func NewClientFileStoreWithDefaults(ctx context.Context, googleCredentials []byte, googleScopes []string, forceNewToken bool) (*http.Client, error) {
	gcfs := GoogleConfigFileStore{
		Scopes:        googleScopes,
		ForceNewToken: forceNewToken}
	if err := gcfs.LoadCredentialsBytes(googleCredentials); err != nil {
		return nil, errorsutil.Wrap(err, "err NewClientFileStoreWithDefaults - LoadCredentialsBytes")
	}
	if err := gcfs.SetDefaultFilepath(); err != nil {
		return nil, errorsutil.Wrap(err, "err NewClientFileStoreWithDefaults - SetDefaultFilepath")
	} else {
		return gcfs.Client(ctx)
	}
}

// NewClientFileStoreWithDefaultsCliEnv instantiates an `*http.Client` for the
// Google API for use from the command line interface (CLI). It will prompt
// the user to open the browser to auth when necessary.
func NewClientFileStoreWithDefaultsCliEnv(ctx context.Context, googleCredentialsEnvVar, googleScopesEnvVar string) (*http.Client, error) {
	googleCredentialsEnvVar = strings.TrimSpace(googleCredentialsEnvVar)
	googleScopesEnvVar = strings.TrimSpace(googleScopesEnvVar)
	if len(googleCredentialsEnvVar) == 0 {
		googleCredentialsEnvVar = EnvGoogleAppCredentials
	}
	if len(googleScopesEnvVar) == 0 {
		googleScopesEnvVar = EnvGoogleAppScopes
	}
	return NewClientFileStoreWithDefaults(
		ctx,
		[]byte(os.Getenv(googleCredentialsEnvVar)),
		stringsutil.SplitCondenseSpace(os.Getenv(googleScopesEnvVar), ","),
		false)
}
