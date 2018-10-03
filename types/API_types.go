package types

// APIUser - public API
// easyjson:json walhalla:
type APIUser struct {
	Username string `json:"username"`
	Password string `json:"password" walhalla:"min:6, max:100"`
}

// APISignUp - public API
// easyjson:json walhalla:
type APISignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   []byte `json:"avatar"`
}

// APIProfile - public API
// easyjson:json walhalla:
type APIProfile struct {
	Username  string `json:"username"`
	AvatarURI string `json:"avatarSource"`
	Score     int    `json:"score"`
}

// APIEditProfile - public API
// easyjson:json walhalla:
type APIEditProfile struct {
	Username    string `json:"username"`
	CurPassword string `json:"curPassword"`
	NewPassword string `json:"newPassword"`
	Avatar      []byte `json:"avatar"`
}

// APILeaderboardRow - public API
// easyjson:json
type APILeaderboardRow struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}

// APILeaderboard - public API
// easyjson:json
type APILeaderboard struct {
	Users []APILeaderboardRow `json:"users"`
	Total int                 `json:"total"`
}

func (pf *APISignUp) AsAPIUser() APIUser {
	return APIUser{
		Username: pf.Username,
		Password: pf.Password,
	}
}
