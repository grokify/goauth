package oauth2util

import (
	"encoding/json"
	"os"
	"os/user"
	"path"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// ReadTokenFile retrieves a Token from a given filepath.
func ReadTokenFile(filepath string) (*oauth2.Token, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	defer f.Close()
	return tok, err
}

// WriteTokenFile writes a token file to the the filepaths.
func WriteTokenFile(filepath string, tok *oauth2.Token) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
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

func NewTokenStoreFile(file string) TokenStoreFile {
	return TokenStoreFile{Filepath: file}
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

func UserCredentialsDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, ".credentials"), nil
}
