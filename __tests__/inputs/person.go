package client

// Contributor represents a person that has made commits to a bundle
type Contributor struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhotoURL string `json:"photoURL"`
	Rank     int    `json:"rank"`
}
