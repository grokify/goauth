package metabase

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/grokify/gotilla/net/urlutil"
	tu "github.com/grokify/gotilla/time/timeutil"
	log "github.com/sirupsen/logrus"
)

type ClientUtil struct {
	HTTPClient *http.Client
	BaseURL    string
}

func (cu *ClientUtil) GetStoreQuestionData(cardId int, filename string, perm os.FileMode) ([]byte, error) {
	data, err := cu.GetQuestionData(cardId)
	if err != nil {
		return data, err
	}
	return data, ioutil.WriteFile(filename, data, perm)
}

func (cu *ClientUtil) GetQuestionData(cardId int) ([]byte, error) {
	cardUrl := cu.BuildMetabaseCardAPI(cardId, "json")

	req, err := http.NewRequest(http.MethodPost, cardUrl, nil)
	resp, err := cu.HTTPClient.Do(req)
	if err != nil {
		return []byte(""), err
	} else if resp.StatusCode >= 300 {
		return []byte(""), fmt.Errorf("Metabase API Error Status: %v", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (cu *ClientUtil) BuildMetabaseCardAPI(cardId int, format string) string {
	relUrl := fmt.Sprintf("api/card/%v/query/%s", cardId, format)
	return urlutil.JoinAbsolute(cu.BaseURL, relUrl)
}

type QuestionsToSlug struct {
	QuestionMap map[string]int
}

func RetrieveQuestions(cu ClientUtil, q2s QuestionsToSlug) (map[string][]byte, error) {
	dt := time.Now()
	dt8 := dt.Format(tu.DT8)
	output := map[string][]byte{}
	for name, cardId := range q2s.QuestionMap {
		filename := fmt.Sprintf("data_%v_%v.json", dt8, name)
		data, err := cu.GetStoreQuestionData(cardId, filename, 0644)
		if err != nil {
			return output, err
		}
		output[name] = data

		log.Info("===")
		log.Info(filename)
		log.Info(string(data))
	}
	return output, nil
}
