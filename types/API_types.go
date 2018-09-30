package types

//easyjson:json
type APIUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//easyjson:json
type APISignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   []byte `json:"avatar"`
}

//easyjson:json
type APIProfile struct {
	Username  string `json:"username"`
	AvatarURI string `json:"avatar"`
}

//easyjson:json
type APIEditProfile struct {
	Password string `json:"password"`
	Avatar   []byte `json:"avatar"`
}

func (pf *APISignUp) AsAPIUser() APIUser {
	return APIUser{
		Username: pf.Username,
		Password: pf.Password,
	}
}
