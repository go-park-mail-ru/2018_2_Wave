package database

import (
	"Wave/server/misc"
	"Wave/server/types"
	"sort"
	"strconv"
	"sync"
)

// DB is a database facade
type DB struct {
	mockTable    map[int]types.User
	cookieToUser map[string]int
	avatarTable  map[int][]byte
	scoreTable   map[int]int
	lastUID      int

	mutex sync.RWMutex
}

// New - create and initialise new database facade
func New() *DB {
	db := &DB{
		lastUID:      0,
		mockTable:    map[int]types.User{},
		cookieToUser: map[string]int{},
		avatarTable:  map[int][]byte{},
		scoreTable:   map[int]int{},
	}
	return db
}

//*****************| Auth

// IsSignedUp - weather the user is signed up
func (db *DB) IsSignedUp(user types.User) bool {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, row := range db.mockTable {
		if row.Username == user.Username {
			return true
		}
	}
	return false
}

// SignUp the user.
// NOTE: each call creates new record with unique uid
func (db *DB) SignUp(profile types.SignUp) (cookie string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.lastUID++

	uid := db.lastUID
	db.mockTable[uid] = profile.AsUser()
	db.avatarTable[uid] = profile.Avatar
	db.scoreTable[uid] = 0
	return db.logIn(uid)
}

// IsLoggedIn - weather the user is logged in
func (db *DB) IsLoggedIn(cookie string) bool {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	_, ok := db.cookieToUser[cookie]
	return ok
}

// LogIn the user if the one is signe up
// NOTE: each call creates new session
func (db *DB) LogIn(user types.User) (cookie string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if uid, ok := db.getUID(user); ok {
		return db.logIn(uid)
	}
	return ""
}

// LogOut a user with the cookie if the one was logged in
func (db *DB) LogOut(cookie string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, ok := db.cookieToUser[cookie]; ok {
		delete(db.cookieToUser, cookie)
	}
}

//*****************| Profile

// GetProfile returns a profile assigned to the cookie
func (db *DB) GetProfile(cookie string) (types.Profile, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if uid, ok := db.cookieToUser[cookie]; ok {
		return types.Profile{
			Username:  db.mockTable[uid].Username,
			AvatarURI: "/img/avatars/" + strconv.Itoa(uid),
			Score:     db.scoreTable[uid],
		}, true
	}
	return types.Profile{}, false
}

// GetAvatar returns avatar's data
func (db *DB) GetAvatar(uid int) ([]byte, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	data, ok := db.avatarTable[uid]
	return data, ok
}

// UpdateProfile updates profile
func (db *DB) UpdateProfile(cookie string, profile types.EditProfile) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

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

//*****************| Leaderboard

func (db *DB) GetUserScore(cookie string) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if uid, ok := db.cookieToUser[cookie]; ok {
		return db.scoreTable[uid]
	}
	return 0
}

func (db *DB) GetTopUsers(start, count int) (board types.Leaderboard) {
	type Pair struct {
		uid   int
		score int
	}
	pairs := []Pair{}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	for key, val := range db.scoreTable {
		pairs = append(pairs, Pair{uid: key, score: val})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].score > pairs[j].score
	})

	board.Total = len(pairs)

	if start >= len(pairs) {
		return board
	}
	end := start + count
	if end > len(pairs) {
		end = len(pairs)
	}

	for _, pair := range pairs[start:end] {
		board.Users = append(board.Users, types.LeaderboardRow{
			Username: db.mockTable[pair.uid].Username,
			Score:    pair.score,
		})
	}
	return board
}

//*****************|

func (db *DB) getUID(user types.User) (uid int, ok bool) {
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
	cookie = misc.GenerateCookie()
	db.cookieToUser[cookie] = uid
	return cookie
}

//*****************|
