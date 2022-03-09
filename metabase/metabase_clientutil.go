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

func NewClientUtil(baseUrl, username, password string, tlsSkipVerify bool) (*ClientUtil, error) {
	httpClient, _, err := NewClientPassword(baseUrl, username, password, tlsSkipVerify)
	if err != nil {
		return nil, err
	}
	return &ClientUtil{
		HTTPClient: httpClient,
		BaseURL:    baseUrl,
	}, nil
}

func (cu *ClientUtil) GetStoreQuestionData(cardId int, filename string, perm os.FileMode) ([]byte, error) {
	data, err := cu.GetQuestionData(cardId)
	if err != nil {
		return data, err
	}
	return data, os.WriteFile(filename, data, perm)
}

func (cu *ClientUtil) GetQuestionData(cardId int) ([]byte, error) {
	cardUrl := cu.BuildMetabaseCardAPI(cardId, "json")

	req, err := http.NewRequest(http.MethodPost, cardUrl, nil)
	if err != nil {
		return []byte(""), err
	}
	resp, err := cu.HTTPClient.Do(req)
	if err != nil {
		return []byte(""), err
	} else if resp.StatusCode >= 300 {
		return []byte(""), fmt.Errorf("Metabase API Error Status: %v", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (cu *ClientUtil) BuildMetabaseCardAPI(cardId int, format string) string {
	relUrl := fmt.Sprintf("api/card/%v/query/%s", cardId, format)
	return urlutil.JoinAbsolute(cu.BaseURL, relUrl)
}

type QuestionsToSlug struct {
	QuestionMap map[string]int
}

func RetrieveQuestions(cu ClientUtil, q2s QuestionsToSlug, dir string) (map[string][]byte, error) {
	dt := time.Now()
	dt8 := dt.Format(timeutil.DT8)
	output := map[string][]byte{}
	for name, cardId := range q2s.QuestionMap {
		filename := fmt.Sprintf("data_%v_%v.json", dt8, name)
		data, err := cu.GetStoreQuestionData(cardId, filename, 0644)
		if err != nil {
			return output, errorsutil.Wrap(err, fmt.Sprintf("Error Retrieving Card #[%v] Name[%v]", cardId, name))
		}
		output[name] = data

		zlog.Info().
			Str("filename", filename).
			Str("data", string(data))
	}
	return output, nil
}

func (apiutil *ClientUtil) GetCurrentUser() (User, *http.Response, error) {
	user := User{}
	apiURL := urlutil.JoinAbsolute(apiutil.BaseURL, CurrentUserURLPath)
	resp, err := apiutil.HTTPClient.Get(apiURL)
	if err != nil {
		return user, nil, err
	} else if resp.StatusCode >= 300 {
		return user, resp, fmt.Errorf("MB_API_ERROR_STATUS_CODE [%v]", resp.StatusCode)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return user, resp, err
	}
	err = json.Unmarshal(bytes, &user)
	return user, resp, err
}

func (apiutil *ClientUtil) GetSCIMUser() (scim.User, error) {
	scimUser := scim.User{}
	mbUser, _, err := apiutil.GetCurrentUser()
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
