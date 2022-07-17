package domain

// Feed is our domain representation of a feed
type Feed struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Articles    []Article `json:"articles"`
}

// Article is our domain representation of an article
type Article struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
	Image       Image  `json:"image,omitempty"`
	URL         string `json:"url,omitempty"`
}

// Image is our domain representation of an image
type Image struct {
	URL   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
}
