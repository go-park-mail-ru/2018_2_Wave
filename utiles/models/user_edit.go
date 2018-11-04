package models

// easyjson:json
type UserEdit struct {
	NewUsername string `json:"newUsername"`
	NewPassword string `json:"newPassword"`
	NewAvatar   string `json:"newAvatar"`
}
