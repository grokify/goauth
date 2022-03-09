package salesforce

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/grokify/go-salesforce/sobjects"
	"github.com/grokify/goauth/credentials"
	"github.com/grokify/mogo/net/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
)

type SalesforceClient struct {
	ClientMore httputilmore.ClientMore
	URLBuilder URLBuilder
}

func NewSalesforceClient(client *http.Client, instanceName string) SalesforceClient {
	return SalesforceClient{
		ClientMore: httputilmore.ClientMore{Client: client},
		URLBuilder: NewURLBuilder(instanceName),
	}
}

func NewSalesforceClientEnv() (SalesforceClient, error) {
	sc := SalesforceClient{
		URLBuilder: NewURLBuilder(os.Getenv("SALESFORCE_INSTANCE_NAME")),
	}
	client, err := NewClientPasswordSalesforceEnv()
	if err != nil {
		return sc, err
	}
	sc.ClientMore = httputilmore.ClientMore{Client: client}
	return sc, nil
}

type OAuth2Credentials struct {
	credentials.CredentialsOAuth2
	InstanceName string
}

func NewSalesforceClientPassword(soc OAuth2Credentials) (SalesforceClient, error) {
	httpClient, err := NewClientPassword(soc.CredentialsOAuth2)
	if err != nil {
		return SalesforceClient{}, err
	}
	return NewSalesforceClient(httpClient, soc.InstanceName), nil
}

func (sc *SalesforceClient) GetServicesData() (*http.Response, error) {
	apiURL := sc.URLBuilder.Build("services/data")
	return sc.ClientMore.Client.Get(apiURL.String())
}

func (sc *SalesforceClient) CreateContact(contact interface{}) (*http.Response, error) {
	//apiURL := sc.URLBuilder.Build("/services/data/v40.0/sobjects/Contact/")
	apiURL := sc.URLBuilder.BuildSobjectURL("Contact")
	return sc.ClientMore.PostToJSON(apiURL.String(), contact)
}

func (sc *SalesforceClient) CreateSobject(sobjectName string, sobject interface{}) (*http.Response, error) {
	apiURL := sc.URLBuilder.BuildSobjectURL(sobjectName)
	return sc.ClientMore.PostToJSON(apiURL.String(), sobject)
}

func (sc *SalesforceClient) ExecSOQL(soql string) (*http.Response, error) {
	//curl https://yourInstance.salesforce.com/services/data/v20.0/query/?q=SELECT+name+from+Account -H "Authorization: Bearer token"
	apiURL := sc.URLBuilder.Build("/services/data/v40.0/query/")
	soql = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(soql), "+")
	qryMap := map[string][]string{"q": {soql}}
	apiURLString := urlutil.URLAddQuery(&apiURL, qryMap).String()
	apiURLString = regexp.MustCompile(`\%2B`).ReplaceAllString(strings.TrimSpace(apiURLString), "+")
	return sc.ClientMore.Client.Get(apiURLString)
}

func (sc *SalesforceClient) GetAccountsAll() (sobjects.AccountSet, error) {
	resp, err := sc.ExecSOQL("SELECT id, name FROM account")
	if err != nil {
		return sobjects.AccountSet{}, err
	}

	err = httputilmore.PrintResponse(resp, true)
	if err != nil {
		return sobjects.AccountSet{}, err
	}

	return sobjects.NewAccountSetFromJSONResponse(resp)
}

func (sc *SalesforceClient) DeleteAccountsAll() error {
	set, err := sc.GetAccountsAll()
	if err != nil {
		return err
	}
	for _, account := range set.Records {
		resp, err := sc.DeleteAccount(account.Id)
		if err != nil {
			continue
		}
		if resp.StatusCode > 299 {
			err := httputilmore.PrintResponse(resp, true)
			if err != nil {
				return err
			}
			fmt.Printf("%v\n", resp.StatusCode)
			continue
		}
	}
	return nil
}

func (sc *SalesforceClient) DeleteAccount(id string) (*http.Response, error) {
	//apiURLString := fmt.Sprintf("/services/data/v40.0/sobjects/%v/%v", "Account", id)
	//apiURL := sc.URLBuilder.Build(apiURLString)
	apiURL := sc.URLBuilder.BuildSobjectURL("Account", id)

	req, err := http.NewRequest("DELETE", apiURL.String(), nil)
	if err != nil {
		return &http.Response{}, err
	}
	return sc.ClientMore.Client.Do(req)
}

func (sc *SalesforceClient) GetContactsAll() (sobjects.ContactSet, error) {
	resp, err := sc.ExecSOQL("SELECT id, name, email FROM contact")
	if err != nil {
		return sobjects.ContactSet{}, err
	}
	return sobjects.NewContactSetFromJSONResponse(resp)
}

func (sc *SalesforceClient) DeleteContactsAll() error {
	set, err := sc.GetContactsAll()
	if err != nil {
		return err
	}
	for _, contact := range set.Records {
		resp, err := sc.DeleteContact(contact.Id)
		if err != nil {
			return err
		}
		if resp.StatusCode > 299 {
			err := httputilmore.PrintResponse(resp, true)
			if err != nil {
				return err
			}
			fmt.Printf("%v\n", resp.StatusCode)
			continue
		}
	}
	return nil
}

func (sc *SalesforceClient) DeleteContact(id string) (*http.Response, error) {
	apiURL := sc.URLBuilder.BuildSobjectURL("Contact", id)

	req, err := http.NewRequest("DELETE", apiURL.String(), nil)
	if err != nil {
		return &http.Response{}, err
	}
	return sc.ClientMore.Client.Do(req)
}

func (sc *SalesforceClient) UserInfo() (User, error) {
	apiURL := "https://login.salesforce.com/services/oauth2/userinfo"
	user := User{}

	req, err := http.NewRequest("GETs", apiURL, nil)
	if err != nil {
		return user, err
	}

	resp, err := sc.ClientMore.Client.Do(req)
	if err != nil {
		return user, err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(bytes, &user)
	return user, err
}

type User struct {
	UserID         string `json:"user_id,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
}

type URLBuilder struct {
	BaseURL url.URL
	Version string
}

func NewURLBuilder(instanceName string) URLBuilder {
	return URLBuilder{
		BaseURL: url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf(HostFormat, instanceName),
		},
		Version: "v40.0",
	}
}

func (b *URLBuilder) Build(path string) url.URL {
	u := b.BaseURL
	u.Path = path
	return u
}

func (b *URLBuilder) BuildSobjectURL(parts ...string) url.URL {
	partsString := path.Join(parts...)
	apiURLString := fmt.Sprintf("/services/data/%v/sobjects/%v", b.Version, partsString)
	return b.Build(apiURLString)
}
