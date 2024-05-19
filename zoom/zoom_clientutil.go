package zoom

import (
	"net/http"
	"strings"
	"time"

	"github.com/grokify/goauth/endpoints"
	"github.com/grokify/goauth/scim"
	"github.com/grokify/mogo/encoding/jsonutil"
)

const (
	ZoomAPIOAuth2AuthzURL = endpoints.ZoomAuthzURL
	ZoomAPIOAuth2TokenURL = endpoints.ZoomTokenURL // #nosec G101
	ZoomAPIURLBase        = "https://api.zoom.us/v2/"
	ZoomAPIURLUsers       = "https://api.zoom.us/v2/users/"
	ZoomAPIURLUsersMe     = "https://api.zoom.us/v2/users/me"
	ZoomAPIUserIDMe       = "me"
)

type ClientUtil struct {
	Client     *http.Client
	UserNative ZoomUser `json:"user,omitempty"`
	UserScim   scim.User
	UserLoaded bool
}

func NewClientUtil(client *http.Client) ClientUtil {
	return ClientUtil{Client: client}
}

func (apiutil *ClientUtil) SetClient(client *http.Client) {
	apiutil.Client = client
}

func (apiutil *ClientUtil) LoadUser() error {
	resp, err := apiutil.Client.Get(ZoomAPIURLUsersMe)
	if err != nil {
		return err
	}

	nativeUser := ZoomUser{}
	_, err = jsonutil.UnmarshalReader(resp.Body, &nativeUser)
	if err != nil {
		return err
	}
	apiutil.UserNative = nativeUser
	if apiutil.UserScim, err = ZoomUserToScimUser(nativeUser); err != nil {
		return err
	} else {
		apiutil.UserLoaded = true
		return nil
	}
}

func (apiutil *ClientUtil) GetSCIMUser() (scim.User, error) {
	if !apiutil.UserLoaded {
		err := apiutil.LoadUser()
		if err != nil {
			return scim.User{}, err
		}
	}
	return apiutil.UserScim, nil
}

func ZoomUserToScimUser(nativeUser ZoomUser) (scim.User, error) {
	scimUser := scim.User{}
	err := scimUser.AddEmail(strings.TrimSpace(nativeUser.Email), true)
	if err != nil {
		return scimUser, err
	}
	scimUser.Name = scim.Name{
		GivenName:  strings.TrimSpace(nativeUser.FirstName),
		FamilyName: strings.TrimSpace(nativeUser.LastName),
		Formatted: strings.TrimSpace(nativeUser.FirstName) +
			" " + strings.TrimSpace(nativeUser.LastName)}
	return scimUser, nil
}

type ZoomUser struct {
	ID                 string    `json:"id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Email              string    `json:"email"`
	Type               int       `json:"type"`
	RoleName           string    `json:"role_name"`
	PMI                int       `json:"pmi"`
	UsePMI             bool      `json:"use_pmi"`
	PersonalMeetingURL string    `json:"personal_meeting_url"`
	Timezone           string    `json:"timezone"`
	Verified           int       `json:"verified"`
	Dept               string    `json:"dept"`
	CreatedAt          time.Time `json:"created_at"`
	LastLoginTime      time.Time `json:"last_login_time"`
	LastClientVersion  string    `json:"last_client_version"`
	PicURL             string    `json:"pic_url"`
	HostKey            string    `json:"host_key"`
	JID                string    `json:"jid"`
	GroupIDs           []string  `json:"group_ids"`
	IMGroupIDs         []string  `json:"im_group_ids"`
	AccountID          string    `json:"account_id"`
	Language           string    `json:"language"`
	PhoneCountry       string    `json:"phone_country"`
	PhoneNumber        string    `json:"phone_number"`
	Status             string    `json:"status"`
}
