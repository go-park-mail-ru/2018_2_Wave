package database

import (
	"Wave/server/misc"
	"Wave/server/types"
	"Wave/utiles"
	"database/sql"
	"regexp"
	"strconv"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//TODO:
// password hashing
// error processing
// validation
// different regexps for username & password (set min/maxlengths!)
// log into file
// curPassword check
// update avatar
// fetch GetUser by any username
// GetAvatar
// avatar on Signup

// ORM!

const (
	UserInfoTable = "userinfo"
	UsernameCol   = "username"
	UID 		  = "uid"

	CookieTable	  = "cookie"
	CookieCol	  = "cookieStr"
)

// Facade
type DB struct {
	dbconf	utiles.DatabaseConfig

	db		*sql.DB
}

func New(dbconf_ utiles.DatabaseConfig) *DB {
	postgr := &DB{
		dbconf : dbconf_,
	}
	
	var err error
	//postgr.db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
	//postgr.dbconf.User, postgr.dbconf.Password, postgr.dbconf.Name))
	postgr.db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
	"waveapp", "surf", "wave"))
	checkErr(err)
	log.Println("postgres connection established")

	return postgr
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (db *DB) present(tableName string, colName string, target string) bool {
	var flag string

	row := db.db.QueryRow("SELECT EXISTS (SELECT true FROM " + tableName + " WHERE " + colName + "='" + target + "')")
	err := row.Scan(&flag)
	checkErr(err)
	
	fl, _ := strconv.ParseBool(flag)
	return fl
}

func validateCredentials(target string) bool {
	// http://regexlib.com/REDetails.aspx?regexp_id=2298
	reg, _ := regexp.Compile("^([a-zA-Z])[a-zA-Z_-]*[\\w_-]*[\\S]$|^([a-zA-Z])[0-9_-]*[\\S]$|^[a-zA-Z]*[\\S]$")

	if reg.MatchString(target) {
		return true
	}
	log.Println("bad username or/and password")

	return false
}

/****************** Authorization block ******************/

func (db *DB) SignUp(profile types.SignUp) (cookie string) {
	if validateCredentials(profile.Username) && validateCredentials(profile.Password) {
		if db.present(UserInfoTable, UsernameCol, profile.Username) {
			log.Println("signup failed: user already exists")

			return ""
		} else {
			cookie := misc.GenerateCookie()
			_, err := db.db.Exec("INSERT INTO userinfo(username,password) VALUES($1, $2)", profile.Username, profile.Password)
			checkErr(err)
			_, err = db.db.Exec("INSERT INTO cookie(uid, cookieStr) VALUES((SELECT uid FROM userinfo WHERE username=$1), $2)", profile.Username, cookie)
			checkErr(err)
			log.Println("signup successful")

			return cookie
		}
	}

	return ""
}

func (db *DB) IsLoggedIn(cookie string) bool {
	foundCookie := true
	if !db.present(CookieTable, CookieCol, cookie) {
		log.Println("is logged in check failed: no such cookie found, returns: false")
		foundCookie = false

		return foundCookie
	} else {
		log.Println("is logged in check successful: cookie found, returns: true")
	}
	
	return foundCookie
}

func (db *DB) LogIn(credentials types.User) (cookie string) {
	if db.present(UserInfoTable, UsernameCol, credentials.Username) {
		var psswd string
		row := db.db.QueryRow("SELECT password FROM userinfo WHERE username=$1", credentials.Username);
		err := row.Scan(&psswd)
		checkErr(err)

		if psswd == credentials.Password {
			cookie := misc.GenerateCookie()
			_, err := db.db.Exec("INSERT INTO cookie(uid, cookieStr) VALUES((SELECT uid FROM userinfo WHERE username=$1), $2);", credentials.Username, cookie)
			checkErr(err)
			log.Println("login successful: cookie set")
			
			return cookie
		} else {
			log.Println("login failed: wrong password")

			return ""
		}
	}

	return ""
}

func (db *DB) LogOut(cookie string) {
	_, err := db.db.Exec("DELETE FROM cookie WHERE cookieStr=$1", cookie);
	checkErr(err)

	log.Println("logout successful");
}

/****************** User profile block ******************/

// current user 
func (db *DB) GetProfile(cookie string) (profile types.Profile) {
	row := db.db.QueryRow("SELECT username,avatar,score FROM userinfo JOIN cookie ON cookie.uid = userinfo.uid AND cookieStr=$1;", cookie)
	err := row.Scan(&profile.Username, &profile.AvatarURI, &profile.Score)
	checkErr(err)
	log.Println("get profile successful");

	return profile
}

func (db *DB) GetAvatar(uid int) (avatarSource string) {
	if !db.present(UserInfoTable, UID, strconv.Itoa(uid)) {
		log.Println("get avatar failed: user doesn't exist")

		return ""
	}
	row := db.db.QueryRow("SELECT avatar FROM userinfo WHERE uid=$1", uid)
	err := row.Scan(&avatarSource)
	checkErr(err)
	log.Println("get avatar successful");

	return avatarSource
}

func (db *DB) UpdateProfile(cookie string, profile types.EditProfile) bool {
	if profile.NewUsername != "" {
		if !db.present(UserInfoTable, UsernameCol, profile.NewUsername) {
			if validateCredentials(profile.NewUsername) {
				// to change!
				db.db.QueryRow("UPDATE userinfo SET username=$1 WHERE userinfo.uid = (SELECT cookie.uid from cookie JOIN userinfo ON cookie.uid = userinfo.uid WHERE cookieStr=$2);", profile.NewUsername, cookie)
				log.Println("update profile successful: username changed")

				return true
			} else {
				log.Println("update profile failed: bad username")

				return false
			}
		} else {
			log.Println("update profile fail: username already in use")

			return false
		}
	}

	if profile.NewPassword != "" {
		if validateCredentials(profile.NewPassword) {
			// to change!
			// curPassword check
			db.db.QueryRow("UPDATE userinfo SET password=$1 WHERE userinfo.uid = (SELECT cookie.uid from cookie JOIN userinfo ON cookie.uid = userinfo.uid WHERE cookieStr=$2);", profile.NewPassword, cookie)
			log.Println("update profile successful: password changed")

			return true
		} else {
			log.Println("update profile failed: bad password")

			return false
		}
	}

	/*
	if profile.Avatar != "" {
		//db.db.Exec("", profile.Avatar, cookie)
		log.Println("update profile successful: avatar changed")
	}
	*/
	return false
}

/****************** Leaderboard block ******************/

func (db *DB) GetTopUsers(limit, offset int) (board types.Leaderboard) {
	row := db.db.QueryRow("SELECT COUNT(*) FROM userinfo")
	err := row.Scan(&board.Total)
	checkErr(err)

	rows, err := db.db.Query("SELECT username,score FROM userinfo ORDER BY score DESC LIMIT $1 OFFSET $2;", limit, offset)
	checkErr(err)

	temp := types.LeaderboardRow{}

	for rows.Next() {
		err = rows.Scan(&temp.Username, &temp.Score)
		board.Users = append(board.Users, temp)
	}
	return board
}