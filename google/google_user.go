package google

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	GoogleApiUrlUserinfo = "https://openidconnect.googleapis.com/v1/userinfo"
)

type GoogleUserinfoOpenIdConnectV2 struct {
	Sub          string `json:"sub,omitempty"`
	Name         string `json:"name,omitempty"`
	GivenName    string `json:"given_name,omitempty"`
	FamilyName   string `json:"family_name,omitempty"`
	Picture      string `json:"picture,omitempty"`
	Email        string `json:"email,omitempty"`
	EmailVerfied bool   `json:"email_verified,omitempty"`
	Locale       string `json:"locale,omitempty"`
}

func GetMeInfo(bearerToken string) (GoogleUserinfoOpenIdConnectV2, error) {
	usr := GoogleUserinfoOpenIdConnectV2{}
	_, bodyBytes, err := HttpGetBearerTokenBody(GoogleApiUrlUserinfo, bearerToken)
	if err != nil {
		return usr, err
	}

	err = json.Unmarshal(bodyBytes, &usr)
	return usr, err
}

func HttpGetBearerTokenBody(url, token string) (*http.Response, []byte, error) {
	resp, err := HttpGetBearerToken(url, token)
	if err != nil {
		return resp, []byte(""), err
	}
	bytes, err := io.ReadAll(resp.Body)
	return resp, bytes, err
}

func HttpGetBearerToken(url, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{}
	return client.Do(req)
}
