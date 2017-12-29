package oauth2more

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var (
	RelCredentialsDir = ".credentials"
)

// ReadTokenFile retrieves a Token from a given filepath.
func ReadTokenFile(fpath string) (*oauth2.Token, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	defer f.Close()
	return tok, err
}

// WriteTokenFile writes a token file to the the filepaths.
func WriteTokenFile(fpath string, tok *oauth2.Token) error {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "Unable to write OAuth token")
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(tok)
}

type TokenStoreFile struct {
	Token    *oauth2.Token
	Filepath string
}

func NewTokenStoreFile(file string) *TokenStoreFile {
	return &TokenStoreFile{Filepath: file}
}

func (ts *TokenStoreFile) Read() error {
	tok, err := ReadTokenFile(ts.Filepath)
	if err != nil {
		return err
	}
	ts.Token = tok
	return nil
}

func (ts *TokenStoreFile) Write() error {
	return WriteTokenFile(ts.Filepath, ts.Token)
}

func (ts *TokenStoreFile) NewTokenFromWeb(cfg *oauth2.Config) (*oauth2.Token, error) {
	tok, err := NewTokenFromWeb(cfg)
	if err != nil {
		return &oauth2.Token{}, err
	}
	ts.Token = tok
	err = ts.Write()
	return tok, err
}

func UserCredentialsDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, RelCredentialsDir), nil
}

func UserCredentialsDirMk(perm os.FileMode) (string, error) {
	dir, err := UserCredentialsDir()
	if err != nil {
		return dir, err
	}
	err = os.MkdirAll(dir, perm)
	return dir, err
}

//
func NewClientWebTokenStore(ctx context.Context, conf *oauth2.Config, tStore *TokenStoreFile, forceNewToken bool) (*http.Client, error) {
	err := tStore.Read()
	client := &http.Client{}

	if err != nil || forceNewToken {
		_, err := tStore.NewTokenFromWeb(conf)
		if err != nil {
			return client, err
		}
	}
	return conf.Client(ctx, tStore.Token), nil
}

func NewTokenStoreFileDefault(tokenPath string, useDefaultDir bool, perm os.FileMode) (*TokenStoreFile, error) {
	tokenPath = strings.TrimSpace(tokenPath)
	tokenFileDefault := "default_credentials.json"
	if tokenPath == "" {
		tokenDir, err := UserCredentialsDirMk(0700)
		if err != nil {
			return &TokenStoreFile{}, err
		}
		tokenPath = filepath.Join(tokenDir, tokenFileDefault)
	} else {
		slashIndex := strings.Index(tokenPath, "/")
		if slashIndex != 0 {
			tokenDir := "."
			if useDefaultDir {
				tokenDirTry, err := UserCredentialsDirMk(0700)
				if err != nil {
					return &TokenStoreFile{}, err
				}
				tokenDir = tokenDirTry
			}
			tokenPath = filepath.Join(tokenDir, tokenPath)
		}
	}
	return NewTokenStoreFile(tokenPath), nil
}
