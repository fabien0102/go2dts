package client

// PageMeta contains the next and previous page tokens that clients can use
// to request pages of data from api.  Next will be empty if no additional pages
// exists, similarly for Prev
type PageMeta struct {
	Next  int `json:"next,omitempty"`
	Prev  int `json:"prev,omitempty"`
	Last  int `json:"last"`
	Count int `json:"count"`
}
