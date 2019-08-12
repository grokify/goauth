package multiservice

import (
	"testing"
)

var providerTests = []struct {
	constantKey OAuth2Provider
	constantVal string
}{
	{Aha, "aha"},
	{RingCentral, "ringcentral"},
	{Visa, "visa"},
}

func TestProviders(t *testing.T) {
	for _, tt := range providerTests {
		try := tt.constantKey.String()
		if try != tt.constantVal {
			t.Errorf("provider.String(%v): want [%v], got [%v]",
				tt.constantKey, tt.constantVal, try)
		}
	}
}

func TestProvidersReverse(t *testing.T) {
	for _, tt := range providerTests {
		try, err := ProviderStringToConst(tt.constantVal)
		if err != nil {
			t.Errorf("provider.ProviderStringToConst(%v): want [%v], err [%v]",
				tt.constantVal, tt.constantKey, err)
		}
		if try != tt.constantKey {
			t.Errorf("provider.ProviderStringToConst(%v): want [%v], got [%v]",
				tt.constantVal, tt.constantKey, try)
		}
	}
}
