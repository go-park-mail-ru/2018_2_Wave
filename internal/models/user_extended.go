package models

// easyjson:json
type UserExtended struct {
	Username string `json:"username"`
	Score    string `json:"score"`
	Avatar   string `json:"avatar"`
	Locale   string `json:"locale"`
}
