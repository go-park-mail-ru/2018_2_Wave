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
}

// APIEditProfile - public API
// easyjson:json walhalla:
type APIEditProfile struct {
	Username    string `json:"username"`
	CurPassword string `json:"curPassword"`
	NewPassword string `json:"newPassword"`
	Avatar      []byte `json:"avatar"`
}

func (pf *APISignUp) AsAPIUser() APIUser {
	return APIUser{
		Username: pf.Username,
		Password: pf.Password,
	}
}
