package multiservice

import (
	"encoding/json"

	"golang.org/x/oauth2"
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
	return &o2cc, err
}

func (c *O2ConfigMore) ProviderType() (OAuth2Provider, error) {
	return ProviderStringToConst(c.Provider)
}

func (c *O2ConfigMore) Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectUris[0],
		Scopes:       c.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AuthUri,
			TokenURL: c.TokenUri,
		},
	}
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
