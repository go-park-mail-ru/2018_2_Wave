package models

// easyjson:json
type UserApplicationInstalled struct {
	Link          string `json:"link"`
	Url           string `json:"url"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	About         string `json:"about"`
	Installations int    `json:"installs"`
	Price         int    `json:"price"`
	Category      string `json:"category"`
	Installed     bool   `json:"installed"`
}

// easyjson:json
type UserApplicationsInstalled struct {
	UserApplications []UserApplication `json:"user_apps_installed"`
}
