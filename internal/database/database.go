package database

import (
	lg "Wave/internal/logger"
	"Wave/internal/misc"
	"Wave/internal/models"
	"os"
	"strconv"

	"github.com/namsral/flag"

	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	UserInfoTable = "userinfo"
	UsernameCol   = "username"

	SessionTable = "session"
	CookieCol    = "cookie"
)

type DatabaseModel struct {
	Database *sqlx.DB
	LG       lg.Logger
}

func New(lg_ *lg.Logger) *DatabaseModel {
	var (
		postgr = &DatabaseModel{
			LG: *lg_,
		}
		dbuser     string
		dbpassword string
		dbname     string
		err        error
	)
	flag.StringVar(&dbuser, "WAVE_DB_USER", "Wave", "")
	flag.StringVar(&dbname, "WAVE_DB_NAME", "Wave", "")
	flag.StringVar(&dbpassword, "WAVE_DB_PASSWORD", "Wave", "")
	flag.Parse()

	// postgr.Database, err = sqlx.Connect("postgres", "user="+dbuser+" password="+dbpassword+" dbname='"+dbname+"' "+"sslmode=disable")

	postgr.Database, err = sqlx.Connect("postgres", "user=waveapp password='surf' dbname='wave' sslmode=disable")

	if err != nil {
		postgr.LG.Info(
			"PostgreSQL connection establishment failed",
			"source", "database.go",
			"who", "New",
		)

		os.Exit(1)
		// exitting
	}

	postgr.LG.Info(
		"PostgreSQL connection establishment succeeded",
		"source", "database.go",
		"who", "New",
	)

	return postgr
}

func (model *DatabaseModel) Present(tableName string, colName string, target string) (fl bool, err error) {
	var exists string
	row := model.Database.QueryRowx("SELECT EXISTS (SELECT true FROM " +
		tableName + " WHERE " + colName + "='" + target + "');")
	err = row.Scan(&exists)

	if err != nil {

		model.LG.Info(
			"Scan failed",
			"source", "database.go",
			"who", "Present",
		)

		return false, err
	}

	fl, err = strconv.ParseBool(exists)

	if err != nil {

		model.LG.Info(
			"strconv.ParseBool failed",
			"source", "database.go",
			"who", "Present",
		)

		return false, err
	}

	return fl, nil
}

func ValidateUname(target string) bool {
	if len(target) < 4 {
		return false
	}

	return true
}

func ValidatePassword(target string) bool {
	if len(target) < 4 {
		return false
	}

	return true
}

func (model *DatabaseModel) Login(credentials models.UserCredentials) (cookie string, err error) {
	if isPresent, problem := model.Present(UserInfoTable, UsernameCol, credentials.Username); isPresent && problem == nil {
		var psswd string

		row := model.Database.QueryRowx(`
			SELECT password
			FROM userinfo
			WHERE username=$1;
		`, credentials.Username)

		err := row.Scan(&psswd)

		if err != nil {

			model.LG.Info(
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

			model.LG.Info(
				"login succeeded, cookie set",
				"source", "database.go",
				"who", "LogIn",
			)

			return cookie, nil
		} else {
			model.LG.Info(
				"login failed, wrong password",
				"source", "database.go",
				"who", "LogIn",
			)

			return "", nil
		}
	} else if problem != nil {

		model.LG.Info(
			"Present failed",
			"source", "database.go",
			"who", "LogIn",
		)

		return "", problem
	}

	model.LG.Info(
		"login failed, no such user",
		"source", "database.go",
		"who", "LogIn",
	)

	return "", nil
}

func (model *DatabaseModel) GetMyProfile(cookie string) (profile models.UserExtended, err error) {
	row := model.Database.QueryRowx(`
		SELECT username, avatar, score, locale
		FROM userinfo
		JOIN session
		USING(uid)
		WHERE cookie=$1;
	`, cookie)
	err = row.Scan(&profile.Username, &profile.Avatar, &profile.Score, &profile.Locale)

	if err != nil {

		model.LG.Info(
			"getmyprofile failed, scan error",
			"source", "database.go",
			"who", "GetMyProfile",
		)

		return models.UserExtended{}, err
	}

	model.LG.Info(
		"getmyprofile succeeded",
		"source", "database.go",
		"who", "GetMyProfile",
	)

	return profile, nil
}

func (model *DatabaseModel) UpdateMyLocale(cookie string, locale string) (err error) {

	model.Database.MustExec(`
		UPDATE userinfo
		SET locale=$1
		WHERE uid=(
			SELECT DISTINCT session.uid FROM session WHERE cookie=$2
		);
	`, locale, cookie)

	if err != nil {

		model.LG.Info(
			"getmyprofile failed, scan error",
			"source", "database.go",
			"who", "GetMyProfile",
		)

		return err
	}

	model.LG.Info(
		"getmyprofile succeeded",
		"source", "database.go",
		"who", "GetMyProfile",
	)

	return nil
}

func (model *DatabaseModel) GetProfile(username string) (profile models.UserExtended, err error) {
	if isPresent, problem := model.Present(UserInfoTable, UsernameCol, username); isPresent && problem == nil {
		row := model.Database.QueryRowx(`
			SELECT username, avatar, score, locale
			FROM userinfo
			WHERE username=$1;
		`, username)
		err = row.Scan(&profile.Username, &profile.Avatar, &profile.Score, &profile.Locale)

		if err != nil {

			model.LG.Info(
				"getprofile failed, scan error",
				"source", "database.go",
				"who", "GetProfile",
			)

			return models.UserExtended{}, err
		}

		model.LG.Info(
			"getprofile succeeded",
			"source", "database.go",
			"who", "GetProfile",
		)

		return profile, nil
	} else if problem != nil {

		model.LG.Info(
			"Present failed",
			"source", "database.go",
			"who", "GetProfile",
		)

		return models.UserExtended{}, err
	}

	model.LG.Info(
		"getprofile failed, user doesn't exist",
		"source", "database.go",
		"who", "GetProfile",
	)

	return models.UserExtended{}, nil
}

func (model *DatabaseModel) UpdateProfile(profile models.UserEdit, cookie string) error {
	if profile.Username != "" {
		isPresent, problem := model.Present(UserInfoTable, UsernameCol, profile.Username)
		if problem != nil {
			model.LG.Info(
				"Present failed",
				"source", "database.go",
				"who", "UpdateProfile",
			)

			return problem
		}
		if !isPresent {
			if ValidateUname(profile.Username) {
				model.Database.MustExec(`
					UPDATE userinfo
					SET username=$1
					WHERE userinfo.uid = (
						SELECT DISCTINCT session.uid from session
						JOIN userinfo
						USING(uid)
						WHERE cookie=$2
					);
				`, profile.Username, cookie)

				model.LG.Info(
					"update profile succeeded, username updated",
					"source", "database.go",
					"who", "UpdateProfile",
				)
			} else {

				model.LG.Info(
					"update profile failed, invalid username",
					"source", "database.go",
					"who", "UpdateProfile",
				)
			}
		}
		if isPresent {
			model.LG.Info(
				"update profile failed, username already in use",
				"source", "database.go",
				"who", "UpdateProfile",
			)
		}
	}

	if profile.Password != "" {
		if ValidatePassword(profile.Password) {
			hashedPsswd := misc.GeneratePasswordHash(profile.Password)
			model.Database.MustExec(`
				UPDATE userinfo
				SET password=$1
				WHERE userinfo.uid = (
					SELECT DISTINCT session.uid from session
					JOIN userinfo
					USING(uid)
					WHERE cookie=$2
				);
			`, hashedPsswd, cookie)

			model.LG.Info(
				"update profile succeeded, password updated",
				"source", "database.go",
				"who", "UpdateProfile",
			)
		} else {

			model.LG.Info(
				"update profile failed, invalid password",
				"source", "database.go",
				"who", "UpdateProfile",
			)
		}
	}

	if profile.Avatar != "" {
		model.Database.MustExec(`
			UPDATE userinfo
			SET avatar=$1
			WHERE userinfo.uid = (
				SELECT DISTINCT session.uid from session
				JOIN userinfo
				USING(uid)
				WHERE cookie=$2
			);
		`, profile.Avatar, cookie)

		model.LG.Info(
			"update profile succeeded, avatar updated",
			"source", "database.go",
			"who", "UpdateProfile",
		)
	}

	return nil
}

func (model *DatabaseModel) Logout(cookie string) bool {
	model.Database.QueryRowx(`
	DELETE
	FROM session
	WHERE cookie=$1;
	`, cookie)

	model.LG.Info(
		"logout succeeded",
		"source", "database.go",
		"who", "Logout",
	)

	return true
}

func (model *DatabaseModel) Register(credentials models.UserEdit) (string, error) {
	if !ValidateUname(credentials.Username) || !ValidatePassword(credentials.Password) {
		return "", fmt.Errorf("non-valid")
	}

	var exists string

	row := model.Database.QueryRowx(`SELECT EXISTS
									(SELECT true FROM userinfo
									WHERE username=$1);
									`, credentials.Username)
	err := row.Scan(&exists)
	err = err
	fl, _ := strconv.ParseBool(exists)

	if fl == true {
		return "", fmt.Errorf("exists")
	}

	cookie := misc.GenerateCookie()
	hashedPsswd := misc.GeneratePasswordHash(credentials.Password)

	model.Database.MustExec(`
			INSERT INTO userinfo(username, password, avatar)
			VALUES($1, $2, $3);
		`, credentials.Username, hashedPsswd, credentials.Avatar)

	model.Database.MustExec(`
		INSERT INTO session(uid, cookie)
		VALUES(
			(SELECT uid FROM userinfo WHERE username=$1),
			$2
		);
	`, credentials.Username, cookie)

	model.LG.Info(
		"signup succeeded",
		"source", "database.go",
		"who", "Register",
	)

	model.AddApp(cookie, "Terminal")
	model.AddApp(cookie, "Snake")
	return cookie, nil
}

func (model *DatabaseModel) Info(cookie string) (string, error) {
	username := ""
	if err := model.Database.QueryRow(`
		SELECT username
		FROM (
			SELECT uid
			FROM session
			WHERE cookie=$1
		) u
		INNER JOIN userinfo USING(uid);
	`, cookie).Scan(&username); err != nil {
		return "", err
	}
	return username, nil
}

/*************************************** App API ***************************************/

func (model *DatabaseModel) GetApps() (apps models.Applications) {
	rows, _ := model.Database.Queryx(`
		SELECT link,url,name,name_de,name_ru,image,about,about_de,about_ru,installs,category
		FROM app
	`)
	defer rows.Close()

	for rows.Next() {
		temp := models.Application{}
		if err := rows.Scan(&temp.Link, &temp.Url, &temp.Name, &temp.NameDE, &temp.NameRU, &temp.Image, &temp.About, &temp.AboutDE, &temp.AboutRU, &temp.Installations, &temp.Category); err != nil {

			model.LG.Info(
				"scan failed",
				"source", "database.go",
				"who", "GetApps",
			)
			return models.Applications{}
		}

		apps.Applications = append(apps.Applications, temp)
	}

	model.LG.Info(
		"GetApps succeeded",
		"source", "database.go",
		"who", "GetApps",
	)

	return apps
}

func (model *DatabaseModel) GetPopularApps() (apps models.Applications) {
	rows, _ := model.Database.Queryx(`
		SELECT link,url,name,name_de,name_ru,image,about,about_de,about_ru,installs,category
			FROM app
			ORDER BY installs DESC;
	`)
	defer rows.Close()

	for rows.Next() {
		temp := models.Application{}
		if err := rows.Scan(&temp.Link, &temp.Url, &temp.Name, &temp.NameDE, &temp.NameRU, &temp.Image, &temp.About, &temp.AboutDE, &temp.AboutRU, &temp.Installations, &temp.Category); err != nil {

			model.LG.Info(
				"scan failed",
				"source", "database.go",
				"who", "GetPopularApps",
			)

			return models.Applications{}
		}

		apps.Applications = append(apps.Applications, temp)
	}

	model.LG.Info(
		"GetPopularApps succeeded",
		"source", "database.go",
		"who", "GetPopularApps",
	)

	return apps
}

func (model *DatabaseModel) GetApp(name string) (app models.Application) {
	if isPresent, problem := model.Present("app", "name", name); isPresent && problem == nil {
		row := model.Database.QueryRowx(`
			SELECT link,url,name,name_de,name_ru,image,about,about_de,about_ru,installs,category
				FROM app
				WHERE name=$1;`, name)
		err := row.Scan(&app.Link, &app.Url, &app.Name, &app.NameDE, &app.NameRU, &app.Image, &app.About, &app.AboutDE, &app.AboutRU, &app.Installations, &app.Category)

		if err != nil {

			model.LG.Info(
				"GetApp failed, scan err: "+err.Error(),
				"source", "database.go",
				"who", "GetApp",
			)

			return models.Application{}
		}

		model.LG.Info(
			"GetApp succeeded",
			"source", "database.go",
			"who", "GetApp",
		)

		return app
	} else if problem != nil {

		model.LG.Info(
			"Present failed",
			"source", "database.go",
			"who", "GetApp",
		)

		return models.Application{}
	}

	model.LG.Info(
		"GetApp failed, app doesn't exist",
		"source", "database.go",
		"who", "GetApp",
	)

	return models.Application{}
}

func (model *DatabaseModel) GetAppPersonal(cookie string, name string) (app models.UserApplicationInstalled) {
	if isPresent, problem := model.Present("app", "name", name); isPresent && problem == nil {
		row := model.Database.QueryRowx(`
			SELECT link,url,name,name_de,name_ru,image,about,about_de,about_ru,installs,category
				FROM app
				WHERE name=$1;`, name)
		err := row.Scan(&app.Link, &app.Url, &app.Name, &app.NameDE, &app.NameRU, &app.Image, &app.About, &app.AboutDE, &app.AboutRU, &app.Installations, &app.Category)

		if err != nil {

			model.LG.Info(
				"GetApp failed, scan error: "+err.Error(),
				"source", "database.go",
				"who", "GetApp",
			)

			return models.UserApplicationInstalled{}
		}

		var exists string
		rowInst := model.Database.QueryRowx(`SELECT EXISTS
											(SELECT true
											FROM userapp
											JOIN app
											USING(appid)
											WHERE userapp.uid=(SELECT DISTINCT session.uid
												FROM session JOIN userinfo
												USING(uid)
												WHERE cookie=$1) AND name=$2);`, cookie, name)
		err = rowInst.Scan(&exists)

		if err != nil {

			model.LG.Info(
				"Scan failed",
				"source", "database.go",
				"who", "Present",
			)

			return models.UserApplicationInstalled{}
		}

		app.Installed, err = strconv.ParseBool(exists)

		if err != nil {

			model.LG.Info(
				"strconv.ParseBool failed",
				"source", "database.go",
				"who", "Present",
			)

			return models.UserApplicationInstalled{}
		}

		model.LG.Info(
			"GetApp succeeded",
			"source", "database.go",
			"who", "GetApp",
		)

		return app
	} else if problem != nil {

		model.LG.Info(
			"Present failed",
			"source", "database.go",
			"who", "GetApp",
		)

		return models.UserApplicationInstalled{}
	}

	model.LG.Info(
		"GetApp failed, app doesn't exist",
		"source", "database.go",
		"who", "GetApp",
	)

	return models.UserApplicationInstalled{}
}

func (model *DatabaseModel) GetMyApps(cookie string) (user_apps models.Applications) {
	rows, _ := model.Database.Queryx(`
	SELECT link,url,name,name_de,name_ru,image,about,about_de,about_ru,installs,category
		FROM app
		JOIN userapp
		USING(appid)
		WHERE userapp.uid=(SELECT DISTINCT session.uid
			FROM session
			JOIN userinfo
			USING(uid)
			WHERE cookie=$1);
	`, cookie)
	defer rows.Close()

	for rows.Next() {
		temp := models.Application{}
		if err := rows.Scan(&temp.Link, &temp.Url, &temp.Name, &temp.NameDE, &temp.NameRU, &temp.Image, &temp.About, &temp.AboutDE, &temp.AboutRU, &temp.Installations, &temp.Category); err != nil {

			model.LG.Info(
				"scan failed",
				"source", "database.go",
				"who", "GetMyApps",
			)

			return models.Applications{}
		}

		user_apps.Applications = append(user_apps.Applications, temp)
	}

	model.LG.Info(
		"GetMyApps succeeded",
		"source", "database.go",
		"who", "GetMyApps",
	)

	return user_apps
}

func (model *DatabaseModel) AddApp(cookie string, appname string) {

	// model.Database.MustExec(`
	// 		UPDATE app
	// 		SET installs=installs+1
	// 		WHERE name=$1;
	// 	`, appname)

	fmt.Println(cookie, appname)
	model.Database.MustExec(`
			INSERT INTO userapp(uid, appid)
			VALUES(
				(SELECT DISTINCT session.uid FROM session
				JOIN userinfo
				USING(uid)
				WHERE cookie=$1)
				,
				(SELECT appid FROM app
				WHERE name=$2)
			)
			ON CONFLICT DO NOTHING;
		`, cookie, appname)

	model.LG.Info(
		"AddApp succeeded",
		"source", "database.go",
		"who", "AddApp",
	)
	return
}

func (model *DatabaseModel) GetAppsByCattegory(category string) (apps models.Applications) {
	rows, _ := model.Database.Queryx(`
		SELECT link,url,name,name_de,name_ru,image,about,about_de,about_ru,installs,category
		FROM app
		WHERE category=$1
		ORDER BY installs DESC;
	`, category)
	defer rows.Close()

	for rows.Next() {
		temp := models.Application{}
		if err := rows.Scan(&temp.Link, &temp.Url, &temp.Name, &temp.Image, &temp.About, &temp.Installations, &temp.Category); err != nil {

			model.LG.Info(
				"scan failed",
				"source", "database.go",
				"who", "GetPopularApps",
			)

			return models.Applications{}
		}

		apps.Applications = append(apps.Applications, temp)
	}

	model.LG.Info(
		"GetPopularApps succeeded",
		"source", "database.go",
		"who", "GetPopularApps",
	)

	return apps
}
