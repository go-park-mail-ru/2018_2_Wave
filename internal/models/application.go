package models

// easyjson:json
type Application struct {
	Name          string `json:"name"`
	Cover         string `json:"cover"`
	Description   string `json:"description"`
	Installations int    `json:"installations"`
	Price         int    `json:"price"`
	Year          string `json:"year"`
}

// easyjson:json
type Applications struct {
	Applications []Application `json:"apps"`
}
