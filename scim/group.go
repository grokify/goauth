package scim

type Group struct {
	Display string `json:"display,omitempty"`
	Ref     string `json:"$ref,omitempty"`
	Value   string `json:"value,omitempty"`
}
