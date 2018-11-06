package models

// easyjson:json
type UserExtended struct {
	Username string `json:"username"`
	Score    string `json:"score"`
	Avatar   []byte `json:"avatar"`
}
