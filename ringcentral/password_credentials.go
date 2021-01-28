package ringcentral

import (
	"net/url"
	"os"
	"strconv"
	"strings"
)

type PasswordCredentials struct {
	GrantType            string `url:"grant_type"`
	AccessTokenTTL       int64  `url:"access_token_ttl"`
	RefreshTokenTTL      int64  `url:"refresh_token_ttl"`
	Username             string `json:"username" url:"username"`
	Extension            string `json:"extension" url:"extension"`
	Password             string `json:"password" url:"password"`
	EndpointId           string `url:"endpoint_id"`
	EngageVoiceAccountId int64  `json:"engageVoiceAccountId"`
}

func NewPasswordCredentialsEnv() PasswordCredentials {
	return PasswordCredentials{
		Username:  os.Getenv(EnvUsername),
		Extension: os.Getenv(EnvExtension),
		Password:  os.Getenv(EnvPassword)}
}

func (pw *PasswordCredentials) URLValues() url.Values {
	v := url.Values{
		"grant_type": {"password"},
		"username":   {pw.Username},
		"password":   {pw.Password}}
	if pw.AccessTokenTTL != 0 {
		v.Set("access_token_ttl", strconv.Itoa(int(pw.AccessTokenTTL)))
	}
	if pw.RefreshTokenTTL != 0 {
		v.Set("refresh_token_ttl", strconv.Itoa(int(pw.RefreshTokenTTL)))
	}
	if len(pw.Extension) > 0 {
		v.Set("extension", pw.Extension)
	}
	if len(pw.EndpointId) > 0 {
		v.Set("endpoint_id", pw.EndpointId)
	}
	return v
}

func (uc *PasswordCredentials) UsernameSimple() string {
	if len(strings.TrimSpace(uc.Extension)) > 0 {
		return strings.Join([]string{uc.Username, uc.Extension}, "*")
	}
	return uc.Username
}
