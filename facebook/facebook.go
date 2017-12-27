package facebook

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/oauth2util/scimutil"
	"golang.org/x/oauth2/facebook"
)

const (
	FacebookAPIMeURL = "https://graph.facebook.com/v2.9/me?locale=en_US&fields=name,email,verified,first_name,middle_name,last_name"
)

func DefaultifyConfig(cfg *oauth2.Config) *oauth2.Config {
	cfg.Endpoint = fb.Endpoint
	return cfg
}

// ClientUtil is a client library to retrieve user info
// from the Facebook API.
type ClientUtil struct {
	Client *http.Client
	User   FacebookUserinfo `json:"user,omitempty"`
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
func (apiutil *ClientUtil) GetUserinfo() (FacebookUserinfo, error) {
	resp, err := apiutil.Client.Get(FacebookAPIMeURL)
	if err != nil {
		return FacebookUserinfo{}, err
	}

	bodyBytes, err := httputilmore.ResponseBody(resp)
	if err != nil {
		return FacebookUserinfo{}, err
	}

	userinfo := FacebookUserinfo{}
	err = json.Unmarshal(bodyBytes, &userinfo)
	if err == nil {
		apiutil.User = userinfo
	}
	return userinfo, err
}

type FacebookUserinfo struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	Verified   bool   `json:"verified,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
}

func (apiutil *ClientUtil) GetSCIMUser() (scimutil.User, error) {
	user := scimutil.User{}

	fbUser, err := apiutil.GetUserinfo()
	if err != nil {
		return user, err
	}

	emailAddr := strings.ToLower(strings.TrimSpace(fbUser.Email))
	if len(emailAddr) > 0 {
		email := scimutil.Item{
			Value:   emailAddr,
			Primary: true}
		user.Emails = []scimutil.Item{email}
	}

	user.Name = scimutil.Name{
		GivenName:  strings.TrimSpace(fbUser.FirstName),
		MiddleName: strings.TrimSpace(fbUser.MiddleName),
		FamilyName: strings.TrimSpace(fbUser.LastName),
		Formatted:  strings.TrimSpace(fbUser.Name)}

	return user, nil
}
