package multiservice

import (
	"fmt"
	"strings"
)

type OAuth2Provider int

// OAuth2Provider a constant list of OAuth2 providers.
// Warning: do not rely on ordering or integer value
// as this will change as additional providers are added.
const (
	Aha OAuth2Provider = iota
	Facebook
	Google
	Instagram
	Lyft
	Metabase
	RingCentral
	SparkPost
	Twitter
	Visa
	Zendesk
)

var providers = [...]string{
	"aha",
	"facebook",
	"google",
	"instagram",
	"lyft",
	"metabase",
	"ringcentral",
	"sparkpost",
	"twitter",
	"visa",
	"zendesk",
}

func (p OAuth2Provider) String() string {
	if Aha <= p && p <= Zendesk {
		return providers[p]
	}
	panic(fmt.Sprintf("Provider not in range [%v]", string(p)))
	return ""
}

func ProviderStringToConst(s string) (OAuth2Provider, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	for i, p := range providers {
		if p == s {
			return OAuth2Provider(i), nil
		}
	}
	return RingCentral, fmt.Errorf("OAuth2 Provider Not Found [%v]", s)
}
