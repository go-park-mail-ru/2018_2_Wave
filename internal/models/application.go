package models

// easyjson:json
type Application struct {
	Link          string `json:"link"`
	Url           string `json:"url"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	About         string `json:"about"`
	Installations int    `json:"installs"`
	Price         int    `json:"price"`
	Category      string `json:"category"`
}

// easyjson:json
type Applications struct {
	Applications []Application `json:"apps"`
}
