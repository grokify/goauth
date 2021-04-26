package zendesk

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/grokify/oauth2more/scim"
	"github.com/grokify/simplego/encoding/jsonutil"
)

// ClientUtil is a client library to retrieve user info
// from the Facebook API.
type ClientUtil struct {
	Client    *http.Client `json:"-"`
	Subdomain string       `json:"subdomain,omitempty"`
	MeUser    Me           `json:"user,omitempty"`
}

func NewClientUtil(client *http.Client, subdomain string) ClientUtil {
	return ClientUtil{Client: client, Subdomain: subdomain}
}

func (cu *ClientUtil) SetClient(client *http.Client) {
	cu.Client = client
}

func (cu *ClientUtil) GetUserinfo() (*Me, error) {
	me, resp, err := GetMe(cu.Client, cu.Subdomain)
	if err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Zendesk API Status Code %v", resp.StatusCode)
	}
	cu.MeUser = *me
	return me, err
}

func (cu *ClientUtil) GetSCIMUser() (scim.User, error) {
	user := scim.User{}

	me, err := cu.GetUserinfo()
	if err != nil {
		return user, err
	}

	err = user.AddEmail(me.Email, true)
	if err != nil {
		return user, err
	}

	user.Name = scim.Name{
		Formatted: strings.TrimSpace(me.Name)}

	return user, nil
}

func MeURL(subdomain string) string {
	return fmt.Sprintf("https://%v.zendesk.com/api/v2/users/me.json", subdomain)
}

type MeResponse struct {
	User Me `json:"user,omitempty"`
}

type Me struct {
	ID                int64  `json:"id,omitempty"`
	URL               string `json:"url,omitempty"`
	Name              string `json:"name,omitempty"`
	Email             string `json:"email,omitempty"`
	Phone             string `json:"phone,omitempty"`
	SharedPhoneNumber string `json:"shared_phone_number,omitempty"`
}

func GetMe(client *http.Client, subdomain string) (*Me, *http.Response, error) {
	meURL := MeURL(subdomain)
	resp, err := client.Get(meURL)
	if err != nil {
		return nil, resp, err
	}
	me := &MeResponse{}
	_, err = jsonutil.UnmarshalReader(resp.Body, me)
	return &me.User, resp, err
}
