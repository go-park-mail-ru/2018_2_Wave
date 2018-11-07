package models

// easyjson:json
type Leaders struct {
	Users []UserScore `json:"users"`
	Total int         `json:"total"`
}
