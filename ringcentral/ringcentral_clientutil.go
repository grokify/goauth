package ringcentral

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/net/urlutil"
	"github.com/grokify/oauth2more/scim"
)

var (
	Hostname = "platform.devtest.ringcentral.com"
)

const (
	ServerURLProduction  = "https://platform.ringcentral.com"
	ServerURLSandbox     = "https://platform.devtest.ringcentral.com"
	ProductionHostname   = "platform.ringcentral.com"
	SandboxHostname      = "platform.devtest.ringcentral.com"
	AuthURLFormat        = "https://%s/restapi/oauth/authorize"
	TokenURLFormat       = "https://%s/restapi/oauth/token"
	AuthURLPart          = "restapi/oauth/authorize"
	TokenURLPart         = "restapi/oauth/token"
	MeURL                = "/restapi/v1.0/account/~/extension/~"
	RestAPI1dot0Fragment = "restapi/v1.0"
)

func NewEndpoint(hostnameOrBasePath string) oauth2.Endpoint {
	hostnameOrBasePath = regexp.
		MustCompile(`/\s*$`).
		ReplaceAllString(hostnameOrBasePath, "")

	m, _ := regexp.MatchString(`^(?i)https?://`, hostnameOrBasePath)

	if m {
		return oauth2.Endpoint{
			AuthURL:  urlutil.JoinAbsolute(hostnameOrBasePath, AuthURLPart),
			TokenURL: urlutil.JoinAbsolute(hostnameOrBasePath, TokenURLPart),
		}
	}
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf(AuthURLFormat, hostnameOrBasePath),
		TokenURL: fmt.Sprintf(TokenURLFormat, hostnameOrBasePath)}
}

// ClientUtil is a client library to retrieve user info
// from the Facebook API.
type ClientUtil struct {
	Client *http.Client             `json:"-"`
	User   RingCentralExtensionInfo `json:"user,omitempty"`
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
func (apiutil *ClientUtil) GetUserinfo() (RingCentralExtensionInfo, error) {
	resp, err := apiutil.Client.Get(
		urlutil.JoinAbsolute(
			fmt.Sprintf("%v://", httputilmore.SchemeHTTPS), Hostname, MeURL))

	if err != nil {
		return RingCentralExtensionInfo{}, err
	}

	bodyBytes, err := httputilmore.ResponseBody(resp)
	if err != nil {
		return RingCentralExtensionInfo{}, err
	}

	userinfo := RingCentralExtensionInfo{}
	err = json.Unmarshal(bodyBytes, &userinfo)
	if err == nil {
		apiutil.User = userinfo
	}
	return userinfo, err
}

type RingCentralExtensionInfo struct {
	ID              int64              `json:"id,omitempty"`
	ExtensionNumber string             `json:"extensionNumber,omitempty"`
	Contact         RingCentralContact `json:"contact,omitempty"`
	Name            string             `json:"name,omitempty"`
	Account         RingCentralAccount `json:"account,omitempty"`
}

type RingCentralContact struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
}

type RingCentralAccount struct {
	URI string `json:"uri,omitempty"`
	ID  string `json:"id,omitempty"`
}

func (apiutil *ClientUtil) GetSCIMUser() (scim.User, error) {
	user := scim.User{}

	rcUser, err := apiutil.GetUserinfo()
	if err != nil {
		return user, err
	}

	err = user.AddEmail(rcUser.Contact.Email, true)
	if err != nil {
		return user, err
	}

	user.Name = scim.Name{
		GivenName:  strings.TrimSpace(rcUser.Contact.FirstName),
		FamilyName: strings.TrimSpace(rcUser.Contact.LastName),
		Formatted:  strings.TrimSpace(rcUser.Name)}

	return user, nil
}

func BuildURL(serverURL, urlFragment string, addRestAPI bool, queryValues url.Values) string {
	apiURL := serverURL
	if addRestAPI {
		apiURL = urlutil.JoinAbsolute(apiURL, RestAPI1dot0Fragment, urlFragment)
	} else {
		apiURL = urlutil.JoinAbsolute(apiURL, urlFragment)
	}
	return urlutil.BuildURL(apiURL, queryValues)
}

func BuildURLPartial(urlFragment string, addRestAPI bool, queryValues url.Values) string {
	apiURL := fmt.Sprintf("%s://%s", httputilmore.SchemeHTTPS, Hostname)
	if addRestAPI {
		apiURL = urlutil.JoinAbsolute(apiURL, RestAPI1dot0Fragment, urlFragment)
	} else {
		apiURL = urlutil.JoinAbsolute(apiURL, urlFragment)
	}
	return urlutil.BuildURL(apiURL, queryValues)
}

func SetHostnameForURL(serverURLString string) error {
	serverURL, err := url.Parse(serverURLString)
	if err != nil {
		return err
	}
	Hostname = strings.TrimSpace(serverURL.Hostname())
	if len(Hostname) < 1 {
		return errors.New("No Hostname")
	}
	return nil
}
