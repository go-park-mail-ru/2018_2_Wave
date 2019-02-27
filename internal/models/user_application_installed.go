package models

// easyjson:json
type UserApplicationInstalled struct {
	Link          string `json:"link"`
	Url           string `json:"url"`
	Name          string `json:"name"`
	NameDE        string `json:"name_de"`
	NameRU        string `json:"name_ru"`
	Image         string `json:"image"`
	About         string `json:"about"`
	AboutDE       string `json:"about_de"`
	AboutRU       string `json:"about_ru"`
	Installations int    `json:"installs"`
	Price         int    `json:"price"`
	Category      string `json:"category"`
	Installed     bool   `json:"installed"`
}

// easyjson:json
type UserApplicationsInstalled struct {
	UserApplications []UserApplication `json:"user_apps_installed"`
}
