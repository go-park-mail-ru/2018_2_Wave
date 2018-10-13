package types

//-----------------| signup |

// SignUp - public API
// easyjson:json walhalla:
type SignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   []byte `json:"avatar"`
}

//-----------------| profile |

// User - public API
// easyjson:json walhalla:
type User struct {
	Username string `json:"username"`
	Password string `json:"password" walhalla:"min:6, max:100"`
}

// Profile - public API
// easyjson:json walhalla:
type Profile struct {
	Username  string `json:"username"`
	AvatarURI string `json:"avatarSource"`
	Score     int    `json:"score"`
}

// EditProfile - public API
// easyjson:json walhalla:
type EditProfile struct {
<<<<<<< HEAD
	NewUsername	string `json:"newUsername"`
=======
	CurUsername string `json:"curUsername"`
	NewUsername string `json:"newUsername"`
>>>>>>> 3e70e414be2bd735a6a17bb8b78af28cebb1a4c9
	CurPassword string `json:"curPassword"`
	NewPassword string `json:"newPassword"`
	Avatar      string `json:"avatar"`
}

//-----------------| pagination |

// Pagination - public API
// easyjson:json walhalla:
type Pagination struct {
	Offset int `json:"offset" walhalla:"min:0"`
	Limit  int `json:"limit"  walhalla:"min:0"`
}

//-----------------| leaderboard |

// LeaderboardRow - public API
// easyjson:json
type LeaderboardRow struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}

// Leaderboard - public API
// easyjson:json
type Leaderboard struct {
	Users []LeaderboardRow `json:"users"`
	Total int              `json:"total"`
}

//-----------------| functions |

func (pf *SignUp) AsUser() User {
	return User{
		Username: pf.Username,
		Password: pf.Password,
	}
}

func Must(bytes []byte, err error) []byte {
	if err != nil {
		return []byte{}
	}
	return bytes
}
