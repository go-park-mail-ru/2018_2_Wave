package main

// todo:
// ports from configs
// error processing
// password hashing
// set default avatar source
// cookie generating and refreshing									:check
// db init
// capital letters - ask Bomkvilt
// defer db.Close()
// validate()														:check
// http integration
// different regexps for username & password (set min/maxlengths!)
// after signup do we automatically log in?							:yes
// if yes then insert cookie into sigup								:check
// cycle imports
// rename types a little bit
// api is an empty cookie!
// avatar processing
// GetAvatar interface is questionable
// no need for CurPassword in types.EditProfile -- excessive
// no need to pass cookie to UpdateProfile()						:fixed
// bug: fix password change if curUser doesn't exist				:fixed
// or maybe i really should use cookie :D							:fixed
// i've change some things in types.go - is it ok?
// leaderbord -> leaderboard

//!!!important!!!
//http
//hash
//avatar

import (
	"2018_2_Wave/server/misc"
	"2018_2_Wave/server/types"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"

	_ "github.com/lib/pq"
)

/****************** DB utility block ******************/

const (
	dbUser     = "postgres"
	dbPassword = "test"
	dbName     = "postgres"
)

const (
	UserInfoTable = "userinfo"
	UsernameCol   = "username"
	PasswordCol	  = "password"
	CookieCol     = "cookie"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// validates the input
func validateCredentials(target string) bool {
	// http://regexlib.com/REDetails.aspx?regexp_id=2298
	reg, _ := regexp.Compile("^([a-zA-Z])[a-zA-Z_-]*[\\w_-]*[\\S]$|^([a-zA-Z])[0-9_-]*[\\S]$|^[a-zA-Z]*[\\S]$")

	if reg.MatchString(target) {
		return true
	}
	log.Println("bad username or/and password")

	return false
}

// identifies if the target record is present in the table
func present(db *sql.DB, colName string, target string) bool {
	var flag string

	//row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM $1 WHERE $2=$3);", UserInfoTable, colName, target)
	row := db.QueryRow("SELECT EXISTS (SELECT true FROM " + UserInfoTable + " WHERE " + colName + "='" + target + "');")
	err := row.Scan(&flag)
	checkErr(err)

	fl, _ := strconv.ParseBool(flag)
	return fl
}

/****************** Authorization block ******************/

// creates a new record and returns a cookie for newly signed up player
func SignUp(db *sql.DB, profile types.SignUp) (cookie string) {
	if validateCredentials(profile.Username) && validateCredentials(profile.Password) {
		if present(db, UsernameCol, profile.Username) {
			log.Println("signup failed: user already exists")

			return ""
		} else {
			cookie := misc.GenerateCookie()
			_, err := db.Exec("INSERT INTO UserInfo(username,password,cookie) VALUES($1, $2, $3)", profile.Username, profile.Password, cookie)
			//_, err := db.Exec("INSERT INTO userinfo($1,$2,$3) VALUES($4, $5, $6)", UsernameCol, PasswordCol, CookieCol, profile.Username, profile.Password, cookie)
			checkErr(err)
			log.Println("signup successful")

			return cookie
		}
	}

	return ""
}

// checks if cookie for the current user exists and returns it
func IsLoggedIn(db *sql.DB, credentials types.User) (cookie string) {
	foundCookie := ""
	if !present(db, UsernameCol, credentials.Username) {
		log.Println("is logged in check failed: no such user found, returns: EMPTY cookie")

		return ""
	}

	row := db.QueryRow("SELECT cookie FROM userinfo WHERE username = '" + credentials.Username + "' AND cookie IS NOT NULL;")
	err := row.Scan(&foundCookie)
	checkErr(err)
	log.Println("is logged in check successful: user's cookie retrieved, returns: EMPTY/NOT EMPTY cookie")

	return foundCookie
}

func LogIn(db *sql.DB, credentials types.User) (cookie string) {
	if present(db, UsernameCol, credentials.Username) {
		var psswd string
		row := db.QueryRow("SELECT password FROM userinfo WHERE username=$1", credentials.Username);
		err := row.Scan(&psswd)
		checkErr(err)

		if psswd == credentials.Password && IsLoggedIn(db, credentials) == "" {
			cookie := misc.GenerateCookie()
			db.Exec("UPDATE userinfo SET cookie=$1 WHERE username=$2", cookie, credentials.Username)
			log.Println("login successful: cookie set")
			
			return cookie
		} else {
			log.Println("login failed: wrong password")

			return ""
		}
	}

	return IsLoggedIn(db, credentials)
}

func LogOut(db *sql.DB, cookie string) {
	db.Exec("UPDATE userinfo SET cookie='' WHERE cookie=$1", cookie);
	log.Println("logout successful");
}

/****************** User profile block ******************/

func GetProfile(db *sql.DB, cookie string) (profile types.Profile) {
	row := db.QueryRow("SELECT username,avatar,score FROM userinfo WHERE cookie=$1", cookie)
	err := row.Scan(&profile.Username, &profile.AvatarURI, &profile.Score)
	checkErr(err)
	log.Println("get profile successful");

	return profile
}

func GetAvatar(db *sql.DB, uid int) (avatarSource string) {
	row := db.QueryRow("SELECT avatar FROM userinfo WHERE uid=$1", uid)
	err := row.Scan(&avatarSource)
	checkErr(err)
	log.Println("get avatar successful");

	return avatarSource
}

func UpdateProfile(db *sql.DB, cookie string, profile types.EditProfile) {
	updateCount := 0
	if profile.NewUsername != "" {
		if !present(db, UsernameCol, profile.NewUsername) {
			if validateCredentials(profile.NewUsername) {
				//db.Exec("UPDATE userinfo SET username=$1 WHERE username=$2", profile.NewUsername, profile.CurUsername)
				db.Exec("UPDATE userinfo SET username=$1 WHERE username=$2 AND cookie=$3", profile.NewUsername, profile.CurUsername, cookie)
				log.Println("update profile successful: username changed");
				updateCount++
			} else {
				log.Println("update profile failed: bad username");
			}
		} else {
			log.Println("update profile fail: username already in use");
		}
	}

	if profile.NewPassword != "" {
		if validateCredentials(profile.NewPassword) {
			if updateCount > 0 {
				//db.Exec("UPDATE userinfo SET password=$1 WHERE username=$2", profile.NewPassword, profile.NewUsername)
				db.Exec("UPDATE userinfo SET password=$1 WHERE username=$2 AND cookie=$3", profile.NewPassword, profile.NewUsername, cookie)
				log.Println("update profile successful: password & username changed");
			} else {
				//db.Exec("UPDATE userinfo SET password=$1 WHERE username=$2", profile.NewPassword, profile.CurUsername)
				db.Exec("UPDATE userinfo SET password=$1 WHERE username=$2 AND cookie=$3", profile.NewPassword, profile.CurUsername, cookie)
				log.Println("update profile successful: password changed");
			}
		} else {
			log.Println("update profile failed: bad password");
		}
	}

	if profile.Avatar != "" {
		//todo
	}
}

/****************** Leaderboard block ******************/

func GetTopUsers(db *sql.DB, offset, limit int) (board types.Leaderboard) {
	row := db.QueryRow("SELECT COUNT(*) FROM userinfo")
	err := row.Scan(&board.Total)
	checkErr(err)

	rows, err := db.Query("SELECT username,score FROM userinfo ORDER BY score ASC;")
	checkErr(err)

	temp := types.LeaderboardRow{}

	for rows.Next() {
		err = rows.Scan(&temp.Username, &temp.Score)
		board.Users = append(board.Users, temp)
	}

	// to be continued
	return board
}

func main() {
	/*
	params := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", params)

	checkErr(err)
	*/
	//defer db.Close()
}
