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
// 
// validation
// different regexps for username & password (set min/maxlengths!)
// log into file
// avatarr when signup

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
	postgr.db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
	postgr.dbconf.User, postgr.dbconf.Password, postgr.dbconf.Name))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("postgres connection established")

	return postgr
}

func (db *DB) present(tableName string, colName string, target string) (bool, error) {
	var flag string

	row := db.db.QueryRow("SELECT EXISTS (SELECT true FROM " + tableName + " WHERE " + colName + "='" + target + "')")
	errScan := row.Scan(&flag)
	if errScan != nil {
		return false, errScan
	}
	
	fl, errParse := strconv.ParseBool(flag)
	if errParse != nil {
		return false, errParse
	}

	return fl, nil
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

func (db *DB) SignUp(profile types.SignUp) (cookie string, err error) {
	if validateCredentials(profile.Username) && validateCredentials(profile.Password) {
		if isPresent, problem := db.present(UserInfoTable, UsernameCol, profile.Username); isPresent && problem == nil {
			log.Println("signup failed: user already exists")

			return "", nil
		} else if !isPresent && problem == nil {
			cookie := misc.GenerateCookie() // = ?
			_, err := db.db.Exec("INSERT INTO userinfo(username,password) VALUES($1, $2)", profile.Username, profile.Password)
			
			if err != nil {
				return "", err
			}

			_, err = db.db.Exec("INSERT INTO cookie(uid, cookieStr) VALUES((SELECT uid FROM userinfo WHERE username=$1), $2)", profile.Username, cookie)
			
			if err != nil {
				return "", err
			}

			log.Println("signup successful")

			return cookie, nil
		} else if problem != nil {
			return "", problem
		}
	}

	return "", nil
}

func (db *DB) IsLoggedIn(cookie string) (bool, error) {
	foundCookie := true
	if isPresent, problem := db.present(CookieTable, CookieCol, cookie); !isPresent && problem == nil {
		log.Println("is logged in check failed: no such cookie found, returns: false")
		foundCookie = false

		return foundCookie, nil
	} else if isPresent && problem == nil {
		log.Println("is logged in check successful: cookie found, returns: true")
	} else if problem != nil {

		return false, problem
	}
	
	return foundCookie, nil
}

func (db *DB) LogIn(credentials types.User) (cookie string, err error) {
	if isPresent, problem := db.present(UserInfoTable, UsernameCol, credentials.Username); isPresent && problem == nil {
		var psswd string
		row := db.db.QueryRow("SELECT password FROM userinfo WHERE username=$1", credentials.Username);
		err := row.Scan(&psswd)
		
		if err != nil {
			return "", err
		}

		if psswd == credentials.Password {
			cookie := misc.GenerateCookie()
			_, err := db.db.Exec("INSERT INTO cookie(uid, cookieStr) VALUES((SELECT uid FROM userinfo WHERE username=$1), $2);", credentials.Username, cookie)
			
			if err != nil {
				return "", err
			}

			log.Println("login successful: cookie set")
			
			return cookie, nil
		} else {
			log.Println("login failed: wrong password")

			return "", nil
		}
	}

	return "", nil
}

func (db *DB) LogOut(cookie string) error {
	_, err := db.db.Exec("DELETE FROM cookie WHERE cookieStr=$1", cookie);
	
	if err != nil {
		return err
	}

	log.Println("logout successful");

	return nil
}

/****************** User profile block ******************/

func (db *DB) GetProfile(cookie string) (profile types.Profile, err error) {
	row := db.db.QueryRow("SELECT username,avatar,score FROM userinfo JOIN cookie ON cookie.uid = userinfo.uid AND cookieStr=$1;", cookie)
	err = row.Scan(&profile.Username, &profile.AvatarURI, &profile.Score)
	
	if err != nil {
		return types.Profile{}, err
	}
	log.Println("get profile successful");

	return profile, nil
}

/*
func (db *DB) GetAvatar(uid int) (avatarSource string, err error) {
	if isPresent, problem := db.present(UserInfoTable, UID, strconv.Itoa(uid)); !isPresent && problem == nil {
		log.Println("get avatar failed: user doesn't exist")

		return "", nil
	} else if problem != nil {
		return "", problem
	}
	row := db.db.QueryRow("SELECT avatar FROM userinfo WHERE uid=$1", uid)
	err := row.Scan(&avatarSource)
	
	if err != nil {
		return "", err
	}

	log.Println("get avatar successful");

	return avatarSource, nil
}
*/

func (db *DB) GetAvatar(cookie string) (avatarSource string, err error) {
	
	row := db.db.QueryRow("SELECT avatar FROM userinfo JOIN cookie ON cookie.uid = userinfo.uid AND cookieStr=$1;", cookie)
	err = row.Scan(&avatarSource)
	
	if err != nil {
		return "", err
	}

	log.Println("get avatar successful");

	return avatarSource, nil
}

func (db *DB) UpdateProfile(cookie string, profile types.EditProfile) (bool, error) {
	if profile.NewUsername != "" {
		isPresent, problem := db.present(UserInfoTable, UsernameCol, profile.NewUsername);
		if problem != nil {
			return false, problem
		} 
		if !isPresent {
			if validateCredentials(profile.NewUsername) {
				// to change!
				db.db.QueryRow("UPDATE userinfo SET username=$1 WHERE userinfo.uid = (SELECT cookie.uid from cookie JOIN userinfo ON cookie.uid = userinfo.uid WHERE cookieStr=$2);", profile.NewUsername, cookie)
				log.Println("update profile successful: username changed")

				return true, nil
			} else {
				log.Println("update profile failed: bad username")

				return false, nil
			}
		}
		if isPresent {
			log.Println("update profile fail: username already in use")

			return false, nil
		}
	}

	if profile.NewPassword != "" {
		if validateCredentials(profile.NewPassword) {
			// to change!
			db.db.QueryRow("UPDATE userinfo SET password=$1 WHERE userinfo.uid = (SELECT cookie.uid from cookie JOIN userinfo ON cookie.uid = userinfo.uid WHERE cookieStr=$2);", profile.NewPassword, cookie)
			log.Println("update profile successful: password changed")

			return true, nil
		} else {
			log.Println("update profile failed: bad password")

			return false, nil
		}
	}

	if profile.Avatar != "" {
		db.db.QueryRow("UPDATE userinfo SET avatar=$1 WHERE userinfo.uid = (SELECT cookie.uid from cookie JOIN userinfo ON cookie.uid = userinfo.uid WHERE cookieStr=$2);", profile.Avatar, cookie)
		log.Println("update profile successful: avatar changed")

		return true, nil
	}
	
	return false, nil
}

/****************** Leaderboard block ******************/

func (db *DB) GetTopUsers(limit int, offset int) (board types.Leaderboard, err error) {
	row := db.db.QueryRow("SELECT COUNT(*) FROM userinfo")
	err = row.Scan(&board.Total)
	
	if err != nil {
		return types.Leaderboard{}, err
	}

	rows, err := db.db.Query("SELECT username,score FROM userinfo ORDER BY score DESC LIMIT $1 OFFSET $2;", limit, offset)
	
	if err != nil {
		return types.Leaderboard{}, err
	}

	temp := types.LeaderboardRow{}

	for rows.Next() {
		err = rows.Scan(&temp.Username, &temp.Score)

		if err != nil {
			return types.Leaderboard{}, err
		}

		board.Users = append(board.Users, temp)
	}

	return board, nil
}