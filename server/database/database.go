package database

import (
	"Wave/utiles/config"
	"Wave/utiles/misc"
	"Wave/utiles/models"
	"fmt"
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseModel struct {
	DBconf   config.DatabaseConfiguration
	Database *sqlx.DB
}

func New(dbconf_ config.DatabaseConfiguration) *DatabaseModel {
	postgr := &DatabaseModel{
		DBconf: dbconf_,
	}

	var err error
	postgr.Database, err = sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname='%s' sslmode=disable", postgr.DBconf.User, postgr.DBconf.Password, postgr.DBconf.DBName))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("postgres connection established")

	return postgr
}

const (
	UserInfoTable = "userinfo"
	UsernameCol   = "username"

	SessionTable = "session"
	CookieCol    = "cookie"
)

func (model *DatabaseModel) present(tableName string, colName string, target string) (fl bool, err error) {
	var exists string
	row := model.Database.QueryRowx("SELECT EXISTS (SELECT true FROM " + tableName + " WHERE " + colName + "='" + target + "');")
	err = row.Scan(&exists)

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	fl, err = strconv.ParseBool(exists)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	log.Println(fl)
	return fl, nil
}

func validateCredentials(target string) bool {
	/*
		// http://regexlib.com/REDetails.aspx?regexp_id=2298
		reg, _ := regexp.Compile("^([a-zA-Z])[a-zA-Z_-]*[\\w_-]*[\\S]$|^([a-zA-Z])[0-9_-]*[\\S]$|^[a-zA-Z]*[\\S]$")

		if reg.MatchString(target) {
			return true
		}
		log.Println("bad username or/and password")

		return false
	*/
	return true
}

/****************************** session block ******************************/

func (model *DatabaseModel) LogIn(credentials models.UserCredentials) (cookie string, err error) {
	if isPresent, problem := model.present(UserInfoTable, UsernameCol, credentials.Username); isPresent && problem == nil {
		var psswd string
		row := model.Database.QueryRowx("SELECT password FROM userinfo WHERE username=$1", credentials.Username)
		err := row.Scan(&psswd)

		if err != nil {
			log.Fatal(err)
			return "", err
		}

		if misc.PasswordsMatched(psswd, credentials.Password) {
			cookie := misc.GenerateCookie()
			model.Database.MustExec("INSERT INTO session(uid, cookie) VALUES((SELECT uid FROM userinfo WHERE username=$1), $2);", credentials.Username, cookie)

			log.Println("login successful: cookie set")

			return cookie, nil
		} else {
			log.Println("login failed: wrong password")

			return "", nil
		}
	} else {
		log.Fatal(problem)
	}

	return "", nil
}

func (model *DatabaseModel) LogOut(cookie string) error {
	log.Println(cookie)
	model.Database.QueryRowx("DELETE FROM session WHERE cookie=$1;", cookie)
	log.Println("logout successful")

	return nil
}

/****************************** user block ******************************/

func (model *DatabaseModel) SignUp(credentials models.UserCredentials) (cookie string, err error) {
	log.Println(credentials.Username)
	if validateCredentials(credentials.Username) && validateCredentials(credentials.Password) {
		if isPresent, problem := model.present(UserInfoTable, UsernameCol, credentials.Username); isPresent && problem == nil {
			log.Println("signup failed: user already exists")

			return "", nil
		} else if !isPresent && problem == nil {
			cookie := misc.GenerateCookie()
			hashedPsswd := misc.GeneratePasswordHash(credentials.Password)
			model.Database.MustExec("INSERT INTO userinfo(username,password) VALUES($1, $2)", credentials.Username, hashedPsswd)
			model.Database.MustExec("INSERT INTO session(uid, cookie) VALUES((SELECT uid FROM userinfo WHERE username=$1), $2)", credentials.Username, cookie)
			log.Println("signup successful")

			return cookie, nil
		} else if problem != nil {
			return "", problem
		}
	}

	return "", nil
}

func (model *DatabaseModel) GetMyProfile(cookie string) (profile models.UserExtended, err error) {
	row := model.Database.QueryRowx("SELECT username,avatar,score FROM userinfo JOIN session ON session.uid = userinfo.uid AND cookie=$1;", cookie)
	err = row.Scan(&profile.Username, &profile.Avatar, &profile.Score)

	if err != nil {
		log.Fatal(err)
		return models.UserExtended{}, err
	}
	log.Println("get my profile successful")

	return profile, nil
}

func (model *DatabaseModel) GetProfile(username string) (profile models.UserExtended, err error) {
	if isPresent, problem := model.present(UserInfoTable, UsernameCol, username); isPresent && problem == nil {
		row := model.Database.QueryRowx("SELECT username,avatar,score FROM userinfo WHERE username=$1;", username)
		err = row.Scan(&profile.Username, &profile.Avatar, &profile.Score)

		if err != nil {
			return models.UserExtended{}, err
		}
		log.Println("get profile successful")

		return profile, nil
	} else if problem != nil {
		log.Println(problem)

		return models.UserExtended{}, err
	} else if !isPresent {
		log.Println("user doesn't exist")
		return models.UserExtended{}, nil
	}

	return models.UserExtended{}, nil
}

func (model *DatabaseModel) UpdateProfile(profile models.UserEdit, cookie string) (bool, error) {
	changedU := false
	changedP := false
	changedA := false

	if profile.NewUsername != "" {
		isPresent, problem := model.present(UserInfoTable, UsernameCol, profile.NewUsername)
		if problem != nil {
			return false, problem
		}
		if !isPresent {
			if validateCredentials(profile.NewUsername) {
				model.Database.MustExec("UPDATE userinfo SET username=$1 WHERE userinfo.uid = (SELECT session.uid from session JOIN userinfo ON session.uid = userinfo.uid WHERE cookie=$2);", profile.NewUsername, cookie)
				log.Println("update profile successful: username changed")

				changedU = true
			} else {
				log.Println("update profile failed: bad username")

				changedU = false
			}
		}
		if isPresent {
			log.Println("update profile fail: username already in use")

			changedU = false
		}
	}

	if profile.NewPassword != "" {
		if validateCredentials(profile.NewPassword) {
			hashedPsswd := misc.GeneratePasswordHash(profile.NewPassword)
			model.Database.MustExec("UPDATE userinfo SET password=$1 WHERE userinfo.uid = (SELECT session.uid from session JOIN userinfo ON session.uid = userinfo.uid WHERE cookie=$2);", hashedPsswd, cookie)
			log.Println("update profile successful: password changed")

			changedP = true
		} else {
			log.Println("update profile failed: bad password")

			changedP = false
		}
	}

	if profile.NewAvatar != "" {
		model.Database.MustExec("UPDATE userinfo SET avatar=$1 WHERE userinfo.uid = (SELECT session.uid from session JOIN userinfo ON userinfo.uid = session.uid WHERE cookie=$2);", profile.NewAvatar, cookie)
		log.Println("update profile successful: avatar changed")

		changedA = true
	}

	if changedU || changedP || changedA {
		return true, nil
	}

	return false, nil
}

func (model *DatabaseModel) GetTopUsers(limit int, offset int) (board models.Leaders, err error) {
	row := model.Database.QueryRowx("SELECT COUNT(*) FROM userinfo")
	if err := row.Scan(&board.Total); err != nil {
		return models.Leaders{}, err
	}

	rows, err := model.Database.Query("SELECT username,score FROM userinfo ORDER BY score DESC LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return models.Leaders{}, err
	}

	for rows.Next() {
		temp := models.UserScore{}
		if err = rows.Scan(&temp.Username, &temp.Score); err != nil {
			return models.Leaders{}, err
		}

		board.Users = append(board.Users, temp)
	}
	return board, nil
}
