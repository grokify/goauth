package credentials

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/grokify/simplego/net/http/httpsimple"
	"github.com/pkg/errors"
)

type CredentialsSet struct {
	Credentials map[string]Credentials
}

func ReadFileCredentialsSet(filename string, inflateEndpoints bool) (CredentialsSet, error) {
	set := CredentialsSet{}
	_, err := jsonutil.ReadFile(filename, &set)
	if err != nil {
		return set, err
	}
	if inflateEndpoints {
		set.Inflate()
	}
	return set, nil
}

func (set *CredentialsSet) Get(key string) (Credentials, error) {
	if creds, ok := set.Credentials[key]; ok {
		return creds, nil
	}
	return Credentials{}, fmt.Errorf("E_CREDS_NOT_FOUND [%s]", key)
}

func (set *CredentialsSet) Inflate() {
	for k, v := range set.Credentials {
		v.Inflate()
		set.Credentials[k] = v
	}
}

func (set *CredentialsSet) NewSimpleClient(key string) (*httpsimple.SimpleClient, error) {
	creds, ok := set.Credentials[key]
	if !ok {
		return nil, fmt.Errorf("client_not_found [%s]", key)
	}
	return creds.NewSimpleClient()
}

func ReadCredentialsFromFile(filename, key string, inclAccountsOnError bool) (Credentials, error) {
	set, err := ReadFileCredentialsSet(filename, true)
	if err != nil {
		return Credentials{}, err
	}
	creds, err := set.Get(key)
	if err != nil {
		if inclAccountsOnError {
			return creds, errors.Wrap(err,
				fmt.Sprintf("validAccounts [%s]", strings.Join(set.Accounts(), ",")))
		}
		return creds, err
	}
	return creds, nil
}

func (set *CredentialsSet) GetClient(key string) (*http.Client, error) {
	creds, ok := set.Credentials[key]
	if !ok {
		return nil, fmt.Errorf("E_CREDS_KEY_NOT_FOUND [%v]", key)
	}
	return creds.NewClient()
}

func (set *CredentialsSet) Accounts() []string { return set.Keys() }

func (set *CredentialsSet) Keys() []string {
	keys := []string{}
	for key := range set.Credentials {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
