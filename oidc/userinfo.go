package oidc

import "strings"

// UserInfo represents an OpenID Connect (OIDC) UserInfo object as described here:
// https://openid.net/specs/openid-connect-core-1_0.html#UserInfoResponse .
type UserInfo struct {
	Audience          string `json:"aud"`
	Issuer            string `json:"iss"`
	Subject           string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	FamilyName        string `json:"family_name"`
	GivenName         string `json:"given_name"`
	Name              string `json:"name"`
	Picture           string `json:"picture"` // URL to picture
	PreferredUsername string `json:"preferred_username"`
	Profile           string `json:"profile"`
}

func (u *UserInfo) AddEmail(email string, addPreferredUsername bool) {
	u.Email = strings.TrimSpace(email)
	pu := ""
	em := strings.TrimSpace(u.Email)
	emParts := strings.Split(em, "@")
	if len(emParts) == 2 {
		pu = strings.TrimSpace(emParts[0])
	}
	u.PreferredUsername = pu
}

func (u *UserInfo) AddName(name string, addFirstAndLastName bool) {
	u.Name = strings.TrimSpace(name)
	dnParts := strings.Fields(u.Name)
	if len(dnParts) == 2 {
		u.GivenName = strings.TrimSpace(dnParts[0])
		u.FamilyName = strings.TrimSpace(dnParts[1])
	}
}
