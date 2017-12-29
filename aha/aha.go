package aha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	hum "github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/strconv/humannameparser"
	ou "github.com/grokify/oauth2more"
	"github.com/grokify/oauth2more/scim"
	"golang.org/x/oauth2"
)

const (
	APIMeURL         = "https://secure.aha.io/api/v1/me"
	AuthURLFormat    = "https://%s.aha.io/oauth/authorize"
	TokenURLFormat   = "https://%s.aha.io/oauth/token"
	AhaAccountHeader = "X-AHA-ACCOUNT"
)

func NewEndpoint(subdomain string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf(AuthURLFormat, subdomain),
		TokenURL: fmt.Sprintf(TokenURLFormat, subdomain)}
}

func NewClient(subdomain, token string) *http.Client {
	client := ou.NewClientAccessToken(token)

	header := http.Header{}
	header.Add(AhaAccountHeader, subdomain)

	client.Transport = hum.TransportWithHeaders{
		Transport: client.Transport,
		Header:    header}
	return client
}

// ClientUtil is a client library to retrieve user info
// from the Facebook API.
type ClientUtil struct {
	Client *http.Client
	User   *AhaUserinfo `json:"user,omitempty"`
}

func NewClientUtil(client *http.Client) ClientUtil {
	return ClientUtil{Client: client}
}

func (apiutil *ClientUtil) SetClient(client *http.Client) {
	apiutil.Client = client
}

// GetUserinfo retrieves the userinfo from the
// https://graph.facebook.com/v2.9/{user-id}
// endpoint.
func (apiutil *ClientUtil) GetUserinfo() (*AhaUserinfo, error) {
	resp, err := apiutil.Client.Get(APIMeURL)
	if err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Aha.io API returned Status Code %v", resp.StatusCode)
	}

	bodyBytes, err := hum.ResponseBody(resp)
	if err != nil {
		return nil, err
	}

	userinfo := AhaUserinfoWrap{}
	err = json.Unmarshal(bodyBytes, &userinfo)
	if err == nil {
		apiutil.User = userinfo.User
	}
	return userinfo.User, err
}

type AhaUserinfoWrap struct {
	User *AhaUserinfo `json:"user,omitempty"`
}

type AhaUserinfo struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func (apiutil *ClientUtil) GetSCIMUser() (scim.User, error) {
	user := scim.User{}

	svcUser, err := apiutil.GetUserinfo()
	if err != nil {
		return user, err
	}

	emailAddr := strings.ToLower(strings.TrimSpace(svcUser.Email))
	if len(emailAddr) > 0 {
		email := scim.Item{
			Value:   emailAddr,
			Primary: true}
		user.Emails = []scim.Item{email}
	}

	user.Name = scim.Name{
		Formatted: strings.TrimSpace(svcUser.Name),
	}

	hn, err := humannameparser.ParseHumanName(user.Name.Formatted)
	if err == nil {
		user.Name.GivenName = hn.FirstName
		user.Name.MiddleName = hn.MiddleName
		user.Name.FamilyName = hn.LastName
	}

	return user, nil
}
