package multiservice

import (
	"testing"
)

var configTests = []struct {
	configJson string
	provider   OAuth2Provider
}{
	{`{"provider":"google","client_id":"1234567890.apps.googleusercontent.com","project_id":"api-project-123456","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://accounts.google.com/o/oauth2/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"1234567890","redirect_uris":["https://example.com/oauth2callback"],"javascript_origins":["https://example.com"],"scopes":["https://www.googleapis.com/auth/bigquery","https://www.googleapis.com/auth/blogger"]}`, Google},
}

func TestConfigs(t *testing.T) {
	for _, tt := range configTests {
		try, err := NewO2ConfigMoreFromJSON([]byte(tt.configJson))
		if err != nil {
			t.Errorf("NewO2ConfigMoreFromJSON(%v): err [%v]",
				tt.configJson, err)
		}
		provider, err := ProviderStringToConst(try.Provider)
		if err != nil {
			t.Errorf("ProviderStringToConst(%v): err [%v]",
				try.Provider, err)
		}
		if provider != tt.provider {
			t.Errorf("NewO2ConfigMoreFromJSON(%v).Provider: want [%v], got [%v]",
				tt.configJson, tt.provider, try.Provider)
		}
	}
}
