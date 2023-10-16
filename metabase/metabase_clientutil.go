package metabase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grokify/goauth/scim"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/time/timeutil"

	zlog "github.com/rs/zerolog/log"
)

const (
	CurrentUserURLPath = "/api/user/current"
)

type ClientUtil struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewClientUtil(baseURL, username, password string, tlsSkipVerify bool) (*ClientUtil, error) {
	httpClient, _, err := NewClientPassword(baseURL, username, password, tlsSkipVerify)
	if err != nil {
		return nil, err
	}
	return &ClientUtil{
		HTTPClient: httpClient,
		BaseURL:    baseURL,
	}, nil
}

func (cu *ClientUtil) GetStoreQuestionData(cardID int, filename string, perm os.FileMode) ([]byte, error) {
	data, err := cu.GetQuestionData(cardID)
	if err != nil {
		return data, err
	}
	return data, os.WriteFile(filename, data, perm)
}

func (cu *ClientUtil) GetQuestionData(cardID int) ([]byte, error) {
	cardURL := cu.BuildMetabaseCardAPI(cardID, "json")

	req, err := http.NewRequest(http.MethodPost, cardURL, nil)
	if err != nil {
		return []byte(""), err
	}
	resp, err := cu.HTTPClient.Do(req)
	if err != nil {
		return []byte(""), err
	} else if resp.StatusCode >= 300 {
		return []byte(""), fmt.Errorf("metabase API Error Status: %v", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (cu *ClientUtil) BuildMetabaseCardAPI(cardID int, format string) string {
	relURL := fmt.Sprintf("api/card/%d/query/%s", cardID, format)
	return urlutil.JoinAbsolute(cu.BaseURL, relURL)
}

type QuestionsToSlug struct {
	QuestionMap map[string]int
}

func RetrieveQuestions(cu ClientUtil, q2s QuestionsToSlug, dir string) (map[string][]byte, error) {
	dt := time.Now()
	dt8 := dt.Format(timeutil.DT8)
	output := map[string][]byte{}
	for name, cardID := range q2s.QuestionMap {
		filename := fmt.Sprintf("data_%v_%v.json", dt8, name)
		data, err := cu.GetStoreQuestionData(cardID, filename, 0600)
		if err != nil {
			return output, errorsutil.Wrap(err, fmt.Sprintf("error retrieving card #(%d) name(%s)", cardID, name))
		}
		output[name] = data

		zlog.Info().
			Str("filename", filename).
			Str("data", string(data))
	}
	return output, nil
}

func (cu *ClientUtil) GetCurrentUser() (User, *http.Response, error) {
	user := User{}
	apiURL := urlutil.JoinAbsolute(cu.BaseURL, CurrentUserURLPath)
	resp, err := cu.HTTPClient.Get(apiURL)
	if err != nil {
		return user, nil, err
	} else if resp.StatusCode >= 300 {
		return user, resp, fmt.Errorf("metabase api error status code (%d)", resp.StatusCode)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return user, resp, err
	}
	err = json.Unmarshal(bytes, &user)
	return user, resp, err
}

func (cu *ClientUtil) GetSCIMUser() (scim.User, error) {
	scimUser := scim.User{}
	mbUser, _, err := cu.GetCurrentUser()
	if err != nil {
		return scimUser, err
	}
	err = scimUser.AddEmail(mbUser.Email, true)
	if err != nil {
		return scimUser, err
	}
	scimUser.Name = scim.Name{
		GivenName:  strings.TrimSpace(mbUser.FirstName),
		FamilyName: strings.TrimSpace(mbUser.LastName),
		Formatted:  strings.TrimSpace(mbUser.CommonName)}
	return scimUser, nil
}

type User struct {
	Email       string    `json:"email,omitempty"`
	LdapAuth    bool      `json:"ldap_auth,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastLogin   time.Time `json:"last_login,omitempty"`
	IsActive    bool      `json:"is_active,omitempty"`
	IsQbnewb    bool      `json:"is_qbnewb,omitempty"`
	IsSuperuser bool      `json:"is_superuser,omitempty"`
	ID          int       `json:"id,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	DateJoined  time.Time `json:"date_joined,omitempty"`
	CommonName  string    `json:"common_name,omitempty"`
	GoogleAuth  bool      `json:"google_auth,omitempty"`
}
