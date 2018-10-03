package database

import (
	"Wave/common"
	"Wave/types"
	"strconv"
)

// DB is a database facade
type DB struct {
	mockTable    map[int]types.APIUser
	cookieToUser map[string]int
	avatarTable  map[int][]byte

	lastUID int
}

// New - create and initialise new database facade
func New() *DB {
	db := &DB{
		lastUID:      0,
		mockTable:    map[int]types.APIUser{},
		cookieToUser: map[string]int{},
		avatarTable:  map[int][]byte{},
	}
	return db
}

//*****************| Auth

// IsSignedUp - weather the user is signed up
func (db *DB) IsSignedUp(user types.APIUser) bool {
	for _, row := range db.mockTable {
		if row.Username == user.Username {
			return true
		}
	}
	return false
}

// SignUp the user.
// NOTE: each call creates new record with unique uid
func (db *DB) SignUp(profile types.APISignUp) (cookie string) {
	db.lastUID++

	uid := db.lastUID
	db.mockTable[uid] = profile.AsAPIUser()
	db.avatarTable[uid] = profile.Avatar
	return db.logIn(uid)
}

// IsLoggedIn - weather the user is logged in
func (db *DB) IsLoggedIn(cookie string) bool {
	_, ok := db.cookieToUser[cookie]
	return ok
}

// LogIn the user if the one is signe up
// NOTE: each call creates new session
func (db *DB) LogIn(user types.APIUser) (cookie string) {
	if uid, ok := db.getUID(user); ok {
		return db.logIn(uid)
	}
	return ""
}

// LogOut a user with the cookie if the one was logged in
func (db *DB) LogOut(cookie string) {
	if _, ok := db.cookieToUser[cookie]; ok {
		delete(db.cookieToUser, cookie)
	}
}

//*****************| Profile

// GetProfile returns a profile assigned to the cookie
func (db *DB) GetProfile(cookie string) (types.APIProfile, bool) {
	if uid, ok := db.cookieToUser[cookie]; ok {
		return types.APIProfile{
			Username:  db.mockTable[uid].Username,
			AvatarURI: "/img/avatars/" + strconv.Itoa(uid),
		}, true
	}
	return types.APIProfile{}, false
}

// GetAvatar returns avatar's data
func (db *DB) GetAvatar(uid int) ([]byte, bool) {
	data, ok := db.avatarTable[uid]
	return data, ok
}

// UpdateProfile updates profile
func (db *DB) UpdateProfile(cookie string, profile types.APIEditProfile) {
	if uid, ok := db.cookieToUser[cookie]; ok {
		user := db.mockTable[uid]

		if newName := profile.Username; newName != "" {
			user.Username = newName
		}
		if newPass := profile.NewPassword; newPass != "" {
			user.Password = newPass
		}
		if len(profile.Avatar) != 0 {
			db.avatarTable[uid] = profile.Avatar
		}

		db.mockTable[uid] = user
	}
}

//*****************|

func (db *DB) getUID(user types.APIUser) (uid int, ok bool) {
	for uid, row := range db.mockTable {
		if row.Username == user.Username {
			if row.Password == user.Password {
				return uid, true
			}
			return 0, false
		}
	}
	return 0, false
}

func (db *DB) logIn(uid int) (cookie string) {
	cookie = common.GenerateCookie()
	db.cookieToUser[cookie] = uid
	return cookie
}

//*****************|
