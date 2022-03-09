package multiservice

import (
	"fmt"
	"strings"
)

// OAuth2Provider a constant list of OAuth2 providers.
// Warning: do not rely on ordering or integer value
// as this will change as additional providers are added.
type OAuth2Provider int

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

// String converts a provider type to a string.
func (p OAuth2Provider) String() string {
	if Aha <= p && p <= Zendesk {
		return providers[p]
	}
	panic(fmt.Sprintf("Provider not in range [%d]", p))
}

// ProviderStringToConst returns an OAuth2Provider type constant
// from a string.
func ProviderStringToConst(s string) (OAuth2Provider, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	for i, p := range providers {
		if p == s {
			return OAuth2Provider(i), nil
		}
	}
	return RingCentral, fmt.Errorf("OAuth2 Provider Not Found [%v]", s)
}
