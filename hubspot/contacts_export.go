package hubspot

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/grokify/goauth/endpoints"
	"github.com/grokify/gohttp/httpsimple"
	"github.com/grokify/mogo/encoding/jsonutil"
)

const (
	ContactsListAPIPathV3        = "/crm/v3/objects/contacts"
	ContactsListAPIPathV1        = "/contacts/v1/lists/all/contacts/all"
	ParamV3Limit                 = "limit"
	ParamV3After                 = "after"
	ParamV3Archived              = "archived"
	ParamV3Properties            = "properties"
	ParamV3PropertiesWithHistory = "propertiesWithHistory"
	ParamV3Associations          = "associations"
	ParamV1Count                 = "count"
	ParamV1VIDOffset             = "vidOffset"
	LimitMax                     = 100
)

type ContactsListV3Opts struct {
	Limit                 int      `url:"limit"`
	After                 string   `url:"after"`
	Properties            []string `url:"properties"`
	PropertiesWithHistory []string `url:"propertiesWithHistory"`
	Associations          []string `url:"associations"`
	Archived              bool     `url:"archived"`
}

// Query generates query string values for the Contacts V3 API.
func (opts *ContactsListV3Opts) Query() url.Values {
	qry := url.Values{}
	if opts.Limit > 0 {
		qry.Add(ParamV3Limit, strconv.Itoa(opts.Limit))
	}
	if len(opts.After) > 0 {
		qry.Add(ParamV3After, opts.After)
	}
	for _, pval := range opts.Properties {
		qry.Add(ParamV3Properties, pval)
	}
	for _, pval := range opts.PropertiesWithHistory {
		qry.Add(ParamV3PropertiesWithHistory, pval)
	}
	for _, assoc := range opts.Associations {
		qry.Add(ParamV3Associations, assoc)
	}
	if opts.Archived {
		qry.Add(ParamV3Archived, "true")
	}
	return qry
}

type ContactsListV1Opts struct {
	Count     int `url:"count"`
	VIDOffset int `url:"vidOffset"`
}

// Query generates query string values for the Contacts V3 API.
func (opts *ContactsListV1Opts) Query() url.Values {
	qry := url.Values{}
	if opts.Count > 0 {
		qry.Add(ParamV1Count, strconv.Itoa(opts.Count))
	}
	if opts.VIDOffset > 0 {
		qry.Add(ParamV1VIDOffset, strconv.Itoa(opts.VIDOffset))
	}
	return qry
}

func ContactsV3ExportWriteFiles(client *http.Client, fileprefix string, opts *ContactsListV3Opts) error {
	if len(fileprefix) == 0 {
		fileprefix = "hubspot_contacts_v3"
	}
	if opts != nil && opts.Limit > LimitMax || opts.Limit < 1 {
		return errors.New("invalid limit - must be between 1 and 100 inclusive")
	}
	sclient := httpsimple.SimpleClient{
		BaseURL:    endpoints.HubspotServerURL,
		HTTPClient: client}

	sreq := httpsimple.SimpleRequest{
		Method: http.MethodGet,
		URL:    ContactsListAPIPathV3}
	if opts != nil {
		qrys := opts.Query().Encode()
		if len(qrys) > 0 {
			sreq.URL = ContactsListAPIPathV3 + "?" + qrys
		}
	}
	pgNum := 1
	for {
		resp, err := sclient.Do(sreq)
		if err != nil {
			return err
		}
		bodyPretty, err := jsonutil.IndentReader(resp.Body, "", "  ")
		if err != nil {
			return err
		}
		filename := fileprefix + "_page-" + strconv.Itoa(pgNum) + ".json"
		err = os.WriteFile(filename, bodyPretty, 0600)
		if err != nil {
			return err
		}
		var pagingResp ResponsePaging
		err = json.Unmarshal(bodyPretty, &pagingResp)
		if err != nil {
			return err
		}
		if len(pagingResp.Paging.Next.After) > 0 {
			if opts == nil {
				opts = &ContactsListV3Opts{}
			}
			opts.After = pagingResp.Paging.Next.After
			sreq.URL = ContactsListAPIPathV3 + "?" + opts.Query().Encode()
		} else {
			break
		}
		pgNum++
	}
	return nil
}

func ContactsV1ExportWriteFiles(client *http.Client, fileprefix string, opts *ContactsListV1Opts) error {
	if len(fileprefix) == 0 {
		fileprefix = "hubspot_contacts_v1"
	}
	if opts != nil && opts.Count > LimitMax || opts.Count < 1 {
		return errors.New("invalid count - must be between 1 and 100 inclusive")
	}
	sclient := httpsimple.SimpleClient{
		BaseURL:    endpoints.HubspotServerURL,
		HTTPClient: client}

	sreq := httpsimple.SimpleRequest{
		Method: http.MethodGet,
		URL:    ContactsListAPIPathV3}
	if opts != nil {
		qrys := opts.Query().Encode()
		if len(qrys) > 0 {
			sreq.URL = ContactsListAPIPathV1 + "?" + qrys
		}
	}
	pgNum := 1
	for {
		resp, err := sclient.Do(sreq)
		if err != nil {
			return err
		}
		bodyPretty, err := jsonutil.IndentReader(resp.Body, "", "  ")
		if err != nil {
			return err
		}
		filename := fileprefix + "_page-" + strconv.Itoa(pgNum) + ".json"
		err = os.WriteFile(filename, bodyPretty, 0600)
		if err != nil {
			return err
		}
		var pagingResp ResponsePaging
		err = json.Unmarshal(bodyPretty, &pagingResp)
		if err != nil {
			return err
		}
		if pagingResp.VIDOffset > 0 {
			if opts == nil {
				opts = &ContactsListV1Opts{}
			}
			opts.VIDOffset = pagingResp.VIDOffset
			sreq.URL = ContactsListAPIPathV1 + "?" + opts.Query().Encode()
		} else {
			break
		}
		pgNum++
	}
	return nil
}
