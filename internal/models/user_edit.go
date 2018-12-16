package models

// easyjson:json
type UserEdit struct {
	Username string `json:"username", omitempty`
	Password string `json:"password", omitempty`
	Avatar   string `json:"avatar", omitempty`
}
