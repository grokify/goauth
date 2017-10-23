package salesforce

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/go-salesforce/sobjects"
	"github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/net/urlutil"
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

func (sc *SalesforceClient) GetServicesData() (*http.Response, error) {
	apiURL := sc.URLBuilder.Build("services/data")
	return sc.ClientMore.Client.Get(apiURL.String())
}

func (sc *SalesforceClient) CreateContact(contact interface{}) (*http.Response, error) {
	apiURL := sc.URLBuilder.Build("/services/data/v40.0/sobjects/Contact/")
	return sc.ClientMore.PostToJSON(apiURL.String(), contact)
}

func (sc *SalesforceClient) ExecSOQL(soql string) (*http.Response, error) {
	//curl https://yourInstance.salesforce.com/services/data/v20.0/query/?q=SELECT+name+from+Account -H "Authorization: Bearer token"
	apiURL := sc.URLBuilder.Build("/services/data/v40.0/query/")
	soql = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(soql), "+")
	q := url.Values{}
	q.Add("q", soql)
	apiURLString := urlutil.BuildURL(apiURL.String(), q)
	apiURLString = regexp.MustCompile(`\%2B`).ReplaceAllString(strings.TrimSpace(apiURLString), "+")
	return sc.ClientMore.Client.Get(apiURLString)
}

func (sc *SalesforceClient) DeleteContactsAll() error {
	resp, err := sc.ExecSOQL("select id,name,email from contact")
	if err != nil {
		return err
	}
	set, err := sobjects.NewContactSetFromJSONResponse(resp)
	if err != nil {
		return err
	}
	for _, contact := range set.Records {
		resp, err := sc.DeleteContact(contact.Id)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode > 299 {
			httputilmore.PrintResponse(resp, true)
			fmt.Printf("%v\n", resp.StatusCode)
			panic("Z")
		}
	}
	return nil
}

func (sc *SalesforceClient) DeleteContact(id string) (*http.Response, error) {
	apiURLString := fmt.Sprintf("/services/data/v40.0/sobjects/%v/%v", "Contact", id)
	apiURL := sc.URLBuilder.Build(apiURLString)

	req, err := http.NewRequest("DELETE", apiURL.String(), nil)
	if err != nil {
		return &http.Response{}, err
	}
	return sc.ClientMore.Client.Do(req)
}

type URLBuilder struct {
	BaseURL url.URL
}

func NewURLBuilder(instanceName string) URLBuilder {
	return URLBuilder{BaseURL: url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf(HostFormat, instanceName)}}
}

func (b *URLBuilder) Build(path string) url.URL {
	u := b.BaseURL
	u.Path = path
	return u
}
