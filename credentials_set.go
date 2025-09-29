package goauth

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/errors/errorsutil"
)

type CredentialsSet struct {
	Credentials map[string]Credentials `json:"credentials,omitempty"`
}

func ReadFileCredentialsSet(filename string, inflateEndpoints bool) (*CredentialsSet, error) {
	var set *CredentialsSet
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := jsonutil.UnmarshalWithLoc(b, &set); err != nil {
		return nil, errorsutil.WrapWithLocation(err)
	} else if inflateEndpoints {
		err := set.Inflate()
		if err != nil {
			return nil, errorsutil.WrapWithLocation(err)
		}
	}
	return set, nil
}

func (set *CredentialsSet) Get(key string) (Credentials, error) {
	if creds, ok := set.Credentials[key]; ok {
		return creds, nil
	}
	return Credentials{}, fmt.Errorf("credentials key not found (%s)", key)
}

func (set *CredentialsSet) Inflate() error {
	for k, v := range set.Credentials {
		err := v.Inflate()
		if err != nil {
			return err
		}
		set.Credentials[k] = v
	}
	return nil
}

func ReadCredentialsFromCLI(inclAccountsOnError bool) (Credentials, error) {
	if opts, err := ParseOptions(); err != nil {
		return Credentials{}, err
	} else {
		return opts.Credentials()
	}
}

func ReadCredentialsFromSetFile(credentialsSetFilename, accountKey string, inclAccountsOnError bool) (Credentials, error) {
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

func (set *CredentialsSet) NewClient(ctx context.Context, key string) (*http.Client, error) {
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

func (set *CredentialsSet) WriteFile(filename, prefix, indent string, perm fs.FileMode) error {
	return jsonutil.MarshalFile(filename, set, prefix, indent, perm)
}
