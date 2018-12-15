package models

// easyjson:json
type Application struct {
	Name          string `json:"name"`
	Thumbnail     string `json:"thumbnail"`
	Description   string `json:"description"`
	Installations int    `json:"installations"`
	Price         int    `json:"price"`
}

// easyjson:json
type Applications struct {
	Applications []Application `json:"apps"`
}
