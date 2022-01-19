package aha

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/grokify/goauth/scim"
	"github.com/grokify/mogo/strconv/humannameparser"
)

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

	bodyBytes, err := io.ReadAll(resp.Body)
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

	err = user.AddEmail(svcUser.Email, true)
	if err != nil {
		return user, err
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
