package models

// easyjson:json
type UserApplication struct {
	Link          string `json:"link"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	About         string `json:"about"`
	Installations int    `json:"installs"`
	Price         int    `json:"price"`
	Category      string `json:"category"`
	TimeTotal     int    `json:"time_total"`
}

// easyjson:json
type UserApplications struct {
	UserApplications []UserApplication `json:"user_apps"`
}
