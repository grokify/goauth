package multiservice

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// O2ConfigCanonical is similar to Google but includes scopes
type O2ConfigMore struct {
	Provider                string   `json:"provider,omitempty"`
	ClientId                string   `json:"client_id,omitempty"`
	ClientSecret            string   `json:"client_secret,omitempty"`
	ProjectId               string   `json:"project_id,omitempty"`
	AuthUri                 string   `json:"auth_uri,omitempty"`
	TokenUri                string   `json:"token_uri,omitempty"`
	AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url,omitempty"`
	RedirectUris            []string `json:"redirect_uris,omitempty"`
	JavaScriptOrigins       []string `json:"javascript_origins,omitempty"`
	Scopes                  []string `json:"scopes,omitempty"`
}

func NewO2ConfigMoreFromJSON(bytes []byte) (*O2ConfigMore, error) {
	o2cc := O2ConfigMore{}
	err := json.Unmarshal(bytes, &o2cc)
	if err != nil {
		return nil, err
	}
	o2cc.Provider = strings.ToLower(strings.TrimSpace(o2cc.Provider))
	switch o2cc.Provider {
	case "facebook":
		if len(strings.TrimSpace(o2cc.AuthUri)) == 0 {
			o2cc.AuthUri = facebook.Endpoint.AuthURL
		}
		if len(strings.TrimSpace(o2cc.TokenUri)) == 0 {
			o2cc.TokenUri = facebook.Endpoint.TokenURL
		}
	}
	return &o2cc, nil
}

func (c *O2ConfigMore) ProviderType() (OAuth2Provider, error) {
	return ProviderStringToConst(c.Provider)
}

func (cm *O2ConfigMore) Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cm.ClientId,
		ClientSecret: cm.ClientSecret,
		RedirectURL:  cm.RedirectURL(),
		Scopes:       cm.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cm.AuthUri,
			TokenURL: cm.TokenUri}}
}

func (cm *O2ConfigMore) AuthURL(state string) string {
	return cm.Config().AuthCodeURL(state)
}

func (cm *O2ConfigMore) RedirectURL() string {
	redirectURL := ""
	for _, try := range cm.RedirectUris {
		try := strings.TrimSpace(try)
		if len(try) > 0 {
			redirectURL = try
			break
		}
	}
	return redirectURL
}

func RandomState(statePrefix string, randomSuffix bool) string {
	parts := []string{}
	if len(statePrefix) > 0 {
		parts = append(parts, statePrefix)
	}
	if randomSuffix {
		parts = append(parts, fmt.Sprintf("%v", rand.Intn(1000000000)))
	}
	return strings.Join(parts, "-")
}

/*
Example:
{
   "web":{
   	  "provider":"google",
      "client_id":"1234567890.apps.googleusercontent.com",
      "project_id":"api-project-123456",
      "auth_uri":"https://accounts.google.com/o/oauth2/auth",
      "token_uri":"https://accounts.google.com/o/oauth2/token",
      "auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
      "client_secret":"1234567890",
      "redirect_uris":[
         "https://example.com/oauth2callback"
      ],
      "javascript_origins":[
         "https://example.com"
      ],
      "scopes":[
         "https://www.googleapis.com/auth/bigquery",
         "https://www.googleapis.com/auth/blogger"
      ]
   }
}
*/
