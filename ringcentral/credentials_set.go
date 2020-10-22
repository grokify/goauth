package ringcentral

import (
	"fmt"
	"net/http"

	"github.com/grokify/gotilla/encoding/jsonutil"
)

type CredentialsSet struct {
	Credentials map[string]Credentials
}

func ReadFileCredentialsSet(filename string) (CredentialsSet, error) {
	set := CredentialsSet{}
	_, err := jsonutil.ReadFile(filename, &set)
	return set, err
}

func (set *CredentialsSet) Get(key string) (Credentials, error) {
	if creds, ok := set.Credentials[key]; ok {
		return creds, nil
	}
	return Credentials{}, fmt.Errorf("E_CREDS_NOT_FOUND [%s]", key)
}

func ReadCredentialsFromFile(filename, key string) (Credentials, error) {
	set, err := ReadFileCredentialsSet(filename)
	if err != nil {
		return Credentials{}, err
	}
	return set.Get(key)
}

func (set *CredentialsSet) GetClient(key string) (*http.Client, error) {
	creds, ok := set.Credentials[key]
	if !ok {
		return nil, fmt.Errorf("E_CREDS_KEY_NOT_FOUND [%v]", key)
	}
	return creds.NewClient()
}
