package facebook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/grokify/gotilla/net/httputil"
	"github.com/grokify/oauth2-util-go/scimutil"
)

const (
	FacebookAPIMeURL = "https://graph.facebook.com/v2.9/me?locale=en_US&fields=name,email,verified,first_name,middle_name,last_name"
)

// FacebookClientUtil is a client library to retrieve user info
// from the Facebook API.
type FacebookClientUtil struct {
	Client *http.Client
}

func NewFacebookClientUtil(client *http.Client) FacebookClientUtil {
	return FacebookClientUtil{Client: client}
}

// GetUserinfo retrieves the userinfo from the
// https://graph.facebook.com/v2.9/{user-id}
// endpoint.
func (apiutil *FacebookClientUtil) GetUserinfo() (FacebookUserinfo, error) {
	resp, err := apiutil.Client.Get(FacebookAPIMeURL)
	if err != nil {
		return FacebookUserinfo{}, err
	}

	bodyBytes, err := httputil.ResponseBody(resp)
	if err != nil {
		return FacebookUserinfo{}, err
	}

	fmt.Printf("%v\n", string(bodyBytes))

	userinfo := FacebookUserinfo{}
	err = json.Unmarshal(bodyBytes, &userinfo)
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

func (apiutil *FacebookClientUtil) GetSCIMUser() (scimutil.User, error) {
	user := scimutil.User{}

	fbUser, err := apiutil.GetUserinfo()
	if err != nil {
		return user, err
	}

	emailAddr := strings.ToLower(strings.TrimSpace(fbUser.Email))
	if len(emailAddr) > 0 {
		email := scimutil.Email{
			Value:   emailAddr,
			Primary: true}
		user.Emails = []scimutil.Email{email}
	}

	user.Name = scimutil.Name{
		GivenName:  strings.TrimSpace(fbUser.FirstName),
		MiddleName: strings.TrimSpace(fbUser.MiddleName),
		FamilyName: strings.TrimSpace(fbUser.LastName),
		Formatted:  strings.TrimSpace(fbUser.Name)}

	return user, nil
}
