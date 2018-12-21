package models

// easyjson:json
type UserApplication struct {
	Name          string `json:"name"`
	Cover         string `json:"cover"`
	Description   string `json:"description"`
	Installations int    `json:"installations"`
	Price         int    `json:"price"`
	Year          string `json:"year"`
	TimeTotal     int    `json:"time_total"`
}

// easyjson:json
type UserApplications struct {
	UserApplications []UserApplication `json:"user_apps"`
}
