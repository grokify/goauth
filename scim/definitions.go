package scim

import (
	"fmt"
	"strings"

	"github.com/grokify/simplego/type/stringsutil"
)

// User is an object from the full user representation
// specified in the SCIM schema:
// http://www.simplecloud.info/specs/draft-scim-core-schema-01.html#anchor7
// https://tools.ietf.org/html/rfc7643
type User struct {
	Schemas           []string  `json:"schemas,omitempty"`
	ID                string    `json:"id,omitempty"`
	ExternalID        string    `json:"externalId,omitempty"`
	UserName          string    `json:"userName,omitempty"`
	Name              Name      `json:"name,omitempty"`
	DisplayName       string    `json:"displayName,omitempty"`
	NickName          string    `json:"nickName,omitempty"`
	ProfileURL        string    `json:"profileUrl,omitempty"`
	PhoneNumbers      []Item    `json:"phoneNumbers,omitempty"`
	Emails            []Item    `json:"emails,omitempty"`
	UserType          string    `json:"userType,omitempty"`
	Title             string    `json:"title,omitempty"`
	PreferredLanguage string    `json:"preferredLanguage,omitempty"`
	Locale            string    `json:"locale,omitempty"`
	Timezone          string    `json:"timezone,omitempty"`
	Active            bool      `json:"active,omitempty"`
	Password          string    `json:"password,omitempty"`
	Addresses         []Address `json:"addresses,omitempty"`
}

func NewUser() User {
	return User{
		Schemas:      []string{},
		PhoneNumbers: []Item{},
		Emails:       []Item{},
		Addresses:    []Address{},
	}
}

func (user *User) DisplayNameAny() string {
	name := strings.TrimSpace(user.DisplayName)
	if len(name) > 0 {
		return name
	}
	name = strings.TrimSpace(user.Name.Formatted)
	if len(name) > 0 {
		return name
	}
	names := stringsutil.SliceCondenseSpace([]string{
		user.Name.GivenName,
		user.Name.MiddleName,
		user.Name.FamilyName}, false, false)
	name = strings.Join(names, " ")
	if len(name) > 0 {
		return name
	}
	name = strings.TrimSpace(user.NickName)
	if len(name) > 0 {
		return name
	}
	return strings.TrimSpace(user.UserName)
}

// AddEmail adds a canonical email address to the user.
// it lowercases and trims preceding and trailing spaces
// from the email address.
func (user *User) AddEmail(emailAddr string, isPrimary bool) error {
	emailAddrCanonical := strings.ToLower(strings.TrimSpace(emailAddr))
	if len(emailAddr) < 1 {
		return fmt.Errorf("Invalid Email Address: %v", emailAddr)
	}
	email := Item{
		Value:   emailAddrCanonical,
		Primary: isPrimary}
	user.Emails = append(user.Emails, email)
	return nil
}

// GetPrimaryEmail returns an email address or empty string.
// It prioritizes primary email addresses and then falls
// back to non-primary email address.
func (user *User) EmailAddress() string {
	one := GetOneItem(user.Emails)
	return one.Value
	/*
		firstSecondaryEmail := ""
		for _, em := range user.Emails {
			if em.Primary && len(strings.TrimSpace(em.Value)) > 0 {
				return strings.TrimSpace(em.Value)
			} else if len(firstSecondaryEmail) == 0 && len(strings.TrimSpace(em.Value)) > 0 {
				firstSecondaryEmail = strings.TrimSpace(em.Value)
			}
		}
		return firstSecondaryEmail*/
}

func (user *User) PhoneNumber() string {
	one := GetOneItem(user.PhoneNumbers)
	return one.Value
}

func GetOneItem(items []Item) Item {
	if len(items) == 0 {
		return Item{}
	}
	havePrimary := false
	haveSecondary := false
	primary := Item{}
	secondary := Item{}
	for _, it := range items {
		it.Value = strings.TrimSpace(it.Value)
		if it.Primary {
			if len(it.Value) > 0 {
				return it
			}
			primary = it
			havePrimary = true
		} else {
			if haveSecondary && len(secondary.Value) > 0 {
				continue
			} else if len(it.Value) > 0 || !haveSecondary {
				secondary = it
				haveSecondary = true
			}
		}
	}
	if havePrimary && len(primary.Value) > 0 {
		return primary
	} else if haveSecondary && len(secondary.Value) > 0 {
		return secondary
	} else if havePrimary {
		return primary
	}
	return secondary
}

// Name is the SCIM user name struct.
type Name struct {
	Formatted       string `json:"formatted,omitempty"`
	FamilyName      string `json:"familyName,omitempty"`
	GivenName       string `json:"givenName,omitempty"`
	MiddleName      string `json:"middleName,omitempty"`
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
}

// Item is a SCIM struct used for email and phone numbers.
type Item struct {
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

// Address is a SCIM struct used for street and mailing addresses.
type Address struct {
	Type          string `json:"Type,omitempty"`
	StreetAddress string `json:"StreetAddress,omitempty"`
	Locality      string `json:"Locality,omitempty"`
	Region        string `json:"Region,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Country       string `json:"country,omitempty"`
	Formatted     string `json:"formatted,omitempty"`
	Primary       bool   `json:"primary,omitempty"`
}

func (addr Address) InflateStreetAddress(force bool) {
	addr.Formatted = strings.TrimSpace(addr.Formatted)
	addr.Locality = strings.TrimSpace(addr.Locality)
	addr.Region = strings.TrimSpace(addr.Region)
	addr.PostalCode = strings.TrimSpace(addr.PostalCode)
	addr.Country = strings.TrimSpace(addr.Country)
	if len(addr.Formatted) > 0 && !force {
		return
	}
	lines := []string{}
	if len(addr.StreetAddress) > 0 {
		lines = append(lines, addr.StreetAddress)
	}
	parts := []string{}
	if len(addr.Locality) > 0 {
		parts = append(parts, addr.Locality+",")
	}
	if len(addr.Region) > 0 {
		parts = append(parts, addr.Region)
	}
	if len(addr.PostalCode) > 0 {
		parts = append(parts, addr.PostalCode)
	}
	if len(addr.Country) > 0 {
		parts = append(parts, addr.Country)
	}
	if len(parts) > 0 {
		lines = append(lines, strings.Join(parts, " "))
	}
	addr.Formatted = strings.Join(lines, "\n")
}
