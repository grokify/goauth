package credentials

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/errors/errorsutil"
)

type CredentialsSet struct {
	Credentials map[string]Credentials
}

func ReadFileCredentialsSet(credentialsSetFilename string, inflateEndpoints bool) (CredentialsSet, error) {
	set := CredentialsSet{}
	_, err := jsonutil.ReadFile(credentialsSetFilename, &set)
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

/*
func (set *CredentialsSet) NewSimpleClient(accountKey string) (*httpsimple.SimpleClient, error) {
	creds, ok := set.Credentials[accountKey]
	if !ok {
		return nil, fmt.Errorf("client_not_found [%s]", accountKey)
	}
	return creds.NewSimpleClient()
}
*/

func ReadCredentialsFromFile(credentialsSetFilename, accountKey string, inclAccountsOnError bool) (Credentials, error) {
	set, err := ReadFileCredentialsSet(credentialsSetFilename, true)
	if err != nil {
		return Credentials{}, err
	}
	creds, err := set.Get(accountKey)
	if err != nil {
		if inclAccountsOnError {
			return creds, errorsutil.Wrap(err,
				fmt.Sprintf("validAccounts [%s]", strings.Join(set.Accounts(), ",")))
		}
		return creds, err
	}
	return creds, nil
}

func (set *CredentialsSet) GetClient(ctx context.Context, key string) (*http.Client, error) {
	creds, ok := set.Credentials[key]
	if !ok {
		return nil, fmt.Errorf("E_CREDS_KEY_NOT_FOUND [%v]", key)
	}
	return creds.NewClient(ctx)
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
