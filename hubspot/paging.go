package hubspot

type ResponsePaging struct {
	Paging    Paging `json:"paging,omitempty"`     // V3 List Contacts API
	VIDOffset int    `json:"vid-offset,omitempty"` // V1 List Contacts API
}

type Paging struct {
	Next PagingNext `json:"next,omitempty"`
}

type PagingNext struct {
	After string `json:"after,omitempty"`
	Link  string `json:"link,omitempty"`
}
