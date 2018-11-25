package database

import (
	lg "Wave/utiles/logger"
	"Wave/utiles/misc"
	"Wave/utiles/models"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgreSQL
)

// Model - psql facade
type Model struct {
	Database *sqlx.DB
	LG       *lg.Logger
}

func envOrDefault(env, def string) string {
	if res := os.Getenv(env); res != "" {
		return res
	}
	return def
}

// New Model
func New(logger *lg.Logger) *Model {
	var (
		postgr = &Model{
			LG: logger,
		}
		dbuser     = envOrDefault("WAVE_DB_USER", "Wave")
		dbpassword = envOrDefault("WAVE_DB_PASSWORD", "Wave")
		dbname     = envOrDefault("WAVE_DB_NAME", "Wave")
		data       = "user="+dbuser+" password="+dbpassword+" dbname='"+dbname+"' "+"sslmode=disable"
		err        error
	)
	postgr.Database, err = sqlx.Connect("postgres", data)
	if err != nil {
		postgr.LG.Sugar.Panicw(
			"PostgreSQL connection establishment failed",
			"source", "database.go",
			"who", "New",
		)
		panic(err)
	}

	postgr.LG.Sugar.Infow(
		"PostgreSQL connection establishment succeded",
		"source", "database.go",
		"who", "New",
	)

	return postgr
}

const ( // we don't need to export them!
	userInfoTable = "userinfo"
	usernameCol   = "username"
	sessionTable  = "session"
	cookieCol     = "cookie"
)

func (model *Model) present(tableName string, colName string, target string) (fl bool, err error) {
	row := model.Database.QueryRowx("SELECT EXISTS (SELECT true FROM " + tableName + " WHERE " + colName + "='" + target + "');")
	
	var exists string
	if err = row.Scan(&exists); err != nil {
		model.LG.Sugar.Infow(
			"Scan failed",
			"source", "database.go",
			"who", "present",
		)
		return false, err
	}

	fl, err = strconv.ParseBool(exists)
	if err != nil {
		model.LG.Sugar.Infow(
			"strconv.ParseBool failed",
			"source", "database.go",
			"who", "present",
		)
		return false, err
	}
	return fl, nil
}

func validateCredentials(target string) bool { 
	return true 
}

/****************************** session block ******************************/

func (model *Model) LogIn(credentials models.UserCredentials) (cookie string, err error) {
	if isPresent, problem := model.present(userInfoTable, usernameCol, credentials.Username); isPresent && problem == nil {

		row := model.Database.QueryRowx(`
		SELECT password
		FROM userinfo
		WHERE username=$1
		`, credentials.Username)
		
		var psswd string
		if err := row.Scan(&psswd); err != nil {
			model.LG.Sugar.Panicw(
				"Scan failed",
				"source", "database.go",
				"who", "LogIn",
			)
			return "", err
		}

		if misc.PasswordsMatched(psswd, credentials.Password) {
			cookie := misc.GenerateCookie()
			model.Database.MustExec(`
				INSERT INTO session(uid, cookie)
				VALUES(
					(SELECT uid FROM userinfo WHERE username=$1),
					$2
				);
			`, credentials.Username, cookie)

			model.LG.Sugar.Infow(
				"login succeded, cookie set",
				"source", "database.go",
				"who", "LogIn",
			)

			return cookie, nil
		} else {
			model.LG.Sugar.Infow(
				"login failed, wrong password",
				"source", "database.go",
				"who", "LogIn",
			)

			return "", nil
		}
	} else if !isPresent && problem == nil {

		model.LG.Sugar.Infow(
			"login failed, no such user",
			"source", "database.go",
			"who", "LogIn",
		)

		return "", nil

	} else if problem != nil {

		model.LG.Sugar.Infow(
			"present failed",
			"source", "database.go",
			"who", "LogIn",
		)

		return "", problem
	}

	model.LG.Sugar.Infow(
		"login failed, no such user",
		"source", "database.go",
		"who", "LogIn",
	)

	return "", nil
}

func (model *Model) LogOut(cookie string) error {
	model.Database.QueryRowx(`
		DELETE
		FROM session
		WHERE cookie=$1;
	`, cookie)

	model.LG.Sugar.Infow(
		"logout succeded",
		"source", "database.go",
		"who", "LogOut",
	)

	return nil
}

/****************************** user block ******************************/

func (model *Model) SignUp(credentials models.UserEdit) (cookie string, err error) {
	if validateCredentials(credentials.Username) && validateCredentials(credentials.Password) {
		if isPresent, problem := model.present(userInfoTable, usernameCol, credentials.Username); isPresent && problem == nil {

			model.LG.Sugar.Infow(
				"signup failed, user already exists",
				"source", "database.go",
				"who", "SignUp",
			)

			return "", nil
		} else if problem != nil {

			model.LG.Sugar.Infow(
				"signup succeded",
				"source", "database.go",
				"who", "SignUp",
			)

			return "", problem
		} else if !isPresent {
			cookie := misc.GenerateCookie()
			hashedPsswd := misc.GeneratePasswordHash(credentials.Password)

			if credentials.Avatar != "" {
				model.Database.MustExec(`
					INSERT INTO userinfo(username,password,avatar)
					VALUES($1, $2, $3)
				`, credentials.Username, hashedPsswd, credentials.Avatar)
			} else {
				model.Database.MustExec(`
					INSERT INTO userinfo(username,password)
					VALUES($1, $2)
				`, credentials.Username, hashedPsswd)
			}

			model.Database.MustExec(`
				INSERT INTO session(uid, cookie)
				VALUES(
					(SELECT uid FROM userinfo WHERE username=$1),
					$2
				)
			`, credentials.Username, cookie)

			model.LG.Sugar.Infow(
				"signup succeded",
				"source", "database.go",
				"who", "SignUp",
			)

			return cookie, nil
		}
	}

	return "", nil
}

func (model *Model) GetMyProfile(cookie string) (profile models.UserExtended, err error) {
	row := model.Database.QueryRowx(`
		SELECT username, avatar, score
		FROM userinfo
		JOIN session
			ON session.uid = userinfo.uid
			AND cookie=$1;
	`, cookie)
	

	if err = row.Scan(&profile.Username, &profile.Avatar, &profile.Score); err != nil {

		model.LG.Sugar.Infow(
			"getmyprofile failed, scan error",
			"source", "database.go",
			"who", "GetMyProfile",
		)

		return models.UserExtended{}, err
	}

	model.LG.Sugar.Infow(
		"getmyprofile succeded",
		"source", "database.go",
		"who", "GetMyProfile",
	)

	return profile, nil
}

func (model *Model) GetProfile(username string) (profile models.UserExtended, err error) {
	if isPresent, problem := model.present(userInfoTable, usernameCol, username); isPresent && problem == nil {
		row := model.Database.QueryRowx(`
			SELECT username, avatar, score
			FROM userinfo
			WHERE username=$1;
		`, username)
		
		if err = row.Scan(&profile.Username, &profile.Avatar, &profile.Score); err != nil {

			model.LG.Sugar.Infow(
				"getprofile failed, scan error",
				"source", "database.go",
				"who", "GetProfile",
			)

			return models.UserExtended{}, err
		}

		model.LG.Sugar.Infow(
			"getprofile succeded",
			"source", "database.go",
			"who", "GetProfile",
		)

		return profile, nil
	} else if problem != nil {

		model.LG.Sugar.Infow(
			"present failed",
			"source", "database.go",
			"who", "GetProfile",
		)

		return models.UserExtended{}, err
	} else if !isPresent {

		model.LG.Sugar.Infow(
			"getprofile failed, user doesn't exist",
			"source", "database.go",
			"who", "GetProfile",
		)

		return models.UserExtended{}, nil
	}

	return models.UserExtended{}, nil
}

func (model *Model) UpdateProfile(profile models.UserEdit, cookie string) (bool, error) {
	changedU := false
	changedP := false
	changedA := false

	if profile.Username != "" {
		isPresent, problem := model.present(userInfoTable, usernameCol, profile.Username)
		if problem != nil {
			model.LG.Sugar.Infow(
				"present failed",
				"source", "database.go",
				"who", "UpdateProfile",
			)

			return false, problem
		}
		if !isPresent {
			if validateCredentials(profile.Username) {
				model.Database.MustExec(`
					UPDATE userinfo
					SET username=$1
					WHERE userinfo.uid = (
						SELECT session.uid from session
						JOIN userinfo
							ON session.uid = userinfo.uid
						WHERE cookie=$2
					);
				`, profile.Username, cookie)

				model.LG.Sugar.Infow(
					"update profile succeded, username updated",
					"source", "database.go",
					"who", "UpdateProfile",
				)
				changedU = true

			} else {

				model.LG.Sugar.Infow(
					"update profile failed, invalid username",
					"source", "database.go",
					"who", "UpdateProfile",
				)

				changedU = false
			}
		}
		if isPresent {
			model.LG.Sugar.Infow(
				"update profile failed, username already in use",
				"source", "database.go",
				"who", "UpdateProfile",
			)

			changedU = false
		}
	}

	if profile.Password != "" {
		if validateCredentials(profile.Password) {
			hashedPsswd := misc.GeneratePasswordHash(profile.Password)
			model.Database.MustExec(`
				UPDATE userinfo
				SET password=$1
				WHERE userinfo.uid = (
					SELECT session.uid from session
					JOIN userinfo ON session.uid = userinfo.uid
					WHERE cookie=$2
				);
			`, hashedPsswd, cookie)

			model.LG.Sugar.Infow(
				"update profile succeded, password updated",
				"source", "database.go",
				"who", "UpdateProfile",
			)

			changedP = true
		} else {

			model.LG.Sugar.Infow(
				"update profile failed, invalid password",
				"source", "database.go",
				"who", "UpdateProfile",
			)

			changedP = false
		}
	}

	if profile.Avatar != "" {
		model.Database.MustExec(`
			UPDATE userinfo
			SET avatar=$1
			WHERE userinfo.uid = (
				SELECT session.uid from session
				JOIN userinfo
				ON userinfo.uid = session.uid
				WHERE cookie=$2
			);
		`, profile.Avatar, cookie)

		model.LG.Sugar.Infow(
			"update profile succeded, avatar updated",
			"source", "database.go",
			"who", "UpdateProfile",
		)

		changedA = true
	}

	if changedU || changedP || changedA {
		return true, nil
	}

	return false, nil
}

func (model *Model) GetTopUsers(limit int, offset int) (board models.Leaders, err error) {
	row := model.Database.QueryRowx(`
		SELECT COUNT(*)
		FROM userinfo
	`)
	if err := row.Scan(&board.Total); err != nil {

		model.LG.Sugar.Infow(
			"scan failed",
			"source", "database.go",
			"who", "GetTopUsers",
		)

		return models.Leaders{}, err
	}

	rows, err := model.Database.Queryx(`
		SELECT username, score
		FROM userinfo
		ORDER BY score DESC LIMIT $1 OFFSET $2;
	`, limit, offset)
	defer rows.Close()

	if err != nil {
		model.LG.Sugar.Infow(
			"queryx failed",
			"source", "database.go",
			"who", "GetTopUsers",
		)

		return models.Leaders{}, err
	}

	for rows.Next() {
		temp := models.UserScore{}
		if err = rows.Scan(&temp.Username, &temp.Score); err != nil {

			model.LG.Sugar.Infow(
				"scan failed",
				"source", "database.go",
				"who", "GetTopUsers",
			)

			return models.Leaders{}, err
		}

		board.Users = append(board.Users, temp)
	}

	return board, nil
}
