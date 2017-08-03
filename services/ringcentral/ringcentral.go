package ringcentral

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"strings"

	"github.com/grokify/gotilla/net/httputil"
	"github.com/grokify/gotilla/net/urlutil"
	"github.com/grokify/oauth2-util-go/scimutil"
)

var (
	Hostname = "platform.devtest.ringcentral.com"
)

const (
	ProductionHostname = "platform.ringcentral.com"
	SandboxHostname    = "platform.devtest.ringcentral.com"
	AuthURLFormat      = "https://%s/restapi/oauth/authorize"
	TokenURLFormat     = "https://%s/restapi/oauth/token"
	APIMeURL           = "/restapi/v1.0/account/~/extension/~"
)

func NewEndpoint(hostname string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf(AuthURLFormat, hostname),
		TokenURL: fmt.Sprintf(TokenURLFormat, hostname)}
}

// FacebookClientUtil is a client library to retrieve user info
// from the Facebook API.
type RingCentralClientUtil struct {
	Client *http.Client
}

func NewRingCentralClientUtil(client *http.Client) RingCentralClientUtil {
	return RingCentralClientUtil{Client: client}
}

// GetUserinfo retrieves the userinfo from the
// https://graph.facebook.com/v2.9/{user-id}
// endpoint.
func (apiutil *RingCentralClientUtil) GetUserinfo() (RingCentralExtensionInfo, error) {
	resp, err := apiutil.Client.Get(
		urlutil.JoinAbsolute(
			fmt.Sprintf("%v://", httputil.SchemeHTTPS), Hostname, APIMeURL))

	if err != nil {
		return RingCentralExtensionInfo{}, err
	}

	bodyBytes, err := httputil.ResponseBody(resp)
	if err != nil {
		return RingCentralExtensionInfo{}, err
	}

	fmt.Printf("%v\n", string(bodyBytes))

	userinfo := RingCentralExtensionInfo{}
	err = json.Unmarshal(bodyBytes, &userinfo)
	return userinfo, err
}

type RingCentralExtensionInfo struct {
	ID              int64              `json:"id,omitempty"`
	ExtensionNumber string             `json:"extensionNumber,omitempty"`
	Contact         RingCentralContact `json:"contact,omitempty"`
	Name            string             `json:"name,omitempty"`
}

type RingCentralContact struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
}

func (apiutil *RingCentralClientUtil) GetSCIMUser() (scimutil.User, error) {
	user := scimutil.User{}

	rcUser, err := apiutil.GetUserinfo()
	if err != nil {
		return user, err
	}

	emailAddr := strings.ToLower(strings.TrimSpace(rcUser.Contact.Email))
	if len(emailAddr) > 0 {
		email := scimutil.Email{
			Value:   emailAddr,
			Primary: true}
		user.Emails = []scimutil.Email{email}
	}

	user.Name = scimutil.Name{
		GivenName:  strings.TrimSpace(rcUser.Contact.FirstName),
		FamilyName: strings.TrimSpace(rcUser.Contact.LastName),
		Formatted:  strings.TrimSpace(rcUser.Name)}

	return user, nil
}
