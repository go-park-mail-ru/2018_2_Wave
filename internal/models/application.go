package models

// easyjson:json
type Application struct {
	Link          string `json:"link"`
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
