package scim

// User is an object from the full user representation
// specified in the SCIM schema:
// http://www.simplecloud.info/specs/draft-scim-core-schema-01.html#anchor7
type User struct {
	Schemas           []string `json:"schemas,omitempty"`
	ExternalId        string   `json:"externalId,omitempty"`
	UserName          string   `json:"userName,omitempty"`
	Name              Name     `json:"name,omitempty"`
	DisplayName       string   `json:"displayName,omitempty"`
	NickName          string   `json:"nickName,omitempty"`
	ProfileUrl        string   `json:"profileUrl,omitempty"`
	PhoneNumbers      []Item   `json:"phoneNumbers,omitempty"`
	Emails            []Item   `json:"emails,omitempty"`
	UserType          string   `json:"userType,omitempty"`
	Title             string   `json:"title,omitempty"`
	PreferredLanguage string   `json:"preferredLanguage,omitempty"`
	Locale            string   `json:"locale,omitempty"`
	Timezone          string   `json:"timezone,omitempty"`
	Active            bool     `json:"active,omitempty"`
	Password          string   `json:"password,omitempty"`
}

type Name struct {
	Formatted       string `json:"formatted,omitempty"`
	FamilyName      string `json:"familyName,omitempty"`
	GivenName       string `json:"givenName,omitempty"`
	MiddleName      string `json:"middleName,omitempty"`
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
}

// Email and PhoneNumber
type Item struct {
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}
